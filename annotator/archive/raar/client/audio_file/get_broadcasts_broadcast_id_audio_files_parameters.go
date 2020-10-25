// Code generated by go-swagger; DO NOT EDIT.

package audio_file

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

// NewGetBroadcastsBroadcastIDAudioFilesParams creates a new GetBroadcastsBroadcastIDAudioFilesParams object
// with the default values initialized.
func NewGetBroadcastsBroadcastIDAudioFilesParams() *GetBroadcastsBroadcastIDAudioFilesParams {
	var ()
	return &GetBroadcastsBroadcastIDAudioFilesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetBroadcastsBroadcastIDAudioFilesParamsWithTimeout creates a new GetBroadcastsBroadcastIDAudioFilesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetBroadcastsBroadcastIDAudioFilesParamsWithTimeout(timeout time.Duration) *GetBroadcastsBroadcastIDAudioFilesParams {
	var ()
	return &GetBroadcastsBroadcastIDAudioFilesParams{

		timeout: timeout,
	}
}

// NewGetBroadcastsBroadcastIDAudioFilesParamsWithContext creates a new GetBroadcastsBroadcastIDAudioFilesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetBroadcastsBroadcastIDAudioFilesParamsWithContext(ctx context.Context) *GetBroadcastsBroadcastIDAudioFilesParams {
	var ()
	return &GetBroadcastsBroadcastIDAudioFilesParams{

		Context: ctx,
	}
}

// NewGetBroadcastsBroadcastIDAudioFilesParamsWithHTTPClient creates a new GetBroadcastsBroadcastIDAudioFilesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetBroadcastsBroadcastIDAudioFilesParamsWithHTTPClient(client *http.Client) *GetBroadcastsBroadcastIDAudioFilesParams {
	var ()
	return &GetBroadcastsBroadcastIDAudioFilesParams{
		HTTPClient: client,
	}
}

/*GetBroadcastsBroadcastIDAudioFilesParams contains all the parameters to send to the API endpoint
for the get broadcasts broadcast ID audio files operation typically these are written to a http.Request
*/
type GetBroadcastsBroadcastIDAudioFilesParams struct {

	/*BroadcastID
	  Id of the broadcast to list the audio files for.

	*/
	BroadcastID int64
	/*PageNumber
	  The page number of the list.

	*/
	PageNumber *int64
	/*PageSize
	  Maximum number of entries that are returned per page. Defaults to 50, maximum is 500.

	*/
	PageSize *int64
	/*Sort
	  Name of the sort field, optionally prefixed with a `-` for descending order.

	*/
	Sort *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithTimeout(timeout time.Duration) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithContext(ctx context.Context) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithHTTPClient(client *http.Client) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBroadcastID adds the broadcastID to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithBroadcastID(broadcastID int64) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetBroadcastID(broadcastID)
	return o
}

// SetBroadcastID adds the broadcastId to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetBroadcastID(broadcastID int64) {
	o.BroadcastID = broadcastID
}

// WithPageNumber adds the pageNumber to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithPageNumber(pageNumber *int64) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetPageNumber(pageNumber)
	return o
}

// SetPageNumber adds the pageNumber to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetPageNumber(pageNumber *int64) {
	o.PageNumber = pageNumber
}

// WithPageSize adds the pageSize to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithPageSize(pageSize *int64) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetPageSize(pageSize *int64) {
	o.PageSize = pageSize
}

// WithSort adds the sort to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WithSort(sort *string) *GetBroadcastsBroadcastIDAudioFilesParams {
	o.SetSort(sort)
	return o
}

// SetSort adds the sort to the get broadcasts broadcast ID audio files params
func (o *GetBroadcastsBroadcastIDAudioFilesParams) SetSort(sort *string) {
	o.Sort = sort
}

// WriteToRequest writes these params to a swagger request
func (o *GetBroadcastsBroadcastIDAudioFilesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param broadcast_id
	if err := r.SetPathParam("broadcast_id", swag.FormatInt64(o.BroadcastID)); err != nil {
		return err
	}

	if o.PageNumber != nil {

		// query param page[number]
		var qrPageNumber int64
		if o.PageNumber != nil {
			qrPageNumber = *o.PageNumber
		}
		qPageNumber := swag.FormatInt64(qrPageNumber)
		if qPageNumber != "" {
			if err := r.SetQueryParam("page[number]", qPageNumber); err != nil {
				return err
			}
		}

	}

	if o.PageSize != nil {

		// query param page[size]
		var qrPageSize int64
		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt64(qrPageSize)
		if qPageSize != "" {
			if err := r.SetQueryParam("page[size]", qPageSize); err != nil {
				return err
			}
		}

	}

	if o.Sort != nil {

		// query param sort
		var qrSort string
		if o.Sort != nil {
			qrSort = *o.Sort
		}
		qSort := qrSort
		if qSort != "" {
			if err := r.SetQueryParam("sort", qSort); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
