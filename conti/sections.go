// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

import (
	"math"
)

type Report struct {
	Balance BalanceSections
	Profit  ProfitSections
}

type BalanceSections struct {
	Assets Tally
	Liabls Tally
	Equity Tally

	// Retained result
	// Accumulated results, i.e. retained earnings (accumulated deficit)
	Retained Tally
}

type ProfitSections struct {
	Revenue Tally
	Expense Tally
	Profit  Tally
}

// specialSection detects if category's section requires special treatment
func specialSection(csec map[string]string, cat string) bool {
	if s, ok := csec[cat]; ok {
		// Found
		// Special sections:
		if s == "Revenues" || s == "Liabilities" || s == "Equity" {
			return true
		} else {
			return false
		}
	}
	return false
}

// startBal calculates the total balance value per section at the beginning
func (this *Report) startBal(sval map[string]float64, csec map[string]string, cat string) {
	if value, ok := sval[cat]; ok {
		if section, ok := csec[cat]; ok {
			switch {
			case section == "Assets":
				this.Balance.Assets.Sta += value

			case section == "Liabilities":
				this.Balance.Liabls.Sta += value

			case section == "Equity":
				this.Balance.Equity.Sta += value

			case section == "Assets" || section == "Liabilities" || section == "Equity":
				this.Balance.Retained.Sta += value

			case section == "Revenues":
				this.Profit.Revenue.Sta += value

			case section == "Expenses":
				this.Profit.Expense.Sta += value

			default:
				// Do nothing
			}
		}
	}
}

// changeBal calculates the difference between the ending and the starting values per section
func (this *Report) changeBal(csec map[string]string, cat string, value float64) {
	if section, ok := csec[cat]; ok {
		switch {
		case section == "Assets":
			this.Balance.Assets.Dif += value

		case section == "Liabilities":
			this.Balance.Liabls.Dif += value

		case section == "Equity":
			this.Balance.Equity.Dif += value

		case section == "Revenues":
			this.Profit.Revenue.Dif += value
			this.Profit.Profit.Dif += value

		case section == "Expenses":
			this.Profit.Expense.Dif += value
			this.Profit.Profit.Dif -= value

		default:
			// Do nothing
		}
	}
}

// finalBal calculates the total balance value per section at the end
func (this *Report) finalBal() {
	this.Balance.Assets.End = this.Balance.Assets.Sta + this.Balance.Assets.Dif
	this.Balance.Liabls.End = this.Balance.Liabls.Sta + this.Balance.Liabls.Dif
	this.Balance.Equity.End = this.Balance.Equity.Sta + this.Balance.Equity.Dif

	this.Profit.Revenue.End = this.Profit.Revenue.Sta + this.Profit.Revenue.Dif
	this.Profit.Expense.End = this.Profit.Expense.Sta + this.Profit.Expense.Dif

	this.Profit.Profit.End = this.Profit.Profit.Sta + this.Profit.Profit.Dif

	this.Balance.Retained.Dif = this.Profit.Profit.Dif
	this.Balance.Retained.End = this.Balance.Retained.Sta + this.Balance.Retained.Dif
}

// roundBal rounds the balance change and the ending balance to the required
// number of decimals
func (this *Report) roundBal() {
	this.Balance.Assets.Sta = math.Round(this.Balance.Assets.Sta*decimals) / decimals
	this.Balance.Liabls.Sta = math.Round(this.Balance.Liabls.Sta*decimals) / decimals
	this.Balance.Equity.Sta = math.Round(this.Balance.Equity.Sta*decimals) / decimals

	this.Profit.Revenue.Sta = math.Round(this.Profit.Revenue.Sta*decimals) / decimals
	this.Profit.Expense.Sta = math.Round(this.Profit.Expense.Sta*decimals) / decimals
	this.Profit.Profit.Sta = math.Round(this.Profit.Profit.Sta*decimals) / decimals

	this.Balance.Assets.End = math.Round(this.Balance.Assets.End*decimals) / decimals
	this.Balance.Liabls.End = math.Round(this.Balance.Liabls.End*decimals) / decimals
	this.Balance.Equity.End = math.Round(this.Balance.Equity.End*decimals) / decimals

	this.Profit.Revenue.End = math.Round(this.Profit.Revenue.End*decimals) / decimals
	this.Profit.Expense.End = math.Round(this.Profit.Expense.End*decimals) / decimals
	this.Profit.Profit.End = math.Round(this.Profit.Profit.End*decimals) / decimals

	this.Balance.Assets.Dif = math.Round(this.Balance.Assets.Dif*decimals) / decimals
	this.Balance.Liabls.Dif = math.Round(this.Balance.Liabls.Dif*decimals) / decimals
	this.Balance.Equity.Dif = math.Round(this.Balance.Equity.Dif*decimals) / decimals

	this.Profit.Revenue.Dif = math.Round(this.Profit.Revenue.Dif*decimals) / decimals
	this.Profit.Expense.Dif = math.Round(this.Profit.Expense.Dif*decimals) / decimals
	this.Profit.Profit.Dif = math.Round(this.Profit.Profit.Dif*decimals) / decimals

	this.Balance.Retained.Dif = math.Round(this.Balance.Retained.Dif*decimals) / decimals
	this.Balance.Retained.End = math.Round(this.Balance.Retained.End*decimals) / decimals
}
