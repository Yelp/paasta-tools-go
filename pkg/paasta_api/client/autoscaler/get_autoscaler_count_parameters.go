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

// NewGetAutoscalerCountParams creates a new GetAutoscalerCountParams object
// with the default values initialized.
func NewGetAutoscalerCountParams() *GetAutoscalerCountParams {
	var ()
	return &GetAutoscalerCountParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAutoscalerCountParamsWithTimeout creates a new GetAutoscalerCountParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAutoscalerCountParamsWithTimeout(timeout time.Duration) *GetAutoscalerCountParams {
	var ()
	return &GetAutoscalerCountParams{

		timeout: timeout,
	}
}

// NewGetAutoscalerCountParamsWithContext creates a new GetAutoscalerCountParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAutoscalerCountParamsWithContext(ctx context.Context) *GetAutoscalerCountParams {
	var ()
	return &GetAutoscalerCountParams{

		Context: ctx,
	}
}

// NewGetAutoscalerCountParamsWithHTTPClient creates a new GetAutoscalerCountParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAutoscalerCountParamsWithHTTPClient(client *http.Client) *GetAutoscalerCountParams {
	var ()
	return &GetAutoscalerCountParams{
		HTTPClient: client,
	}
}

/*GetAutoscalerCountParams contains all the parameters to send to the API endpoint
for the get autoscaler count operation typically these are written to a http.Request
*/
type GetAutoscalerCountParams struct {

	/*Instance
	  Instance name

	*/
	Instance string
	/*Service
	  Service name

	*/
	Service string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get autoscaler count params
func (o *GetAutoscalerCountParams) WithTimeout(timeout time.Duration) *GetAutoscalerCountParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get autoscaler count params
func (o *GetAutoscalerCountParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get autoscaler count params
func (o *GetAutoscalerCountParams) WithContext(ctx context.Context) *GetAutoscalerCountParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get autoscaler count params
func (o *GetAutoscalerCountParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get autoscaler count params
func (o *GetAutoscalerCountParams) WithHTTPClient(client *http.Client) *GetAutoscalerCountParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get autoscaler count params
func (o *GetAutoscalerCountParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithInstance adds the instance to the get autoscaler count params
func (o *GetAutoscalerCountParams) WithInstance(instance string) *GetAutoscalerCountParams {
	o.SetInstance(instance)
	return o
}

// SetInstance adds the instance to the get autoscaler count params
func (o *GetAutoscalerCountParams) SetInstance(instance string) {
	o.Instance = instance
}

// WithService adds the service to the get autoscaler count params
func (o *GetAutoscalerCountParams) WithService(service string) *GetAutoscalerCountParams {
	o.SetService(service)
	return o
}

// SetService adds the service to the get autoscaler count params
func (o *GetAutoscalerCountParams) SetService(service string) {
	o.Service = service
}

// WriteToRequest writes these params to a swagger request
func (o *GetAutoscalerCountParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param instance
	if err := r.SetPathParam("instance", o.Instance); err != nil {
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