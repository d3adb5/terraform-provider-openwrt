package forwarding

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	destAttribute            = "dest"
	destAttributeDescription = "Name of the destination zone."
	destUCIOption            = "dest"

	familyAttribute            = "family"
	familyAttributeDescription = "Protocol family for this forwarding rule (ipv4, ipv6, any)."
	familyUCIOption            = "family"

	srcAttribute            = "src"
	srcAttributeDescription = "Name of the source zone."
	srcUCIOption            = "src"

	schemaDescription = "Allows traffic to flow from one firewall zone to another."

	uciConfig = "firewall"
	uciType   = "forwarding"
)

var (
	familyValidators = []validator.String{
		stringvalidator.OneOf("ipv4", "ipv6", "any"),
	}

	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators:        familyValidators,
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		destAttribute:           destSchemaAttribute,
		familyAttribute:         familySchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		srcAttribute:            srcSchemaAttribute,
	}
)

func NewDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

func NewResource() resource.Resource {
	return lucirpcglue.NewResource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	Id     types.String `tfsdk:"id"`
	Src    types.String `tfsdk:"src"`
	Dest   types.String `tfsdk:"dest"`
	Family types.String `tfsdk:"family"`
}

func modelGetId(m model) types.String     { return m.Id }
func modelGetSrc(m model) types.String    { return m.Src }
func modelGetDest(m model) types.String   { return m.Dest }
func modelGetFamily(m model) types.String { return m.Family }

func modelSetId(m *model, value types.String)     { m.Id = value }
func modelSetSrc(m *model, value types.String)    { m.Src = value }
func modelSetDest(m *model, value types.String)   { m.Dest = value }
func modelSetFamily(m *model, value types.String) { m.Family = value }
