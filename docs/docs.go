/*
Package classification Money Processing Service

Documentation Money Processing Service API


Schemes: http
BasePath: /api/v1
Version: 2.0.0

Consumes:
- application/json

Produces:
- application/json
swagger:meta
*/
package docs

import "github.com/brxyxn/go_mps_redcage/data"

// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

/* Creation */

/*
Created item returns the details of the client
swagger:response clientResponse
*/
type clientResponse struct {
	// Response in body
	// in: body
	// required: true
	Body data.Client
}

/*
Created item returns the details of the account
swagger:response accountResponse
*/
type accountResponse struct {
	// in: body
	// required: true
	Body data.Account
}

/*
Created item returns the details of the transaction
swagger:response transactionResponse
*/
type transactionResponse struct {
	// in: body
	// required: true
	Body data.Transaction
}

/*-------------*/

/*
Returns the details of the requested account
swagger:response accountGetResponse
*/
type accountGetResponse struct {
	// in: body
	// required: true
	Body data.Account
}

/*-------------*/

// Data structure representing a list of accounts
// swagger:response accountsResponse
type accountsResponseWrapper struct {
	// All existing accounts
	// in: body
	Body data.Accounts
}

// Data structure representing a list of transactions
// swagger:response transactionsResponse
type transactionsResponseWrapper struct {
	// All existing transactions
	// in: body
	Body data.Transactions
}

/*---------*/

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Bad Request Error Response
// swagger:response badRequestErrorResponse
type badRequestErrorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Not Found Error Response
// swagger:response notFoundErrorResponse
type notFoundErrorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Internal Error Response
// swagger:response internalErrorResponse
type internalErrorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Generic Error is a generic error message returned by a server
type GenericError struct {
	// "error":"message"
	Error string `json:"error"`
}
