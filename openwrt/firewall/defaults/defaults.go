package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"github.com/ORFops/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	inputAttribute            = "input"
	inputAttributeDescription = "Default policy for traffic entering the router (ACCEPT, DROP, REJECT)."
	inputUCIOption            = "input"

	outputAttribute            = "output"
	outputAttributeDescription = "Default policy for traffic leaving the router (ACCEPT, DROP, REJECT)."
	outputUCIOption            = "output"

	forwardAttribute            = "forward"
	forwardAttributeDescription = "Default policy for traffic forwarded through the router (ACCEPT, DROP, REJECT)."
	forwardUCIOption            = "forward"

	dropInvalidAttribute            = "drop_invalid"
	dropInvalidAttributeDescription = "Drop packets with an invalid conntrack state."
	dropInvalidUCIOption            = "drop_invalid"

	synfloodProtectAttribute            = "synflood_protect"
	synfloodProtectAttributeDescription = "Enable SYN flood protection."
	synfloodProtectUCIOption            = "synflood_protect"

	synfloodRateAttribute            = "synflood_rate"
	synfloodRateAttributeDescription = `Rate limit for SYN flood protection (e.g. "25/second").`
	synfloodRateUCIOption            = "synflood_rate"

	synfloodBurstAttribute            = "synflood_burst"
	synfloodBurstAttributeDescription = "Burst limit for SYN flood protection."
	synfloodBurstUCIOption            = "synflood_burst"

	tcpSynCookiesAttribute            = "tcp_syncookies"
	tcpSynCookiesAttributeDescription = "Enable TCP SYN cookie protection."
	tcpSynCookiesUCIOption            = "tcp_syncookies"

	tcpEcnAttribute            = "tcp_ecn"
	tcpEcnAttributeDescription = "Enable TCP Explicit Congestion Notification (ECN)."
	tcpEcnUCIOption            = "tcp_ecn"

	tcpWindowScalingAttribute            = "tcp_window_scaling"
	tcpWindowScalingAttributeDescription = "Enable TCP window scaling."
	tcpWindowScalingUCIOption            = "tcp_window_scaling"

	acceptRedirectsAttribute            = "accept_redirects"
	acceptRedirectsAttributeDescription = "Accept ICMP redirect messages."
	acceptRedirectsUCIOption            = "accept_redirects"

	acceptSourceRouteAttribute            = "accept_source_route"
	acceptSourceRouteAttributeDescription = "Accept source-routed packets."
	acceptSourceRouteUCIOption            = "accept_source_route"

	autoHelperAttribute            = "auto_helper"
	autoHelperAttributeDescription = "Automatically assign conntrack helpers to connections based on destination port."
	autoHelperUCIOption            = "auto_helper"

	customChainsAttribute            = "custom_chains"
	customChainsAttributeDescription = "Create per-zone custom iptables/nftables chains to allow user-defined rules."
	customChainsUCIOption            = "custom_chains"

	flowOffloadingAttribute            = "flow_offloading"
	flowOffloadingAttributeDescription = "Enable software flow offloading for forwarded traffic (fw4/nftables only)."
	flowOffloadingUCIOption            = "flow_offloading"

	flowOffloadingHwAttribute            = "flow_offloading_hw"
	flowOffloadingHwAttributeDescription = "Enable hardware flow offloading for forwarded traffic (requires flow_offloading, fw4/nftables only)."
	flowOffloadingHwUCIOption            = "flow_offloading_hw"

	schemaDescription = "Manages global firewall settings that apply to all zones and rules."

	uciConfig = "firewall"
	uciType   = "defaults"
)

