package system

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"github.com/ORFops/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	conLogLevelAttribute            = "conloglevel"
	conLogLevelAttributeDescription = "The maximum log level for kernel messages to be logged to the console."
	conLogLevelUCIOption            = "conloglevel"

	cronLogLevelAttribute            = "cronloglevel"
	cronLogLevelAttributeDescription = "The minimum level for cron messages to be logged to syslog."
	cronLogLevelUCIOption            = "cronloglevel"

	descriptionAttribute            = "description"
	descriptionAttributeDescription = "A short description for the system."
	descriptionUCIOption            = "description"

	hostnameAttribute            = "hostname"
	hostnameAttributeDescription = "The hostname for the system."
	hostnameUCIOption            = "hostname"

	klogConLogLevelAttribute            = "klogconloglevel"
	klogConLogLevelAttributeDescription = "The maximum log level for kernel messages to be logged to the console. Only used when klogd is running."
	klogConLogLevelUCIOption            = "klogconloglevel"

	logFileAttribute            = "log_file"
	logFileAttributeDescription = "File path for local log output."
	logFileUCIOption            = "log_file"

	logHostnameAttribute            = "log_hostname"
	logHostnameAttributeDescription = "Hostname to send with syslog messages."
	logHostnameUCIOption            = "log_hostname"

	logIPAttribute            = "log_ip"
	logIPAttributeDescription = "IP address of a remote syslog server."
	logIPUCIOption            = "log_ip"

	logPortAttribute            = "log_port"
	logPortAttributeDescription = "Port of the remote syslog server (default 514)."
	logPortUCIOption            = "log_port"

	logProtoAttribute            = "log_proto"
	logProtoAttributeDescription = `Transport protocol for remote syslog. Must be one of: "tcp", "udp".`
	logProtoTCP                  = "tcp"
	logProtoUDP                  = "udp"
	logProtoUCIOption            = "log_proto"

	logRemoteAttribute            = "log_remote"
	logRemoteAttributeDescription = "Enable sending syslog to a remote server."
	logRemoteUCIOption            = "log_remote"

	logSizeAttribute            = "log_size"
	logSizeAttributeDescription = "Size of the file based log buffer in KiB."
	logSizeUCIOption            = "log_size"

	notesAttribute            = "notes"
	notesAttributeDescription = "Multi-line free-form text about the system."
	notesUCIOption            = "notes"

	schemaDescription = "Provides system data about an OpenWrt device"

	timezoneAttribute            = "timezone"
	timezoneAttributeDescription = "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua."
	timezoneUCIOption            = "timezone"

	ttyLoginAttribute            = "ttylogin"
	ttyLoginAttributeDescription = "Require authentication for local users to log in the system."
	ttyLoginUCIOption            = "ttylogin"

	uciConfig = "system"
	uciType   = "system"

	zonenameAttribute            = "zonename"
	zonenameAttributeDescription = "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua."
	zonenameUCIOption            = "zonename"
)

