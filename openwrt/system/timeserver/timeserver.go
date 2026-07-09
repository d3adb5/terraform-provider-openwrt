package timeserver

import (
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
	enabledAttribute            = "enabled"
	enabledAttributeDescription = "Enable the NTP client."
	enabledUCIOption            = "enabled"

	enableServerAttribute            = "enable_server"
	enableServerAttributeDescription = "Enable the local NTP server, making this device available as an NTP source for other hosts."
	enableServerUCIOption            = "enable_server"

	schemaDescription = "NTP client and server configuration."

	serverAttribute            = "server"
	serverAttributeDescription = "List of NTP server hostnames or addresses."
	serverUCIOption            = "server"

	uciConfig = "system"
	uciType   = "timeserver"
)

var (
	enabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnabled, enabledAttribute, enabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnabled, enabledAttribute, enabledUCIOption),
	}

	enableServerSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enableServerAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnableServer, enableServerAttribute, enableServerUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnableServer, enableServerAttribute, enableServerUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		enabledAttribute:        enabledSchemaAttribute,
		enableServerAttribute:   enableServerSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		serverAttribute:         serverSchemaAttribute,
	}

	serverSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       serverAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetServer, serverAttribute, serverUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetServer, serverAttribute, serverUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.ValueStringsAre(
				stringvalidator.LengthAtLeast(1),
			),
		},
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
	Enabled      types.Bool   `tfsdk:"enabled"`
	EnableServer types.Bool   `tfsdk:"enable_server"`
	Id           types.String `tfsdk:"id"`
	Server       types.List   `tfsdk:"server"`
}

func modelGetEnabled(m model) types.Bool      { return m.Enabled }
func modelGetEnableServer(m model) types.Bool { return m.EnableServer }
func modelGetId(m model) types.String         { return m.Id }
func modelGetServer(m model) types.List       { return m.Server }

func modelSetEnabled(m *model, value types.Bool)      { m.Enabled = value }
func modelSetEnableServer(m *model, value types.Bool) { m.EnableServer = value }
func modelSetId(m *model, value types.String)         { m.Id = value }
func modelSetServer(m *model, value types.List)       { m.Server = value }
