// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Categories represents The Chart of Accounts, a sctructure of accounts.
type Categories struct {
	// Unique (!) Category id
	Cat string

	// Section, such as Assets, Liabilities, Equity, Revenues, Expenses etc
	Sect string

	// Category name
	Name string

	// Balance value, e.g. opening (starting) or closing balance
	// Note: it's optional for revenue and expense accounts
	Bal Tally
}

type Tally struct {
	Sta float64 // Starting
	Dif float64 // Change
	End float64 // Ending
}

// Transactions for accounting book entries
type Transactions struct {
	// Amount
	Amount float64

	// The source of funds, i.e. an id of account (category) from which money is
	// spent (paid / invested / lent) or received (earned / borrowed)
	Source string

	// The purpose (or reason) for paying for smth or the use of receipts,
	// i.e. an id of account (category) to which the money is purposed
	Purpose string
}

// gatherTransactions reads records from CSV data files into a slice of
// Transactions objects. No data validation.
func gatherTransactions(q Schema, headers bool) ([]Transactions, NoticeOfError) {
	var (
		file      string
		alert     NoticeOfError
		rec, recs []Transactions
		raw       [][]string
	)

	for _, record := range q.Records {
		// Exclude files if Include == 0
		if record.Include == 0 {
			// Go to the next iteration: skip record, don't count this file
			continue
		}

		// Exclude unnamed files
		file, alert = fileType(record.Id)

		if alert.Error != nil {
			alert.Trace.Crumbs("gatherTransactions")
			return recs, alert
		}

		if alert.Code == CaseUnnamedFile {
			// Go to the next iteration: skip record, don't count this file,
			// reset alert
			alert = NoticeOfError{}
			continue
		}

		raw, alert = file2mx(filepath.Join(q.Path, file), headers)
		// Note: Alternatively, use Join() from path/filepath
		// raw, alert = file2mx(q.Path + file, headers)
		if alert.Error != nil {
			alert.Trace.Crumbs("gatherTransactions")
			return recs, alert
		}

		rec, alert = readTransactions(raw)
		if alert.Error != nil {
			alert.Trace.Crumbs("gatherTransactions")
			return recs, alert
		}

		recs = append(recs, rec...)
	}

	return recs, alert
}

// gatherCategories reads records from CSV data files (arranged by preset
// Sections) into a slice of Transactions objects containing Sections.
// No data validation.
func gatherCategories(q Schema, headers bool) ([]Categories, NoticeOfError) {
	var (
		file  string
		alert NoticeOfError
		cats  []Categories
		cat   []Categories
		raw   [][]string
	)

	// using Join() from path/filepath
	sections := map[string]string{
		"Assets":      filepath.Join(q.Path, q.Chart.Assets),
		"Liabilities": filepath.Join(q.Path, q.Chart.Liabilities),
		"Equity":      filepath.Join(q.Path, q.Chart.Equity),
		"Revenues":    filepath.Join(q.Path, q.Chart.Revenues),
		"Expenses":    filepath.Join(q.Path, q.Chart.Expenses),
		/*
			"Assets":      q.Path + q.Chart.Assets,
			"Liabilities": q.Path + q.Chart.Liabilities,
			"Equity":      q.Path + q.Chart.Equity,
			"Revenues":    q.Path + q.Chart.Revenues,
			"Expenses":    q.Path + q.Chart.Expenses,
		*/
	}

	for i := range sections {
		file, alert = fileType(sections[i])
		/*
			if alert.Error == nil || alert.Code != "" {
				alert.Trace.Crumbs("gatherCategories")
				return cats, alert
			}
		*/

		if alert.Error != nil || alert.Code == CaseUnnamedFile {
			alert.Trace.Crumbs("gatherCategories")
			return cats, alert
		}

		raw, alert = file2mx(file, headers)
		if alert.Error != nil {
			alert.Trace.Crumbs("gatherCategories")
			return cats, alert
		}
		cat, alert = mx2cats(raw, i)
		if alert.Error != nil {
			alert.Trace.Crumbs("gatherCategories")
			return cats, alert
		}

		cats = append(cats, cat...)
	}

	// Sort categories keeping original order of equal elements.
	sort.SliceStable(cats, func(i, j int) bool {
		return cats[i].Cat < cats[j].Cat
	})

	return cats, alert
}