var (
	policyValidators = []validator.String{
		stringvalidator.OneOf("ACCEPT", "DROP", "REJECT"),
	}

	inputSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       inputAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetInput, inputAttribute, inputUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetInput, inputAttribute, inputUCIOption),
		Validators:        policyValidators,
	}

	outputSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       outputAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetOutput, outputAttribute, outputUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetOutput, outputAttribute, outputUCIOption),
		Validators:        policyValidators,
	}

	forwardSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       forwardAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetForward, forwardAttribute, forwardUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetForward, forwardAttribute, forwardUCIOption),
		Validators:        policyValidators,
	}

	dropInvalidSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dropInvalidAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDropInvalid, dropInvalidAttribute, dropInvalidUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDropInvalid, dropInvalidAttribute, dropInvalidUCIOption),
	}

	synfloodProtectSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       synfloodProtectAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetSynfloodProtect, synfloodProtectAttribute, synfloodProtectUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetSynfloodProtect, synfloodProtectAttribute, synfloodProtectUCIOption),
	}

	synfloodRateSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       synfloodRateAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSynfloodRate, synfloodRateAttribute, synfloodRateUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSynfloodRate, synfloodRateAttribute, synfloodRateUCIOption),
	}

	synfloodBurstSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       synfloodBurstAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetSynfloodBurst, synfloodBurstAttribute, synfloodBurstUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetSynfloodBurst, synfloodBurstAttribute, synfloodBurstUCIOption),
	}

	tcpSynCookiesSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       tcpSynCookiesAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetTcpSynCookies, tcpSynCookiesAttribute, tcpSynCookiesUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetTcpSynCookies, tcpSynCookiesAttribute, tcpSynCookiesUCIOption),
	}

	tcpEcnSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       tcpEcnAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetTcpEcn, tcpEcnAttribute, tcpEcnUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetTcpEcn, tcpEcnAttribute, tcpEcnUCIOption),
	}

	tcpWindowScalingSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       tcpWindowScalingAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetTcpWindowScaling, tcpWindowScalingAttribute, tcpWindowScalingUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetTcpWindowScaling, tcpWindowScalingAttribute, tcpWindowScalingUCIOption),
	}

	acceptRedirectsSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       acceptRedirectsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetAcceptRedirects, acceptRedirectsAttribute, acceptRedirectsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetAcceptRedirects, acceptRedirectsAttribute, acceptRedirectsUCIOption),
	}

	acceptSourceRouteSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       acceptSourceRouteAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetAcceptSourceRoute, acceptSourceRouteAttribute, acceptSourceRouteUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetAcceptSourceRoute, acceptSourceRouteAttribute, acceptSourceRouteUCIOption),
	}

	autoHelperSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       autoHelperAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetAutoHelper, autoHelperAttribute, autoHelperUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetAutoHelper, autoHelperAttribute, autoHelperUCIOption),
	}

	customChainsSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       customChainsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetCustomChains, customChainsAttribute, customChainsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetCustomChains, customChainsAttribute, customChainsUCIOption),
	}

	flowOffloadingSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       flowOffloadingAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetFlowOffloading, flowOffloadingAttribute, flowOffloadingUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetFlowOffloading, flowOffloadingAttribute, flowOffloadingUCIOption),
	}

	flowOffloadingHwSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       flowOffloadingHwAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetFlowOffloadingHw, flowOffloadingHwAttribute, flowOffloadingHwUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetFlowOffloadingHw, flowOffloadingHwAttribute, flowOffloadingHwUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		acceptRedirectsAttribute:   acceptRedirectsSchemaAttribute,
		acceptSourceRouteAttribute: acceptSourceRouteSchemaAttribute,
		autoHelperAttribute:        autoHelperSchemaAttribute,
		customChainsAttribute:      customChainsSchemaAttribute,
		dropInvalidAttribute:       dropInvalidSchemaAttribute,
		flowOffloadingAttribute:    flowOffloadingSchemaAttribute,
		flowOffloadingHwAttribute:  flowOffloadingHwSchemaAttribute,
		forwardAttribute:           forwardSchemaAttribute,
		lucirpcglue.IdAttribute:    lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		inputAttribute:             inputSchemaAttribute,
		outputAttribute:            outputSchemaAttribute,
		synfloodProtectAttribute:   synfloodProtectSchemaAttribute,
		synfloodBurstAttribute:     synfloodBurstSchemaAttribute,
		synfloodRateAttribute:      synfloodRateSchemaAttribute,
		tcpEcnAttribute:            tcpEcnSchemaAttribute,
		tcpSynCookiesAttribute:     tcpSynCookiesSchemaAttribute,
		tcpWindowScalingAttribute:  tcpWindowScalingSchemaAttribute,
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
	Id                types.String `tfsdk:"id"`
	Input             types.String `tfsdk:"input"`
	Output            types.String `tfsdk:"output"`
	Forward           types.String `tfsdk:"forward"`
	DropInvalid       types.Bool   `tfsdk:"drop_invalid"`
	SynfloodProtect   types.Bool   `tfsdk:"synflood_protect"`
	SynfloodRate      types.String `tfsdk:"synflood_rate"`
	SynfloodBurst     types.Int64  `tfsdk:"synflood_burst"`
	TcpSynCookies     types.Bool   `tfsdk:"tcp_syncookies"`
	TcpEcn            types.Bool   `tfsdk:"tcp_ecn"`
	TcpWindowScaling  types.Bool   `tfsdk:"tcp_window_scaling"`
	AcceptRedirects   types.Bool   `tfsdk:"accept_redirects"`
	AcceptSourceRoute types.Bool   `tfsdk:"accept_source_route"`
	AutoHelper        types.Bool   `tfsdk:"auto_helper"`
	CustomChains      types.Bool   `tfsdk:"custom_chains"`
	FlowOffloading    types.Bool   `tfsdk:"flow_offloading"`
	FlowOffloadingHw  types.Bool   `tfsdk:"flow_offloading_hw"`
}

