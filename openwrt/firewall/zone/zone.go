package zone

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
	nameAttributeDescription = "The name of the zone."
	nameUCIOption            = "name"

	forwardAttribute            = "forward"
	forwardAttributeDescription = "Default policy for forwarded traffic in this zone (ACCEPT, REJECT, DROP)."
	forwardUCIOption            = "forward"

	inputAttribute            = "input"
	inputAttributeDescription = "Default policy for traffic entering this zone (ACCEPT, REJECT, DROP)."
	inputUCIOption            = "input"

	outputAttribute            = "output"
	outputAttributeDescription = "Default policy for traffic leaving this zone (ACCEPT, REJECT, DROP)."
	outputUCIOption            = "output"

	networkAttribute            = "network"
	networkAttributeDescription = "List of interfaces or networks belonging to this zone."
	networkUCIOption            = "network"

	masqAttribute            = "masq"
	masqAttributeDescription = "Enable masquerading (NAT) for outbound traffic from this zone."
	masqUCIOption            = "masq"

	masqSrcAttribute            = "masq_src"
	masqSrcAttributeDescription = "Restrict masquerading to traffic originating from these source prefixes."
	masqSrcUCIOption            = "masq_src"

	masqDestAttribute            = "masq_dest"
	masqDestAttributeDescription = "Restrict masquerading to traffic destined for these prefixes."
	masqDestUCIOption            = "masq_dest"

	masqAllowInvalidAttribute            = "masq_allow_invalid"
	masqAllowInvalidAttributeDescription = "Allow masquerading of packets with an invalid conntrack state."
	masqAllowInvalidUCIOption            = "masq_allow_invalid"

	mtuFixAttribute            = "mtu_fix"
	mtuFixAttributeDescription = "Enable MSS clamping for outgoing TCP connections in this zone."
	mtuFixUCIOption            = "mtu_fix"

	logAttribute            = "log"
	logAttributeDescription = "Enable logging for packets rejected or dropped by this zone's default policy."
	logUCIOption            = "log"

	logLimitAttribute            = "log_limit"
	logLimitAttributeDescription = `Rate limit for log messages, e.g. "10/minute".`
	logLimitUCIOption            = "log_limit"

	familyAttribute            = "family"
	familyAttributeDescription = "Protocol family to apply the zone to (ipv4, ipv6, any)."
	familyUCIOption            = "family"

	conntrackAttribute            = "conntrack"
	conntrackAttributeDescription = "Force connection tracking for all traffic in this zone."
	conntrackUCIOption            = "conntrack"

	autoHelperAttribute            = "auto_helper"
	autoHelperAttributeDescription = "Automatically assign conntrack helpers based on destination port."
	autoHelperUCIOption            = "auto_helper"

	schemaDescription = "Manages a firewall zone. Zones group network interfaces and define default policies for traffic."

	uciConfig = "firewall"
	uciType   = "zone"

	typeAccept = "ACCEPT"
	typeReject = "REJECT"
	typeDrop   = "DROP"
)

