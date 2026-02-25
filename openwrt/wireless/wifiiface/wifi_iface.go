package wifiiface

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	disabledAttribute            = "disabled"
	disabledAttributeDescription = "Disables this wireless interface without removing it."
	disabledUCIOption            = "disabled"

	deviceAttribute            = "device"
	deviceAttributeDescription = "Name of the physical device. This name is what the device is known as in LuCI/UCI, or the `id` field in Terraform."
	deviceUCIOption            = "device"

	encryptionMethodAttribute            = "encryption"
	encryptionMethodAttributeDescription = `Encryption method. Currently, only PSK encryption methods are supported. Must be one of: "none", "psk", "psk2", "psk2+aes", "psk2+ccmp", "psk2+tkip", "psk2+tkip+aes", "psk2+tkip+ccmp", "psk+aes", "psk+ccmp", "psk-mixed", "psk-mixed+aes", "psk-mixed+ccmp", "psk-mixed+tkip", "psk-mixed+tkip+aes", "psk-mixed+tkip+ccmp", "psk+tkip", "psk+tkip+aes", "psk+tkip+ccmp", "sae", "sae-mixed".`
	encryptionMethodNone                 = "none"
	encryptionMethodPSK                  = "psk"
	encryptionMethodPSK2                 = "psk2"
	encryptionMethodPSK2AES              = "psk2+aes"
	encryptionMethodPSK2CCMP             = "psk2+ccmp"
	encryptionMethodPSK2TKIP             = "psk2+tkip"
	encryptionMethodPSK2TKIPAES          = "psk2+tkip+aes"
	encryptionMethodPSK2TKIPCCMP         = "psk2+tkip+ccmp"
	encryptionMethodPSKAES               = "psk+aes"
	encryptionMethodPSKCCMP              = "psk+ccmp"
	encryptionMethodPSKMixed             = "psk-mixed"
	encryptionMethodPSKMixedAES          = "psk-mixed+aes"
	encryptionMethodPSKMixedCCMP         = "psk-mixed+ccmp"
	encryptionMethodPSKMixedTKIP         = "psk-mixed+tkip"
	encryptionMethodPSKMixedTKIPAES      = "psk-mixed+tkip+aes"
	encryptionMethodPSKMixedTKIPCCMP     = "psk-mixed+tkip+ccmp"
	encryptionMethodPSKTKIP              = "psk+tkip"
	encryptionMethodPSKTKIPAES           = "psk+tkip+aes"
	encryptionMethodPSKTKIPCCMP          = "psk+tkip+ccmp"
	encryptionMethodSAE                  = "sae"
	encryptionMethodSAEMixed             = "sae-mixed"
	encryptionMethodUCIOption            = "encryption"

	hiddenAttribute            = "hidden"
	hiddenAttributeDescription = "Suppress SSID broadcast (hidden network)."
	hiddenUCIOption            = "hidden"

	ieee80211rAttribute            = "ieee80211r"
	ieee80211rAttributeDescription = "Enable 802.11r fast BSS transition (roaming)."
	ieee80211rUCIOption            = "ieee80211r"

	ieee80211wAttribute            = "ieee80211w"
	ieee80211wAttributeDescription = `Management frame protection (802.11w). Must be one of: 0 = disabled, 1 = optional, 2 = required.`
	ieee80211wUCIOption            = "ieee80211w"

	isolateClientsAttribute            = "isolate"
	isolateClientsAttributeDescription = "Isolate wireless clients from each other."
	isolateClientsUCIOption            = "isolate"

	keyAttribute            = "key"
	keyAttributeDescription = "The pre-shared passphrase from which the pre-shared key will be derived. The clear text key has to be 8-63 characters long."
	keyUCIOption            = "key"

	krackWorkaroundAttribute            = "wpa_disable_eapol_key_retries"

	krackWorkaroundAttributeDescription = "Enable WPA key reinstallation attack (KRACK) workaround. This should be `true` to enable KRACK workaround (you almost surely want this enabled)."
	krackWorkaroundUCIOption            = "wpa_disable_eapol_key_retries"

	macFilterAttribute            = "macfilter"
	macFilterAttributeDescription = `MAC address filter mode. Must be one of: "disable", "allow", "deny".`
	macFilterAllow                = "allow"
	macFilterDeny                 = "deny"
	macFilterDisable              = "disable"
	macFilterUCIOption            = "macfilter"

	macListAttribute            = "maclist"
	macListAttributeDescription = "List of MAC addresses for the MAC filter."
	macListUCIOption            = "maclist"

	maxAssocAttribute            = "maxassoc"
	maxAssocAttributeDescription = "Maximum number of associated clients."
	maxAssocUCIOption            = "maxassoc"

	modeAP                   = "ap"
	modeAttribute            = "mode"
	modeAttributeDescription = `The operation mode of the wireless network interface controller.. Currently only "ap" is supported.`
	modeUCIOption            = "mode"

	networkAttribute            = "network"
	networkAttributeDescription = "Network interface to attach the wireless network. This name is what the interface is known as in UCI, or the `id` field in Terraform."
	networkUCIOption            = "network"

	schemaDescription = "A wireless network."

	ssidAttribute            = "ssid"
	ssidAttributeDescription = "The broadcasted SSID of the wireless network. This is what actual clients will see the network as."
	ssidUCIOption            = "ssid"

	uciConfig = "wireless"
	uciType   = "wifi-iface"
)

