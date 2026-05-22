package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

// validServiceAccountVaultPermissions is the closed set of permission strings
// the op CLI accepts in `op service-account create --vault NAME:perm,perm`.
var validServiceAccountVaultPermissions = []string{
	string(model.ServiceAccountVaultPermissionReadItems),
	string(model.ServiceAccountVaultPermissionWriteItems),
	string(model.ServiceAccountVaultPermissionShareItems),
}

var (
	_ resource.Resource                   = &serviceAccountResource{}
	_ resource.ResourceWithConfigure      = &serviceAccountResource{}
	_ resource.ResourceWithImportState    = &serviceAccountResource{}
	_ resource.ResourceWithValidateConfig = &serviceAccountResource{}
)

func NewServiceAccountResource() resource.Resource {
	return &serviceAccountResource{}
}

type serviceAccountResource struct {
	repo storage.Repository
}

func (r *serviceAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account"
}

func (r *serviceAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Service Account resource.

A 1Password service account is a non-human account used for automated access to 1Password. Upon creation,
a bearer token is generated once and stored in Terraform state. This token cannot be retrieved again after
creation, so it must be captured from state immediately after the resource is created.

~> **Note:** Changing the name of a service account forces the creation of a new resource because 1Password
does not support renaming service accounts via the CLI.

~> **Note:** The ` + "`vault_access`" + ` map is configurable only at service account creation time; the
1Password CLI does not expose a command to add, remove, or modify a service account's vault access after the
fact. Any change to ` + "`vault_access`" + ` therefore destroys and recreates the service account, which
generates a **new bearer token**. Downstream consumers of the token must be updated accordingly.

~> **Note:** The token attribute is sensitive and stored in Terraform state. Protect your state file accordingly.
If you import an existing service account, the token will be empty because it cannot be retrieved.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the service account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the service account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The bearer token for the service account. Only available at creation time; preserved in state across subsequent plans.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_access": schema.MapAttribute{
				Optional: true,
				Description: fmt.Sprintf(
					"Vaults this service account can access, configured at creation time. "+
						"Keys are vault names or UUIDs (as accepted by `op service-account create --vault`); "+
						"values are the set of permissions to grant on that vault. "+
						"Allowed permissions: %s. `write_items` and `share_items` each require `read_items`. "+
						"Vault identifiers must not contain `:` or `,` (used as delimiters by the CLI). "+
						"Any change to this attribute forces resource replacement and rotates the bearer token.",
					strings.Join(validServiceAccountVaultPermissions, ", ")),
				ElementType: types.SetType{ElemType: types.StringType},
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// ValidateConfig enforces vault-identifier escaping and the
// write_items/share_items → read_items dependency rule before plan/apply.
func (r *serviceAccountResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg ServiceAccount
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.VaultAccess.IsNull() || cfg.VaultAccess.IsUnknown() {
		return
	}

	vaultAccess, err := tfMapToVaultAccess(ctx, cfg.VaultAccess)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("vault_access"), "Invalid vault_access", err.Error())
		return
	}

	for vault, perms := range vaultAccess {
		if strings.ContainsAny(vault, ":,") {
			resp.Diagnostics.AddAttributeError(
				path.Root("vault_access"),
				"Invalid vault identifier",
				fmt.Sprintf("vault identifier %q must not contain ':' or ',' (used as delimiters by `op service-account create --vault`)", vault),
			)
		}

		seen := map[string]bool{}
		for _, p := range perms {
			ps := string(p)
			if !contains(validServiceAccountVaultPermissions, ps) {
				resp.Diagnostics.AddAttributeError(
					path.Root("vault_access"),
					"Invalid service account vault permission",
					fmt.Sprintf("permission %q on vault %q is not one of: %s", ps, vault, strings.Join(validServiceAccountVaultPermissions, ", ")),
				)
			}
			seen[ps] = true
		}
		if (seen[string(model.ServiceAccountVaultPermissionWriteItems)] ||
			seen[string(model.ServiceAccountVaultPermissionShareItems)]) &&
			!seen[string(model.ServiceAccountVaultPermissionReadItems)] {
			resp.Diagnostics.AddAttributeError(
				path.Root("vault_access"),
				"Missing required service account vault permission",
				fmt.Sprintf("vault %q has write_items or share_items but is missing read_items; both write_items and share_items require read_items", vault),
			)
		}
	}
}

func (r *serviceAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r *serviceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var tfSA ServiceAccount
	diags := req.Plan.Get(ctx, &tfSA)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sa, err := mapTfToModelServiceAccount(ctx, tfSA)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping service account", err.Error())
		return
	}
	newSA, err := r.repo.CreateServiceAccount(ctx, *sa)
	if err != nil {
		resp.Diagnostics.AddError("Error creating service account", "Could not create service account, unexpected error: "+err.Error())
		return
	}

	newTfSA, err := mapModelToTfServiceAccount(*newSA)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping service account", err.Error())
		return
	}

	diags = resp.State.Set(ctx, newTfSA)
	resp.Diagnostics.Append(diags...)
}

