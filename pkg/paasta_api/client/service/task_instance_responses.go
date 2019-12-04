// Code generated by go-swagger; DO NOT EDIT.

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/Yelp/paasta-tools-go/pkg/paasta_api/models"
)

// TaskInstanceReader is a Reader for the TaskInstance structure.
type TaskInstanceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *TaskInstanceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewTaskInstanceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewTaskInstanceBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewTaskInstanceNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewTaskInstanceInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewTaskInstanceOK creates a TaskInstanceOK with default headers values
func NewTaskInstanceOK() *TaskInstanceOK {
	return &TaskInstanceOK{}
}

/*TaskInstanceOK handles this case with default header values.

Task associated with an instance with specified ID
*/
type TaskInstanceOK struct {
	Payload models.InstanceTask
}

func (o *TaskInstanceOK) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/tasks/{task_id}][%d] taskInstanceOK  %+v", 200, o.Payload)
}

func (o *TaskInstanceOK) GetPayload() models.InstanceTask {
	return o.Payload
}

func (o *TaskInstanceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewTaskInstanceBadRequest creates a TaskInstanceBadRequest with default headers values
func NewTaskInstanceBadRequest() *TaskInstanceBadRequest {
	return &TaskInstanceBadRequest{}
}

/*TaskInstanceBadRequest handles this case with default header values.

Bad request
*/
type TaskInstanceBadRequest struct {
}

func (o *TaskInstanceBadRequest) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/tasks/{task_id}][%d] taskInstanceBadRequest ", 400)
}

func (o *TaskInstanceBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewTaskInstanceNotFound creates a TaskInstanceNotFound with default headers values
func NewTaskInstanceNotFound() *TaskInstanceNotFound {
	return &TaskInstanceNotFound{}
}

/*TaskInstanceNotFound handles this case with default header values.

Task with ID not found
*/
type TaskInstanceNotFound struct {
}

func (o *TaskInstanceNotFound) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/tasks/{task_id}][%d] taskInstanceNotFound ", 404)
}

func (o *TaskInstanceNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewTaskInstanceInternalServerError creates a TaskInstanceInternalServerError with default headers values
func NewTaskInstanceInternalServerError() *TaskInstanceInternalServerError {
	return &TaskInstanceInternalServerError{}
}

/*TaskInstanceInternalServerError handles this case with default header values.

Instance failure
*/
type TaskInstanceInternalServerError struct {
}

func (o *TaskInstanceInternalServerError) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/tasks/{task_id}][%d] taskInstanceInternalServerError ", 500)
}

func (o *TaskInstanceInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
