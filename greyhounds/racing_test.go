package greyhounds

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

func TestBulkParseGreyhoundRacing(t *testing.T) {
	dirs := []string{
		"testdata/Crayford",
		"testdata/Nottingham",
		"testdata/Perry Barr",
		"testdata/The Meadows",
		"testdata/Wheeling Island",
		"testdata/feed",
	}

	// Assertsions is a list of custom checks for successfuly parsing some
	// fields to their non default value.
	assertions := map[string]bool{}
	//catch is a helper function for marking custom test check as passed.
	catch := func(name string, val bool) {
		assertions[name] = assertions[name] || val
	}

	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err, dir)
		for _, f := range files {
			path := path.Join(dir, f.Name())
			t.Log(fmt.Sprintf("checking: %s", path))
			blob, err := ioutil.ReadFile(path)
			require.NoError(t, err, path)
			obj, err := ParseFile(blob)
			assert.NoError(t, err, path)

			assert.True(t, len(obj.Meetings) == 1, "always exactly one meeting perfile")
			for _, m := range obj.Meetings {
				assert.True(t, len(m.Races) >= 1, "always at least one race per meeting")
				for _, r := range m.Races {
					t.Log(fmt.Sprintf("time: %s", r.Time))
					assert.False(t, r.Time.Year() == 0, "Race time always has a date")

					catch("parsed Race with status Dormant", r.State == RaceDormant)
					// catch("parsed Race with status Delayed", r.State == RaceDelayed)
					catch("parsed Race with status Parading", r.State == RaceParading)
					catch("parsed Race with status Approaching", r.State == RaceApproaching)
					catch("parsed Race with status Going in traps", r.State == RaceGoingInTraps)
					catch("parsed Race with status Hare Running", r.State == RaceHareRunning)
					catch("parsed Race with status Off", r.State == RaceOff)
					// catch("parsed Race with status Blanket Finish", r.State == RaceBlanketFinish)
					// catch("parsed Race with status Result", r.State == RaceResult)
					catch("parsed Race with status Final Result", r.State == RaceFinalResult)
					// catch("parsed Race with status Race Void", r.State == RaceVoid)
					// catch("parsed Race with status No Race", r.State == RaceNoRace)
					// catch("parsed Race with status Rerun", r.State == RaceRerun)
					// catch("parsed Race with status Stewards Inquiry", r.State == RaceStewardsInquiry)
					// catch("parsed Race with status Stopped for Safety", r.State == RaceStoppedForSafety)
					// catch("parsed Race with status Abandoned", r.State == RaceAbandoned)
				}
			}
		}
	}

	for check, value := range assertions {
		assert.True(t, value, check)
	}
}

