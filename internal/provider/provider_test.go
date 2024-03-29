package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/fake"
)

// Acceptance tests don't run against onepassword, they use a fake file based storage.
// This way we test terraform from the user perspective without the need to use one password
// CLI and authentication.
//
// This is done because storage repository is abstracted and we have use a file  as if
// onepassword resources where being served.
//
// The we can test onepassword integration tests independently from terraform.
//
// To run the tests you will need to set `OP_FAKE_STORAGE_PATH` pointing to the fake storage file.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"onepasswordorg": providerserver.NewProtocol6WithError(provider.New()),
}

func testAccPreCheck(t *testing.T) {
	_ = getFakePath(t)
}

func getFakePath(t *testing.T) string {
	fakePath := os.Getenv(provider.EnvVarOpFakeStoragePath)
	if fakePath == "" {
		t.Fatalf("%q env var must be set for acceptance tests", provider.EnvVarOpFakeStoragePath)
	}

	return fakePath
}

func getFakeRepository(t *testing.T) storage.Repository {
	fakePath := getFakePath(t)
	r, err := fake.NewRepository(fakePath)
	if err != nil {
		t.Fatalf("could not get fake repository: %s", err)
	}

	return r
}
