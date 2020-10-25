// Code generated by go-swagger; DO NOT EDIT.

package broadcast

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/radiorabe/annotation-agent/annotator/archive/raar/models"
)

// GetBroadcastsYearMonthDayReader is a Reader for the GetBroadcastsYearMonthDay structure.
type GetBroadcastsYearMonthDayReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetBroadcastsYearMonthDayReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetBroadcastsYearMonthDayOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetBroadcastsYearMonthDayOK creates a GetBroadcastsYearMonthDayOK with default headers values
func NewGetBroadcastsYearMonthDayOK() *GetBroadcastsYearMonthDayOK {
	return &GetBroadcastsYearMonthDayOK{}
}

/*GetBroadcastsYearMonthDayOK handles this case with default header values.

successfull operation
*/
type GetBroadcastsYearMonthDayOK struct {
	Payload *GetBroadcastsYearMonthDayOKBody
}

func (o *GetBroadcastsYearMonthDayOK) Error() string {
	return fmt.Sprintf("[GET /broadcasts/{year}/{month}/{day}][%d] getBroadcastsYearMonthDayOK  %+v", 200, o.Payload)
}

func (o *GetBroadcastsYearMonthDayOK) GetPayload() *GetBroadcastsYearMonthDayOKBody {
	return o.Payload
}

func (o *GetBroadcastsYearMonthDayOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetBroadcastsYearMonthDayOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetBroadcastsYearMonthDayOKBody get broadcasts year month day o k body
swagger:model GetBroadcastsYearMonthDayOKBody
*/
type GetBroadcastsYearMonthDayOKBody struct {

	// data
	Data []*models.Broadcast `json:"data"`

	// included
	Included []*models.Show `json:"included"`
}

// Validate validates this get broadcasts year month day o k body
func (o *GetBroadcastsYearMonthDayOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateIncluded(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetBroadcastsYearMonthDayOKBody) validateData(formats strfmt.Registry) error {

	if swag.IsZero(o.Data) { // not required
		return nil
	}

	for i := 0; i < len(o.Data); i++ {
		if swag.IsZero(o.Data[i]) { // not required
			continue
		}

		if o.Data[i] != nil {
			if err := o.Data[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getBroadcastsYearMonthDayOK" + "." + "data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *GetBroadcastsYearMonthDayOKBody) validateIncluded(formats strfmt.Registry) error {

	if swag.IsZero(o.Included) { // not required
		return nil
	}

	for i := 0; i < len(o.Included); i++ {
		if swag.IsZero(o.Included[i]) { // not required
			continue
		}

		if o.Included[i] != nil {
			if err := o.Included[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getBroadcastsYearMonthDayOK" + "." + "included" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetBroadcastsYearMonthDayOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetBroadcastsYearMonthDayOKBody) UnmarshalBinary(b []byte) error {
	var res GetBroadcastsYearMonthDayOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
