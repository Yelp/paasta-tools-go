/*
 * Paasta API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package paastaapi

import (
	_context "context"
	_ioutil "io/ioutil"
	_nethttp "net/http"
	_neturl "net/url"
	"strings"
)

// Linger please
var (
	_ _context.Context
)

// ServiceApiService ServiceApi service
type ServiceApiService service

type ApiDelayInstanceRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
	service string
	instance string
}


func (r ApiDelayInstanceRequest) Execute() (map[string]interface{}, *_nethttp.Response, error) {
	return r.ApiService.DelayInstanceExecute(r)
}

/*
 * DelayInstance Get the possible reasons for a deployment delay for a marathon service.instance
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param service Service name
 * @param instance Instance name
 * @return ApiDelayInstanceRequest
 */
func (a *ServiceApiService) DelayInstance(ctx _context.Context, service string, instance string) ApiDelayInstanceRequest {
	return ApiDelayInstanceRequest{
		ApiService: a,
		ctx: ctx,
		service: service,
		instance: instance,
	}
}

/*
 * Execute executes the request
 * @return map[string]interface{}
 */
func (a *ServiceApiService) DelayInstanceExecute(r ApiDelayInstanceRequest) (map[string]interface{}, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.DelayInstance")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services/{service}/{instance}/delay"
	localVarPath = strings.Replace(localVarPath, "{"+"service"+"}", _neturl.PathEscape(parameterToString(r.service, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instance"+"}", _neturl.PathEscape(parameterToString(r.instance, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiInstanceSetStateRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
	service string
	instance string
	desiredState string
}


func (r ApiInstanceSetStateRequest) Execute() (*_nethttp.Response, error) {
	return r.ApiService.InstanceSetStateExecute(r)
}

/*
 * InstanceSetState Change state of service_name.instance_name
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param service Service name
 * @param instance Instance name
 * @param desiredState Desired state
 * @return ApiInstanceSetStateRequest
 */
func (a *ServiceApiService) InstanceSetState(ctx _context.Context, service string, instance string, desiredState string) ApiInstanceSetStateRequest {
	return ApiInstanceSetStateRequest{
		ApiService: a,
		ctx: ctx,
		service: service,
		instance: instance,
		desiredState: desiredState,
	}
}

/*
 * Execute executes the request
 */
func (a *ServiceApiService) InstanceSetStateExecute(r ApiInstanceSetStateRequest) (*_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPost
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.InstanceSetState")
	if err != nil {
		return nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services/{service}/{instance}/state/{desired_state}"
	localVarPath = strings.Replace(localVarPath, "{"+"service"+"}", _neturl.PathEscape(parameterToString(r.service, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instance"+"}", _neturl.PathEscape(parameterToString(r.instance, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"desired_state"+"}", _neturl.PathEscape(parameterToString(r.desiredState, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiListInstancesRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
	service string
}


func (r ApiListInstancesRequest) Execute() (InlineResponse2001, *_nethttp.Response, error) {
	return r.ApiService.ListInstancesExecute(r)
}

/*
 * ListInstances List instances of service_name
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param service Service name
 * @return ApiListInstancesRequest
 */
func (a *ServiceApiService) ListInstances(ctx _context.Context, service string) ApiListInstancesRequest {
	return ApiListInstancesRequest{
		ApiService: a,
		ctx: ctx,
		service: service,
	}
}

/*
 * Execute executes the request
 * @return InlineResponse2001
 */
func (a *ServiceApiService) ListInstancesExecute(r ApiListInstancesRequest) (InlineResponse2001, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  InlineResponse2001
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.ListInstances")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services/{service}"
	localVarPath = strings.Replace(localVarPath, "{"+"service"+"}", _neturl.PathEscape(parameterToString(r.service, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiListServicesForClusterRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
}


func (r ApiListServicesForClusterRequest) Execute() (InlineResponse200, *_nethttp.Response, error) {
	return r.ApiService.ListServicesForClusterExecute(r)
}

/*
 * ListServicesForCluster List service names and service instance names on the current cluster
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @return ApiListServicesForClusterRequest
 */
func (a *ServiceApiService) ListServicesForCluster(ctx _context.Context) ApiListServicesForClusterRequest {
	return ApiListServicesForClusterRequest{
		ApiService: a,
		ctx: ctx,
	}
}

/*
 * Execute executes the request
 * @return InlineResponse200
 */
func (a *ServiceApiService) ListServicesForClusterExecute(r ApiListServicesForClusterRequest) (InlineResponse200, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  InlineResponse200
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.ListServicesForCluster")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiStatusInstanceRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
	service string
	instance string
	verbose *int32
	includeSmartstack *bool
	includeEnvoy *bool
	includeMesos *bool
}

func (r ApiStatusInstanceRequest) Verbose(verbose int32) ApiStatusInstanceRequest {
	r.verbose = &verbose
	return r
}
func (r ApiStatusInstanceRequest) IncludeSmartstack(includeSmartstack bool) ApiStatusInstanceRequest {
	r.includeSmartstack = &includeSmartstack
	return r
}
func (r ApiStatusInstanceRequest) IncludeEnvoy(includeEnvoy bool) ApiStatusInstanceRequest {
	r.includeEnvoy = &includeEnvoy
	return r
}
func (r ApiStatusInstanceRequest) IncludeMesos(includeMesos bool) ApiStatusInstanceRequest {
	r.includeMesos = &includeMesos
	return r
}

func (r ApiStatusInstanceRequest) Execute() (InstanceStatus, *_nethttp.Response, error) {
	return r.ApiService.StatusInstanceExecute(r)
}

/*
 * StatusInstance Get status of service_name.instance_name
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param service Service name
 * @param instance Instance name
 * @return ApiStatusInstanceRequest
 */
func (a *ServiceApiService) StatusInstance(ctx _context.Context, service string, instance string) ApiStatusInstanceRequest {
	return ApiStatusInstanceRequest{
		ApiService: a,
		ctx: ctx,
		service: service,
		instance: instance,
	}
}

/*
 * Execute executes the request
 * @return InstanceStatus
 */
func (a *ServiceApiService) StatusInstanceExecute(r ApiStatusInstanceRequest) (InstanceStatus, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  InstanceStatus
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.StatusInstance")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services/{service}/{instance}/status"
	localVarPath = strings.Replace(localVarPath, "{"+"service"+"}", _neturl.PathEscape(parameterToString(r.service, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instance"+"}", _neturl.PathEscape(parameterToString(r.instance, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.verbose != nil {
		localVarQueryParams.Add("verbose", parameterToString(*r.verbose, ""))
	}
	if r.includeSmartstack != nil {
		localVarQueryParams.Add("include_smartstack", parameterToString(*r.includeSmartstack, ""))
	}
	if r.includeEnvoy != nil {
		localVarQueryParams.Add("include_envoy", parameterToString(*r.includeEnvoy, ""))
	}
	if r.includeMesos != nil {
		localVarQueryParams.Add("include_mesos", parameterToString(*r.includeMesos, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiTaskInstanceRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
	service string
	instance string
	taskId string
	verbose *bool
}

func (r ApiTaskInstanceRequest) Verbose(verbose bool) ApiTaskInstanceRequest {
	r.verbose = &verbose
	return r
}

func (r ApiTaskInstanceRequest) Execute() (map[string]interface{}, *_nethttp.Response, error) {
	return r.ApiService.TaskInstanceExecute(r)
}

/*
 * TaskInstance Get mesos task of service_name.instance_name by task_id
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param service Service name
 * @param instance Instance name
 * @param taskId mesos task id
 * @return ApiTaskInstanceRequest
 */
func (a *ServiceApiService) TaskInstance(ctx _context.Context, service string, instance string, taskId string) ApiTaskInstanceRequest {
	return ApiTaskInstanceRequest{
		ApiService: a,
		ctx: ctx,
		service: service,
		instance: instance,
		taskId: taskId,
	}
}

/*
 * Execute executes the request
 * @return map[string]interface{}
 */
func (a *ServiceApiService) TaskInstanceExecute(r ApiTaskInstanceRequest) (map[string]interface{}, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.TaskInstance")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services/{service}/{instance}/tasks/{task_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"service"+"}", _neturl.PathEscape(parameterToString(r.service, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instance"+"}", _neturl.PathEscape(parameterToString(r.instance, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"task_id"+"}", _neturl.PathEscape(parameterToString(r.taskId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.verbose != nil {
		localVarQueryParams.Add("verbose", parameterToString(*r.verbose, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiTasksInstanceRequest struct {
	ctx _context.Context
	ApiService *ServiceApiService
	service string
	instance string
	slaveHostname *string
	verbose *bool
}

func (r ApiTasksInstanceRequest) SlaveHostname(slaveHostname string) ApiTasksInstanceRequest {
	r.slaveHostname = &slaveHostname
	return r
}
func (r ApiTasksInstanceRequest) Verbose(verbose bool) ApiTasksInstanceRequest {
	r.verbose = &verbose
	return r
}

func (r ApiTasksInstanceRequest) Execute() ([]map[string]interface{}, *_nethttp.Response, error) {
	return r.ApiService.TasksInstanceExecute(r)
}

/*
 * TasksInstance Get mesos tasks of service_name.instance_name
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param service Service name
 * @param instance Instance name
 * @return ApiTasksInstanceRequest
 */
func (a *ServiceApiService) TasksInstance(ctx _context.Context, service string, instance string) ApiTasksInstanceRequest {
	return ApiTasksInstanceRequest{
		ApiService: a,
		ctx: ctx,
		service: service,
		instance: instance,
	}
}

/*
 * Execute executes the request
 * @return []map[string]interface{}
 */
func (a *ServiceApiService) TasksInstanceExecute(r ApiTasksInstanceRequest) ([]map[string]interface{}, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  []map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ServiceApiService.TasksInstance")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/services/{service}/{instance}/tasks"
	localVarPath = strings.Replace(localVarPath, "{"+"service"+"}", _neturl.PathEscape(parameterToString(r.service, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instance"+"}", _neturl.PathEscape(parameterToString(r.instance, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.slaveHostname != nil {
		localVarQueryParams.Add("slave_hostname", parameterToString(*r.slaveHostname, ""))
	}
	if r.verbose != nil {
		localVarQueryParams.Add("verbose", parameterToString(*r.verbose, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
