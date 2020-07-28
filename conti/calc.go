// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"fmt"
	"math"
)

const decimals float64 = 10000

// postTransactionsToAccounts calculates the balance (for balance categories and the cumulative
// sums for P/L categories) per category on the basis of records of transactions.
func postTransactionsToAccounts(catsIn []Categories, records []Transactions) (catsOut []Categories, err error) {
	var (
		sVal   map[string]float64
		cVal   map[string]float64
		pVal   *map[string]float64
		cSec   map[string]string
		Result Report
	)

	// Create a map of category-value pairs
	cVal = catVal(catsIn)
	// Create a map of category-value pairs with opening balance values
	sVal = staVal(catsIn)

	// Create a pointer to this map, i.e. a pointer to category's value
	pVal = &cVal

	// Create a map of category-section pairs to detect category's section
	cSec = catSec(catsIn)

	// Calculate a starting balance per section
	for j0 := range catsIn {
		Result.startBal(sVal, cSec, catsIn[j0].Cat)
	}

	// Calculate a starting profit (loss)
	// Note: this result doesn't reflect the actual starting P/L value,
	// as the retained profit or the loss carried forward is in
	// 'Equity, P&L Account' with a user-defined key
	Result.Profit.Profit.Sta = Result.Profit.Revenue.Sta - Result.Profit.Expense.Sta

	// Allocate space for a slice of categories
	catsOut = make([]Categories, len(catsIn))

	for i := range records {
		// Update a Source-account value. Set value through the pointer to
		// category's value.
		// Note: Some sections ('Revenues', 'Liabilities' and 'Equity') require
		// special treatment. The applied sign depends on the type (section) of the source
		// category.
		if specialSection(cSec, records[i].Source) {
			*pVal = addBal(cVal, records[i].Source, +records[i].Amount)
		} else {
			*pVal = addBal(cVal, records[i].Source, -records[i].Amount)
		}

		// Update a Purpose-account value. Set value through the pointer to
		// category's value.
		// Note: some sections ('Revenues', 'Liabilities' and 'Equity') require
		// special treatment. The applied sign depends on the type (section) of the source
		// category.
		if specialSection(cSec, records[i].Purpose) {
			*pVal = addBal(cVal, records[i].Purpose, -records[i].Amount)
		} else {
			*pVal = addBal(cVal, records[i].Purpose, +records[i].Amount)
		}
	}

	// Assign values to each element of the slice of categories
	for j := range catsIn {
		catsOut[j].Cat = catsIn[j].Cat
		catsOut[j].Sect = catsIn[j].Sect
		catsOut[j].Name = catsIn[j].Name
		catsOut[j].Bal.Sta = catsIn[j].Bal.Sta

		catsOut[j].Bal.Dif = cVal[catsIn[j].Cat]

		catsOut[j].Bal.End = catsOut[j].Bal.Sta + catsOut[j].Bal.Dif

		// Round the balance change and the ending balance to the required number
		// of decimals
		catsOut[j].Bal.Dif = math.Round(catsOut[j].Bal.Dif*decimals) / decimals
		catsOut[j].Bal.End = math.Round(catsOut[j].Bal.End*decimals) / decimals

		// Calculate a difference between the balance at the end and the balance at
		// the beginning per section
		Result.changeBal(cSec, catsOut[j].Cat, catsOut[j].Bal.Dif)
	}

	Result.Profit.Profit.End = Result.Profit.Revenue.End - Result.Profit.Expense.End

	// Calculate an ending balance per section
	Result.finalBal()

	Result.roundBal()

	fmt.Println("Balance, Assets:         ", Result.Balance.Assets)
	fmt.Println("Balance, Liabilities:    ", Result.Balance.Liabls)
	fmt.Println("Balance, Equity:         ", Result.Balance.Equity)
	fmt.Println("Balance, Retained Result:", Result.Balance.Retained)
	fmt.Println("P&L, Revenues:", Result.Profit.Revenue)
	fmt.Println("P&L, Expenses:", Result.Profit.Expense)
	fmt.Println("P&L, Profit:  ", Result.Profit.Profit)

	return
}
