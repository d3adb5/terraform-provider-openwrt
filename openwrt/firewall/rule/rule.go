package rule

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
	nameAttributeDescription = "Name of the rule."
	nameUCIOption            = "name"

	targetAttribute            = "target"
	targetAttributeDescription = "Action to take for matching traffic (ACCEPT, REJECT, DROP, NOTRACK)."
	targetUCIOption            = "target"

	srcAttribute            = "src"
	srcAttributeDescription = "Source zone. If omitted, the rule matches traffic from any zone."
	srcUCIOption            = "src"

	srcIpAttribute            = "src_ip"
	srcIpAttributeDescription = "Match traffic from this source IP address or CIDR range."
	srcIpUCIOption            = "src_ip"

	srcMacAttribute            = "src_mac"
	srcMacAttributeDescription = "Match traffic from this source MAC address."
	srcMacUCIOption            = "src_mac"

	srcPortAttribute            = "src_port"
	srcPortAttributeDescription = "Match traffic from this source port or port range (e.g. \"80\" or \"80:443\")."
	srcPortUCIOption            = "src_port"

	destAttribute            = "dest"
	destAttributeDescription = "Destination zone. If omitted, the rule matches traffic to any destination."
	destUCIOption            = "dest"

	destIpAttribute            = "dest_ip"
	destIpAttributeDescription = "Match traffic to this destination IP address or CIDR range."
	destIpUCIOption            = "dest_ip"

	destPortAttribute            = "dest_port"
	destPortAttributeDescription = "Match traffic to this destination port or port range (e.g. \"80\" or \"80:443\")."
	destPortUCIOption            = "dest_port"

	protoAttribute            = "proto"
	protoAttributeDescription = "Match traffic using this protocol (e.g. tcp, udp, tcpudp, icmp, all)."
	protoUCIOption            = "proto"

	familyAttribute            = "family"
	familyAttributeDescription = "Protocol family to match (ipv4, ipv6, any)."
	familyUCIOption            = "family"

	enabledAttribute            = "enabled"
	enabledAttributeDescription = "Enable or disable this rule."
	enabledUCIOption            = "enabled"

	schemaDescription = "Defines a firewall rule to accept, drop, or reject traffic matching the given criteria."

	uciConfig = "firewall"
	uciType   = "rule"
)

var (
	targetValidators = []validator.String{
		stringvalidator.OneOf(
			"ACCEPT",
			"REJECT",
			"DROP",
			"NOTRACK",
		),
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

	targetSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       targetAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTarget, targetAttribute, targetUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTarget, targetAttribute, targetUCIOption),
		Validators:        targetValidators,
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
	}

	srcIpSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcIp, srcIpAttribute, srcIpUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcIp, srcIpAttribute, srcIpUCIOption),
	}

	srcMacSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcMacAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcMac, srcMacAttribute, srcMacUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcMac, srcMacAttribute, srcMacUCIOption),
	}

	srcPortSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcPort, srcPortAttribute, srcPortUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcPort, srcPortAttribute, srcPortUCIOption),
	}

	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
	}

	destIpSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDestIp, destIpAttribute, destIpUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDestIp, destIpAttribute, destIpUCIOption),
	}

	destPortSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDestPort, destPortAttribute, destPortUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDestPort, destPortAttribute, destPortUCIOption),
	}

	protoSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protoAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetProto, protoAttribute, protoUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetProto, protoAttribute, protoUCIOption),
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators:        familyValidators,
	}

	enabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnabled, enabledAttribute, enabledUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnabled, enabledAttribute, enabledUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		destAttribute:           destSchemaAttribute,
		destIpAttribute:         destIpSchemaAttribute,
		destPortAttribute:       destPortSchemaAttribute,
		enabledAttribute:        enabledSchemaAttribute,
		familyAttribute:         familySchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		nameAttribute:           nameSchemaAttribute,
		protoAttribute:          protoSchemaAttribute,
		srcAttribute:            srcSchemaAttribute,
		srcIpAttribute:          srcIpSchemaAttribute,
		srcMacAttribute:         srcMacSchemaAttribute,
		srcPortAttribute:        srcPortSchemaAttribute,
		targetAttribute:         targetSchemaAttribute,
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
	Target   types.String `tfsdk:"target"`
	Src      types.String `tfsdk:"src"`
	SrcIp    types.String `tfsdk:"src_ip"`
	SrcMac   types.String `tfsdk:"src_mac"`
	SrcPort  types.String `tfsdk:"src_port"`
	Dest     types.String `tfsdk:"dest"`
	DestIp   types.String `tfsdk:"dest_ip"`
	DestPort types.String `tfsdk:"dest_port"`
	Proto    types.String `tfsdk:"proto"`
	Family   types.String `tfsdk:"family"`
	Enabled  types.Bool   `tfsdk:"enabled"`
}

func modelGetId(m model) types.String      { return m.Id }
func modelGetName(m model) types.String    { return m.Name }
func modelGetTarget(m model) types.String  { return m.Target }
func modelGetSrc(m model) types.String     { return m.Src }
func modelGetSrcIp(m model) types.String   { return m.SrcIp }
func modelGetSrcMac(m model) types.String  { return m.SrcMac }
func modelGetSrcPort(m model) types.String { return m.SrcPort }
func modelGetDest(m model) types.String    { return m.Dest }
func modelGetDestIp(m model) types.String  { return m.DestIp }
func modelGetDestPort(m model) types.String { return m.DestPort }
func modelGetProto(m model) types.String   { return m.Proto }
func modelGetFamily(m model) types.String  { return m.Family }
func modelGetEnabled(m model) types.Bool   { return m.Enabled }

func modelSetId(m *model, value types.String)       { m.Id = value }
func modelSetName(m *model, value types.String)     { m.Name = value }
func modelSetTarget(m *model, value types.String)   { m.Target = value }
func modelSetSrc(m *model, value types.String)      { m.Src = value }
func modelSetSrcIp(m *model, value types.String)    { m.SrcIp = value }
func modelSetSrcMac(m *model, value types.String)   { m.SrcMac = value }
func modelSetSrcPort(m *model, value types.String)  { m.SrcPort = value }
func modelSetDest(m *model, value types.String)     { m.Dest = value }
func modelSetDestIp(m *model, value types.String)   { m.DestIp = value }
func modelSetDestPort(m *model, value types.String) { m.DestPort = value }
func modelSetProto(m *model, value types.String)    { m.Proto = value }
func modelSetFamily(m *model, value types.String)   { m.Family = value }
func modelSetEnabled(m *model, value types.Bool)    { m.Enabled = value }
