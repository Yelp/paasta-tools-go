// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DeleteServiceAutoscalerPause unpauses the autoscaler
*/
func (a *Client) DeleteServiceAutoscalerPause(params *DeleteServiceAutoscalerPauseParams) (*DeleteServiceAutoscalerPauseOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteServiceAutoscalerPauseParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "delete_service_autoscaler_pause",
		Method:             "DELETE",
		PathPattern:        "/service_autoscaler/pause",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/x-www-form-urlencoded"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteServiceAutoscalerPauseReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteServiceAutoscalerPauseOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for delete_service_autoscaler_pause: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetServiceAutoscalerPause gets autoscaling pause time
*/
func (a *Client) GetServiceAutoscalerPause(params *GetServiceAutoscalerPauseParams) (*GetServiceAutoscalerPauseOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetServiceAutoscalerPauseParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "get_service_autoscaler_pause",
		Method:             "GET",
		PathPattern:        "/service_autoscaler/pause",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/x-www-form-urlencoded"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetServiceAutoscalerPauseReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetServiceAutoscalerPauseOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for get_service_autoscaler_pause: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
Metastatus gets metastatus
*/
func (a *Client) Metastatus(params *MetastatusParams) (*MetastatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewMetastatusParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "metastatus",
		Method:             "GET",
		PathPattern:        "/metastatus",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/x-www-form-urlencoded"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &MetastatusReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*MetastatusOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for metastatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
ShowVersion versions of paasta tools package
*/
func (a *Client) ShowVersion(params *ShowVersionParams) (*ShowVersionOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewShowVersionParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "showVersion",
		Method:             "GET",
		PathPattern:        "/version",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/x-www-form-urlencoded"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ShowVersionReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ShowVersionOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for showVersion: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
UpdateServiceAutoscalerPause update service autoscaler pause API
*/
func (a *Client) UpdateServiceAutoscalerPause(params *UpdateServiceAutoscalerPauseParams) (*UpdateServiceAutoscalerPauseOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateServiceAutoscalerPauseParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "update_service_autoscaler_pause",
		Method:             "POST",
		PathPattern:        "/service_autoscaler/pause",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/x-www-form-urlencoded"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateServiceAutoscalerPauseReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*UpdateServiceAutoscalerPauseOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for update_service_autoscaler_pause: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
