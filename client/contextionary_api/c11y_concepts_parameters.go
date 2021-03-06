//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2020 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

// Code generated by go-swagger; DO NOT EDIT.

package contextionary_api

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

// NewC11yConceptsParams creates a new C11yConceptsParams object
// with the default values initialized.
func NewC11yConceptsParams() *C11yConceptsParams {
	var ()
	return &C11yConceptsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewC11yConceptsParamsWithTimeout creates a new C11yConceptsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewC11yConceptsParamsWithTimeout(timeout time.Duration) *C11yConceptsParams {
	var ()
	return &C11yConceptsParams{

		timeout: timeout,
	}
}

// NewC11yConceptsParamsWithContext creates a new C11yConceptsParams object
// with the default values initialized, and the ability to set a context for a request
func NewC11yConceptsParamsWithContext(ctx context.Context) *C11yConceptsParams {
	var ()
	return &C11yConceptsParams{

		Context: ctx,
	}
}

// NewC11yConceptsParamsWithHTTPClient creates a new C11yConceptsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewC11yConceptsParamsWithHTTPClient(client *http.Client) *C11yConceptsParams {
	var ()
	return &C11yConceptsParams{
		HTTPClient: client,
	}
}

/*C11yConceptsParams contains all the parameters to send to the API endpoint
for the c11y concepts operation typically these are written to a http.Request
*/
type C11yConceptsParams struct {

	/*Concept
	  CamelCase list of words to validate.

	*/
	Concept string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the c11y concepts params
func (o *C11yConceptsParams) WithTimeout(timeout time.Duration) *C11yConceptsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the c11y concepts params
func (o *C11yConceptsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the c11y concepts params
func (o *C11yConceptsParams) WithContext(ctx context.Context) *C11yConceptsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the c11y concepts params
func (o *C11yConceptsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the c11y concepts params
func (o *C11yConceptsParams) WithHTTPClient(client *http.Client) *C11yConceptsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the c11y concepts params
func (o *C11yConceptsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithConcept adds the concept to the c11y concepts params
func (o *C11yConceptsParams) WithConcept(concept string) *C11yConceptsParams {
	o.SetConcept(concept)
	return o
}

// SetConcept adds the concept to the c11y concepts params
func (o *C11yConceptsParams) SetConcept(concept string) {
	o.Concept = concept
}

// WriteToRequest writes these params to a swagger request
func (o *C11yConceptsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param concept
	if err := r.SetPathParam("concept", o.Concept); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
