// Code generated by go-swagger; DO NOT EDIT.

package autoscaler

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetAutoscalerCountReader is a Reader for the GetAutoscalerCount structure.
type GetAutoscalerCountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAutoscalerCountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAutoscalerCountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetAutoscalerCountNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAutoscalerCountInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAutoscalerCountOK creates a GetAutoscalerCountOK with default headers values
func NewGetAutoscalerCountOK() *GetAutoscalerCountOK {
	return &GetAutoscalerCountOK{}
}

/*GetAutoscalerCountOK handles this case with default header values.

Get desired instance count for a service instance
*/
type GetAutoscalerCountOK struct {
	Payload *GetAutoscalerCountOKBody
}

func (o *GetAutoscalerCountOK) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/autoscaler][%d] getAutoscalerCountOK  %+v", 200, o.Payload)
}

func (o *GetAutoscalerCountOK) GetPayload() *GetAutoscalerCountOKBody {
	return o.Payload
}

func (o *GetAutoscalerCountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetAutoscalerCountOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAutoscalerCountNotFound creates a GetAutoscalerCountNotFound with default headers values
func NewGetAutoscalerCountNotFound() *GetAutoscalerCountNotFound {
	return &GetAutoscalerCountNotFound{}
}

/*GetAutoscalerCountNotFound handles this case with default header values.

Deployment key not found
*/
type GetAutoscalerCountNotFound struct {
}

func (o *GetAutoscalerCountNotFound) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/autoscaler][%d] getAutoscalerCountNotFound ", 404)
}

func (o *GetAutoscalerCountNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetAutoscalerCountInternalServerError creates a GetAutoscalerCountInternalServerError with default headers values
func NewGetAutoscalerCountInternalServerError() *GetAutoscalerCountInternalServerError {
	return &GetAutoscalerCountInternalServerError{}
}

/*GetAutoscalerCountInternalServerError handles this case with default header values.

Instance failure
*/
type GetAutoscalerCountInternalServerError struct {
}

func (o *GetAutoscalerCountInternalServerError) Error() string {
	return fmt.Sprintf("[GET /services/{service}/{instance}/autoscaler][%d] getAutoscalerCountInternalServerError ", 500)
}

func (o *GetAutoscalerCountInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*GetAutoscalerCountOKBody get autoscaler count o k body
swagger:model GetAutoscalerCountOKBody
*/
type GetAutoscalerCountOKBody struct {

	// desired instances
	DesiredInstances int64 `json:"desired_instances,omitempty"`
}

// Validate validates this get autoscaler count o k body
func (o *GetAutoscalerCountOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetAutoscalerCountOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetAutoscalerCountOKBody) UnmarshalBinary(b []byte) error {
	var res GetAutoscalerCountOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}