var (
	deviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       deviceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDevice, deviceAttribute, deviceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDevice, deviceAttribute, deviceUCIOption),
	}

	disabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       disabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDisabled, disabledAttribute, disabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDisabled, disabledAttribute, disabledUCIOption),
	}

	encryptionMethodSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       encryptionMethodAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetEncryptionMethod, encryptionMethodAttribute, encryptionMethodUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetEncryptionMethod, encryptionMethodAttribute, encryptionMethodUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				encryptionMethodNone,
				encryptionMethodPSK,
				encryptionMethodPSK2,
				encryptionMethodPSK2AES,
				encryptionMethodPSK2CCMP,
				encryptionMethodPSK2TKIP,
				encryptionMethodPSK2TKIPAES,
				encryptionMethodPSK2TKIPCCMP,
				encryptionMethodPSKAES,
				encryptionMethodPSKCCMP,
				encryptionMethodPSKMixed,
				encryptionMethodPSKMixedAES,
				encryptionMethodPSKMixedCCMP,
				encryptionMethodPSKMixedTKIP,
				encryptionMethodPSKMixedTKIPAES,
				encryptionMethodPSKMixedTKIPCCMP,
				encryptionMethodPSKTKIP,
				encryptionMethodPSKTKIPAES,
				encryptionMethodPSKTKIPCCMP,
				encryptionMethodSAE,
				encryptionMethodSAEMixed,
			),
		},
	}

	hiddenSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       hiddenAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetHidden, hiddenAttribute, hiddenUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetHidden, hiddenAttribute, hiddenUCIOption),
	}

	ieee80211rSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ieee80211rAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetIEEE80211r, ieee80211rAttribute, ieee80211rUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetIEEE80211r, ieee80211rAttribute, ieee80211rUCIOption),
	}

	ieee80211wSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ieee80211wAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetIEEE80211w, ieee80211wAttribute, ieee80211wUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetIEEE80211w, ieee80211wAttribute, ieee80211wUCIOption),
		Validators: []validator.Int64{
			int64validator.OneOf(0, 1, 2),
		},
	}

	isolateClientsSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       isolateClientsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetIsolateClients, isolateClientsAttribute, isolateClientsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetIsolateClients, isolateClientsAttribute, isolateClientsUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(modeAttribute),
				modeAP,
			),
		},
	}

	keySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       keyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetKey, keyAttribute, keyUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		Sensitive:         true,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetKey, keyAttribute, keyUCIOption),
		Validators: []validator.String{
			stringvalidator.AlsoRequires(path.MatchRoot(encryptionMethodAttribute)),
			stringvalidator.LengthBetween(8, 63),
		},
	}

	krackWorkaroundSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       krackWorkaroundAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetKRACKWorkaround, krackWorkaroundAttribute, krackWorkaroundUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetKRACKWorkaround, krackWorkaroundAttribute, krackWorkaroundUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(modeAttribute),
				modeAP,
			),
		},
	}

	macFilterSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       macFilterAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMACFilter, macFilterAttribute, macFilterUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMACFilter, macFilterAttribute, macFilterUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				macFilterDisable,
				macFilterAllow,
				macFilterDeny,
			),
		},
	}

	macListSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       macListAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetMACList, macListAttribute, macListUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetMACList, macListAttribute, macListUCIOption),
		Validators: []validator.List{
			listvalidator.ValueStringsAre(
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$`),
					`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
				),
			),
		},
	}

	maxAssocSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       maxAssocAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMaxAssoc, maxAssocAttribute, maxAssocUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMaxAssoc, maxAssocAttribute, maxAssocUCIOption),
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
	}

	modeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       modeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMode, modeAttribute, modeUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMode, modeAttribute, modeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				modeAP,
			),
		},
	}

	networkSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       networkAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNetwork, networkAttribute, networkUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNetwork, networkAttribute, networkUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		deviceAttribute:           deviceSchemaAttribute,
		disabledAttribute:         disabledSchemaAttribute,
		encryptionMethodAttribute: encryptionMethodSchemaAttribute,
		hiddenAttribute:           hiddenSchemaAttribute,
		ieee80211rAttribute:       ieee80211rSchemaAttribute,
		ieee80211wAttribute:       ieee80211wSchemaAttribute,
		isolateClientsAttribute:   isolateClientsSchemaAttribute,
		keyAttribute:              keySchemaAttribute,
		krackWorkaroundAttribute:  krackWorkaroundSchemaAttribute,
		lucirpcglue.IdAttribute:   lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		macFilterAttribute:        macFilterSchemaAttribute,
		macListAttribute:          macListSchemaAttribute,
		maxAssocAttribute:         maxAssocSchemaAttribute,
		modeAttribute:             modeSchemaAttribute,
		networkAttribute:          networkSchemaAttribute,
		ssidAttribute:             ssidSchemaAttribute,
	}

	ssidSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ssidAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSSID, ssidAttribute, ssidUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSSID, ssidAttribute, ssidUCIOption),
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
	Device           types.String `tfsdk:"device"`
	Disabled         types.Bool   `tfsdk:"disabled"`
	EncryptionMethod types.String `tfsdk:"encryption"`
	Hidden           types.Bool   `tfsdk:"hidden"`
	Id               types.String `tfsdk:"id"`
	IEEE80211r       types.Bool   `tfsdk:"ieee80211r"`
	IEEE80211w       types.Int64  `tfsdk:"ieee80211w"`
	IsolateClients   types.Bool   `tfsdk:"isolate"`
	Key              types.String `tfsdk:"key"`
	KRACKWorkaround  types.Bool   `tfsdk:"wpa_disable_eapol_key_retries"`
	MACFilter        types.String `tfsdk:"macfilter"`
	MACList          types.List   `tfsdk:"maclist"`
	MaxAssoc         types.Int64  `tfsdk:"maxassoc"`
	Mode             types.String `tfsdk:"mode"`
	Network          types.String `tfsdk:"network"`
	SSID             types.String `tfsdk:"ssid"`
}

