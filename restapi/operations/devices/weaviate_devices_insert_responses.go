/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 Weaviate. All rights reserved.
 * LICENSE: https://github.com/weaviate/weaviate/blob/master/LICENSE
 * AUTHOR: Bob van Luijt (bob@weaviate.com)
 * See www.weaviate.com for details
 * See package.json for author and maintainer info
 * Contact: @weaviate_iot / yourfriends@weaviate.com
 */
 package devices

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/weaviate/weaviate/models"
)

// WeaviateDevicesInsertCreatedCode is the HTTP code returned for type WeaviateDevicesInsertCreated
const WeaviateDevicesInsertCreatedCode int = 201

/*WeaviateDevicesInsertCreated Successful created.

swagger:response weaviateDevicesInsertCreated
*/
type WeaviateDevicesInsertCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Device `json:"body,omitempty"`
}

// NewWeaviateDevicesInsertCreated creates WeaviateDevicesInsertCreated with default headers values
func NewWeaviateDevicesInsertCreated() *WeaviateDevicesInsertCreated {
	return &WeaviateDevicesInsertCreated{}
}

// WithPayload adds the payload to the weaviate devices insert created response
func (o *WeaviateDevicesInsertCreated) WithPayload(payload *models.Device) *WeaviateDevicesInsertCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the weaviate devices insert created response
func (o *WeaviateDevicesInsertCreated) SetPayload(payload *models.Device) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WeaviateDevicesInsertCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WeaviateDevicesInsertNotImplementedCode is the HTTP code returned for type WeaviateDevicesInsertNotImplemented
const WeaviateDevicesInsertNotImplementedCode int = 501

/*WeaviateDevicesInsertNotImplemented Not (yet) implemented.

swagger:response weaviateDevicesInsertNotImplemented
*/
type WeaviateDevicesInsertNotImplemented struct {
}

// NewWeaviateDevicesInsertNotImplemented creates WeaviateDevicesInsertNotImplemented with default headers values
func NewWeaviateDevicesInsertNotImplemented() *WeaviateDevicesInsertNotImplemented {
	return &WeaviateDevicesInsertNotImplemented{}
}

// WriteResponse to the client
func (o *WeaviateDevicesInsertNotImplemented) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(501)
}