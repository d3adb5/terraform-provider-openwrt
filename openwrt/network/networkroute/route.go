package networkroute

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	disabledAttribute            = "disabled"
	disabledAttributeDescription = "Disable this route without removing it."
	disabledUCIOption            = "disabled"

	gatewayAttribute            = "gateway"
	gatewayAttributeDescription = "Next hop IP address for this route."
	gatewayUCIOption            = "gateway"

	interfaceAttribute            = "interface"
	interfaceAttributeDescription = "Logical network interface this route is associated with."
	interfaceUCIOption            = "interface"

	metricAttribute            = "metric"
	metricAttributeDescription = "Route metric (priority). Lower values are preferred."
	metricUCIOption            = "metric"

	mtuAttribute            = "mtu"
	mtuAttributeDescription = "MTU for this route."
	mtuUCIOption            = "mtu"

	netmaskAttribute            = "netmask"
	netmaskAttributeDescription = "Subnet mask for the target network (e.g. \"255.255.255.0\")."
	netmaskUCIOption            = "netmask"

	onlinkAttribute            = "onlink"
	onlinkAttributeDescription = "Treat the gateway as directly reachable on this link even if it does not match any interface prefix."
	onlinkUCIOption            = "onlink"

	sourceAttribute            = "source"
	sourceAttributeDescription = "Preferred source address when sending to this destination."
	sourceUCIOption            = "source"

	tableAttribute            = "table"
	tableAttributeDescription = "Routing table to add this route to (e.g. \"main\", \"local\", or a numeric table ID)."
	tableUCIOption            = "table"

	targetAttribute            = "target"
	targetAttributeDescription = "Destination network or host address (e.g. \"192.168.2.0\")."
	targetUCIOption            = "target"

	schemaDescription = "Manages a static IPv4 route."

	uciConfig = "network"
	uciType   = "route"
)

var (
	disabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       disabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDisabled, disabledAttribute, disabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDisabled, disabledAttribute, disabledUCIOption),
	}

	gatewaySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       gatewayAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetGateway, gatewayAttribute, gatewayUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetGateway, gatewayAttribute, gatewayUCIOption),
	}

	interfaceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetInterface, interfaceAttribute, interfaceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetInterface, interfaceAttribute, interfaceUCIOption),
	}

	metricSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       metricAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMetric, metricAttribute, metricUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMetric, metricAttribute, metricUCIOption),
	}

	mtuSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtuAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMtu, mtuAttribute, mtuUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMtu, mtuAttribute, mtuUCIOption),
	}

	netmaskSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       netmaskAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNetmask, netmaskAttribute, netmaskUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNetmask, netmaskAttribute, netmaskUCIOption),
	}

	onlinkSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       onlinkAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetOnlink, onlinkAttribute, onlinkUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetOnlink, onlinkAttribute, onlinkUCIOption),
	}

	sourceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       sourceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSource, sourceAttribute, sourceUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSource, sourceAttribute, sourceUCIOption),
	}

	tableSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       tableAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTable, tableAttribute, tableUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTable, tableAttribute, tableUCIOption),
	}

	targetSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       targetAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTarget, targetAttribute, targetUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTarget, targetAttribute, targetUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		disabledAttribute:       disabledSchemaAttribute,
		gatewayAttribute:        gatewaySchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		interfaceAttribute:      interfaceSchemaAttribute,
		metricAttribute:         metricSchemaAttribute,
		mtuAttribute:            mtuSchemaAttribute,
		netmaskAttribute:        netmaskSchemaAttribute,
		onlinkAttribute:         onlinkSchemaAttribute,
		sourceAttribute:         sourceSchemaAttribute,
		tableAttribute:          tableSchemaAttribute,
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
	Disabled  types.Bool   `tfsdk:"disabled"`
	Gateway   types.String `tfsdk:"gateway"`
	Id        types.String `tfsdk:"id"`
	Interface types.String `tfsdk:"interface"`
	Metric    types.Int64  `tfsdk:"metric"`
	Mtu       types.Int64  `tfsdk:"mtu"`
	Netmask   types.String `tfsdk:"netmask"`
	Onlink    types.Bool   `tfsdk:"onlink"`
	Source    types.String `tfsdk:"source"`
	Table     types.String `tfsdk:"table"`
	Target    types.String `tfsdk:"target"`
}

func modelGetDisabled(m model) types.Bool   { return m.Disabled }
func modelGetGateway(m model) types.String  { return m.Gateway }
func modelGetId(m model) types.String       { return m.Id }
func modelGetInterface(m model) types.String { return m.Interface }
func modelGetMetric(m model) types.Int64    { return m.Metric }
func modelGetMtu(m model) types.Int64       { return m.Mtu }
func modelGetNetmask(m model) types.String  { return m.Netmask }
func modelGetOnlink(m model) types.Bool     { return m.Onlink }
func modelGetSource(m model) types.String   { return m.Source }
func modelGetTable(m model) types.String    { return m.Table }
func modelGetTarget(m model) types.String   { return m.Target }

func modelSetDisabled(m *model, value types.Bool)   { m.Disabled = value }
func modelSetGateway(m *model, value types.String)  { m.Gateway = value }
func modelSetId(m *model, value types.String)       { m.Id = value }
func modelSetInterface(m *model, value types.String) { m.Interface = value }
func modelSetMetric(m *model, value types.Int64)    { m.Metric = value }
func modelSetMtu(m *model, value types.Int64)       { m.Mtu = value }
func modelSetNetmask(m *model, value types.String)  { m.Netmask = value }
func modelSetOnlink(m *model, value types.Bool)     { m.Onlink = value }
func modelSetSource(m *model, value types.String)   { m.Source = value }
func modelSetTable(m *model, value types.String)    { m.Table = value }
func modelSetTarget(m *model, value types.String)   { m.Target = value }
