package horses

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/big"
	"path"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/advbet/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDecimal(t *testing.T, s string) decimal.Number {
	d, err := decimal.FromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func makeTime(t *testing.T, s string) time.Time {
	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

func makeRat(t *testing.T, s string) big.Rat {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		t.Fatal("error converting", s, "to big.Rat")
	}
	return *r
}

func TestBulkParseHorseRacing(t *testing.T) {
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
			// try parsing only HorceRacing documents
			if !IsRacingFile(f.Name()) {
				continue
			}
			path := path.Join(dir, f.Name())
			t.Log(fmt.Sprintf("checking: %s", path))
			blob, err := ioutil.ReadFile(path)
			require.NoError(t, err, path)
			obj, err := ParseRacingFile(blob)
			require.NoError(t, err, path)

			assert.True(t, len(obj.Meetings) == 1, "always exactly one meeting perfile")
			for _, m := range obj.Meetings {
				assert.True(t, len(m.Races) <= 1, "always at least one race per meeting")
				for _, r := range m.Races {
					catch("parsed Race with status Dormant", r.Status == RaceDormant)
					catch("parsed Race with status Delayed", r.Status == RaceDelayed)
					catch("parsed Race with status Parading", r.Status == RaceParading)
					catch("parsed Race with status GoingDown", r.Status == RaceGoingDown)
					catch("parsed Race with status AtThePost", r.Status == RaceAtThePost)
					catch("parsed Race with status GoingBehind", r.Status == RaceGoingBehind)
					//catch("parsed Race with status GoingInStalls", r.Status == RaceStatusGoingInStalls)
					//catch("parsed Race with status UnderOrders", r.Status == RaceStatusUnderOrders)
					catch("parsed Race with status Off", r.Status == RaceOff)
					catch("parsed Race with status Finished", r.Status == RaceFinished)
					catch("parsed Race with status FalseStart", r.Status == RaceFalseStart)
					//catch("parsed Race with status Photograph", r.Status == RaceStatusPhotograph)
					catch("parsed Race with status Result", r.Status == RaceResult)
					catch("parsed Race with status WeighedIn", r.Status == RaceWeighedIn)
					catch("parsed Race with status Void", r.Status == RaceRaceVoid)
					catch("parsed Race with status Abandoned", r.Status == RaceAbandoned)

					catch("parsed Race with stewards status None", r.Stewards == StewardsNone)
					catch("parsed Race with stewards status Inquiry", r.Stewards == StewardsInquiry)
					catch("parsed Race with stewards status Objection", r.Stewards == StewardsObjection)
					catch("parsed Race with stewards status InquiryAndObjection", r.Stewards == StewardsInquiryAndObjection)
					catch("parsed Race with stewards status AmendedResult", r.Stewards == StewardsAmendedResult)
					catch("parsed Race with stewards status ResultStands", r.Stewards == StewardsResultStands)
					if r.Stewards == StewardsInquiry {
						assert.True(t, r.StewardsInquiry != "", "if stewards status is Inquiry then StewardsInquiry must be non empty")
					}
					if r.Stewards == StewardsObjection {
						assert.True(t, r.StewardsObjection != "", "if stewards status is Objection then StewardsObjection must be non empty")
					}
					if r.Stewards == StewardsInquiryAndObjection {
						assert.True(t, r.StewardsObjection != "", "if stewards status is InquiryAndObjection then StewardsObjection must be non empty")
						assert.True(t, r.StewardsObjection != "", "if stewards status is InquiryAndObjection then StewardsObjection must be non empty")
					}

					catch("parsed Race Stewards Inquiry", r.StewardsInquiry != "")
					catch("parsed Race Stewards Objection", r.StewardsObjection != "")
					for _, h := range r.Horses {
						catch("parsed Horse with status Runner", h.Status == HorseRunner)
						catch("parsed Horse with status NonRunner", h.Status == HorseNonRunner)
						catch("parsed Horse with status Withdrawn", h.Status == HorseWithdrawn)
						//catch("parsed Horse with status Reserve", h.Status == HorseStatusReserve)
						if h.Result != nil {
							catch("parsedHorseResultDisqualified", h.Result.Disqualified)
							catch("parsedHorseResultAmendedPos", h.Result.AmendedPos != 0)
						}
						/*
							if h.Status == HorseStatusWithdrawn {
								assert.True(t, !h.WithdrawnTime.IsZero(), "Withdrawn time is always present if horse status is withdrawn")
								assert.True(t, h.WithdrawnBetMarket != 0, "Withdrawn bet market is always non zero if horse status is withdrawn")
							}*/
						catch("parsed Horse WithdrawnTime", !h.WithdrawnTime.IsZero())
						catch("parsed Horse WithdrawnMarketNumber", h.WithdrawnBetMarket != 0)
					}
					for _, m := range r.BetMarkets {
						catch("parsedBetMarketSuspended", !m.Suspended.IsZero())
					}
				}
			}
		}
	}

	for check, value := range assertions {
		assert.True(t, value, check)
	}
}

