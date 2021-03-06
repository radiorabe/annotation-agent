// Code generated by go-swagger; DO NOT EDIT.

package broadcast

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/radiorabe/annotation-agent/annotator/archive/raar/models"
)

// GetBroadcastsIDReader is a Reader for the GetBroadcastsID structure.
type GetBroadcastsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetBroadcastsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetBroadcastsIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetBroadcastsIDOK creates a GetBroadcastsIDOK with default headers values
func NewGetBroadcastsIDOK() *GetBroadcastsIDOK {
	return &GetBroadcastsIDOK{}
}

/*GetBroadcastsIDOK handles this case with default header values.

successfull operation
*/
type GetBroadcastsIDOK struct {
	Payload *GetBroadcastsIDOKBody
}

func (o *GetBroadcastsIDOK) Error() string {
	return fmt.Sprintf("[GET /broadcasts/{id}][%d] getBroadcastsIdOK  %+v", 200, o.Payload)
}

func (o *GetBroadcastsIDOK) GetPayload() *GetBroadcastsIDOKBody {
	return o.Payload
}

func (o *GetBroadcastsIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetBroadcastsIDOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetBroadcastsIDOKBody get broadcasts ID o k body
swagger:model GetBroadcastsIDOKBody
*/
type GetBroadcastsIDOKBody struct {

	// data
	Data *models.Broadcast `json:"data,omitempty"`
}

// Validate validates this get broadcasts ID o k body
func (o *GetBroadcastsIDOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetBroadcastsIDOKBody) validateData(formats strfmt.Registry) error {

	if swag.IsZero(o.Data) { // not required
		return nil
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getBroadcastsIdOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetBroadcastsIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetBroadcastsIDOKBody) UnmarshalBinary(b []byte) error {
	var res GetBroadcastsIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