func TestParseGreyhoundRacing(t *testing.T) {
	tests := []struct {
		file string
		obj  DogRacing
		err  error
	}{
		{
			file: "testdata/Crayford/b2018041433736119270007.xml",
			obj: DogRacing{
				Type: MessageRace,
				Meetings: []Meeting{
					Meeting{
						MeetingID: 337361,
						Track:     "Crayford",
						Date:      makeTime(t, "2018-04-14T00:00:00Z"),
						State:     MeetingDormant,
						Races: []Race{
							Race{
								Revision:   7,
								RaceNumber: 1,
								Time:       makeTime(t, "2018-04-14T19:27:00+01:00"),
								Type:       RaceTypeFlat,
								Handicap:   false,
								Class:      "A7",
								Distance:   380,
								State:      RaceDormant,
								Traps: []Trap{
									Trap{
										TrapNo:  1,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   478812,
											Name: "Clonmannon Lady",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T19:21:56+01:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(6, 1),
												},
											},
										},
									},
									Trap{
										TrapNo:  2,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   497334,
											Name: "Kelva Matty",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T19:21:59+01:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(2, 1),
												},
											},
										},
									},
									Trap{
										TrapNo:  3,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   504096,
											Name: "Galtee Blue",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T19:22:04+01:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(2, 1),
												},
											},
										},
									},
									Trap{
										TrapNo:  4,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   482241,
											Name: "Cromac Terror",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T19:22:08+01:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(3, 1),
												},
											},
										},
									},
									Trap{
										TrapNo:  5,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   507585,
											Name: "Pesky Pigeon",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T19:22:13+01:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(5, 1),
												},
											},
										},
									},
									Trap{
										TrapNo:  6,
										Vacant:  false,
										Wide:    true,
										Reserve: false,
										Dog: &Dog{
											ID:   476879,
											Name: "Aoifes Speedy",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T19:22:17+01:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(6, 1),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "testdata/The Meadows/b201804143181110023.xml",
			obj: DogRacing{
				Type: MessageRace,
				Meetings: []Meeting{
					Meeting{
						MeetingID: 3181,
						Track:     "The Meadows",
						Country:   "Australia",
						Date:      makeTime(t, "2018-04-14T00:00:00Z"),
						State:     MeetingActive,
						Races: []Race{
							Race{
								Revision:   23,
								RaceNumber: 11,
								Time:       makeTime(t, "2018-04-14T12:45:00+00:00"),
								Type:       RaceTypeFlat,
								Handicap:   false,
								Class:      "A",
								Distance:   525,
								OffTime:    makeTime(t, "2018-04-14T12:46:27+00:00"),
								State:      RaceFinalResult,
								Traps: []Trap{
									Trap{
										TrapNo:  1,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   27251,
											Name: "Oh Jay Korie",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:25+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(8, 1),
												},
											},
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:45:46+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(9, 1),
												},
											},
										},
										Result: &Result{
											Position: "0",
											StartingPrice: &Price{
												Fractional: *big.NewRat(10, 1),
											},
										},
									},
									Trap{
										TrapNo:  2,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   20969,
											Name: "Like A Rocket",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:32+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(10, 3),
												},
											},
										},
										Result: &Result{
											Position: "0",
											StartingPrice: &Price{
												Fractional: *big.NewRat(3, 1),
											},
										},
									},
									Trap{
										TrapNo:  3,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   19531,
											Name: "One Plus Two",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:35+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(17, 2),
												},
											},
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:45:50+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(10, 1),
												},
											},
										},
										Result: &Result{
											Position: "2",
											StartingPrice: &Price{
												Fractional: *big.NewRat(12, 1),
											},
										},
									},
									Trap{
										TrapNo:  4,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   21013,
											Name: "Dyna Benny",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:41+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(11, 1),
												},
											},
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:45:52+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(12, 1),
												},
											},
										},
										Result: &Result{
											Position: "0",
											StartingPrice: &Price{
												Fractional: *big.NewRat(12, 1),
											},
										},
									},
									Trap{
										TrapNo:  5,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   23556,
											Name: "Orazzi Ohh",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:45+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(3, 1),
												},
											},
										},
										Result: &Result{
											Position: "0",
											StartingPrice: &Price{
												Fractional: *big.NewRat(10, 3),
											},
										},
									},
									Trap{
										TrapNo:  6,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   21020,
											Name: "He's Loaded",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:48+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(6, 4),
												},
											},
										},
										Result: &Result{
											Position: "1",
											StartingPrice: &Price{
												Fractional: *big.NewRat(5, 4),
											},
										},
									},
									Trap{
										TrapNo:  7,
										Vacant:  false,
										Wide:    false,
										Reserve: true,
										Dog: &Dog{
											ID:   21014,
											Name: "Fat Rhino",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:51+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(33, 1),
												},
											},
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:45:57+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(40, 1),
												},
											},
										},
										Result: &Result{
											Position: "0",
											StartingPrice: &Price{
												Fractional: *big.NewRat(50, 1),
											},
										},
									},
									Trap{
										TrapNo:  8,
										Vacant:  false,
										Wide:    false,
										Reserve: false,
										Dog: &Dog{
											ID:   26416,
											Name: "Dorrigo Bale",
										},
										Shows: []Show{
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:44:54+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(16, 1),
												},
											},
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:45:59+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(18, 1),
												},
											},
											Show{
												TimeStamp:    makeTime(t, "2018-04-14T12:46:03+00:00"),
												MarketNumber: 1,
												Price: &Price{
													Fractional: *big.NewRat(20, 1),
												},
											},
										},
										Result: &Result{
											Position: "3",
											StartingPrice: &Price{
												Fractional: *big.NewRat(25, 1),
											},
										},
									},
								},
								NonRunners: []NonRunner{
									NonRunner{
										Trap: 7,
										Dog: &Dog{
											ID:   35155,
											Name: "Elevated",
										},
									},
								},
								Dividends: &Dividends{
									Forecast: []Forecast{
										Forecast{
											Trap1:    6,
											Trap2:    3,
											Dividend: makeDecimal(t, "18.70"),
										},
									},
									Tricast: []Tricast{
										Tricast{
											Trap1:    6,
											Trap2:    3,
											Trap3:    8,
											Dividend: makeDecimal(t, "156.99"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		blob, err := ioutil.ReadFile(test.file)
		require.NoError(t, err)
		var obj DogRacing
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

func TestAddDate(t *testing.T) {
	tests := []struct {
		timePart time.Time
		datePart time.Time
		expected time.Time
	}{
		{
			timePart: time.Date(0, 0, 0, 15, 05, 30, 0, time.UTC),
			datePart: time.Date(2019, time.January, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2019, time.January, 15, 15, 05, 30, 0, time.UTC),
		},
		{
			timePart: time.Date(0, 0, 0, 6, 10, 0, 0, time.UTC),
			datePart: time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2020, time.March, 1, 6, 10, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		add := addDate(test.timePart, test.datePart)
		assert.Equal(t, test.expected, add)
	}
}
