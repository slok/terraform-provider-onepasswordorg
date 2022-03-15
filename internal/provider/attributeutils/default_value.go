package attributeutils

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type defaultValue struct {
	value attr.Value
}

func (d defaultValue) Description(ctx context.Context) string { return "" }

func (d defaultValue) MarkdownDescription(ctx context.Context) string { return "" }

func (d defaultValue) Modify(ctx context.Context, request tfsdk.ModifyAttributePlanRequest, response *tfsdk.ModifyAttributePlanResponse) {
	result, _, _ := tftypes.WalkAttributePath(request.Config.Raw, tftypes.NewAttributePathWithSteps(request.AttributePath.Steps()))

	isNull := false
	if result.(tftypes.Value).IsNull() {
		isNull = true
	} else if result.(tftypes.Value).Type().Is(tftypes.List{}) {
		if request.AttributeConfig.(types.List).Null {
			isNull = true
		}
	} else if result.(tftypes.Value).Type().Is(tftypes.Map{}) {
		if request.AttributeConfig.(types.Map).Null {
			isNull = true
		}
	} else if result.(tftypes.Value).Type().Is(tftypes.Set{}) {
		if request.AttributeConfig.(types.Set).Null {
			isNull = true
		}
	}

	if isNull {
		response.AttributePlan = d.value
	}
}

func DefaultNumber(v float64) tfsdk.AttributePlanModifier {
	return defaultValue{value: types.Number{Value: big.NewFloat(v)}}
}

func DefaultString(v string) tfsdk.AttributePlanModifier {
	return defaultValue{value: types.String{Value: v}}
}

func DefaultBool(v bool) tfsdk.AttributePlanModifier {
	return defaultValue{value: types.Bool{Value: v}}
}

func DefaultObject(t map[string]attr.Type, v map[string]attr.Value) tfsdk.AttributePlanModifier {
	return defaultValue{
		value: types.Object{
			AttrTypes: t,
			Attrs:     v,
		},
	}
}

func DefaultListOfNumbers(items ...float64) tfsdk.AttributePlanModifier {
	values := make([]attr.Value, 0)

	for _, v := range items {
		values = append(values, types.Number{Value: big.NewFloat(v)})
	}
	return defaultValue{
		value: types.List{
			ElemType: types.NumberType,
			Elems:    values,
		},
	}
}

func DefaultListOfStrings(items ...string) tfsdk.AttributePlanModifier {
	values := make([]attr.Value, 0)

	for _, v := range items {
		values = append(values, types.String{Value: v})
	}
	return defaultValue{
		value: types.List{
			ElemType: types.StringType,
			Elems:    values,
		},
	}
}
