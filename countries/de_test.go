package countries

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestCanConvertStringToBundesbankEntry(t *testing.T) {
	data := "100000001Bundesbank                                                10591Berlin                             BBk Berlin                 20100MARKDEF110009011380U000000000"
	result := BundesbankStringToEntry(data)

	assert.Equal(t, result.Bankcode, "10000000", "Couldn't parse bank code.")
	assert.Equal(t, result.M, 1, "Couldn't parse M.")
	assert.Equal(t, result.Name, "Bundesbank", "Couldn't parse name.")
	assert.Equal(t, result.Zip, "10591", "Couldn't parse zip.")
	assert.Equal(t, result.City, "Berlin", "Couldn't parse city.")
	assert.Equal(t, result.ShortName, "BBk Berlin", "Couldn't parse short name.")
	assert.Equal(t, result.Pan, 20100, "Couldn't parse short pan.")
	assert.Equal(t, result.Bic, "MARKDEF1100", "Couldn't parse bic.")
	assert.Equal(t, result.CheckAlgo, "09", "Couldn't parse check algo.")
	assert.Equal(t, result.Id, "01138", "Couldn't parse internal id.")
	assert.Equal(t, result.Change, "U", fmt.Sprintf("Couldn't parse change: %v", result.Change))
	assert.Equal(t, result.ToBeDeleted, 0, "Couldn't parse to be deleted.")
	assert.Equal(t, result.NewBankCode, "00000000", "Couldn't parse new bank code.")
}
