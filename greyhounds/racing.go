package greyhounds

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/advbet/decimal"
)

// MessageType is the type of the message (Card or Race)
type MessageType string

// CardState is the state of the card message (Advance or Final)
type CardState string

// MeetingState is the state of the meeting
type MeetingState string

// RaceType is the type of the race
type RaceType string

// RaceState is the state of the race
type RaceState string

// TrapSeeding is the seeding of the trap
type TrapSeeding string

// DogSex is the sex of a dog
type DogSex string

// xmlTimeElement is a date value with cusom XML unmarshaler that reads ISO 8601:1988
// date value.
type xmlTimeElement time.Time

// xmlYesNo is typed boolean with custom XML unmarshaler that converts Yes/No
// string values to boolean value.
type xmlYesNo bool

// DogRacing is the main object sent via PA greyhound racing feed. It holds all
// the information for a single (or more) greyhound races.
type DogRacing struct {
	Type  MessageType // Message type
	State string      // Applicable only if type is Card indicates if card is an advance or final card

	Meetings []Meeting // The meeting(s)
}

// Meeting object describes a greyhound racing meeting with information on zero
// or more races.
type Meeting struct {
	MeetingID int          // The unique identifier for this meeting
	Track     string       // The track at which the meeting is being held
	Country   string       // The country where the track is located
	Date      time.Time    // The date when the meeting is started
	State     MeetingState // The current state of the meeting

	Races       []Race // Race details for this meeting
	ReserveDogs []Dog  // List of reserve dogs for meeting
}

type xmlMeeting Meeting

// Race holds the details of a single race.
type Race struct {
	Revision   int           // Incremental revision number used in race messages
	RaceNumber int           // The number of this race within the meeting
	Time       time.Time     // The time that the race is scheduled to start
	Type       RaceType      // Type of race (Flat or Hurdles)
	Handicap   bool          // Whether the race is a handicap
	Class      string        // The class of the race
	Distance   int           // The distance of the race (metres)
	Title      string        // The title of the race.
	Prizes     string        // The prizes awarded for the race
	OffTime    time.Time     // The time that the race actually started
	Going      string        // The going allowance for this race
	WinTime    time.Duration // The time taken to complete the race
	State      RaceState     // The current state of this race
	Bags       bool          // Whether the race uses bags
	Tricast    bool          // Indicates whether a tricast will be returned

	Comments   []Comment   // The comments on this race
	Traps      []Trap      // Trap Details for this race.
	NonRunners []NonRunner // Any dogs that were withdrawn from this race
	Dividends  *Dividends  // Dividends paid on the result (forecast, tricast etc)
}

type xmlRace Race

// Trap holds the information of a single race trap.
type Trap struct {
	TrapNo   int         // The number of the trap this dog was due to start from
	Vacant   bool        // Whether the trap is vacant (empty)
	Wide     bool        // Whether this trap contains a wide runner
	Seeding  TrapSeeding // One of Wide, Mid or Rails
	Handicap string      // The handicap start for this trap
	Reserve  bool        // Whether this trap contains one of the reserve dogs
	Photo    int         // If this trap is involved in a photo finish this attribute contains the position the photo is for.

	Dog    *Dog    // Details of the dog racing from this trap
	Shows  []Show  // Betting shows associated with this trap
	Result *Result // Result details for this trap
}

type xmlTrap Trap

// Comment contains information of a comment
type Comment struct {
	Source string // Source description e.g. PA, Timeform
	Type   string // Description of the comment type
	Text   string // Comment text
}

type xmlComment Comment

// NonRunner contains details of the trap that is a nonrunner
type NonRunner struct {
	Trap   int    // The number of the trap this dog was due to start from
	Reason string // Reason for dog withdrawal

	Dog *Dog // Details of the dog that is a nonrunner
}

type xmlNonRunner NonRunner

// Dog contains all the information of a single dog
type Dog struct {
	ID                  int    // The id number of this dog
	Name                string // The name of this dog
	Origin              string // The origin of this dog
	ForecastPriceSource string // The source of forecast price

	ForecastPrice *Price         // Forecast price of the dog
	BestTime      *BestTime      // Best time for this dog
	ExpectedTimes []ExpectedTime // Expected sectional and final times for this dog
	Breeding      *Breeding      // Breeding related details for this dog
	Trainer       Trainer        // Trainer details for this dog
	Owner         Owner          // Owner details for this dog
	Ratings       []Rating       // Rating(s) for the dog
	Comments      []Comment      // Comment(s) for the dog
	FormRaces     []FormRace     // The form (list of previous runs) for this dog
}