func (r *serviceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state (including the token we stored at creation).
	var state ServiceAccount
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		// Corrupted state (e.g. from a failed prior apply) — remove so Terraform can recreate.
		resp.State.RemoveResource(ctx)
		return
	}
	sa, err := r.repo.GetServiceAccountByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Service account not found remotely, removing from state", map[string]interface{}{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading service account", fmt.Sprintf("Could not get service account %q, unexpected error: %s", id, err.Error()))
		return
	}

	// The API never returns the token after creation, so preserve it from state.
	// op CLI v2.34+ also no longer returns the name on read; preserve it from state too.
	// vault_access is also unreadable post-creation (the CLI has no command to list
	// a service account's vault grants), so preserve from state and rely on
	// RequiresReplace to surface drift via config changes.
	readTfSA, err := mapModelToTfServiceAccount(*sa)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping service account", err.Error())
		return
	}
	readTfSA.Token = state.Token
	readTfSA.Name = state.Name
	readTfSA.VaultAccess = state.VaultAccess

	diags = resp.State.Set(ctx, readTfSA)
	resp.Diagnostics.Append(diags...)
}

// Update is intentionally a no-op: all mutable attributes use RequiresReplace,
// so any change triggers a destroy+create cycle instead.
func (r *serviceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Service accounts cannot be updated", "All changes to service accounts require replacement. This is an internal error.")
}

func (r *serviceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var tfSA ServiceAccount
	diags := req.State.Get(ctx, &tfSA)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := tfSA.ID.ValueString()
	err := r.repo.DeleteServiceAccount(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Service account already deleted remotely, removing from state", map[string]interface{}{"id": id})
		} else {
			resp.Diagnostics.AddError("Error deleting service account", fmt.Sprintf("Could not delete service account %q, unexpected error: %s", id, err.Error()))
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

func (r *serviceAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// tfMapToVaultAccess converts a Terraform Map[String, Set[String]] into the
// model representation. Permissions in each set are sorted for deterministic
// downstream command construction.
func tfMapToVaultAccess(ctx context.Context, m types.Map) (map[string][]model.ServiceAccountVaultPermission, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}

	raw := map[string]types.Set{}
	if diags := m.ElementsAs(ctx, &raw, false); diags.HasError() {
		return nil, fmt.Errorf("could not decode vault_access: %s", diags)
	}

	out := make(map[string][]model.ServiceAccountVaultPermission, len(raw))
	for vault, set := range raw {
		var perms []string
		if diags := set.ElementsAs(ctx, &perms, false); diags.HasError() {
			return nil, fmt.Errorf("could not decode vault_access[%q] permissions: %s", vault, diags)
		}
		sort.Strings(perms)
		modelPerms := make([]model.ServiceAccountVaultPermission, 0, len(perms))
		for _, p := range perms {
			modelPerms = append(modelPerms, model.ServiceAccountVaultPermission(p))
		}
		out[vault] = modelPerms
	}
	return out, nil
}

// vaultAccessToTfMap is the inverse of tfMapToVaultAccess.
func vaultAccessToTfMap(va map[string][]model.ServiceAccountVaultPermission) (types.Map, error) {
	if len(va) == 0 {
		return types.MapNull(types.SetType{ElemType: types.StringType}), nil
	}

	elements := make(map[string]attr.Value, len(va))
	for vault, perms := range va {
		strVals := make([]attr.Value, 0, len(perms))
		for _, p := range perms {
			strVals = append(strVals, types.StringValue(string(p)))
		}
		set, diags := types.SetValue(types.StringType, strVals)
		if diags.HasError() {
			return types.Map{}, fmt.Errorf("could not encode permissions for vault %q: %s", vault, diags)
		}
		elements[vault] = set
	}
	m, diags := types.MapValue(types.SetType{ElemType: types.StringType}, elements)
	if diags.HasError() {
		return types.Map{}, fmt.Errorf("could not encode vault_access: %s", diags)
	}
	return m, nil
}

func mapTfToModelServiceAccount(ctx context.Context, sa ServiceAccount) (*model.ServiceAccount, error) {
	va, err := tfMapToVaultAccess(ctx, sa.VaultAccess)
	if err != nil {
		return nil, err
	}
	return &model.ServiceAccount{
		ID:          sa.ID.ValueString(),
		Name:        sa.Name.ValueString(),
		Token:       sa.Token.ValueString(),
		VaultAccess: va,
	}, nil
}

func mapModelToTfServiceAccount(sa model.ServiceAccount) (ServiceAccount, error) {
	va, err := vaultAccessToTfMap(sa.VaultAccess)
	if err != nil {
		return ServiceAccount{}, err
	}
	return ServiceAccount{
		ID:          types.StringValue(sa.ID),
		Name:        types.StringValue(sa.Name),
		Token:       types.StringValue(sa.Token),
		VaultAccess: va,
	}, nil
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// isNotFoundError returns true if the error indicates the remote resource was
// not found (HTTP 404 / op CLI exit status 4). This allows Read and Delete to
// handle already-removed resources gracefully.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Not Found") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "isn't a service account") ||
		strings.Contains(msg, "isn't a user")
}