func modelGetDevice(m model) types.String           { return m.Device }
func modelGetDisabled(m model) types.Bool           { return m.Disabled }
func modelGetEncryptionMethod(m model) types.String { return m.EncryptionMethod }
func modelGetHidden(m model) types.Bool             { return m.Hidden }
func modelGetId(m model) types.String               { return m.Id }
func modelGetIEEE80211r(m model) types.Bool         { return m.IEEE80211r }
func modelGetIEEE80211w(m model) types.Int64        { return m.IEEE80211w }
func modelGetIsolateClients(m model) types.Bool     { return m.IsolateClients }
func modelGetKey(m model) types.String              { return m.Key }
func modelGetKRACKWorkaround(m model) types.Bool    { return m.KRACKWorkaround }
func modelGetMACFilter(m model) types.String        { return m.MACFilter }
func modelGetMACList(m model) types.List            { return m.MACList }
func modelGetMaxAssoc(m model) types.Int64          { return m.MaxAssoc }
func modelGetMode(m model) types.String             { return m.Mode }
func modelGetNetwork(m model) types.String          { return m.Network }
func modelGetSSID(m model) types.String             { return m.SSID }

func modelSetDevice(m *model, value types.String)           { m.Device = value }
func modelSetDisabled(m *model, value types.Bool)           { m.Disabled = value }
func modelSetEncryptionMethod(m *model, value types.String) { m.EncryptionMethod = value }
func modelSetHidden(m *model, value types.Bool)             { m.Hidden = value }
func modelSetId(m *model, value types.String)               { m.Id = value }
func modelSetIEEE80211r(m *model, value types.Bool)         { m.IEEE80211r = value }
func modelSetIEEE80211w(m *model, value types.Int64)        { m.IEEE80211w = value }
func modelSetIsolateClients(m *model, value types.Bool)     { m.IsolateClients = value }
func modelSetKey(m *model, value types.String)              { m.Key = value }
func modelSetKRACKWorkaround(m *model, value types.Bool)    { m.KRACKWorkaround = value }
func modelSetMACFilter(m *model, value types.String)        { m.MACFilter = value }
func modelSetMACList(m *model, value types.List)            { m.MACList = value }
func modelSetMaxAssoc(m *model, value types.Int64)          { m.MaxAssoc = value }
func modelSetMode(m *model, value types.String)             { m.Mode = value }
func modelSetNetwork(m *model, value types.String)          { m.Network = value }
func modelSetSSID(m *model, value types.String)             { m.SSID = value }
