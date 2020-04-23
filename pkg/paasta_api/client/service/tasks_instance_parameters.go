// Code generated by go-swagger; DO NOT EDIT.

package service

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
	"github.com/go-openapi/swag"
)

// NewTasksInstanceParams creates a new TasksInstanceParams object
// with the default values initialized.
func NewTasksInstanceParams() *TasksInstanceParams {
	var ()
	return &TasksInstanceParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewTasksInstanceParamsWithTimeout creates a new TasksInstanceParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewTasksInstanceParamsWithTimeout(timeout time.Duration) *TasksInstanceParams {
	var ()
	return &TasksInstanceParams{

		timeout: timeout,
	}
}

// NewTasksInstanceParamsWithContext creates a new TasksInstanceParams object
// with the default values initialized, and the ability to set a context for a request
func NewTasksInstanceParamsWithContext(ctx context.Context) *TasksInstanceParams {
	var ()
	return &TasksInstanceParams{

		Context: ctx,
	}
}

// NewTasksInstanceParamsWithHTTPClient creates a new TasksInstanceParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewTasksInstanceParamsWithHTTPClient(client *http.Client) *TasksInstanceParams {
	var ()
	return &TasksInstanceParams{
		HTTPClient: client,
	}
}

/*TasksInstanceParams contains all the parameters to send to the API endpoint
for the tasks instance operation typically these are written to a http.Request
*/
type TasksInstanceParams struct {

	/*Instance
	  Instance name

	*/
	Instance string
	/*Service
	  Service name

	*/
	Service string
	/*SlaveHostname
	  slave hostname to filter tasks by

	*/
	SlaveHostname *string
	/*Verbose
	  Return slave and executor for task

	*/
	Verbose *bool

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the tasks instance params
func (o *TasksInstanceParams) WithTimeout(timeout time.Duration) *TasksInstanceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the tasks instance params
func (o *TasksInstanceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the tasks instance params
func (o *TasksInstanceParams) WithContext(ctx context.Context) *TasksInstanceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the tasks instance params
func (o *TasksInstanceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the tasks instance params
func (o *TasksInstanceParams) WithHTTPClient(client *http.Client) *TasksInstanceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the tasks instance params
func (o *TasksInstanceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithInstance adds the instance to the tasks instance params
func (o *TasksInstanceParams) WithInstance(instance string) *TasksInstanceParams {
	o.SetInstance(instance)
	return o
}

// SetInstance adds the instance to the tasks instance params
func (o *TasksInstanceParams) SetInstance(instance string) {
	o.Instance = instance
}

// WithService adds the service to the tasks instance params
func (o *TasksInstanceParams) WithService(service string) *TasksInstanceParams {
	o.SetService(service)
	return o
}

// SetService adds the service to the tasks instance params
func (o *TasksInstanceParams) SetService(service string) {
	o.Service = service
}

// WithSlaveHostname adds the slaveHostname to the tasks instance params
func (o *TasksInstanceParams) WithSlaveHostname(slaveHostname *string) *TasksInstanceParams {
	o.SetSlaveHostname(slaveHostname)
	return o
}

// SetSlaveHostname adds the slaveHostname to the tasks instance params
func (o *TasksInstanceParams) SetSlaveHostname(slaveHostname *string) {
	o.SlaveHostname = slaveHostname
}

// WithVerbose adds the verbose to the tasks instance params
func (o *TasksInstanceParams) WithVerbose(verbose *bool) *TasksInstanceParams {
	o.SetVerbose(verbose)
	return o
}

// SetVerbose adds the verbose to the tasks instance params
func (o *TasksInstanceParams) SetVerbose(verbose *bool) {
	o.Verbose = verbose
}

// WriteToRequest writes these params to a swagger request
func (o *TasksInstanceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	if o.SlaveHostname != nil {

		// query param slave_hostname
		var qrSlaveHostname string
		if o.SlaveHostname != nil {
			qrSlaveHostname = *o.SlaveHostname
		}
		qSlaveHostname := qrSlaveHostname
		if qSlaveHostname != "" {
			if err := r.SetQueryParam("slave_hostname", qSlaveHostname); err != nil {
				return err
			}
		}

	}

	if o.Verbose != nil {

		// query param verbose
		var qrVerbose bool
		if o.Verbose != nil {
			qrVerbose = *o.Verbose
		}
		qVerbose := swag.FormatBool(qrVerbose)
		if qVerbose != "" {
			if err := r.SetQueryParam("verbose", qVerbose); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}