package networkrule6

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
	actionAttribute            = "action"
	actionAttributeDescription = "Rule action when matched: blackhole (drop silently), prohibit (drop with EACCES), or unreachable (drop with ENETUNREACH)."
	actionUCIOption            = "action"

	destAttribute            = "dest"
	destAttributeDescription = "Destination IPv6 subnet to match in CIDR notation (e.g. \"2001:db8:1::/48\")."
	destUCIOption            = "dest"

	disabledAttribute            = "disabled"
	disabledAttributeDescription = "Disable this rule without removing it."
	disabledUCIOption            = "disabled"

	gotoAttribute            = "goto"
	gotoAttributeDescription = "Jump to the rule with the specified priority instead of performing an action."
	gotoUCIOption            = "goto"

	inAttribute            = "in"
	inAttributeDescription = "Incoming network interface to match."
	inUCIOption            = "in"

	invertAttribute            = "invert"
	invertAttributeDescription = "Invert the sense of the match."
	invertUCIOption            = "invert"

	lookupAttribute            = "lookup"
	lookupAttributeDescription = "Routing table to use when this rule matches (e.g. \"main\", or a numeric table ID)."
	lookupUCIOption            = "lookup"

	markAttribute            = "mark"
	markAttributeDescription = "Firewall mark to match, optionally with a mask (e.g. \"0x100\" or \"0x100/0x100\")."
	markUCIOption            = "mark"

	outAttribute            = "out"
	outAttributeDescription = "Outgoing network interface to match."
	outUCIOption            = "out"

	priorityAttribute            = "priority"
	priorityAttributeDescription = "Rule priority. Rules are evaluated in ascending order."
	priorityUCIOption            = "priority"

	srcAttribute            = "src"
	srcAttributeDescription = "Source IPv6 subnet to match in CIDR notation (e.g. \"2001:db8::/32\")."
	srcUCIOption            = "src"

	tosAttribute            = "tos"
	tosAttributeDescription = "Type of Service (TOS) value to match (hexadecimal, e.g. \"0x10\")."
	tosUCIOption            = "tos"

	schemaDescription = "Manages an IPv6 policy routing rule. Rules select which routing table to use based on packet attributes."

	uciConfig = "network"
	uciType   = "rule6"
)

var (
	actionValidators = []validator.String{
		stringvalidator.OneOf("blackhole", "prohibit", "unreachable"),
	}

	actionSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       actionAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetAction, actionAttribute, actionUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetAction, actionAttribute, actionUCIOption),
		Validators:        actionValidators,
	}

	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
	}

	disabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       disabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDisabled, disabledAttribute, disabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDisabled, disabledAttribute, disabledUCIOption),
	}

	gotoSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       gotoAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetGoto, gotoAttribute, gotoUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetGoto, gotoAttribute, gotoUCIOption),
	}

	inSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       inAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIn, inAttribute, inUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIn, inAttribute, inUCIOption),
	}

	invertSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       invertAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetInvert, invertAttribute, invertUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetInvert, invertAttribute, invertUCIOption),
	}

	lookupSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       lookupAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLookup, lookupAttribute, lookupUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLookup, lookupAttribute, lookupUCIOption),
	}

	markSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       markAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMark, markAttribute, markUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMark, markAttribute, markUCIOption),
	}

	outSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       outAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetOut, outAttribute, outUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetOut, outAttribute, outUCIOption),
	}

	prioritySchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       priorityAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetPriority, priorityAttribute, priorityUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetPriority, priorityAttribute, priorityUCIOption),
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
	}

	tosSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       tosAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTos, tosAttribute, tosUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTos, tosAttribute, tosUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		actionAttribute:         actionSchemaAttribute,
		destAttribute:           destSchemaAttribute,
		disabledAttribute:       disabledSchemaAttribute,
		gotoAttribute:           gotoSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		inAttribute:             inSchemaAttribute,
		invertAttribute:         invertSchemaAttribute,
		lookupAttribute:         lookupSchemaAttribute,
		markAttribute:           markSchemaAttribute,
		outAttribute:            outSchemaAttribute,
		priorityAttribute:       prioritySchemaAttribute,
		srcAttribute:            srcSchemaAttribute,
		tosAttribute:            tosSchemaAttribute,
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
	Action   types.String `tfsdk:"action"`
	Dest     types.String `tfsdk:"dest"`
	Disabled types.Bool   `tfsdk:"disabled"`
	Goto     types.Int64  `tfsdk:"goto"`
	Id       types.String `tfsdk:"id"`
	In       types.String `tfsdk:"in"`
	Invert   types.Bool   `tfsdk:"invert"`
	Lookup   types.String `tfsdk:"lookup"`
	Mark     types.String `tfsdk:"mark"`
	Out      types.String `tfsdk:"out"`
	Priority types.Int64  `tfsdk:"priority"`
	Src      types.String `tfsdk:"src"`
	Tos      types.String `tfsdk:"tos"`
}

func modelGetAction(m model) types.String  { return m.Action }
func modelGetDest(m model) types.String    { return m.Dest }
func modelGetDisabled(m model) types.Bool  { return m.Disabled }
func modelGetGoto(m model) types.Int64     { return m.Goto }
func modelGetId(m model) types.String      { return m.Id }
func modelGetIn(m model) types.String      { return m.In }
func modelGetInvert(m model) types.Bool    { return m.Invert }
func modelGetLookup(m model) types.String  { return m.Lookup }
func modelGetMark(m model) types.String    { return m.Mark }
func modelGetOut(m model) types.String     { return m.Out }
func modelGetPriority(m model) types.Int64 { return m.Priority }
func modelGetSrc(m model) types.String     { return m.Src }
func modelGetTos(m model) types.String     { return m.Tos }

func modelSetAction(m *model, value types.String)  { m.Action = value }
func modelSetDest(m *model, value types.String)    { m.Dest = value }
func modelSetDisabled(m *model, value types.Bool)  { m.Disabled = value }
func modelSetGoto(m *model, value types.Int64)     { m.Goto = value }
func modelSetId(m *model, value types.String)      { m.Id = value }
func modelSetIn(m *model, value types.String)      { m.In = value }
func modelSetInvert(m *model, value types.Bool)    { m.Invert = value }
func modelSetLookup(m *model, value types.String)  { m.Lookup = value }
func modelSetMark(m *model, value types.String)    { m.Mark = value }
func modelSetOut(m *model, value types.String)     { m.Out = value }
func modelSetPriority(m *model, value types.Int64) { m.Priority = value }
func modelSetSrc(m *model, value types.String)     { m.Src = value }
func modelSetTos(m *model, value types.String)     { m.Tos = value }