type xmlDog Dog

// BestTime is the fastest runtime for the greyhound at the same course
// and over the same distance as the race the dog is currently entered for
// within the period of the last 3 months.
//
// If no runtime is available for the specified period, then the best run
// time ever run by the dog at that course and distance is used.
//
// If no run time is still available then the time will be 0.00
type BestTime struct {
	AdjustedTime time.Duration // The finishing time adjusted for going and handicap.
	Date         time.Time     // The date of the race in which the best time occurred
	RaceNumber   int           // The racenumber in the meeting where the best time occurred
	MeetingID    int           // The PA meeting id of the meeting where the best time occurred
	Class        string        // The class of the race where the best time occurred
}

type xmlBestTime BestTime

// Breeding contains information of a single dogs breed
type Breeding struct {
	Sire   string    // The sire of the dog
	Dam    string    // The dam of the dog
	Born   time.Time // The date on which the dog was born
	Colour string    // The colour of the dog
	Sex    DogSex    // The sex of the dog
	Season string    // The season of the dog (bitches only)
}

type xmlBreeding Breeding

// Trainer contains information for a dogs trainer
type Trainer struct {
	ID    int    // The trainers id
	Name  string // The trainers name
	Track string // The track this trainer is "local" to
}

type xmlTrainer Trainer

// Owner holds the information of a dogs owner
type Owner struct {
	ID   int    // The owners id
	Name string // The owners name
}

type xmlOwner Owner

// Rating contains the information of a single rating
type Rating struct {
	Source string // Source description e.g. PA, Timeform
	Type   string // The type of rating (e.g. star)
	Value  string // The value of the data
}

type xmlRating Rating

// ExpectedTime contains data of a single expected time for a dog
type ExpectedTime struct {
	Source string // The data supplier
	Type   string // The type of time (e.g. final, sectional)
	Value  string // The value of the data
}

type xmlExpectedTime ExpectedTime

// FormRace contains information of a single form race
type FormRace struct {
	MeetingID   int           // The meeting identified for this race
	Track       string        // The track where the race was held
	RaceNumber  int           // The number of this race within the meeting
	Going       string        // The going allowance for this race
	Date        time.Time     // The time that the race started
	Type        RaceType      // Type of race (Flat or Hurdles)
	Class       string        // The class of the race
	Distance    int           // The distance of the race
	WinningTime time.Duration // The winning dog's time

	FormTraps []FormTrap // Details of race result
}

type xmlFormRace FormRace

// FormTrap contains information on a trap of a form race
type FormTrap struct {
	Trap     int         // The number of the trap this dog was due to start from.
	Wide     bool        // Whether this trap contained a wide runner.
	Seeding  TrapSeeding // One of Wide, Mid or Rails
	Handicap string      // The handicap start for this trap (if appropriate)

	Dog    *Dog    // The dog that started from this trap.
	Result *Result // The result details for this dog.
}

type xmlFormTrap FormTrap

// Result contains the result information of a single trap
type Result struct {
	Position      string        // The finish position of the dog
	BtnDistance   string        // If the Result element is contained within a FormTrap Element this is the distance between this dog and the winner. If the Result element is contained within a Trap Element this is the distance between this dog and the dog in front
	SectionalTime time.Duration // The time taken to reach the first bend.
	BendPosition  string        // The dog's position at each bend.
	RunComment    string        // The description of how the dog ran
	RunTime       time.Duration // The time taken by this dog to run the race
	Weight        float64       // The weight of the dog (kilograms)
	AdjustedTime  time.Duration // The finishing time adjusted for going and handicap.

	StartingPrice *Price // The price returned for this dog
}

type xmlResult Result

// Show contains the details of a single show
type Show struct {
	TimeStamp    time.Time // The time at which the show was available
	MarketNumber int       // When more than one betting market has been formed, this attribute indicates which market the show is applicable to, otherwise it will be absent.
	NoOffers     bool      // If no show price is currently being offered then this will be true

	Price *Price // Show price. Absent only if noOffers attribute is true.
}

type xmlShow Show

// Price contains the single price value for a dog
type Price struct {
	Decimal    decimal.Number // Decimal representation of the price (empty or in HK format)
	Fractional big.Rat        // Fractional representation of the price
}

type xmlPrice Price

// Dividends contains forecast, tricast information
type Dividends struct {
	Forecast []Forecast // The forecast dividends
	Tricast  []Tricast  // The tricast dividends
}

