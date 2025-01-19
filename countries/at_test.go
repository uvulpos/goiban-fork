package countries

import (
	"fmt"
	"testing"

	countryValidationRules "github.com/fourcube/goiban/countries/validation-rules"
	"gotest.tools/assert"
)

func TestCanConvertStringToAustriaBankEntry(t *testing.T) {
	data := "Hauptanstalt;\"10050973\";\"52300\";\"KI\";\"Joint stock banks and private banks\";\"350921k\";\"Addiko Bank AG\";\"Wipplingerstraße 34/4\";\"1010\";\"Wien\";\"\";\"\";\"\";\"\";\"\";\"Wien\";\"050232\";\"050232/3000\";\"holding@addiko.com\";\"HSEEAT2KXXX\";\"www.addiko.com\";\"20130621\";;"
	result := AustriaBankStringToEntry(data, countryValidationRules.COUNTRY_CODE_TO_BANK_CODE_LENGTH)

	assert.Equal(t, result.Name, "Addiko Bank AG", fmt.Sprintf("Couldn't parse name: %v", result.Name))
	assert.Equal(t, result.Bic, "HSEEAT2KXXX", fmt.Sprintf("Couldn't parse bic: %v", result.Bic))
	assert.Equal(t, result.Bankcode, "52300", fmt.Sprintf("Couldn't parse bank code: %v", result.Bankcode))
}
