package countries

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestCanConvertSliceToSwitzerlandBankEntry(t *testing.T) {
	countryCodeToBankCodeMap := map[string]int{
		"CH": 5,
	}

	data := []string{"4", "4750", "0000", "4835", "047501", "4835", "1", "20061020", "1", "1", "2", "CS (Schweiz) AG", "Credit Suisse (Schweiz) AG", "Rue du  Simplon 50", "Case postale 210", "1800", "Vevey 1", "021 925 01 11", "021 921 90 87", "", "", "*12-35-2", "CRESCHZZ18A"}
	entry := SwitzerlandRowToEntry(data, countryCodeToBankCodeMap)

	assert.Equal(t, entry.Bankcode, "04835", fmt.Sprintf("expected 04835 as bankcode, got %v", entry.Bankcode))
	assert.Equal(t, entry.Bic, "CRESCHZZ18A", fmt.Sprintf("expected CRESCHZZ18A as bic, got %v", entry.Bic))
	assert.Equal(t, entry.Name, "Credit Suisse (Schweiz) AG", fmt.Sprintf("Credit Suisse (Schweiz) AG as name, got %v", entry.Name))
}