type xmlDividends Dividends

// Forecast contains the dividend paid for correctly predicting first two dogs in race
type Forecast struct {
	Trap1    int            // The trap number of the 1st placed dog
	Trap2    int            // The trap number of the 2nd placed dog
	Dividend decimal.Number // The amount paid (to a unit stake)
}

type xmlForecast Forecast

// Tricast contains the dividend paid for correctly predicting first three dogs in race
type Tricast struct {
	Trap1    int            // The trap number of the 1st placed dog
	Trap2    int            // The trap number of the 2nd placed dog
	Trap3    int            // The trap number of the 3rd placed dog
	Dividend decimal.Number // The amount paid (to a unit stake)
}

type xmlTricast Tricast

// List of allowed MessageType values.
const (
	MessageCard MessageType = "Card" // Greyhound Card information
	MessageRace MessageType = "Race" // individual race shows and result
)

// List of allowed CardState values.
const (
	CardAdvance CardState = "Advance" // The card is advance
	CardFinal   CardState = "Final"   // The card is final
)

// List of allowed MeetingState values.
const (
	MeetingDormant   MeetingState = "Dormant"   // The meeting has not started yet
	MeetingActive    MeetingState = "Active"    // The meeting has started
	MeetingDelayed   MeetingState = "Delayed"   // The meeting is currently delayed
	MeetingFinished  MeetingState = "Finished"  // The meeting has finished
	MeetingAbandoned MeetingState = "Abandoned" // The meeting has been abandoned
)

// List of allowed RaceType values.
const (
	RaceTypeFlat    RaceType = "Flat"    // The race is flat
	RaceTypeHurdles RaceType = "Hurdles" // The race has hurdles
)

// List of allowed RaceState values.
const (
	RaceDormant          RaceState = "Dormant"            // The race has not yet started
	RaceDelayed          RaceState = "Delayed"            // The race is currently delayed
	RaceParading         RaceState = "Parading"           // The dogs are parading
	RaceApproaching      RaceState = "Approaching"        // The dogs are approaching the traps
	RaceGoingInTraps     RaceState = "Going in traps"     // The dogs are being put into the traps
	RaceHareRunning      RaceState = "Hare Running"       // The hare has been started
	RaceOff              RaceState = "Off"                // The race has started
	RaceBlanketFinish    RaceState = "Blanket Finish"     // The race has finished with a blanket finish (many dogs involved in photo finish)
	RaceResult           RaceState = "Result"             // Result message
	RaceFinalResult      RaceState = "Final Result"       // Final version of result
	RaceVoid             RaceState = "Race Void"          // The race has been declared void
	RaceNoRace           RaceState = "No Race"            // The race has been declared a "no race" (followed by either "Rerun" or "Race Void")
	RaceRerun            RaceState = "Rerun"              // The race will be rerun (this usually occurs after the last race)
	RaceStewardsInquiry  RaceState = "Stewards Inquiry"   // There is a stewards inquiry (Note: greyhound results are not subject to amendment)
	RaceStoppedForSafety RaceState = "Stopped for Safety" // The race has been stopped in the interest of safety (followed by either "Rerun" or "Race Void")
	RaceAbandoned        RaceState = "Abandoned"          // The race has been abandoned
	RaceMeetingAbandoned RaceState = "Meeting Abandoned"  // The meeting has been abandoned
	RaceFinished         RaceState = "Finished"           // The race is finished
	RacePhotoSecond      RaceState = "Photo Second"
	RacePhotoThird       RaceState = "Photo Third"
	RaceTrapFailure      RaceState = "Trap Failure"
	RaceHareFailure      RaceState = "Hare Failure"
)

// List of allowed TrapSeeding values.
const (
	SeedingWide  TrapSeeding = "Wide"  // Wide trap seeding
	SeedingMid   TrapSeeding = "Mid"   // Mid trap seeding
	SeedingRails TrapSeeding = "Rails" // Rails trap seeding
)

