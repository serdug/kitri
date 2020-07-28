// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"encoding/json"
	"net/http"
)

const Headers bool = true

// Schema represents the structure from JSON received in the client request body
type Schema struct {
	// Working directory
	Path string `json:"path"`

	// Names of CSV files containing the Chart of Accounts, a file per section
	Chart struct {
		Assets      string `json:"assets"`
		Liabilities string `json:"liabilities"`
		Equity      string `json:"equity"`
		Revenues    string `json:"revenues"`
		Expenses    string `json:"expenses"`
	}

	// A list of CSV files containing records of transactions to be processed
	Records []Record `json:"records"`
}

type Record struct {
	Include int
	Id      string
}

// DecodeSchema parses a JSON payload from request body.
func DecodeSchema(r *http.Request) (s Schema) {
	dec := json.NewDecoder(r.Body)

	// Decode a schema
	err := dec.Decode(&s)
	if err != nil {
		warning("JSON unmarshal failed!", err)
	}
	return
}
