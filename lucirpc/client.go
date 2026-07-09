package lucirpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ErrSectionNotFound is returned by GetSection when the requested section does not exist.
var ErrSectionNotFound = errors.New("section not found")

const (
	humanReadableCommitChanges   = "commit changes"
	humanReadableCreateSection   = "create section"
	humanReadableDeleteSection   = "delete section"
	humanReadableGetSection      = "get section"
	humanReadableGetSections     = "get sections"
	humanReadableLogin           = "login"
	humanReadableReorderSections = "reorder sections"
	humanReadableRevertChanges   = "revert changes"
	humanReadableShowChanges     = "show changes"
	humanReadableUpdateSection   = "update section"

	methodChanges = "changes"
	methodCommit  = "commit"
	methodDelete  = "delete"
	methodForeach = "foreach"
	methodGetAll  = "get_all"
	methodLogin   = "login"
	methodRevert  = "revert"
	methodSection = "section"
	methodTSet    = "tset"

	pathAuth = "/cgi-bin/luci/rpc/auth"
	pathUCI  = "/cgi-bin/luci/rpc/uci"

	queryKeyAuth = "auth"
)

type Client struct {
	jsonRPCClientUCI jsonRPCClient
}

func (c *Client) CommitChanges(
	ctx context.Context,
	config string,
) (bool, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return false, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableCommitChanges, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodCommit,
		Params: []json.RawMessage{
			marshalledConfig,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableCommitChanges,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableCommitChanges, err)
	}

	// The result can be `true` to indicate success,
	// or `null` to indicate failure.
	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) CreateSection(
	ctx context.Context,
	config string,
	sectionType string,
	section string,
	options Options,
) (bool, error) {
	set, unset := partitionOptions(options)
	err := c.stageCreate(ctx, config, sectionType, section, set, unset)
	if err != nil {
		// Discard staged changes so a later commit cannot apply a partial create.
		_, _ = c.invokeBoolean(ctx, humanReadableRevertChanges, methodRevert, config)
		return false, err
	}

	result, err := c.CommitChanges(
		ctx,
		config,
	)
	if err != nil {
		return false, fmt.Errorf("was able to %s, but could not %s: %w", humanReadableCreateSection, humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) stageCreate(
	ctx context.Context,
	config string,
	sectionType string,
	section string,
	set Options,
	unset []string,
) error {
	if len(unset) > 0 {
		current, err := c.GetSection(ctx, config, section)
		switch {
		case errors.Is(err, ErrSectionNotFound):
			// The section is new; there is nothing to delete.

		case err != nil:
			return fmt.Errorf("unable to %s: %w", humanReadableCreateSection, err)

		default:
			err = c.deleteOptions(ctx, humanReadableCreateSection, config, section, unset, current)
			if err != nil {
				return err
			}
		}
	}

	// The section method appends list options to a pre-existing section instead of replacing them,
	// so it only establishes the section and its type; tset then sets the options.
	ok, err := c.invokeBoolean(ctx, humanReadableCreateSection, methodSection, config, sectionType, section, Options{})
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableCreateSection)
	}

	// tset fails on an empty set of options.
	if len(set) == 0 {
		return nil
	}

	ok, err = c.invokeBoolean(ctx, humanReadableCreateSection, methodTSet, config, section, set)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableCreateSection)
	}

	return nil
}

// deleteOptions stages the deletion of each unset option that is present in the section.
func (c *Client) deleteOptions(
	ctx context.Context,
	humanReadableMethod string,
	config string,
	section string,
	unset []string,
	current Options,
) error {
	for _, option := range unset {
		if _, ok := current[option]; !ok {
			continue
		}

		ok, err := c.invokeBoolean(ctx, humanReadableMethod, methodDelete, config, section, option)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("unable to %s: could not delete option %s.%s.%s", humanReadableMethod, config, section, option)
		}
	}

	return nil
}

func (c *Client) DeleteSection(
	ctx context.Context,
	config string,
	section string,
) (bool, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return false, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableDeleteSection, err)
	}

	marshalledSection, err := json.Marshal(section)
	if err != nil {
		return false, fmt.Errorf("unable to serialize section %q for %s: %w", section, humanReadableDeleteSection, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodDelete,
		Params: []json.RawMessage{
			marshalledConfig,
			marshalledSection,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableDeleteSection,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableDeleteSection, err)
	}

	// The result can be `true` to indicate success,
	// or `null` to indicate failure.
	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableDeleteSection, err)
	}

	if !result {
		return false, fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableDeleteSection)
	}

	result, err = c.CommitChanges(
		ctx,
		config,
	)
	if err != nil {
		return false, fmt.Errorf("was able to %s, but could not %s: %w", humanReadableDeleteSection, humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) GetSection(
	ctx context.Context,
	config string,
	section string,
) (Options, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableGetSection, err)
	}

	marshalledSection, err := json.Marshal(section)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize section %q for %s: %w", section, humanReadableGetSection, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodGetAll,
		Params: []json.RawMessage{
			marshalledConfig,
			marshalledSection,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableGetSection,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableGetSection, err)
	}

	if responseBody == nil {
		return nil, fmt.Errorf("%w: %s.%s", ErrSectionNotFound, config, section)
	}

	// Depending on the `config` and `section`,
	// this method can return a response that is an array instead of an object.
	// We have to handle that case as well.
	var unknownResult any
	err = json.Unmarshal(*responseBody, &unknownResult)
	if err != nil {
		return nil, fmt.Errorf("unable to determine type of %s response: %w", humanReadableGetSection, err)
	}

	_, ok := unknownResult.([]any)
	if ok {
		return nil, fmt.Errorf("incorrect config (%q) and/or section (%q): result from LuCI: %s", config, section, *responseBody)
	}

	var result Options
	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableGetSection, err)
	}

	return result, nil
}