// List of allowed DogSex values.
const (
	SexDog   DogSex = "d" // Dog
	SexBitch DogSex = "b" // Bitch
)

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *DogRacing) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Type  MessageType `xml:"type,attr"`
		State string      `xml:"state,attr"`

		Meetings []xmlMeeting `xml:"Meeting"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	if !data.Type.isValid() {
		return fmt.Errorf("invalid Message type attibute value: %s", data.Type)
	}
	var meetings []Meeting
	for _, m := range data.Meetings {
		meetings = append(meetings, Meeting(m))
	}
	*r = DogRacing{
		Type:     data.Type,
		State:    data.State,
		Meetings: meetings,
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (m *xmlMeeting) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initial values for fields having non empty defaults
	data := struct {
		MeetingID int            `xml:"meetingId,attr"`
		Track     string         `xml:"track,attr"`
		Country   string         `xml:"country,attr"`
		Date      xmlTimeElement `xml:"date,attr"`
		State     MeetingState   `xml:"state,attr"`

		Races       []xmlRace `xml:"Race"`
		ReserveDogs struct {
			Dogs []xmlDog `xml:"Dog"`
		} `xml:"ReserveDogs"`
	}{
		State: MeetingDormant,
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	if !data.State.isValid() {
		return fmt.Errorf("invalid Meeting state attibute value: %s", data.State)
	}
	var races []Race
	for _, r := range data.Races {
		if r.Time.Year() == 0 { // Get full date
			r.Time = addDate(r.Time, time.Time(data.Date))
		}
		if r.OffTime.Year() == 0 { // Get full date
			r.OffTime = addDate(r.OffTime, time.Time(data.Date))
		}

		for i, t := range r.Traps {
			for j, s := range t.Shows {
				if s.TimeStamp.Year() == 0 { // Get full date
					s.TimeStamp = addDate(s.TimeStamp, time.Time(data.Date))
					t.Shows[j] = s
				}
			}
			r.Traps[i] = t
		}

		races = append(races, Race(r))
	}
	var reserveDogs []Dog
	for _, r := range data.ReserveDogs.Dogs {
		reserveDogs = append(reserveDogs, Dog(r))
	}
	*m = xmlMeeting{
		MeetingID:   data.MeetingID,
		Track:       data.Track,
		Country:     data.Country,
		Date:        time.Time(data.Date),
		State:       data.State,
		Races:       races,
		ReserveDogs: reserveDogs,
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *xmlRace) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		Revision   int            `xml:"revision,attr"`
		RaceNumber int            `xml:"raceNumber,attr"`
		Time       xmlTimeElement `xml:"time,attr"`
		Type       RaceType       `xml:"type,attr"`
		Handicap   xmlYesNo       `xml:"handicap,attr"`
		Class      string         `xml:"class,attr"`
		Distance   int            `xml:"distance,attr"`
		Title      string         `xml:"title,attr"`
		Prizes     string         `xml:"prizes,attr"`
		OffTime    xmlTimeElement `xml:"offTime,attr"`
		Going      string         `xml:"going,attr"`
		WinTime    string         `xml:"winTime,attr"`
		State      RaceState      `xml:"state,attr"`
		Bags       xmlYesNo       `xml:"Bags,attr"`
		Tricast    xmlYesNo       `xml:"tricast,attr"`
		Comments   struct {
			Comments []xmlComment `xml:"Comment"`
		} `xml:"Comments"`
		Traps      []xmlTrap      `xml:"Trap"`
		NonRunners []xmlNonRunner `xml:"NonRunner"`
		Dividends  *xmlDividends  `xml:"Dividends"`
	}{
		State: RaceDormant,
	}

	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	if !data.Type.isValid() {
		return fmt.Errorf("invalid Race type attibute value: %s", data.Type)
	}
	if !data.State.isValid() {
		return fmt.Errorf("invalid Race state attibute value: %s", data.State)
	}

	winTime, err := parseDuration(data.WinTime)
	if err != nil {
		return err
	}
	var comments []Comment
	for _, c := range data.Comments.Comments {
		comments = append(comments, Comment(c))
	}
	var traps []Trap
	for _, t := range data.Traps {
		traps = append(traps, Trap(t))
	}
	var nonRunners []NonRunner
	for _, nr := range data.NonRunners {
		nonRunners = append(nonRunners, NonRunner(nr))
	}
	*r = xmlRace{
		Revision:   data.Revision,
		RaceNumber: data.RaceNumber,
		Time:       time.Time(data.Time),
		Type:       data.Type,
		Handicap:   bool(data.Handicap),
		Class:      data.Class,
		Distance:   data.Distance,
		Title:      data.Title,
		Prizes:     data.Prizes,
		OffTime:    time.Time(data.OffTime),
		Going:      data.Going,
		WinTime:    winTime,
		State:      data.State,
		Bags:       bool(data.Bags),
		Tricast:    bool(data.Tricast),
		Comments:   comments,
		Traps:      traps,
		NonRunners: nonRunners,
		Dividends:  (*Dividends)(data.Dividends),
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *xmlTrap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		TrapNo   int         `xml:"trap,attr"`
		Vacant   xmlYesNo    `xml:"vacant,attr"`
		Wide     xmlYesNo    `xml:"wide,attr"`
		Seeding  TrapSeeding `xml:"seeding,attr"`
		Handicap string      `xml:"handicap,attr"`
		Reserve  xmlYesNo    `xml:"reserve,attr"`
		Photo    int         `xml:"photo,attr"`

		Dog    *xmlDog    `xml:"Dog"`
		Shows  []xmlShow  `xml:"Show"`
		Result *xmlResult `xml:"Result"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	var shows []Show
	for _, s := range data.Shows {
		shows = append(shows, Show(s))
	}
	*t = xmlTrap{
		TrapNo:   data.TrapNo,
		Vacant:   bool(data.Vacant),
		Wide:     bool(data.Wide),
		Seeding:  data.Seeding,
		Handicap: data.Handicap,
		Reserve:  bool(data.Reserve),
		Photo:    data.Photo,
		Dog:      (*Dog)(data.Dog),
		Shows:    shows,
		Result:   (*Result)(data.Result),
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *xmlDog) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		ID     int    `xml:"id,attr"`
		Name   string `xml:"name,attr"`
		Origin string `xml:"origin,attr"`

		BestTime      *xmlBestTime `xml:"BestTime"`
		ExpectedTimes struct {
			ExpectedTimes []xmlExpectedTime `xml:"ExpectedTime"`
		} `xml:"ExpectedTimes"`
		Breeding      *xmlBreeding `xml:"Breeding"`
		Trainer       xmlTrainer   `xml:"Trainer"`
		Owner         xmlOwner     `xml:"Owner"`
		Ratings       []xmlRating  `xml:"Rating"`
		Comments      []xmlComment `xml:"Comment"`
		ForecastPrice struct {
			Source string    `xml:"source,attr"`
			Price  *xmlPrice `xml:"Price"`
		} `xml:"ForecastPrice"`
		Form struct {
			FormRaces []xmlFormRace `xml:"FormRace"`
		} `xml:"Form"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	var expectedTimes []ExpectedTime
	for _, t := range data.ExpectedTimes.ExpectedTimes {
		expectedTimes = append(expectedTimes, ExpectedTime(t))
	}
	var ratings []Rating
	for _, r := range data.Ratings {
		ratings = append(ratings, Rating(r))
	}
	var comments []Comment
	for _, c := range data.Comments {
		comments = append(comments, Comment(c))
	}
	var formRaces []FormRace
	for _, r := range data.Form.FormRaces {
		formRaces = append(formRaces, FormRace(r))
	}
	*e = xmlDog{
		ID:                  data.ID,                            // The id number of this dog
		Name:                data.Name,                          // The name of this dog
		Origin:              data.Origin,                        // The origin of this dog
		ForecastPriceSource: data.ForecastPrice.Source,          // The source of forecast price
		ForecastPrice:       (*Price)(data.ForecastPrice.Price), // Forecast price of the dog
		BestTime:            (*BestTime)(data.BestTime),         // Best time for this dog
		ExpectedTimes:       expectedTimes,                      // Expected sectional and final times for this dog
		Breeding:            (*Breeding)(data.Breeding),         // Breeding related details for this dog
		Trainer:             Trainer(data.Trainer),              // Trainer details for this dog
		Owner:               Owner(data.Owner),                  // Owner details for this dog
		Ratings:             ratings,                            // Rating(s) for the dog
		Comments:            comments,                           // Comment(s) for the dog
		FormRaces:           formRaces,                          // The form (list of previous runs) for this dog
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (c *xmlComment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Source  string `xml:"source,attr"`
		Type    string `xml:"type,attr"`
		Comment string `xml:",innerxml"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*c = xmlComment{
		Source: data.Source,  // Source description e.g. PA, Timeform
		Type:   data.Type,    // Description of the comment type
		Text:   data.Comment, // Comment text
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *xmlNonRunner) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Trap   int     `xml:"trap,attr"`
		Reason string  `xml:"reasonForWithdrawal,attr"`
		Dog    *xmlDog `xml:"Dog"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*e = xmlNonRunner{
		Trap:   data.Trap,        // The number of the trap this dog was due to start from
		Reason: data.Reason,      // Reason for dog withdrawal
		Dog:    (*Dog)(data.Dog), // Details of the dog that is a nonrunner
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *xmlBestTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		AdjustedTime string         `xml:"adjustedTime,attr"`
		Date         xmlTimeElement `xml:"date,attr"`
		RaceNumber   int            `xml:"raceNumber,attr"`
		MeetingID    int            `xml:"meetingId,attr"`
		Class        string         `xml:"class,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	adjustedTime, err := parseDuration(data.AdjustedTime)
	if err != nil {
		return err
	}
	*t = xmlBestTime{
		AdjustedTime: adjustedTime,         // The finishing time adjusted for going and handicap.
		Date:         time.Time(data.Date), // The date of the race in which the best time occurred
		RaceNumber:   data.RaceNumber,      // The racenumber in the meeting where the best time occurred
		MeetingID:    data.MeetingID,       // The PA meeting id of the meeting where the best time occurred
		Class:        data.Class,           // The class of the race where the best time occurred
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (b *xmlBreeding) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Sire   string         `xml:"sire,attr"`
		Dam    string         `xml:"dam,attr"`
		Born   xmlTimeElement `xml:"born,attr"`
		Colour string         `xml:"colour,attr"`
		Sex    DogSex         `xml:"sex,attr"`
		Season string         `xml:"season,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*b = xmlBreeding{
		Sire:   data.Sire,            // The sire of the dog
		Dam:    data.Dam,             // The dam of the dog
		Born:   time.Time(data.Born), // The date on which the dog was born
		Colour: data.Colour,          // The colour of the dog
		Sex:    data.Sex,             // The sex of the dog
		Season: data.Season,          // The season of the dog (bitches only)
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *xmlTrainer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		ID    int    `xml:"id,attr"`
		Name  string `xml:"name,attr"`
		Track string `xml:"track,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*t = xmlTrainer{
		ID:    data.ID,    // The trainers id
		Name:  data.Name,  // The trainers name
		Track: data.Track, // The track this trainer is "local" to
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (o *xmlOwner) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		ID   int    `xml:"id,attr"`
		Name string `xml:"name,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*o = xmlOwner{
		ID:   data.ID,   // The owners id
		Name: data.Name, // The owners name
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *xmlRating) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Source string `xml:"source,attr"`
		Type   string `xml:"type,attr"`
		Value  string `xml:"value,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*r = xmlRating{
		Source: data.Source, // Source description e.g. PA, Timeform
		Type:   data.Type,   // The type of rating (e.g. star)
		Value:  data.Value,  // The value of the data
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *xmlExpectedTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Source string `xml:"source,attr"`
		Type   string `xml:"type,attr"`
		Value  string `xml:"value,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*t = xmlExpectedTime{
		Source: data.Source, // Source description e.g. PA, Timeform
		Type:   data.Type,   // The type of rating (e.g. star)
		Value:  data.Value,  // The value of the data
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *xmlFormRace) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		MeetingID   int            `xml:"meetingId,attr"`
		Track       string         `xml:"track,attr"`
		Date        xmlTimeElement `xml:"date,attr"`
		RaceNumber  int            `xml:"raceNumber,attr"`
		Going       string         `xml:"going,attr"`
		Time        xmlTimeElement `xml:"time,attr"`
		Type        RaceType       `xml:"type,attr"`
		Class       string         `xml:"class,attr"`
		Distance    int            `xml:"distance,attr"`
		WinningTime string         `xml:"winningTime,attr"`

		FormTraps []xmlFormTrap `xml:"FormTrap"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	winningTime, err := parseDuration(data.WinningTime)
	if err != nil {
		return err
	}
	startTime := time.Time(data.Time)
	if startTime.Year() == 0 { // Get full date
		startTime = addDate(startTime, time.Time(data.Date))
	}
	var traps []FormTrap
	for _, t := range data.FormTraps {
		traps = append(traps, FormTrap(t))
	}
	*r = xmlFormRace{
		MeetingID:   data.MeetingID,  // The meeting identified for this race
		Track:       data.Track,      // The track where the race was held
		RaceNumber:  data.RaceNumber, // The number of this race within the meeting
		Going:       data.Going,      // The going allowance for this race
		Date:        startTime,       // The time that the race started
		Type:        data.Type,       // Type of race (Flat or Hurdles)
		Class:       data.Class,      // The class of the race
		Distance:    data.Distance,   // The distance of the race
		WinningTime: winningTime,     // The winning dog's time
		FormTraps:   traps,           // Details of race result
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *xmlFormTrap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Trap     int         `xml:"trap,attr"`
		Wide     xmlYesNo    `xml:"wide,attr"`
		Seeding  TrapSeeding `xml:"seeding,attr"`
		Handicap string      `xml:"handicap,attr"`
		Dog      *xmlDog     `xml:"Dog"`
		Result   *xmlResult  `xml:"Result"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*t = xmlFormTrap{
		Trap:     data.Trap,              // The number of the trap this dog was due to start from.
		Wide:     bool(data.Wide),        // Whether this trap contained a wide runner.
		Seeding:  data.Seeding,           // One of Wide, Mid or Rails
		Handicap: data.Handicap,          // The handicap start for this trap (if appropriate)
		Dog:      (*Dog)(data.Dog),       // The dog that started from this trap.
		Result:   (*Result)(data.Result), // The result details for this dog.
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *xmlResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Position      string  `xml:"position,attr"`
		BtnDistance   string  `xml:"btnDistance,attr"`
		SectionalTime string  `xml:"sectionalTime,attr"`
		BendPosition  string  `xml:"bendPosition,attr"`
		RunComment    string  `xml:"runComment,attr"`
		RunTime       string  `xml:"runTime,attr"`
		Weight        float64 `xml:"weight,attr"`
		AdjustedTime  string  `xml:"adjustedTime,attr"`

		StartingPrice struct {
			MarketPos int `xml:"marketPos,attr"`
			MarketCnt int `xml:"marketCnt,attr"`

			Price *xmlPrice `xml:"Price"`
		} `xml:"StartingPrice"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	sectionalTime, err := parseDuration(data.SectionalTime)
	if err != nil {
		return err
	}
	runTime, err := parseDuration(data.RunTime)
	if err != nil {
		return err
	}
	adjustedTime, err := parseDuration(data.AdjustedTime)
	if err != nil {
		return err
	}
	*r = xmlResult{
		Position:      data.Position,                      // The finish position of the dog
		BtnDistance:   data.BtnDistance,                   // If the Result element is contained within a FormTrap Element this is the distance between this dog and the winner. If the Result element is contained within a Trap Element this is the distance between this dog and the dog in front
		SectionalTime: sectionalTime,                      // The time taken to reach the first bend.
		BendPosition:  data.BendPosition,                  // The dog's position at each bend.
		RunComment:    data.RunComment,                    // The description of how the dog ran
		RunTime:       runTime,                            // The time taken by this dog to run the race
		Weight:        data.Weight,                        // The weight of the dog (kilograms)
		AdjustedTime:  adjustedTime,                       // The finishing time adjusted for going and handicap.
		StartingPrice: (*Price)(data.StartingPrice.Price), // The price returned for this dog
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (s *xmlShow) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		TimeStamp    xmlTimeElement `xml:"timeStamp,attr"`
		MarketNumber int            `xml:"marketNumber,attr"`
		NoOffers     xmlYesNo       `xml:"noOffers,attr"`

		Price *xmlPrice `xml:"Price"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*s = xmlShow{
		TimeStamp:    time.Time(data.TimeStamp), // The time at which the show was available
		MarketNumber: data.MarketNumber,         // When more than one betting market has been formed, this attribute indicates which market the show is applicable to, otherwise it will be absent.
		NoOffers:     bool(data.NoOffers),       // If no show price is currently being offered then this will be true

		Price: (*Price)(data.Price), // Show price. Absent only if noOffers attribute is true.
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (p *xmlPrice) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Decimal     decimal.Number `xml:"decimal,attr"`
		Numerator   int            `xml:"numerator,attr"`
		Denominator int            `xml:"denominator,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	var fraction big.Rat
	if data.Denominator != 0 {
		fraction = *big.NewRat(int64(data.Numerator), int64(data.Denominator))
	}
	*p = xmlPrice{
		Decimal:    data.Decimal, // Decimal representation of the price
		Fractional: fraction,     // Fractional representation of the price
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *xmlDividends) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Forecast []xmlForecast `xml:"Forecast"`
		Tricast  []xmlTricast  `xml:"Tricast"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	var forecast []Forecast
	for _, f := range data.Forecast {
		forecast = append(forecast, Forecast(f))
	}
	var tricast []Tricast
	for _, t := range data.Tricast {
		tricast = append(tricast, Tricast(t))
	}
	*e = xmlDividends{
		Forecast: forecast, // The forecast dividends
		Tricast:  tricast,  // The tricast dividends
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (f *xmlForecast) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Trap1    int            `xml:"trap1,attr"`
		Trap2    int            `xml:"trap2,attr"`
		Dividend decimal.Number `xml:"dividend,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*f = xmlForecast{
		Trap1:    data.Trap1,    // The trap number of the 1st placed dog
		Trap2:    data.Trap2,    // The trap number of the 2nd placed dog
		Dividend: data.Dividend, // The amount paid (to a unit stake)
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *xmlTricast) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		Trap1    int            `xml:"trap1,attr"`
		Trap2    int            `xml:"trap2,attr"`
		Trap3    int            `xml:"trap3,attr"`
		Dividend decimal.Number `xml:"dividend,attr"`
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}

	*t = xmlTricast{
		Trap1:    data.Trap1,    // The trap number of the 1st placed dog
		Trap2:    data.Trap2,    // The trap number of the 2nd placed dog
		Trap3:    data.Trap3,    // The trap number of the 3rd placed dog
		Dividend: data.Dividend, // The amount paid (to a unit stake)
	}
	return nil
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr intrface.
func (t *xmlTimeElement) UnmarshalXMLAttr(attr xml.Attr) error {
	var tm time.Time
	var err error
	switch len(attr.Value) {
	case 8:
		tm, err = time.Parse("20060102", attr.Value)
	case 9:
		tm, err = time.Parse("1504-0700", attr.Value)
	case 10:
		if len(strings.Split(attr.Value, "/")) == 3 {
			tm, err = time.Parse("02/01/2006", attr.Value)
		} else {
			tm, err = time.Parse("02-01-2006", attr.Value)
		}
	case 11:
		tm, err = time.Parse("150405-0700", attr.Value)
	case 17:
		tm, err = time.Parse("200601021504-0700", attr.Value)
	case 18:
		tm, err = time.Parse("20060102T1504-0700", attr.Value)
	case 19:
		tm, err = time.Parse("20060102150405-0700", attr.Value)
	case 20:
		tm, err = time.Parse("20060102T150405-0700", attr.Value)
	}
	if err != nil {
		return fmt.Errorf("parsing %v attribute (%s): %v", attr.Name, attr.Value, err)
	}
	*t = xmlTimeElement(tm)
	return nil
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr intrface.
func (b *xmlYesNo) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "yes", "Yes":
		*b = true
		return nil
	case "no", "No":
		*b = false
		return nil
	default:
		return fmt.Errorf("parsing %s attribute as Yes/No field, unexpected value: %s", attr.Name, attr.Value)
	}
}

