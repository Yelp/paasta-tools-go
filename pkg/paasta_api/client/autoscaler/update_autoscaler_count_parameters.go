// Code generated by go-swagger; DO NOT EDIT.

package autoscaler

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewUpdateAutoscalerCountParams creates a new UpdateAutoscalerCountParams object
// with the default values initialized.
func NewUpdateAutoscalerCountParams() *UpdateAutoscalerCountParams {
	var ()
	return &UpdateAutoscalerCountParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateAutoscalerCountParamsWithTimeout creates a new UpdateAutoscalerCountParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateAutoscalerCountParamsWithTimeout(timeout time.Duration) *UpdateAutoscalerCountParams {
	var ()
	return &UpdateAutoscalerCountParams{

		timeout: timeout,
	}
}

// NewUpdateAutoscalerCountParamsWithContext creates a new UpdateAutoscalerCountParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateAutoscalerCountParamsWithContext(ctx context.Context) *UpdateAutoscalerCountParams {
	var ()
	return &UpdateAutoscalerCountParams{

		Context: ctx,
	}
}

// NewUpdateAutoscalerCountParamsWithHTTPClient creates a new UpdateAutoscalerCountParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateAutoscalerCountParamsWithHTTPClient(client *http.Client) *UpdateAutoscalerCountParams {
	var ()
	return &UpdateAutoscalerCountParams{
		HTTPClient: client,
	}
}

/*UpdateAutoscalerCountParams contains all the parameters to send to the API endpoint
for the update autoscaler count operation typically these are written to a http.Request
*/
type UpdateAutoscalerCountParams struct {

	/*Instance
	  Instance name

	*/
	Instance string
	/*JSONBody*/
	JSONBody UpdateAutoscalerCountBody
	/*Service
	  Service name

	*/
	Service string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) WithTimeout(timeout time.Duration) *UpdateAutoscalerCountParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) WithContext(ctx context.Context) *UpdateAutoscalerCountParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) WithHTTPClient(client *http.Client) *UpdateAutoscalerCountParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithInstance adds the instance to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) WithInstance(instance string) *UpdateAutoscalerCountParams {
	o.SetInstance(instance)
	return o
}

// SetInstance adds the instance to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) SetInstance(instance string) {
	o.Instance = instance
}

// WithJSONBody adds the jSONBody to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) WithJSONBody(jSONBody UpdateAutoscalerCountBody) *UpdateAutoscalerCountParams {
	o.SetJSONBody(jSONBody)
	return o
}

// SetJSONBody adds the jsonBody to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) SetJSONBody(jSONBody UpdateAutoscalerCountBody) {
	o.JSONBody = jSONBody
}

// WithService adds the service to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) WithService(service string) *UpdateAutoscalerCountParams {
	o.SetService(service)
	return o
}

// SetService adds the service to the update autoscaler count params
func (o *UpdateAutoscalerCountParams) SetService(service string) {
	o.Service = service
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateAutoscalerCountParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param instance
	if err := r.SetPathParam("instance", o.Instance); err != nil {
		return err
	}

	if err := r.SetBodyParam(o.JSONBody); err != nil {
		return err
	}

	// path param service
	if err := r.SetPathParam("service", o.Service); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
