// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package handlers provides functions serving client requests
package handlers

import (
	"encoding/json"
	// "fmt"
	"io/ioutil" // to read files

	"gopkg.in/yaml.v2"

	"github.com/serdug/kitri/conti"
)

// ReadSchemaJSON parses a JSON schema file
func ReadSchemaJSON(filename string) conti.Schema {
	var s conti.Schema
	// fmt.Println("\nConfig file read:", filename)

	// read the config file
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		// Handle error
	}

	err = json.Unmarshal(dat, &s)
	if err != nil {
		// Handle error
	}
	return s
}

// ReadSchemaYAML parses a YAML schema file
func ReadSchemaYAML(filename string) conti.Schema {
	var s conti.Schema
	// fmt.Println("\nConfig file read:", filename)

	// read the config file
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		// Handle error
	}

	err = yaml.Unmarshal(dat, &s)
	if err != nil {
		// Handle error
	}
	return s
}
