package countries

import (
	"fmt"
	"testing"

	countryValidationRules "github.com/fourcube/goiban/countries/validation-rules"
	"gotest.tools/assert"
)

func TestCanConvertSliceToLiechtensteinBankEntry(t *testing.T) {
	data := []string{"Bank Alpinum AG", "BALPLI22", "8801"}
	entry := LiechtensteinRowToEntry(data, countryValidationRules.COUNTRY_CODE_TO_BANK_CODE_LENGTH)

	assert.Equal(t, entry.Bankcode, "08801", fmt.Sprintf("expected 08801 as bankcode, got %v", entry.Bankcode))
	assert.Equal(t, entry.Bic, "BALPLI22", fmt.Sprintf("expected BALPLI22 as bic, got %v", entry.Bankcode))
	assert.Equal(t, entry.Name, "Bank Alpinum AG", fmt.Sprintf("expected Bank Alpinum AG as name, got %v", entry.Bankcode))
}
