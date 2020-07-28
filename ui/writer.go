// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package ui for GUI (front end)
package ui

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"fyne.io/fyne"

	"github.com/serdug/kitri/conti"
)

// fileNamed extracts the file name from fyne's 'URIWriteCloser' object,
// removes the 'file://' added by fyne from the writer.URI() string and
// returns a string containing a refined file name
func fileNamed(writer fyne.URIWriteCloser) string {
	p := fmt.Sprintf("%s", writer.URI())

	// Remove "file://" from writer.URI() added by fyne
	return strings.Replace(p, "file://", "", 1)
}

// schemaWriter saves schema template as a '.yaml' file
func schemaWriter(kit kitri, name string) {
	ext := filepath.Ext(name)
	ext = strings.ToLower(ext)

	switch {
	case ext == ".":
		writeSchemaYAML(kit, name+"yaml")
		fmt.Println("Template saved to", name+"yaml")

	case len(ext) == 0:
		writeSchemaYAML(kit, name+".yaml")
		fmt.Println("Template saved to", name+".yaml")

	case ext == ".yaml" || ext == ".yml":
		// Save with a user-typed name
		writeSchemaYAML(kit, name)
		fmt.Println("Template saved to", name)

	default:
		fmt.Printf("File '" + name +
			"' has an unacceptable extension '" + ext +
			"'\nTemplates are only saved as '.yaml' (or '.yml') files.\nPlease set a file name without extension or type it with a YAML extension.")
		return
	}
}

// outputWriter recalculates and saves results as a '.csv' file
func outputWriter(kit kitri, name string) {
	schema := templateSchema(kit)
	cats, _ := conti.Accounts(schema)

	ext := filepath.Ext(name)
	ext = strings.ToLower(ext)

	switch {
	case ext == ".":
		conti.ExportAccountsToCsv(cats, name+"csv")
		fmt.Println("Output saved to", name+"csv")

	case len(ext) == 0:
		conti.ExportAccountsToCsv(cats, name+".csv")
		fmt.Println("Output saved to", name+".csv")

	case ext == ".csv":
		// Save with a user-typed name
		conti.ExportAccountsToCsv(cats, name)
		fmt.Println("Output saved to", name)

	default:
		fmt.Printf("File '" + name +
			"' has an unacceptable extension '" + ext +
			"'\nResults are only saved as '.csv' files.\nPlease set a file name without extension or type it with a CSV extension.")
		return
	}
}

// writeSchemaYAML serializes and writes a schema template into a YAML file
// with a provided name
func writeSchemaYAML(kit kitri, filename string) {
	s := templateSchema(kit)

	// Serialize schema into a YAML document
	data, err := yaml.Marshal(s)

	if err != nil {
		fmt.Println("YAML Marshal error:", err)
		return
	}

	erw := ioutil.WriteFile(filename, data, 0644)
	if erw != nil {
		fmt.Println("Writing error:", erw)
		return
	}
}
