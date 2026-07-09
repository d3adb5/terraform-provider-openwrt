package redirect

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"github.com/ORFops/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	destAttribute            = "dest"
	destAttributeDescription = "Destination zone the traffic is forwarded into (e.g. \"lan\"). If omitted, matches any destination zone."
	destUCIOption            = "dest"

	destIpAttribute            = "dest_ip"
	destIpAttributeDescription = "Internal IP address to forward matching traffic to (the LAN-side server). Required for DNAT to be useful."
	destIpUCIOption            = "dest_ip"

	destPortAttribute            = "dest_port"
	destPortAttributeDescription = "Internal port to forward matching traffic to (the LAN-side port, e.g. \"80\"). If omitted, the same port as src_dport is used."
	destPortUCIOption            = "dest_port"

	enabledAttribute            = "enabled"
	enabledAttributeDescription = "Enable or disable this redirect."
	enabledUCIOption            = "enabled"

	familyAttribute            = "family"
	familyAttributeDescription = "Protocol family to match (ipv4, ipv6, any)."
	familyUCIOption            = "family"

	nameAttribute            = "name"
	nameAttributeDescription = "Name of the redirect."
	nameUCIOption            = "name"

	orderingSchemaDescription = "Relative order of firewall redirects (port forwards). OpenWrt applies redirects in the order they appear in the firewall config, and the first matching redirect wins, so overlapping redirects behave differently depending on their order."

	protoAttribute            = "proto"
	protoAttributeDescription = "List of protocols to match (e.g. [\"tcp\"], [\"udp\"], [\"tcp\", \"udp\"])."
	protoUCIOption            = "proto"

	reflectionAttribute            = "reflection"
	reflectionAttributeDescription = "Enable NAT reflection (hairpin NAT) so LAN clients can reach forwarded ports via the WAN IP."
	reflectionUCIOption            = "reflection"

	schemaDescription = "Firewall redirect (port forwarding / NAT) configuration."

	srcAttribute            = "src"
	srcAttributeDescription = "Source zone where traffic originates (e.g. \"wan\"). If omitted, matches traffic from any zone."
	srcUCIOption            = "src"

	srcDipAttribute            = "src_dip"
	srcDipAttributeDescription = "Match traffic whose pre-DNAT destination IP equals this address (i.e. the WAN/external IP the client connected to). Useful when a host has multiple WAN IPs."
	srcDipUCIOption            = "src_dip"

	srcDportAttribute            = "src_dport"
	srcDportAttributeDescription = "External port (or range) to match on the source zone interface — the port the outside client connects to (e.g. \"8080\" or \"8080-8090\")."
	srcDportUCIOption            = "src_dport"

	srcIpAttribute            = "src_ip"
	srcIpAttributeDescription = "Match traffic from this source (client) IP address or CIDR range."
	srcIpUCIOption            = "src_ip"

	srcMacAttribute            = "src_mac"
	srcMacAttributeDescription = "Match traffic from this source MAC address."
	srcMacUCIOption            = "src_mac"

	targetAttribute            = "target"
	targetAttributeDescription = "NAT action to perform (DNAT for port forwarding, SNAT for source NAT)."
	targetUCIOption            = "target"

	uciConfig = "firewall"
	uciType   = "redirect"
)