// GetSections returns all sections of the given type,
// in the order they appear in the config.
func (c *Client) GetSections(
	ctx context.Context,
	config string,
	sectionType string,
) ([]Options, error) {
	params, err := marshalParams(humanReadableGetSections, config, sectionType)
	if err != nil {
		return nil, err
	}

	requestBody := jsonRPCRequestBody{
		Method: methodForeach,
		Params: params,
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableGetSections,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableGetSections, err)
	}

	// The result is `false` when no sections match.
	result := []Options{}
	if responseBody == nil {
		return result, nil
	}

	var noSections bool
	if json.Unmarshal(*responseBody, &noSections) == nil {
		return result, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableGetSections, err)
	}

	return result, nil
}

// ReorderSections rewrites the given sections so they appear in the given order.
// The LuCI JSON-RPC API does not expose UCI's reorder,
// so each section is deleted and recreated with its existing options,
// which moves it to the end of the config.
// All changes are staged and committed once;
// staged changes are reverted if any step fails.
func (c *Client) ReorderSections(
	ctx context.Context,
	config string,
	sectionType string,
	sections []string,
) (bool, error) {
	optionsBySection := map[string]Options{}
	for _, section := range sections {
		options, err := c.GetSection(ctx, config, section)
		if err != nil {
			return false, fmt.Errorf("unable to %s: %w", humanReadableReorderSections, err)
		}

		actualType, err := options.GetString(".type")
		if err != nil || actualType != sectionType {
			return false, fmt.Errorf("unable to %s: section %s.%s is not of type %q", humanReadableReorderSections, config, section, sectionType)
		}

		for option := range options {
			if strings.HasPrefix(option, ".") {
				delete(options, option)
			}
		}

		optionsBySection[section] = options
	}

	err := c.stageReorder(ctx, config, sectionType, sections, optionsBySection)
	if err != nil {
		// Discard staged changes so a later commit cannot apply a partial reorder.
		_, _ = c.invokeBoolean(ctx, humanReadableRevertChanges, methodRevert, config)
		return false, err
	}

	return c.CommitChanges(ctx, config)
}

func (c *Client) stageReorder(
	ctx context.Context,
	config string,
	sectionType string,
	sections []string,
	optionsBySection map[string]Options,
) error {
	for _, section := range sections {
		ok, err := c.invokeBoolean(ctx, humanReadableReorderSections, methodDelete, config, section)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("unable to %s: could not delete section %s.%s", humanReadableReorderSections, config, section)
		}

		ok, err = c.invokeBoolean(ctx, humanReadableReorderSections, methodSection, config, sectionType, section, optionsBySection[section])
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("unable to %s: could not recreate section %s.%s", humanReadableReorderSections, config, section)
		}
	}

	return nil
}

func (c *Client) ShowChanges(
	ctx context.Context,
	config string,
) ([][]string, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableShowChanges, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodChanges,
		Params: []json.RawMessage{
			marshalledConfig,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableShowChanges,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableShowChanges, err)
	}

	result := [][]string{}
	if responseBody == nil {
		return result, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableShowChanges, err)
	}

	return result, nil
}

func (c *Client) UpdateSection(
	ctx context.Context,
	config string,
	section string,
	options Options,
) (bool, error) {
	set, unset := partitionOptions(options)
	err := c.stageUpdate(ctx, config, section, set, unset)
	if err != nil {
		// Discard staged changes so a later commit cannot apply a partial update.
		_, _ = c.invokeBoolean(ctx, humanReadableRevertChanges, methodRevert, config)
		return false, err
	}

	result, err := c.CommitChanges(
		ctx,
		config,
	)
	if err != nil {
		return false, fmt.Errorf("was able to %s, but could not %s: %w", humanReadableUpdateSection, humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) stageUpdate(
	ctx context.Context,
	config string,
	section string,
	set Options,
	unset []string,
) error {
	if len(set) == 0 && len(unset) == 0 {
		return fmt.Errorf("unable to %s: no options provided", humanReadableUpdateSection)
	}

	if len(unset) > 0 {
		current, err := c.GetSection(ctx, config, section)
		if err != nil {
			return fmt.Errorf("unable to %s: %w", humanReadableUpdateSection, err)
		}

		err = c.deleteOptions(ctx, humanReadableUpdateSection, config, section, unset, current)
		if err != nil {
			return err
		}
	}

	// tset fails on an empty set of options.
	if len(set) == 0 {
		return nil
	}

	ok, err := c.invokeBoolean(ctx, humanReadableUpdateSection, methodTSet, config, section, set)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableUpdateSection)
	}

	return nil
}

