// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"fmt"
)

// Accounts runs ending category [and Balance and P/L] calculations
// based on the records passed in CSV data files.
func Accounts(q Schema) ([]Categories, NoticeOfError) {
	var (
		alert NoticeOfError
		cats  []Categories
		recs  []Transactions
	)

	// Note: const Headers bool = true
	cats, alert = gatherCategories(q, Headers)
	if alert.Error != nil {
		alert.Trace.Crumbs("Accounts")
		fmt.Printf("Trail (%v): %v\n", len(alert.Trace.x), alert.Trace)
		return nil, alert
	}
	// fmt.Println("Total categories read:", len(cats))

	recs, alert = gatherTransactions(q, Headers)
	if alert.Error != nil {
		alert.Trace.Crumbs("Accounts")
		fmt.Printf("Trail (%v): %v\n", len(alert.Trace.x), alert.Trace)
		return nil, alert
	}
	// fmt.Println("Total records read:", len(recs))

	conti, err := postTransactionsToAccounts(cats, recs)
	if err != nil {
		alert = NoticeOfError{
			Code:  CaseInnerError,
			Hint:  "Send this error to the program developer",
			Error: err,
		}
		alert.Trace.Crumbs("postTransactionsToAccounts")
		return nil, alert
	}

	// fmt.Println("Categories in results:", len(conti))

	return conti, alert
}