func TestParseHorseRacing(t *testing.T) {
	tests := []struct {
		file string
		obj  RacingFile
		err  error
	}{
		{
			file: "testdata/Lingfield/b20180414lin17400007.xml",
			obj: RacingFile{
				Timestamp: makeTime(t, "2018-04-14T17:34:37+01:00"),
				Meetings: []Meeting{{
					ID:       97192,
					Revision: 4,
					Country:  "England",
					Course:   "Lingfield",
					Date:     makeTime(t, "2018-04-14T00:00:00Z"),
					Status:   MeetingDormant,
					//Abandoned
					//Delayed
					Weather:    "Sunny",
					GoingBrief: "Standard",
					GoingFull:  "Standard",
					Races: []Race{{
						ID:         798361,
						Revision:   7,
						StartTime:  makeTime(t, "2018-04-14T17:40:00+01:00"),
						Runners:    10,
						Handicap:   true,
						Showcase:   false,
						Trifecta:   true,
						Stewards:   StewardsNone,
						Status:     RaceDormant,
						Weather:    "Sunny",
						GoingBrief: "Standard",
						GoingFull:  "Standard",
						//OffTime
						//WinTime
						//StewardsInquiry
						//StewardsObjection
						BetMarkets: []BetMarket{
							{
								MarketNumber:  1,
								Formed:        makeTime(t, "2018-04-14T17:33:20+01:00"),
								DeductionType: DeductionNone,
							},
						},
						//LackFinishers   TODO // Used if not enough finishers to fill result
						//Message         TODO // Any other information about the race
						Horses: []Horse{
							{
								ID:          1961454,
								Name:        "Officer Drivel",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 1,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 135,
									Text:  "9st 9lbs",
								},
								Jockey: Jockey{
									ID:   1150396,
									Name: "Harry Burns",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 7,
									// },
								},
								Trainer: Trainer{
									ID:   131079,
									Name: "Suzi Best",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "14/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-04-14T17:34:35+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "12/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2102945,
								Name:        "Zephyros",
								Bred:        "GER",
								Status:      HorseRunner,
								ClothNumber: 2,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 135,
									Text:  "9st 9lbs",
								},
								Jockey: Jockey{
									ID:   1156790,
									Name: "Poppy Bridgwater",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 7,
									// },
								},
								Trainer: Trainer{
									ID:   12102,
									Name: "D G Bridgwater",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "7/2"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2166053,
								Name:        "Hatem",
								Bred:        "FR",
								Status:      HorseRunner,
								ClothNumber: 3,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 133,
									Text:  "9st 7lbs",
								},
								Jockey: Jockey{
									ID:   14394,
									Name: "Fran Berry",
								},
								Trainer: Trainer{
									ID:   9237,
									Name: "N P Littmoden",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "25/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          1713481,
								Name:        "Ready",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 4,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 131,
									Text:  "9st 5lbs",
								},
								Jockey: Jockey{
									ID:   76191,
									Name: "Kieren Fox",
								},
								Trainer: Trainer{
									ID:   129102,
									Name: "Mark Pattinson",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "4/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-04-14T17:34:35+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "9/2"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2214710,
								Name:        "Oceanus",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 5,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 130,
									Text:  "9st 4lbs",
								},
								Jockey: Jockey{
									ID:   1140493,
									Name: "Shelley Birkett",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 3,
									// },
								},
								Trainer: Trainer{
									ID:   14707,
									Name: "Miss J Feilden",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "10/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2234307,
								Name:        "Presence Process",
								Bred:        "GB",
								Status:      HorseRunner,
								ClothNumber: 6,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 129,
									Text:  "9st 3lbs",
								},
								Jockey: Jockey{
									ID:   1149121,
									Name: "Paddy Bradley",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 5,
									// },
								},
								Trainer: Trainer{
									ID:   3897,
									Name: "P Phelan",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "3/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          1556290,
								Name:        "Karam Albaari",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 7,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 127,
									Text:  "9st 1lbs",
								},
								Jockey: Jockey{
									ID:   1157192,
									Name: "Tom Marquand",
								},
								Trainer: Trainer{
									ID:   103,
									Name: "J R Jenkins",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "7/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2033610,
								Name:        "Maraakib",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 8,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 121,
									Text:  "8st 9lbs",
								},
								Jockey: Jockey{
									ID:   31657,
									Name: "K T O'Neill",
								},
								Trainer: Trainer{
									ID:   111091,
									Name: "A Dunn",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "16/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2234195,
								Name:        "Amadeus Rox",
								Bred:        "FR",
								Status:      HorseRunner,
								ClothNumber: 9,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 120,
									Text:  "8st 8lbs",
								},
								Jockey: Jockey{
									ID:   74413,
									Name: "J P Fahy",
								},
								Trainer: Trainer{
									ID:   111091,
									Name: "A Dunn",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "22/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-04-14T17:34:35+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "25/1"),
									},
								},
								//StartingPrice
							},
							{
								ID:          2196756,
								Name:        "Feel The Vibes",
								Bred:        "GB",
								Status:      HorseRunner,
								ClothNumber: 10,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 119,
									Text:  "8st 7lbs",
								},
								Jockey: Jockey{
									ID:   65282,
									Name: "David Probert",
								},
								Trainer: Trainer{
									ID:   8436,
									Name: "M Blanshard",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-04-14T17:33:20+01:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "11/1"),
									},
								},
								//StartingPrice
							},
						},
						//WinningDistance TODO // The distances between the runners on completing the course
						//Returns         TODO // The returns generated by the result of the race
						//SellingDetails  TODO // Details of horses sold or claimed after the result
					}},
					//Messages
					//MultiBet
				}},
			},
		},
		{
			file: "testdata/feed/b20181128wth12150045.xml",
			obj: RacingFile{
				Timestamp: makeTime(t, "2018-11-28T12:50:14+00:00"),
				Meetings: []Meeting{{
					ID:       104930,
					Revision: 3,
					Country:  "England",
					Course:   "Wetherby",
					Date:     makeTime(t, "2018-11-28T00:00:00Z"),
					Status:   MeetingDormant,
					//Abandoned
					//Delayed
					Weather:    "Overcast & Showers",
					GoingBrief: "Good to Soft",
					GoingFull:  "Good to Soft",
					Races: []Race{{
						ID:         854412,
						Revision:   45,
						StartTime:  makeTime(t, "2018-11-28T12:15:00+00:00"),
						Runners:    10,
						Handicap:   false,
						Showcase:   false,
						Trifecta:   true,
						Stewards:   StewardsNone,
						Status:     RaceWeighedIn,
						Weather:    "Overcast & Showers",
						GoingBrief: "Good to Soft",
						GoingFull:  "Good to Soft",
						OffTime:    makeTime(t, "2018-11-28T12:15:49+00:00"),
						WinTime:    4*time.Minute + 3*time.Second + 100*time.Millisecond,
						//StewardsInquiry
						//StewardsObjection
						BetMarkets: []BetMarket{
							{
								MarketNumber:  1,
								Formed:        makeTime(t, "2018-11-28T12:06:15+00:00"),
								DeductionType: DeductionNone,
							},
						},
						//LackFinishers   TODO // Used if not enough finishers to fill result
						//Message         TODO // Any other information about the race
						Horses: []Horse{
							{
								ID:          2358957,
								Name:        "Alexanderthegreat",
								Bred:        "FR",
								Status:      HorseRunner,
								ClothNumber: 1,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 152,
									Text:  "10st 12lbs",
								},
								Jockey: Jockey{
									ID:   41547,
									Name: "B S Hughes",
								},
								Trainer: Trainer{
									ID:   9243,
									Name: "J J Quinn",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "5/4"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:11:40+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "11/8"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:13:08+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "6/4"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:14:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "13/8"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:15:14+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "7/4"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "13/8"),
									FavouritePosition: 1,
									FavouriteJoint:    1,
								},
								CasualtyReason:      UnseatedRider,
								CloseUpComment:      "tracked leaders, tracked winner soon after 6th, ridden 3 out, weakened next, stumbled on landing and unseated rider last",
								BetMovementsComment: "op 5/4 tchd 7/4",
							},
							{
								ID:          2342636,
								Name:        "Alliteration",
								Bred:        "GB",
								Status:      HorseRunner,
								ClothNumber: 2,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 152,
									Text:  "10st 12lbs",
								},
								Jockey: Jockey{
									ID:   83305,
									Name: "Danny Cook",
								},
								Trainer: Trainer{
									ID:   107851,
									Name: "J Hughes",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "11/2"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:11:57+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "5/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:14:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "9/2"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:15:05+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "4/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "4/1"),
									FavouritePosition: 3,
									FavouriteJoint:    1,
								},
								Result: &Result{
									FinishPos:    2,
									Disqualified: false,
									//AmendedPos
									BetweenDistance: "17 lengths",
								},
								CasualtyReason:      NoCasualty,
								CloseUpComment:      "held up, headway 6th, chased winner when hit 3 out, plugged on",
								BetMovementsComment: "op 11/2",
							},
							{
								ID:          2298225,
								Name:        "Burnieboozle",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 3,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 152,
									Text:  "10st 12lbs",
								},
								Jockey: Jockey{
									ID:   1148952,
									Name: "C R King",
								},
								Trainer: Trainer{
									ID:   9243,
									Name: "J J Quinn",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "33/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:14:03+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "25/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:14:55+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "20/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:15:14+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "16/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "16/1"),
									FavouritePosition: 6,
									FavouriteJoint:    1,
								},
								CasualtyReason:      Fell,
								CloseUpComment:      "keen, held up, over jumped and fell 4th",
								BetMovementsComment: "op 33/1",
							},
							{
								ID:          2295288,
								Name:        "Keynote",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 4,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 152,
									Text:  "10st 12lbs",
								},
								Jockey: Jockey{
									ID:   1164129,
									Name: "Mr P Armson",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 7,
									// },
								},
								Trainer: Trainer{
									ID:   9194,
									Name: "R J Armson",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "200/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "200/1"),
									FavouritePosition: 10,
									FavouriteJoint:    1,
								},
								Result: &Result{
									FinishPos:    4,
									Disqualified: false,
									//AmendedPos
									BetweenDistance: "30 lengths",
								},
								CasualtyReason: NoCasualty,
								CloseUpComment: "towards rear, ridden 6th, never on terms",
								//BetMovementsComment
							},
							{
								ID:          2310027,
								Name:        "Astrofire",
								Bred:        "GB",
								Status:      HorseRunner,
								ClothNumber: 5,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 145,
									Text:  "10st 5lbs",
								},
								Jockey: Jockey{
									ID:   1154755,
									Name: "Mr Alex Chadwick",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 7,
									// },
									Overweight: UnitsValue{
										Units: "lbs",
										Value: 2,
									},
								},
								Trainer: Trainer{
									ID:   163,
									Name: "M H Tompkins",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "200/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:10:12+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "150/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "150/1"),
									FavouritePosition: 9,
									FavouriteJoint:    1,
								},
								Result: &Result{
									FinishPos:    6,
									Disqualified: false,
									//AmendedPos
									BetweenDistance: "13 lengths",
								},
								CasualtyReason:      NoCasualty,
								CloseUpComment:      "keen headway to lead 2nd, soon clear, reduced lead and headed soon after 6th, weakened quickly",
								BetMovementsComment: "op 200/1",
							},
							{
								ID:          2402973,
								Name:        "Don't Fence Me In",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 6,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 145,
									Text:  "10st 5lbs",
								},
								Jockey: Jockey{
									ID:   75251,
									Name: "R P McLernon",
								},
								Trainer: Trainer{
									ID:   9710,
									Name: "P R Webber",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "14/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:11:40+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "16/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:13:08+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "14/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "14/1"),
									FavouritePosition: 5,
									FavouriteJoint:    1,
								},
								CasualtyReason:      PulledUp,
								CloseUpComment:      "green in rear and not fluent, blundered and nearly unseated rider and lost irons 5th, pulled up next",
								BetMovementsComment: "tchd 16/1",
							},
							{
								ID:          2338916,
								Name:        "Fabianski",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 7,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 145,
									Text:  "10st 5lbs",
								},
								Jockey: Jockey{
									ID:   55809,
									Name: "C O'Farrell",
								},
								Trainer: Trainer{
									ID:   118739,
									Name: "Rebecca Menzies",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "33/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:08:50+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "28/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:11:40+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "25/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:13:29+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "20/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:14:03+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "18/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:15:18+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "20/1"),
									},
								},
								Result: &Result{
									FinishPos:    1,
									Disqualified: false,
									//AmendedPos
									//BetweenDistance
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "20/1"),
									FavouritePosition: 7,
									FavouriteJoint:    1,
								},
								CasualtyReason:      NoCasualty,
								CloseUpComment:      "led and bumped 1st, headed 2nd, led again soon after 6th, clear 2 out, ridden and ran on",
								BetMovementsComment: "op 33/1 tchd 18/1",
							},
							{
								ID:          2279152,
								Name:        "Kheleyf's Girl",
								Bred:        "GB",
								Status:      HorseRunner,
								ClothNumber: 8,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 145,
									Text:  "10st 5lbs",
								},
								Jockey: Jockey{
									ID:   1150129,
									Name: "Harrison Beswick",
									// Allowance: UnitsValue{
									// 	Units: "lbs",
									// 	Value: 7,
									// },
								},
								Trainer: Trainer{
									ID:   125032,
									Name: "Clare Ellam",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "150/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:07:19+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "100/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "100/1"),
									FavouritePosition: 8,
									FavouriteJoint:    1,
								},
								CasualtyReason:      PulledUp,
								CloseUpComment:      "keen, tracked winner when bumped 1st, weakened 6th, tailed off when pulled up next",
								BetMovementsComment: "op 150/1",
							},
							{
								ID:          2298615,
								Name:        "Pepper Street",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 9,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 145,
									Text:  "10st 5lbs",
								},
								Jockey: Jockey{
									ID:   80831,
									Name: "Jack Quinlan",
								},
								Trainer: Trainer{
									ID:   128567,
									Name: "Miss Amy Murphy",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "9/4"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:12:30+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "2/1"),
									},
								},
								Result: &Result{
									FinishPos:    3,
									Disqualified: false,
									//AmendedPos
									BetweenDistance: "1 1/4 length",
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "2/1"),
									FavouritePosition: 2,
									FavouriteJoint:    1,
								},
								CasualtyReason:      NoCasualty,
								CloseUpComment:      "keen close up, tracked leaders when ridden 3 out, weakened next",
								BetMovementsComment: "op 9/4",
							},
							{
								ID:          2370895,
								Name:        "Sweet Marmalade",
								Bred:        "IRE",
								Status:      HorseRunner,
								ClothNumber: 10,
								Weight: UnitsValueText{
									Units: "lbs",
									Value: 145,
									Text:  "10st 5lbs",
								},
								Jockey: Jockey{
									ID:   1142033,
									Name: "Jamie Hamilton",
								},
								Trainer: Trainer{
									ID:   61138,
									Name: "L A Mullaney",
								},
								Shows: []Show{
									{
										Timestamp:    makeTime(t, "2018-11-28T12:06:15+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "12/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:08:50+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "11/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:10:12+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "12/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:10:56+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "11/1"),
									},
									{
										Timestamp:    makeTime(t, "2018-11-28T12:11:40+00:00"),
										MarketNumber: 1,
										Price:        makeRat(t, "12/1"),
									},
								},
								StartingPrice: StartingPrice{
									Price:             makeRat(t, "12/1"),
									FavouritePosition: 4,
									FavouriteJoint:    1,
								},
								Result: &Result{
									FinishPos:    5,
									Disqualified: false,
									//AmendedPos
									BetweenDistance: "33 lengths",
								},
								CasualtyReason:      NoCasualty,
								CloseUpComment:      "tracked leaders, ridden and lost place 6th",
								BetMovementsComment: "tchd 11/1",
							},
						},
						//WinningDistance TODO // The distances between the runners on completing the course
						Returns: &Returns{
							Tote: []Tote{
								{
									Type:     ToteWin,
									Currency: "GBP",
									Dividend: makeDecimal(t, "19.20"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
									},
								},
								{
									Type:     TotePlace,
									Currency: "GBP",
									Dividend: makeDecimal(t, "3.70"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
									},
								},
								{
									Type:     TotePlace,
									Currency: "GBP",
									Dividend: makeDecimal(t, "1.60"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2342636,
											Name: "Alliteration",
											Bred: "GB",
										},
									},
								},
								{
									Type:     TotePlace,
									Currency: "GBP",
									Dividend: makeDecimal(t, "1.20"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2298615,
											Name: "Pepper Street",
											Bred: "IRE",
										},
									},
								},
								{
									Type:     ToteExacta,
									Currency: "GBP",
									Dividend: makeDecimal(t, "137.70"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
										{
											ID:   2342636,
											Name: "Alliteration",
											Bred: "GB",
										},
									},
								},
								{
									Type:     ToteTrifecta,
									Currency: "GBP",
									Dividend: makeDecimal(t, "336.50"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
										{
											ID:   2342636,
											Name: "Alliteration",
											Bred: "GB",
										},
										{
											ID:   2298615,
											Name: "Pepper Street",
											Bred: "IRE",
										},
									},
								},
								{
									Type:     ToteSwinger,
									Currency: "GBP",
									Dividend: makeDecimal(t, "4.10"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2342636,
											Name: "Alliteration",
											Bred: "GB",
										},
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
									},
								},
								{
									Type:     ToteSwinger,
									Currency: "GBP",
									Dividend: makeDecimal(t, "2.00"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2342636,
											Name: "Alliteration",
											Bred: "GB",
										},
										{
											ID:   2298615,
											Name: "Pepper Street",
											Bred: "IRE",
										},
									},
								},
								{
									Type:     ToteSwinger,
									Currency: "GBP",
									Dividend: makeDecimal(t, "4.20"),
									Stake:    1,
									HorseRef: []HorseRef{
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
										{
											ID:   2298615,
											Name: "Pepper Street",
											Bred: "IRE",
										},
									},
								},
							},
							Bet: []Bet{
								{
									Type:     BetTypeCSF,
									Currency: "GBP",
									Dividend: makeDecimal(t, "101.18"),
									HorseRef: []HorseRef{
										{
											ID:   2338916,
											Name: "Fabianski",
											Bred: "IRE",
										},
										{
											ID:   2342636,
											Name: "Alliteration",
											Bred: "GB",
										},
									},
								},
							},
						},
						//SellingDetails  TODO // Details of horses sold or claimed after the result
					}},
					//Messages
					//MultiBet
				}},
			},
		},
	}

	for _, test := range tests {
		blob, err := ioutil.ReadFile(test.file)
		require.NoError(t, err)
		var obj RacingFile
		err = xml.Unmarshal(blob, &obj)
		assert.Equal(t, test.err, err, test.file)
		assert.Equal(t, test.obj, obj, test.file)
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		s   string
		d   time.Duration
		err error
	}{
		{
			s: "hi",
			err: &strconv.NumError{
				Func: "ParseInt",
				Num:  "hi",
				Err:  strconv.ErrSyntax,
			},
		},
		{
			s: "",
			d: 0,
		},
		{
			s: "0",
			d: 0,
		},
		{
			s: "0.1",
			d: time.Millisecond * 100,
		},
		{
			s: "0.100",
			d: time.Millisecond * 100,
		},
		{
			s: "0.01",
			d: time.Millisecond * 10,
		},
		{
			s: "0.001",
			d: time.Millisecond * 1,
		},
		{
			// <1ms values gets truncated to 0ms
			s: "0.0001",
			d: 0,
		},
		{
			s: "5",
			d: time.Second * 5,
		},
		{
			s: "80",
			d: time.Second * 80,
		},
		{
			s: "500",
			d: time.Minute * 5,
		},
		{
			s: "9900",
			d: time.Minute * 99,
		},
		{
			s: "10000",
			d: time.Hour * 1,
		},
		{
			s: "0102.003",
			d: time.Minute + time.Second*2 + time.Millisecond*3,
		},
	}

	for _, test := range tests {
		d, err := parseDuration(test.s)
		assert.Equal(t, test.d, d, fmt.Sprintf("input: %s", test.s))
		assert.Equal(t, test.err, err, fmt.Sprintf("input: %s", test.s))
	}
}