func modelGetId(m model) types.String                { return m.Id }
func modelGetInput(m model) types.String             { return m.Input }
func modelGetOutput(m model) types.String            { return m.Output }
func modelGetForward(m model) types.String           { return m.Forward }
func modelGetDropInvalid(m model) types.Bool         { return m.DropInvalid }
func modelGetSynfloodProtect(m model) types.Bool     { return m.SynfloodProtect }
func modelGetSynfloodRate(m model) types.String      { return m.SynfloodRate }
func modelGetSynfloodBurst(m model) types.Int64      { return m.SynfloodBurst }
func modelGetTcpSynCookies(m model) types.Bool       { return m.TcpSynCookies }
func modelGetTcpEcn(m model) types.Bool              { return m.TcpEcn }
func modelGetTcpWindowScaling(m model) types.Bool    { return m.TcpWindowScaling }
func modelGetAcceptRedirects(m model) types.Bool     { return m.AcceptRedirects }
func modelGetAcceptSourceRoute(m model) types.Bool   { return m.AcceptSourceRoute }
func modelGetAutoHelper(m model) types.Bool          { return m.AutoHelper }
func modelGetCustomChains(m model) types.Bool        { return m.CustomChains }
func modelGetFlowOffloading(m model) types.Bool      { return m.FlowOffloading }
func modelGetFlowOffloadingHw(m model) types.Bool    { return m.FlowOffloadingHw }

func modelSetId(m *model, value types.String)                { m.Id = value }
func modelSetInput(m *model, value types.String)             { m.Input = value }
func modelSetOutput(m *model, value types.String)            { m.Output = value }
func modelSetForward(m *model, value types.String)           { m.Forward = value }
func modelSetDropInvalid(m *model, value types.Bool)         { m.DropInvalid = value }
func modelSetSynfloodProtect(m *model, value types.Bool)     { m.SynfloodProtect = value }
func modelSetSynfloodRate(m *model, value types.String)      { m.SynfloodRate = value }
func modelSetSynfloodBurst(m *model, value types.Int64)      { m.SynfloodBurst = value }
func modelSetTcpSynCookies(m *model, value types.Bool)       { m.TcpSynCookies = value }
func modelSetTcpEcn(m *model, value types.Bool)              { m.TcpEcn = value }
func modelSetTcpWindowScaling(m *model, value types.Bool)    { m.TcpWindowScaling = value }
func modelSetAcceptRedirects(m *model, value types.Bool)     { m.AcceptRedirects = value }
func modelSetAcceptSourceRoute(m *model, value types.Bool)   { m.AcceptSourceRoute = value }
func modelSetAutoHelper(m *model, value types.Bool)          { m.AutoHelper = value }
func modelSetCustomChains(m *model, value types.Bool)        { m.CustomChains = value }
func modelSetFlowOffloading(m *model, value types.Bool)      { m.FlowOffloading = value }
func modelSetFlowOffloadingHw(m *model, value types.Bool)    { m.FlowOffloadingHw = value }
