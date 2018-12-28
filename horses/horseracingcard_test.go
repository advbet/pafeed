package horses

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulkParseHorseRacingCard(t *testing.T) {
	dirs := []string{
		"testdata/NewcastleRule4AllBets",
		"testdata/WindsorRule4BoardPrices",
		"testdata/Aintree",
		"testdata/Lingfield",
		"testdata/VaalLateWithdrawal",
		"testdata/Abandoned",
		"testdata/Abandoned/Hexham",
		"testdata/GreyvilleJockeyChanges",
		"testdata/EdgeCases",
		"testdata/feed",
	}

	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err, dir)
		for _, f := range files {
			// try parsing only HorceRacingCard documents
			if !strings.HasPrefix(f.Name(), "c") {
				continue
			}
			path := path.Join(dir, f.Name())
			t.Log(fmt.Sprintf("checking: %s", path))
			blob, err := ioutil.ReadFile(path)
			require.NoError(t, err, path)
			var obj RacingCard
			assert.NoError(t, xml.Unmarshal(blob, &obj), path)
		}
	}
}