// invokeBoolean invokes a method whose result is `true` to indicate success,
// or `null`/`false` to indicate failure.
func (c *Client) invokeBoolean(
	ctx context.Context,
	humanReadableMethod string,
	method string,
	params ...any,
) (bool, error) {
	marshalledParams, err := marshalParams(humanReadableMethod, params...)
	if err != nil {
		return false, err
	}

	requestBody := jsonRPCRequestBody{
		Method: method,
		Params: marshalledParams,
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableMethod,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableMethod, err)
	}

	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableMethod, err)
	}

	return result, nil
}

func marshalParams(
	humanReadableMethod string,
	params ...any,
) ([]json.RawMessage, error) {
	marshalled := make([]json.RawMessage, 0, len(params))
	for _, param := range params {
		value, err := json.Marshal(param)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize parameter %v for %s: %w", param, humanReadableMethod, err)
		}

		marshalled = append(marshalled, value)
	}

	return marshalled, nil
}

func NewClient(
	ctx context.Context,
	scheme string,
	hostname string,
	port uint16,
	username string,
	password string,
) (*Client, error) {
	host := hostname
	if port != 0 {
		host = fmt.Sprintf("%s:%d", host, port)
	}

	address := url.URL{
		Host:   host,
		Path:   pathAuth,
		Scheme: scheme,
	}
	httpClient := &http.Client{}
	marshalledUsername, err := json.Marshal(username)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize username for %s: %w", humanReadableLogin, err)
	}

	marshalledPassword, err := json.Marshal(password)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize password for %s: %w", humanReadableLogin, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodLogin,
		Params: []json.RawMessage{
			marshalledUsername,
			marshalledPassword,
		},
	}
	jsonRPCClient := jsonRPCNewClient(
		*httpClient,
		address,
	)
	responseBody, err := jsonRPCClient.InvokeNotNull(
		ctx,
		humanReadableLogin,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableLogin, err)
	}

	var authToken string
	err = json.Unmarshal(responseBody, &authToken)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableLogin, err)
	}

	query := url.Values{}
	query.Add(queryKeyAuth, authToken)
	addressUCI := url.URL{
		Host:     host,
		Path:     pathUCI,
		RawQuery: query.Encode(),
		Scheme:   scheme,
	}
	jsonRPCClientUCI := jsonRPCNewClient(
		*httpClient,
		addressUCI,
	)
	client := &Client{
		jsonRPCClientUCI: jsonRPCClientUCI,
	}
	return client, nil
}

type jsonRPCClient struct {
	address url.URL
	client  http.Client
}

func (c jsonRPCClient) InvokeNotNull(
	ctx context.Context,
	humanReadableMethod string,
	requestBody jsonRPCRequestBody,
) (json.RawMessage, error) {
	result, err := c.Invoke(
		ctx,
		humanReadableMethod,
		requestBody,
	)
	if err != nil {
		return json.RawMessage{}, err
	}

	if result == nil {
		return nil, fmt.Errorf("invalid %s response: expected either an error or a result, got neither", humanReadableMethod)
	}

	return *result, nil
}

func (c jsonRPCClient) Invoke(
	ctx context.Context,
	humanReadableMethod string,
	requestBody jsonRPCRequestBody,
) (*json.RawMessage, error) {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("problem encoding %s request: %w", humanReadableMethod, err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.address.String(),
		&buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("problem creating %s request: %w", humanReadableMethod, err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("problem sending request to %s: %w", humanReadableMethod, err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected %s to respond with a 200: got %s", humanReadableMethod, response.Status)
	}

	var responseBody jsonRPCResponseBody
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableMethod, err)
	}

	if responseBody.Error != nil {
		return nil, fmt.Errorf("%s error: %s", humanReadableMethod, *responseBody.Error)
	}

	return responseBody.Result, nil
}

func jsonRPCNewClient(
	httpClient http.Client,
	address url.URL,
) jsonRPCClient {
	return jsonRPCClient{
		address: address,
		client:  httpClient,
	}
}

type jsonRPCRequestBody struct {
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

type jsonRPCResponseBody struct {
	Error  *string          `json:"error"`
	Result *json.RawMessage `json:"result"`
}