// fileType detects whether the provided file name has a '.csv' extension and
// either:
// (1) adds '.csv' to the file name if no extension is provided or
// (2) returns a warning of a wrong file type
func fileType(filename string) (string, NoticeOfError) {
	var (
		nameFull string
		alert    NoticeOfError
	)
	// Split path to distinguish the file name
	// Note: The values returned by Split() have the property
	// that path=dir+file
	// Otherwise, it's possible to correctly join elements of the path
	// using Join() from path/filepath
	d, f := filepath.Split(filename)
	ext := strings.ToLower(filepath.Ext(filename))

	switch {
	case ext == ".csv":
		return filename, alert

	case f != "" && ext == "":
		// No extension, add '.csv'
		nameFull = filename + ".csv"
		fmt.Println("Extended: ", nameFull)
		return nameFull, alert

	case f != "" && ext != ".csv":
		fmt.Println("Unrecognized extension! ", filename)
		alert = NoticeOfError{
			Code: CaseWrongFileType,
			Hint: "File '" + filename + "' has an unacceptable extension '" + ext + "'",
		}
		alert.Trace.Crumbs("fileType")
		return filename, alert

	case d != "" && f == "":
		// Empty file name provided
		alert = NoticeOfError{
			Code: CaseUnnamedFile,
		}
		fmt.Println("Unnamed records file, skipped it!")
		return filename, alert

	case d == "" && f == "":
		// No path and no file name provided
		alert = NoticeOfError{
			Code: CaseNoData,
			Hint: "No path and no file name provided",
		}
		fmt.Println("No path and no file name provided!")
		return filename, alert

	default:
		// Anything else
		fmt.Println("Oops... ", filename)
		return filename, alert
	}
}

// file2mx reads file and puts CSV data into a [][]string matrix (raws-columns)
func file2mx(filename string, headers bool) ([][]string, NoticeOfError) {
	var (
		alert NoticeOfError
		mx    [][]string
	)
	mx, alert = readFileCsv(filename)
	if alert.Error != nil {
		alert.Trace.Crumbs("file2mx")
		return mx, alert
	}

	switch headers {
	case true:

		// Delete the title row (first element from the slice)
		mx = append(mx[:0], mx[1:]...)
		return mx, alert

	default:
		return mx, alert
	}
}

// mx2cats puts data from a matrix of read input into a slice of
// Categories objects. No data validation.
func mx2cats(raw [][]string, section string) ([]Categories, NoticeOfError) {
	var (
		alert NoticeOfError
		one   Categories
		all   []Categories
	)
	all = make([]Categories, len(raw))
	for i, each := range raw {
		bal, err := strconv.ParseFloat(each[2], 64) //col3
		if err != nil {
			alert = NoticeOfError{
				Code:  CaseWrongFormat,
				Error: err,
				Hint:  "WARNING! Balance '" + each[2] + "' read as '" + fmt.Sprintf("%f", bal) + "'",
			}
			alert.Trace.Crumbs("mx2cats")
		}
		one = Categories{
			Cat:  each[0], //col1
			Sect: section,
			Name: each[1], //col2
			Bal: Tally{
				Sta: bal, //col3
			},
		}
		all[i] = one
	}
	return all, alert
}

// readTransactions puts data from a matrix of read input into a slice of
// Transactions objects. No data validation.
func readTransactions(mx [][]string) ([]Transactions, NoticeOfError) {
	var (
		alert NoticeOfError
		one   Transactions
		all   []Transactions
	)
	all = make([]Transactions, len(mx))
	for i, each := range mx {
		// The monetary value of transaction must be in the 1st column
		amount, err := strconv.ParseFloat(each[0], 64) //col1
		if err != nil {
			alert = NoticeOfError{
				Code:  CaseWrongFormat,
				Error: err,
				Hint:  "WARNING! Amount '" + each[0] + "' read as '" + fmt.Sprintf("%f", amount) + "'",
			}
			alert.Trace.Crumbs("readTransactions")
		}

		one = Transactions{
			Amount:  amount,  //col1
			Source:  each[1], //col2
			Purpose: each[2], //col3
		}
		all[i] = one
	}
	return all, alert
}
