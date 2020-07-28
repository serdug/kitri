// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package conti provides business logic of trial account calculation
package conti

// catVal builds a category-value map
func catVal(cats []Categories) map[string]float64 {
	cval := map[string]float64{}
	for i := range cats {
		cval[cats[i].Cat] = cats[i].Bal.Dif
	}
	return cval
}

// staVal builds a category-value map of starting values
func staVal(cats []Categories) map[string]float64 {
	sval := map[string]float64{}
	for i := range cats {
		sval[cats[i].Cat] = cats[i].Bal.Sta // incorrect
	}
	return sval
}

// catSec builds a category-section map
func catSec(cats []Categories) map[string]string {
	csec := map[string]string{}
	for i := range cats {
		csec[cats[i].Cat] = cats[i].Sect
	}
	return csec
}

// addBal modifies the balance value for the category
func addBal(cval map[string]float64, cat string, val float64) map[string]float64 {
	if v, ok := cval[cat]; ok {
		// Found
		cval[cat] = v + val
		return cval
	}
	return map[string]float64{}
}
