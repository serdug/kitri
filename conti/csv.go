// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings" // to split strings
)

var (
	attributes string = "Cat+Sect+Name+Starting+Change+Ending"
)

// TO DO: TRACE errors!!!
// ExportAccountsToCsv writes results of value-by-category
// calculations in a CSV file
func ExportAccountsToCsv(conti []Categories, filename string) {
	var (
		err1  error
		field []string
	)

	csvNewFile, err := os.OpenFile(filename, os.O_RDWR, 0666) // os.O_RDONLY, os.O_WRONLY, os.O_RDWR
	//csvNewFile, err := os.Open(filename) // READ-ONLY!!!
	if err != nil {
		//fmt.Println("Open error:", err)
		csvNewFile, err1 = os.Create(filename)
		if err1 != nil {
			fmt.Printf("Error (creating file %s): %s\n", filename, err1)
		}
	}
	defer csvNewFile.Close()

	writer := csv.NewWriter(csvNewFile)

	headers := strings.Split(attributes, "+")
	writer.Write(headers)

	for _, one := range conti {
		field = make([]string, len(headers))

		field[0] = one.Cat
		field[1] = one.Sect
		field[2] = one.Name
		// Converting response to a single string
		field[3] = fmt.Sprintf("%f", one.Bal.Sta)
		field[4] = fmt.Sprintf("%f", one.Bal.Dif)
		field[5] = fmt.Sprintf("%f", one.Bal.End)

		writer.Write(field)
	}

	// remember to flush!
	writer.Flush()
	return
}

// readFileCsv reads data from a csv file into a [][]string matrix
func readFileCsv(filename string) ([][]string, NoticeOfError) {
	var (
		alert NoticeOfError
		mx    [][]string
	)

	f, errOpen := os.Open(filename)
	if errOpen != nil {
		alert = NoticeOfError{
			Code:     CaseNotFound,
			Resource: filename,
			Hint:     "File not found: " + filename,
			Error:    errOpen,
		}
		alert.Trace.Crumbs("readFileCsv")
		return mx, alert
	}

	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	mx, errRead := reader.ReadAll()
	if errRead != nil {
		alert = NoticeOfError{
			Code:     CaseUnreadable,
			Resource: filename,
			Hint:     "Failed to read: " + filename,
			Error:    errRead,
		}
		alert.Trace.Crumbs("readFileCsv")
		os.Exit(1)
		return mx, alert
	}
	return mx, alert
}