var (
	conLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       conLogLevelAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetConLogLevel, conLogLevelAttribute, conLogLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetConLogLevel, conLogLevelAttribute, conLogLevelUCIOption),
	}

	cronLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       cronLogLevelAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetCronLogLevel, cronLogLevelAttribute, cronLogLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetCronLogLevel, cronLogLevelAttribute, cronLogLevelUCIOption),
	}

	descriptionSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       descriptionAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDescription, descriptionAttribute, descriptionUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDescription, descriptionAttribute, descriptionUCIOption),
	}

	hostnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       hostnameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetHostname, hostnameAttribute, hostnameUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetHostname, hostnameAttribute, hostnameUCIOption),
	}

	klogConLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       klogConLogLevelAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetKLogConLogLevel, klogConLogLevelAttribute, klogConLogLevelUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetKLogConLogLevel, klogConLogLevelAttribute, klogConLogLevelUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(1, 8),
		},
	}

	logFileSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logFileAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLogFile, logFileAttribute, logFileUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLogFile, logFileAttribute, logFileUCIOption),
	}

	logHostnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logHostnameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLogHostname, logHostnameAttribute, logHostnameUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLogHostname, logHostnameAttribute, logHostnameUCIOption),
	}

	logIPSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logIPAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLogIP, logIPAttribute, logIPUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLogIP, logIPAttribute, logIPUCIOption),
	}

	logPortSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetLogPort, logPortAttribute, logPortUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetLogPort, logPortAttribute, logPortUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(1, 65535),
		},
	}

	logProtoSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logProtoAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLogProto, logProtoAttribute, logProtoUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLogProto, logProtoAttribute, logProtoUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(logProtoTCP, logProtoUDP),
		},
	}

	logRemoteSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logRemoteAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetLogRemote, logRemoteAttribute, logRemoteUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetLogRemote, logRemoteAttribute, logRemoteUCIOption),
	}

	logSizeSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logSizeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetLogSize, logSizeAttribute, logSizeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetLogSize, logSizeAttribute, logSizeUCIOption),
	}

	notesSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       notesAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNotes, notesAttribute, notesUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNotes, notesAttribute, notesUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		conLogLevelAttribute:     conLogLevelSchemaAttribute,
		cronLogLevelAttribute:    cronLogLevelSchemaAttribute,
		descriptionAttribute:     descriptionSchemaAttribute,
		hostnameAttribute:        hostnameSchemaAttribute,
		klogConLogLevelAttribute: klogConLogLevelSchemaAttribute,
		logFileAttribute:         logFileSchemaAttribute,
		logHostnameAttribute:     logHostnameSchemaAttribute,
		logIPAttribute:           logIPSchemaAttribute,
		logPortAttribute:         logPortSchemaAttribute,
		logProtoAttribute:        logProtoSchemaAttribute,
		logRemoteAttribute:       logRemoteSchemaAttribute,
		logSizeAttribute:         logSizeSchemaAttribute,
		lucirpcglue.IdAttribute:  lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		notesAttribute:           notesSchemaAttribute,
		timezoneAttribute:        timezoneSchemaAttribute,
		ttyLoginAttribute:        ttyLoginSchemaAttribute,
		zonenameAttribute:        zonenameSchemaAttribute,
	}

	timezoneSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       timezoneAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTimezone, timezoneAttribute, timezoneUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTimezone, timezoneAttribute, timezoneUCIOption),
	}

	ttyLoginSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ttyLoginAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetTTYLogin, ttyLoginAttribute, ttyLoginUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetTTYLogin, ttyLoginAttribute, ttyLoginUCIOption),
	}

	zonenameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       zonenameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetZonename, zonenameAttribute, zonenameUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetZonename, zonenameAttribute, zonenameUCIOption),
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
	ConLogLevel     types.Int64  `tfsdk:"conloglevel"`
	CronLogLevel    types.Int64  `tfsdk:"cronloglevel"`
	Description     types.String `tfsdk:"description"`
	Hostname        types.String `tfsdk:"hostname"`
	Id              types.String `tfsdk:"id"`
	KLogConLogLevel types.Int64  `tfsdk:"klogconloglevel"`
	LogFile         types.String `tfsdk:"log_file"`
	LogHostname     types.String `tfsdk:"log_hostname"`
	LogIP           types.String `tfsdk:"log_ip"`
	LogPort         types.Int64  `tfsdk:"log_port"`
	LogProto        types.String `tfsdk:"log_proto"`
	LogRemote       types.Bool   `tfsdk:"log_remote"`
	LogSize         types.Int64  `tfsdk:"log_size"`
	Notes           types.String `tfsdk:"notes"`
	Timezone        types.String `tfsdk:"timezone"`
	TTYLogin        types.Bool   `tfsdk:"ttylogin"`
	Zonename        types.String `tfsdk:"zonename"`
}

func modelGetConLogLevel(m model) types.Int64     { return m.ConLogLevel }
func modelGetCronLogLevel(m model) types.Int64    { return m.CronLogLevel }
func modelGetDescription(m model) types.String    { return m.Description }
func modelGetHostname(m model) types.String       { return m.Hostname }
func modelGetId(m model) types.String             { return m.Id }
func modelGetKLogConLogLevel(m model) types.Int64 { return m.KLogConLogLevel }
func modelGetLogFile(m model) types.String        { return m.LogFile }
func modelGetLogHostname(m model) types.String    { return m.LogHostname }
func modelGetLogIP(m model) types.String          { return m.LogIP }
func modelGetLogPort(m model) types.Int64         { return m.LogPort }
func modelGetLogProto(m model) types.String       { return m.LogProto }
func modelGetLogRemote(m model) types.Bool        { return m.LogRemote }
func modelGetLogSize(m model) types.Int64         { return m.LogSize }
func modelGetNotes(m model) types.String          { return m.Notes }
func modelGetTimezone(m model) types.String       { return m.Timezone }
func modelGetTTYLogin(m model) types.Bool         { return m.TTYLogin }
func modelGetZonename(m model) types.String       { return m.Zonename }

func modelSetConLogLevel(m *model, value types.Int64)     { m.ConLogLevel = value }
func modelSetCronLogLevel(m *model, value types.Int64)    { m.CronLogLevel = value }
func modelSetDescription(m *model, value types.String)    { m.Description = value }
func modelSetHostname(m *model, value types.String)       { m.Hostname = value }
func modelSetId(m *model, value types.String)             { m.Id = value }
func modelSetKLogConLogLevel(m *model, value types.Int64) { m.KLogConLogLevel = value }
func modelSetLogFile(m *model, value types.String)        { m.LogFile = value }
func modelSetLogHostname(m *model, value types.String)    { m.LogHostname = value }
func modelSetLogIP(m *model, value types.String)          { m.LogIP = value }
func modelSetLogPort(m *model, value types.Int64)         { m.LogPort = value }
func modelSetLogProto(m *model, value types.String)       { m.LogProto = value }
func modelSetLogRemote(m *model, value types.Bool)        { m.LogRemote = value }
func modelSetLogSize(m *model, value types.Int64)         { m.LogSize = value }
func modelSetNotes(m *model, value types.String)          { m.Notes = value }
func modelSetTimezone(m *model, value types.String)       { m.Timezone = value }
func modelSetTTYLogin(m *model, value types.Bool)         { m.TTYLogin = value }
func modelSetZonename(m *model, value types.String)       { m.Zonename = value }