var (
	// ipCidrValidators accepts an IPv4 or IPv6 address with an optional CIDR
	// suffix, optionally prefixed with ! for negation.
	ipCidrValidators = []validator.String{
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^!?[0-9a-fA-F.:][0-9a-fA-F.:/]*$`),
			`must be an IPv4 or IPv6 address or CIDR range (e.g. "192.168.1.0/24"), optionally prefixed with ! to negate`,
		),
	}

	// macAddressValidators accepts standard colon-separated MAC addresses,
	// optionally prefixed with ! for negation.
	macAddressValidators = []validator.String{
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^!?([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`),
			`must be a MAC address in the format XX:XX:XX:XX:XX:XX, optionally prefixed with ! to negate`,
		),
	}

	// portValidators accepts a single port, a range (80-443 or 80:443), or a
	// comma-separated list thereof, each optionally prefixed with !.
	portValidators = []validator.String{
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^!?\d+([:-]\d+)?(,!?\d+([:-]\d+)?)*$`),
			`must be a port number (e.g. "80"), a range (e.g. "8080-8090"), or a comma-separated list, optionally prefixed with ! to negate`,
		),
	}

	// protoValidators accepts named protocols and numeric protocol numbers.
	protoValidators = []validator.String{
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^(tcp|udp|tcpudp|udplite|icmp|icmpv6|esp|ah|sctp|all|\d+)$`),
			`must be a protocol name (tcp, udp, tcpudp, udplite, icmp, icmpv6, esp, ah, sctp, all) or a numeric protocol number`,
		),
	}

	zoneValidators = []validator.String{
		stringvalidator.LengthAtLeast(1),
	}

	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
		Validators:        zoneValidators,
	}

	destIpSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDestIp, destIpAttribute, destIpUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDestIp, destIpAttribute, destIpUCIOption),
		Validators:        ipCidrValidators,
	}

	destPortSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDestPort, destPortAttribute, destPortUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDestPort, destPortAttribute, destPortUCIOption),
		Validators:        portValidators,
	}

	enabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnabled, enabledAttribute, enabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnabled, enabledAttribute, enabledUCIOption),
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf("ipv4", "ipv6", "any"),
		},
	}

	nameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       nameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetName, nameAttribute, nameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetName, nameAttribute, nameUCIOption),
		Validators:        zoneValidators,
	}

	protoSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protoAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetProto, protoAttribute, protoUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetProto, protoAttribute, protoUCIOption),
		Validators:        []validator.List{listvalidator.ValueStringsAre(protoValidators...)},
	}

	reflectionSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       reflectionAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetReflection, reflectionAttribute, reflectionUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetReflection, reflectionAttribute, reflectionUCIOption),
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
		Validators:        zoneValidators,
	}

	srcDipSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcDipAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcDip, srcDipAttribute, srcDipUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcDip, srcDipAttribute, srcDipUCIOption),
		Validators:        ipCidrValidators,
	}

	srcDportSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcDportAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcDport, srcDportAttribute, srcDportUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcDport, srcDportAttribute, srcDportUCIOption),
		Validators:        portValidators,
	}

	srcIpSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcIp, srcIpAttribute, srcIpUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcIp, srcIpAttribute, srcIpUCIOption),
		Validators:        ipCidrValidators,
	}

	srcMacSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcMacAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrcMac, srcMacAttribute, srcMacUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrcMac, srcMacAttribute, srcMacUCIOption),
		Validators:        macAddressValidators,
	}

	targetSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       targetAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTarget, targetAttribute, targetUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTarget, targetAttribute, targetUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf("DNAT", "SNAT"),
		},
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

func NewOrderingResource() resource.Resource {
	return lucirpcglue.NewOrderingResource(
		orderingSchemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	Dest       types.String `tfsdk:"dest"`
	DestIp     types.String `tfsdk:"dest_ip"`
	DestPort   types.String `tfsdk:"dest_port"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	Family     types.String `tfsdk:"family"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Proto      types.List   `tfsdk:"proto"`
	Reflection types.Bool   `tfsdk:"reflection"`
	Src        types.String `tfsdk:"src"`
	SrcDip     types.String `tfsdk:"src_dip"`
	SrcDport   types.String `tfsdk:"src_dport"`
	SrcIp      types.String `tfsdk:"src_ip"`
	SrcMac     types.String `tfsdk:"src_mac"`
	Target     types.String `tfsdk:"target"`
}

func modelGetDest(m model) types.String      { return m.Dest }
func modelGetDestIp(m model) types.String    { return m.DestIp }
func modelGetDestPort(m model) types.String  { return m.DestPort }
func modelGetEnabled(m model) types.Bool     { return m.Enabled }
func modelGetFamily(m model) types.String    { return m.Family }
func modelGetId(m model) types.String        { return m.Id }
func modelGetName(m model) types.String      { return m.Name }
func modelGetProto(m model) types.List       { return m.Proto }
func modelGetReflection(m model) types.Bool  { return m.Reflection }
func modelGetSrc(m model) types.String       { return m.Src }
func modelGetSrcDip(m model) types.String    { return m.SrcDip }
func modelGetSrcDport(m model) types.String  { return m.SrcDport }
func modelGetSrcIp(m model) types.String     { return m.SrcIp }
func modelGetSrcMac(m model) types.String    { return m.SrcMac }
func modelGetTarget(m model) types.String    { return m.Target }

func modelSetDest(m *model, value types.String)      { m.Dest = value }
func modelSetDestIp(m *model, value types.String)    { m.DestIp = value }
func modelSetDestPort(m *model, value types.String)  { m.DestPort = value }
func modelSetEnabled(m *model, value types.Bool)     { m.Enabled = value }
func modelSetFamily(m *model, value types.String)    { m.Family = value }
func modelSetId(m *model, value types.String)        { m.Id = value }
func modelSetName(m *model, value types.String)      { m.Name = value }
func modelSetProto(m *model, value types.List)       { m.Proto = value }
func modelSetReflection(m *model, value types.Bool)  { m.Reflection = value }
func modelSetSrc(m *model, value types.String)       { m.Src = value }
func modelSetSrcDip(m *model, value types.String)    { m.SrcDip = value }
func modelSetSrcDport(m *model, value types.String)  { m.SrcDport = value }
func modelSetSrcIp(m *model, value types.String)     { m.SrcIp = value }
func modelSetSrcMac(m *model, value types.String)    { m.SrcMac = value }
func modelSetTarget(m *model, value types.String)    { m.Target = value }
