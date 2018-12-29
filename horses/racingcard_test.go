package horses

import (
	"fmt"
	"io/ioutil"
	"path"
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

	// Assertsions is a list of custom checks for successfuly parsing some
	// fields to their non default value.
	assertions := map[string]bool{}
	// catch is a helper function for marking custom test check as passed.
	catch := func(name string, val bool) {
		assertions[name] = assertions[name] || val
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
			cards, err := ParseRacingCard(blob)
			require.NoError(t, err, path)

			assert.True(t, len(*cards) == 1, "always exactly one meeting card per file")
			for _, m := range *cards {
				assert.True(t, len(m.Races) >= 1, "always at least one race per meeting")
				for _, race := range m.Races {
					for _, h := range race.Horses {
						catch("parsed CardHorse Jockey ID", h.Jockey.ID != 0)
						catch("parsed CardHorse Jockey Name", h.Jockey.Name != "")
						catch("parsed CardHorse Jockey Allowance", h.Jockey.Allowance.Units != "")
						catch("parsed CardHorse Trainer ID", h.Trainer.ID != 0)
						catch("parsed CardHorse Trainer Name", h.Trainer.Name != "")
						catch("parsed CardHorse Trainer Nationality", h.Trainer.Nationality != "")
						catch("parsed CardHorse Trainer Location", h.Trainer.Location != "")
						catch("parsed CardHorse Breeding", len(h.Breeding) > 0)
						for _, b := range h.Breeding {
							catch("parsed Breeding Relation", b.Relation != "")
							catch("parsed Breeding Name", b.Name != "")
							catch("parsed Breeding Bred", b.Bred != "")
							catch("parsed Breeding YearBorn", b.YearBorn != 0)
						}
						catch("parsed CardHorse Ratings", len(h.Ratings) > 0)
						for _, r := range h.Ratings {
							catch("parsed Ratings Type", r.Type != "")
							catch("parsed Ratings Value", r.Value != 0)
						}
					}
				}
			}
		}
	}

	for check, value := range assertions {
		assert.True(t, value, check)
	}
}
