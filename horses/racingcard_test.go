package horses

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

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
			if !IsRacingCardFile(f.Name()) {
				continue
			}
			path := path.Join(dir, f.Name())
			t.Log(fmt.Sprintf("checking: %s", path))
			blob, err := ioutil.ReadFile(path)
			require.NoError(t, err, path)
			_, err = ParseRacingCard(blob)
			require.NoError(t, err, path)
		}
	}
}