var (
	TypeValidators = []validator.String{
		stringvalidator.OneOf(
			typeAccept,
			typeReject,
			typeDrop,
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

	forwardSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       forwardAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetForward, forwardAttribute, forwardUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetForward, forwardAttribute, forwardUCIOption),
		Validators:        TypeValidators,
	}

	inputSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       inputAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetInput, inputAttribute, inputUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetInput, inputAttribute, inputUCIOption),
		Validators:        TypeValidators,
	}

	outputSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       outputAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetOutput, outputAttribute, outputUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetOutput, outputAttribute, outputUCIOption),
		Validators:        TypeValidators,
	}

	networkSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       networkAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetNetwork, networkAttribute, networkUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetNetwork, networkAttribute, networkUCIOption),
	}

	masqSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       masqAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetMasq, masqAttribute, masqUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetMasq, masqAttribute, masqUCIOption),
	}

	masqSrcSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       masqSrcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetMasqSrc, masqSrcAttribute, masqSrcUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetMasqSrc, masqSrcAttribute, masqSrcUCIOption),
	}

	masqDestSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       masqDestAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetMasqDest, masqDestAttribute, masqDestUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetMasqDest, masqDestAttribute, masqDestUCIOption),
	}

	masqAllowInvalidSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       masqAllowInvalidAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetMasqAllowInvalid, masqAllowInvalidAttribute, masqAllowInvalidUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetMasqAllowInvalid, masqAllowInvalidAttribute, masqAllowInvalidUCIOption),
	}

	mtuFixSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtuFixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetMtuFix, mtuFixAttribute, mtuFixUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetMtuFix, mtuFixAttribute, mtuFixUCIOption),
	}

	logSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetLog, logAttribute, logUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetLog, logAttribute, logUCIOption),
	}

	logLimitSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logLimitAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLogLimit, logLimitAttribute, logLimitUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLogLimit, logLimitAttribute, logLimitUCIOption),
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators:        familyValidators,
	}

	conntrackSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       conntrackAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetConntrack, conntrackAttribute, conntrackUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetConntrack, conntrackAttribute, conntrackUCIOption),
	}

	autoHelperSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       autoHelperAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetAutoHelper, autoHelperAttribute, autoHelperUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetAutoHelper, autoHelperAttribute, autoHelperUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		autoHelperAttribute:       autoHelperSchemaAttribute,
		conntrackAttribute:        conntrackSchemaAttribute,
		familyAttribute:           familySchemaAttribute,
		forwardAttribute:          forwardSchemaAttribute,
		lucirpcglue.IdAttribute:   lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		inputAttribute:            inputSchemaAttribute,
		logAttribute:              logSchemaAttribute,
		logLimitAttribute:         logLimitSchemaAttribute,
		masqAttribute:             masqSchemaAttribute,
		masqAllowInvalidAttribute: masqAllowInvalidSchemaAttribute,
		masqDestAttribute:         masqDestSchemaAttribute,
		masqSrcAttribute:          masqSrcSchemaAttribute,
		mtuFixAttribute:           mtuFixSchemaAttribute,
		nameAttribute:             nameSchemaAttribute,
		networkAttribute:          networkSchemaAttribute,
		outputAttribute:           outputSchemaAttribute,
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
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Forward          types.String `tfsdk:"forward"`
	Input            types.String `tfsdk:"input"`
	Output           types.String `tfsdk:"output"`
	Network          types.List   `tfsdk:"network"`
	Masq             types.Bool   `tfsdk:"masq"`
	MasqSrc          types.List   `tfsdk:"masq_src"`
	MasqDest         types.List   `tfsdk:"masq_dest"`
	MasqAllowInvalid types.Bool   `tfsdk:"masq_allow_invalid"`
	MtuFix           types.Bool   `tfsdk:"mtu_fix"`
	Log              types.Bool   `tfsdk:"log"`
	LogLimit         types.String `tfsdk:"log_limit"`
	Family           types.String `tfsdk:"family"`
	Conntrack        types.Bool   `tfsdk:"conntrack"`
	AutoHelper       types.Bool   `tfsdk:"auto_helper"`
}

func modelGetId(m model) types.String             { return m.Id }
func modelGetName(m model) types.String           { return m.Name }
func modelGetForward(m model) types.String        { return m.Forward }
func modelGetInput(m model) types.String          { return m.Input }
func modelGetOutput(m model) types.String         { return m.Output }
func modelGetNetwork(m model) types.List          { return m.Network }
func modelGetMasq(m model) types.Bool             { return m.Masq }
func modelGetMasqSrc(m model) types.List          { return m.MasqSrc }
func modelGetMasqDest(m model) types.List         { return m.MasqDest }
func modelGetMasqAllowInvalid(m model) types.Bool { return m.MasqAllowInvalid }
func modelGetMtuFix(m model) types.Bool           { return m.MtuFix }
func modelGetLog(m model) types.Bool              { return m.Log }
func modelGetLogLimit(m model) types.String       { return m.LogLimit }
func modelGetFamily(m model) types.String         { return m.Family }
func modelGetConntrack(m model) types.Bool        { return m.Conntrack }
func modelGetAutoHelper(m model) types.Bool       { return m.AutoHelper }

func modelSetId(m *model, value types.String)             { m.Id = value }
func modelSetName(m *model, value types.String)           { m.Name = value }
func modelSetForward(m *model, value types.String)        { m.Forward = value }
func modelSetInput(m *model, value types.String)          { m.Input = value }
func modelSetOutput(m *model, value types.String)         { m.Output = value }
func modelSetNetwork(m *model, value types.List)          { m.Network = value }
func modelSetMasq(m *model, value types.Bool)             { m.Masq = value }
func modelSetMasqSrc(m *model, value types.List)          { m.MasqSrc = value }
func modelSetMasqDest(m *model, value types.List)         { m.MasqDest = value }
func modelSetMasqAllowInvalid(m *model, value types.Bool) { m.MasqAllowInvalid = value }
func modelSetMtuFix(m *model, value types.Bool)           { m.MtuFix = value }
func modelSetLog(m *model, value types.Bool)              { m.Log = value }
func modelSetLogLimit(m *model, value types.String)       { m.LogLimit = value }
func modelSetFamily(m *model, value types.String)         { m.Family = value }
func modelSetConntrack(m *model, value types.Bool)        { m.Conntrack = value }
func modelSetAutoHelper(m *model, value types.Bool)       { m.AutoHelper = value }
