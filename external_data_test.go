/*
The MIT License (MIT)

Copyright (c) 2014 Chris Grieger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package goiban

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	data "github.com/uvulpos/goiban-fork/data"
	"gotest.tools/assert"
)

var (
	repo = data.NewInMemoryStore()
)

//go:embed test/austria.csv
var testDataAustria string

//go:embed test/bundesbank.txt
var testDataGermanBundesbank string

func TestMain(m *testing.M) {
	retCode := m.Run()

	os.Exit(retCode)
}

func TestCanReadFromAustriaFile(t *testing.T) {
	data := ReadAustriaBankFileEntry(testDataAustria)
	assert.Check(t, len(data) > 0, "Failed to read file. No rows in slice")
}

func TestCanReadFromBundesbankFile(t *testing.T) {
	data := ReadGermanBankFileEntry(testDataGermanBundesbank)
	assert.Check(t, data[0].Name != "", "Failed to read file.")
}

func TestSpecialRuleForCommerzbankBic(t *testing.T) {
	dataRepo := data.NewInMemoryStore()
	input := "DE06200400000052065002"
	iban := ParseToIban(input)
	result := NewValidationResult(true, "", input)

	data := data.BankInfo{Bankcode: "20040000", Country: "DE", Source: "Foo"}
	dataRepo.Store(data)

	result = GetBic(iban, result, dataRepo)
	assert.Equal(t, result.BankData.Bankcode, "20040000", "BLZ is wrong")
	assert.Equal(t, result.BankData.Bic, "COBADEFFXXX", fmt.Sprintf("Expected Bic COBADEFFXXX, was %v", result.BankData.Bic))
}
