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

	srcIpAttribute            = "src_ip"
	srcIpAttributeDescription = "Match traffic from this source IP address."
	srcIpUCIOption            = "src_ip"

	srcDipAttribute            = "src_dip"
	srcDipAttributeDescription = "Match traffic whose original destination IP matches this address (before DNAT)."
	srcDipUCIOption            = "src_dip"

	srcDportAttribute            = "src_dport"
	srcDportAttributeDescription = "Match incoming destination port(s) on the source zone interface."
	srcDportUCIOption            = "src_dport"

	srcMacAttribute            = "src_mac"
	srcMacAttributeDescription = "Match traffic from this source MAC address."
	srcMacUCIOption            = "src_mac"

	destAttribute            = "dest"
	destAttributeDescription = "Destination zone for the redirect."
	destUCIOption            = "dest"

	destIpAttribute            = "dest_ip"
	destIpAttributeDescription = "Destination IP address to redirect matching traffic to."
	destIpUCIOption            = "dest_ip"

	destPortAttribute            = "dest_port"
	destPortAttributeDescription = "Destination port to redirect matching traffic to."
	destPortUCIOption            = "dest_port"

	protoAttribute            = "proto"
	protoAttributeDescription = "Protocol to match (e.g. tcp, udp, tcpudp)."
	protoUCIOption            = "proto"

	targetAttribute            = "target"
	targetAttributeDescription = "NAT action to perform (DNAT for port forwarding, SNAT for source NAT)."
	targetUCIOption            = "target"

	familyAttribute            = "family"
	familyAttributeDescription = "Protocol family to match (ipv4, ipv6, any)."
	familyUCIOption            = "family"

	enabledAttribute            = "enabled"
	enabledAttributeDescription = "Enable or disable this redirect."
	enabledUCIOption            = "enabled"

	reflectionAttribute            = "reflection"
	reflectionAttributeDescription = "Enable NAT reflection (hairpin NAT) so LAN clients can reach forwarded ports via the WAN IP."
	reflectionUCIOption            = "reflection"

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

	srcIpSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcIp, srcIpAttribute, srcIpUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcIp, srcIpAttribute, srcIpUCIOption),
	}

	srcDipSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcDipAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcDip, srcDipAttribute, srcDipUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcDip, srcDipAttribute, srcDipUCIOption),
	}

	srcDportSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcDportAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcDport, srcDportAttribute, srcDportUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcDport, srcDportAttribute, srcDportUCIOption),
	}

	srcMacSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcMacAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcMac, srcMacAttribute, srcMacUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcMac, srcMacAttribute, srcMacUCIOption),
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

	reflectionSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       reflectionAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetReflection, reflectionAttribute, reflectionUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetReflection, reflectionAttribute, reflectionUCIOption),
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
		reflectionAttribute:     reflectionSchemaAttribute,
		srcAttribute:            srcSchemaAttribute,
		srcDipAttribute:         srcDipSchemaAttribute,
		srcDportAttribute:       srcDportSchemaAttribute,
		srcIpAttribute:          srcIpSchemaAttribute,
		srcMacAttribute:         srcMacSchemaAttribute,
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
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Src        types.String `tfsdk:"src"`
	SrcIp      types.String `tfsdk:"src_ip"`
	SrcDip     types.String `tfsdk:"src_dip"`
	SrcDport   types.String `tfsdk:"src_dport"`
	SrcMac     types.String `tfsdk:"src_mac"`
	Dest       types.String `tfsdk:"dest"`
	DestIp     types.String `tfsdk:"dest_ip"`
	DestPort   types.String `tfsdk:"dest_port"`
	Proto      types.String `tfsdk:"proto"`
	Target     types.String `tfsdk:"target"`
	Family     types.String `tfsdk:"family"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	Reflection types.Bool   `tfsdk:"reflection"`
}

func modelGetId(m model) types.String         { return m.Id }
func modelGetName(m model) types.String       { return m.Name }
func modelGetSrc(m model) types.String        { return m.Src }
func modelGetSrcIp(m model) types.String      { return m.SrcIp }
func modelGetSrcDip(m model) types.String     { return m.SrcDip }
func modelGetSrcDport(m model) types.String   { return m.SrcDport }
func modelGetSrcMac(m model) types.String     { return m.SrcMac }
func modelGetDest(m model) types.String       { return m.Dest }
func modelGetDestIp(m model) types.String     { return m.DestIp }
func modelGetDestPort(m model) types.String   { return m.DestPort }
func modelGetProto(m model) types.String      { return m.Proto }
func modelGetTarget(m model) types.String     { return m.Target }
func modelGetFamily(m model) types.String     { return m.Family }
func modelGetEnabled(m model) types.Bool      { return m.Enabled }
func modelGetReflection(m model) types.Bool   { return m.Reflection }

func modelSetId(m *model, value types.String)         { m.Id = value }
func modelSetName(m *model, value types.String)       { m.Name = value }
func modelSetSrc(m *model, value types.String)        { m.Src = value }
func modelSetSrcIp(m *model, value types.String)      { m.SrcIp = value }
func modelSetSrcDip(m *model, value types.String)     { m.SrcDip = value }
func modelSetSrcDport(m *model, value types.String)   { m.SrcDport = value }
func modelSetSrcMac(m *model, value types.String)     { m.SrcMac = value }
func modelSetDest(m *model, value types.String)       { m.Dest = value }
func modelSetDestIp(m *model, value types.String)     { m.DestIp = value }
func modelSetDestPort(m *model, value types.String)   { m.DestPort = value }
func modelSetProto(m *model, value types.String)      { m.Proto = value }
func modelSetTarget(m *model, value types.String)     { m.Target = value }
func modelSetFamily(m *model, value types.String)     { m.Family = value }
func modelSetEnabled(m *model, value types.Bool)      { m.Enabled = value }
func modelSetReflection(m *model, value types.Bool)   { m.Reflection = value }
