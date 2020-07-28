// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"fmt"
)

// Status codes
const (
	CaseNotFound         = "Resource not found"
	CaseUnreadable       = "Unexpected type of or damaged file"
	CaseCategoryNotKnown = "Unrecognizable category"
	CaseNoData           = "No data provided"
	CaseUnnamedFile      = "No file name"
	CaseWrongFileType    = "Wrong file type"
	CaseWrongFormat      = "Wrong data format"
	CaseInnerError       = "Internal program error"
)

// NoticeOfError provides a structure for user guidance if calculation has gone not as
// expected
type NoticeOfError struct {
	// A code of the state of calculation returned to client
	Code string

	// Trace collects breadcrumbs; it's an array of executed process names
	// (i.e. functions, methods etc), so that an output error can be traced back
	// to the process that owns of the fault.
	// Note: The service is being designed to help users identify what could be
	// changed to attain better result. So it would be helpful to tell what
	// exactly has lead to the problem.
	Trace Trail

	// ID of the user-provided resource where the problem occurred
	Resource string

	// Note: Insight is expected. Practical information or handy advice should
	// be provided
	Hint string

	// Standard error message
	Error error
}

// Trail represents an array of marks
type Trail struct {
	x []string
}

// warning prints an error message; it does not cause a process to end.
func warning(msg string, e error) {
	if e != nil {
		fmt.Printf("%s [%v]\n", msg, e)
	}
}

// ***************************************************************************
// * METHODS
// ***************************************************************************

// Crumbs collects marks of execution, i.e. executed process' IDs
func (t *Trail) Crumbs(mark string) {
	t.x = append(t.x, mark)
}