func (s MeetingState) isValid() bool {
	switch s {
	case MeetingDormant,
		MeetingActive,
		MeetingDelayed,
		MeetingFinished,
		MeetingAbandoned:
		return true
	default:
		return false
	}
}

func (s MessageType) isValid() bool {
	switch s {
	case MessageCard,
		MessageRace:
		return true
	default:
		return false
	}
}

func (s RaceState) isValid() bool {
	switch s {
	case RaceDormant,
		RaceDelayed,
		RaceParading,
		RaceApproaching,
		RaceGoingInTraps,
		RaceHareRunning,
		RaceOff,
		RaceBlanketFinish,
		RaceResult,
		RaceFinalResult,
		RaceVoid,
		RaceNoRace,
		RaceRerun,
		RaceStewardsInquiry,
		RaceStoppedForSafety,
		RaceAbandoned,
		RaceMeetingAbandoned,
		RaceFinished,
		RacePhotoSecond,
		RacePhotoThird,
		RaceTrapFailure,
		RaceHareFailure:
		return true
	default:
		return false
	}
}

func (s RaceType) isValid() bool {
	switch s {
	case RaceTypeFlat,
		RaceTypeHurdles:
		return true
	default:
		return false
	}
}

// parseDudation converts ISO 8601:1988 mmss.ss formated string to golang
// time.Duration value.
func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	n, err := decimal.FromString(s)
	if err != nil {
		return 0, err
	}
	hour := n.ScaledVal(4) % 100
	mins := n.ScaledVal(2) % 100
	secs := n.ScaledVal(0) % 100
	mils := n.ScaledVal(-3) % 1000
	return time.Hour*time.Duration(hour) +
		time.Minute*time.Duration(mins) +
		time.Second*time.Duration(secs) +
		time.Millisecond*time.Duration(mils), err
}

func addDate(t time.Time, date time.Time) time.Time {
	hour, min, sec := t.Clock()
	year, month, day := date.Date()

	return time.Date(year, month, day, hour, min, sec, 0, t.Location())
}
