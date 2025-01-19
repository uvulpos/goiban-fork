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
	"database/sql"
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
	co "github.com/uvulpos/goiban-fork/countries"
	countryValidationRules "github.com/uvulpos/goiban-fork/countries/validation-rules"
	"github.com/uvulpos/goiban-fork/data"
)

var (
	SELECT_BIC                   = "SELECT bic FROM BANK_DATA WHERE bankcode = ? AND country = ?;"
	SELECT_BIC_STMT              *sql.Stmt
	SELECT_BANK_INFORMATION      = "SELECT bankcode, name, zip, city, bic FROM BANK_DATA WHERE bankcode = ? AND country = ?;"
	SELECT_BANK_INFORMATION_STMT *sql.Stmt
)

func GetBic(iban *Iban, intermediateResult *ValidationResult, repo data.BankDataRepository) *ValidationResult {
	length, ok := countryValidationRules.COUNTRY_CODE_TO_BANK_CODE_LENGTH[(iban.countryCode)]

	if !ok {
		intermediateResult.Messages = append(intermediateResult.Messages, "Cannot get BIC. No information available.")
		return intermediateResult
	}

	if len(iban.bban) < length {
		intermediateResult.Messages = append(intermediateResult.Messages, "Cannot get BIC for BBAN "+iban.bban)
		return intermediateResult
	}

	bankCode := iban.bban[0:length]
	bankData := GetBankInformationByCountryAndBankCodeFromDb(iban.countryCode, bankCode, repo)

	if bankData == nil {
		intermediateResult.Messages = append(intermediateResult.Messages, "No BIC found for bank code: "+bankCode)
		return intermediateResult
	}

	// issue #17 - Custom Rule for Commerzbank
	//
	// See https://www.eckd-kigst.de/fileadmin/user_upload/eckd/Downloads_KFM/Deutsche_Bundesbank_Uebersicht_der_IBAN_Regeln_Stand_Juni_2013.pdf <-- broken link
	// See GitHub Issue: https://github.com/apilayer/goiban-service/issues/17
	if iban.countryCode == "DE" && isCommerzbank(bankData) {
		bankData.Bic = "COBADEFFXXX"
	}

	intermediateResult.BankData = *bankData

	return intermediateResult
}

func isCommerzbank(bd *data.BankInfo) bool {
	return len(bd.Bankcode) > 6 && bd.Bankcode[3:6] == "400"
}

func prepareSelectBankInformationStatement(db *sql.DB) {
	var err error

	SELECT_BANK_INFORMATION_STMT, err = db.Prepare(SELECT_BANK_INFORMATION)
	if err != nil {
		panic("Couldn't prepare statement: " + SELECT_BANK_INFORMATION)
	}

}

func GetBankInformationByCountryAndBankCodeFromDb(countryCode string, bankCode string, repo data.BankDataRepository) *data.BankInfo {

	// if SELECT_BANK_INFORMATION_STMT == nil {
	// 	prepareSelectBankInformationStatement(db)
	// }

	// var dbBankcode, dbName, dbZip, dbCity, dbBic string

	//bankCode = strings.TrimLeft(bankCode, "0")

	data, err := repo.Find(countryCode, bankCode)
	// err := SELECT_BANK_INFORMATION_STMT.QueryRow(bankCode, countryCode).Scan(&dbBankcode, &dbName, &dbZip, &dbCity, &dbBic)

	if err != nil {
		panic("Failed to load bank info from db.")
	}

	return data
}

func prepareSelectBicStatement(db *sql.DB) {
	var err error
	SELECT_BIC_STMT, err = db.Prepare(SELECT_BIC)
	if err != nil {
		panic("Couldn't prepare statement: " + SELECT_BIC)
	}
}

func ReadAustriaBankFileEntry(fileContent string) (result []*co.AustriaBankFileEntry) {
	for _, line := range strings.Split(fileContent, "\n")[5:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		result = append(result, co.AustriaBankStringToEntry(line, countryValidationRules.COUNTRY_CODE_TO_BANK_CODE_LENGTH))
	}
	return result
}

func ReadGermanBankFileEntry(fileContent string) (result []*co.GermanBankFileEntry) {
	for _, line := range strings.Split(fileContent, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		result = append(result, co.BundesbankStringToEntry(line))
	}
	return result
}

func xlsxFileToSlice(fileContent string) ([][][]string, error) {
	xlsxFile, xlsxFileErr := xlsx.OpenBinary([]byte(fileContent))
	if xlsxFileErr != nil {
		return [][][]string{}, fmt.Errorf("Couldn't read xlsx file content, %v", xlsxFileErr)
	}

	file, fileErr := xlsxFile.ToSlice()
	if fileErr != nil {
		return [][][]string{}, fmt.Errorf("Couldn't read xlsx file content, %v", fileErr)
	}

	return file, nil
}
