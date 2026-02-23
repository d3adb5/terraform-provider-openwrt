package redirect

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
	nameAttribute            = "name"
	nameAttributeDescription = "Name of the redirect."
	nameUCIOption            = "name"

	srcAttribute            = "src"
	srcAttributeDescription = "Source zone for the redirect."
	srcUCIOption            = "src"

	srcDportAttribute            = "src_dport"
	srcDportAttributeDescription = "Match incoming destination port(s) on the source zone interface."
	srcDportUCIOption            = "src_dport"

	destAttribute            = "dest"
	destAttributeDescription = "Destination zone for the redirect."
	destUCIOption            = "dest"

	destIpAttribute            = "dest_ip"
	destIpAttributeDescription = "Destination IP address for the redirect."
	destIpUCIOption            = "dest_ip"

	destPortAttribute            = "dest_port"
	destPortAttributeDescription = "Destination port for the redirect."
	destPortUCIOption            = "dest_port"

	protoAttribute            = "proto"
	protoAttributeDescription = "Protocol of the redirect."
	protoUCIOption            = "proto"

	targetAttribute            = "target"
	targetAttributeDescription = "Target of the redirect."
	targetUCIOption            = "target"

	familyAttribute            = "family"
	familyAttributeDescription = "IP address family."
	familyUCIOption            = "family"

	schemaDescription = "Firewall redirect configuration."

	uciConfig = "firewall"
	uciType   = "redirect"
)

var (
	targetValidators = []validator.String{
		stringvalidator.OneOf("DNAT", "SNAT"),
	}

	familyValidators = []validator.String{
		stringvalidator.OneOf("ipv4", "ipv6", "any"),
	}

	nameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       nameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetName, nameAttribute, nameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetName, nameAttribute, nameUCIOption),
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
	}

	srcDportSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcDportAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcDport, srcDportAttribute, srcDportUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcDport, srcDportAttribute, srcDportUCIOption),
	}

	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
	}

	destIpSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDestIp, destIpAttribute, destIpUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDestIp, destIpAttribute, destIpUCIOption),
	}

	destPortSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDestPort, destPortAttribute, destPortUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDestPort, destPortAttribute, destPortUCIOption),
	}

	protoSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protoAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetProto, protoAttribute, protoUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetProto, protoAttribute, protoUCIOption),
	}

	targetSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       targetAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTarget, targetAttribute, targetUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTarget, targetAttribute, targetUCIOption),
		Validators:        targetValidators,
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators:        familyValidators,
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		nameAttribute:           nameSchemaAttribute,
		srcAttribute:            srcSchemaAttribute,
		srcDportAttribute:       srcDportSchemaAttribute,
		destAttribute:           destSchemaAttribute,
		destIpAttribute:         destIpSchemaAttribute,
		destPortAttribute:       destPortSchemaAttribute,
		protoAttribute:          protoSchemaAttribute,
		targetAttribute:         targetSchemaAttribute,
		familyAttribute:         familySchemaAttribute,
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
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Src      types.String `tfsdk:"src"`
	SrcDport types.String `tfsdk:"src_dport"`
	Dest     types.String `tfsdk:"dest"`
	DestIp   types.String `tfsdk:"dest_ip"`
	DestPort types.String `tfsdk:"dest_port"`
	Proto    types.String `tfsdk:"proto"`
	Target   types.String `tfsdk:"target"`
	Family   types.String `tfsdk:"family"`
}

func modelGetId(m model) types.String       { return m.Id }
func modelGetName(m model) types.String     { return m.Name }
func modelGetSrc(m model) types.String      { return m.Src }
func modelGetSrcDport(m model) types.String { return m.SrcDport }
func modelGetDest(m model) types.String     { return m.Dest }
func modelGetDestIp(m model) types.String   { return m.DestIp }
func modelGetDestPort(m model) types.String { return m.DestPort }
func modelGetProto(m model) types.String    { return m.Proto }
func modelGetTarget(m model) types.String   { return m.Target }
func modelGetFamily(m model) types.String   { return m.Family }

func modelSetId(m *model, value types.String)       { m.Id = value }
func modelSetName(m *model, value types.String)     { m.Name = value }
func modelSetSrc(m *model, value types.String)      { m.Src = value }
func modelSetSrcDport(m *model, value types.String) { m.SrcDport = value }
func modelSetDest(m *model, value types.String)     { m.Dest = value }
func modelSetDestIp(m *model, value types.String)   { m.DestIp = value }
func modelSetDestPort(m *model, value types.String) { m.DestPort = value }
func modelSetProto(m *model, value types.String)    { m.Proto = value }
func modelSetTarget(m *model, value types.String)   { m.Target = value }
func modelSetFamily(m *model, value types.String)   { m.Family = value }
