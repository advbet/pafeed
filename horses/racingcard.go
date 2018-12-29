package horses

import (
	"encoding/xml"
	"fmt"
	"time"

	"bitbucket.org/advbet/decimal"
)

// RacingCard is main data structure in the PA horse racing feed, that
// describe upcoming horse racing meetings.
type RacingCard []CardMeeting

// CardMeetingStatus is an enum for meeting status in horse racing cards.
type CardMeetingStatus string

// CardMeeting describes a single horse racing meeting. It is similar to
// Meeting, main difference is meeting cards get sent before Meeting.
type CardMeeting struct {
	ID      int               // Meeting internal database ID
	Country string            // The country where the meeting is being held
	Course  string            // The course where the meeting is being held
	Date    time.Time         // Date when the meeting starts (format ISO 8601:1988 yyyymmdd)
	Status  CardMeetingStatus // Meeting status, one of Dormant, Inspection, Abandoned
	//DeclarationStage UNUSED // Declaration stage of races at the meeting (summarised), one of Early, Final, Mixed
	WeatherForecast string    // Forecasted weather for meeting
	Inspection      time.Time // Present if the meeting is subject to an inspection
	AbandonedReason string    // Gives the reason for a meeting being abandoned
	DrawAdvantage   string    // Generalised comment about advantage gained from stalls position
	AdvancedGoing   string    // Indication of expected going at the meeting
	//Messages      UNUSED    // Other textual messages associated with meeting
	Races []CardRace // Meeting races
}

// CardRace describes a single race in the horse racing card meeting.
type CardRace struct {
	ID        int       // The internal identifier for the race
	StartTime time.Time // The date of the race (format ISO 8601:1988 yyyymmdd)
	RaceType  RaceType  //Type of race (Flat, Hurdle, Chase, National Hunt Flat)
	TrackType TrackType // The type of surface being raced on
	Handicap  bool      // Whether or not this race is a handicap
	Trifecta  bool      // Whether or not this race has a trifecta associated with it
	Showcase  bool      // Whether or not this is a showcase race
	Class     int       // The class of the race
	//DeclarationStage UNUSED //Declaration stage of the race. Early - used for early declarations (fourday etc). Final - used for final declarations (overnight etc)
	MaxRunners    int                    // Maximum field size
	NumFences     int                    // Number of fences to be jumped
	Title         string                 // The title of the race
	AddedMoney    *MoneyValue            // The money added to the prize fund for the race
	PenaltyValue  *MoneyValue            // The prize money awarded to the winner
	PrizeCurrency string                 // The currency of the prize money
	Prizes        map[int]decimal.Number // Map from finishing position to prize amount
	//Fees        UNUSED // Fees associated with the race
	Eligibility string         // The type of horses eligible for the race. Example: 3yo plus.
	Distance    UnitsValueText // The distance of the race
	//WeightsRaised   UNUSED // Amount weights raised (at overnight stage)
	//LastWinner      *TODO  // The winner of corresponding race last year
	//Conditions      UNUSED // The conditions for the race (penalty weights etc)
	//Televised       UNUSED // Television coverage details
	//RaceFlags       UNUSED // Optional extra info breaking down type of race etc.
	//PreviewComments UNUSED // Preview text comment(s)
	//Selections      UNUSED // Selections (tips) for race
	//DrawBias        UNUSED // The effect of the draw in this race (Flat races only)
	Ratings []Rating // Race ratings
	//Messages UNUSED `xml:"Message"`        // Other textual messages associated with race
	//Totes  []TODO `xml:"Tote"`  // Tote bets applicable to this race
	Horses []CardHorse `xml:"Horse"` // The horse(s)

}

// CardHorse contains data about a single horse participating in a race. This
// object is sent in race cards and include more details then general Horse
// object that is in normal racing messages.
type CardHorse struct {
	ID          int             // The internal identifier for the horse
	Name        string          // The name of the horse
	Bred        string          // The country of breeding of the horse
	Status      CardHorseStatus // Horse status - Runner, Doubtful
	ClothNumber int             // The saddlecloth number for the horse
	DrawnStall  int             // The stall the horse starts from (Flat races only)
	//FormFigures     []struct{}      // Recent form (figures) for the horse
	//LastRunDays     []struct{}      // Number of days since the horse last ran
	//RaceHistoryStat []struct{}      // The race history for the horse
	AgeInYears        int            // The age of the horse (in years)
	Weight            UnitsValueText // The weight carried by the horse
	WeightPenalty     UnitsValue     // Extra weight incurred through recent win
	Trainer           CardTrainer    // Details of the trainer of the horse
	OwnerName         string         // Details of the owner of the horse
	BreederName       string         // Details of the breeder of the horse
	Jockey            CardJockey     // Details of the jockey of the horse
	JockeyColours     string         // Textual description of the jockey's colours (silks)
	JockeyColoursFile string         // Name of the graphics file which represents the the jockey's colours (silks)
	//Tackle          []struct{}      // The tackle which the horse will be wearing
	//Career          []struct{}      // The career performance for the horse
	Colours  []string   // The colour(s) of the horse
	Sex      Sex        // The sex of the horse
	Breeding []Breeding // The lineage of the horse
	//Lineage         *struct{}       // Lineage comment for horse
	//FoalDate        *struct{}       // Date horse was foaled
	//Comment         *struct{}       // Textual comment for the horse
	//ForecastPrice   *struct{}       // The betting forecast price for the horse
	//StartingPrice   *struct{}       // Starting price of horse (used in LastWinner context)
	//Rating          []struct{}      // Ratings associated with this horse
	//Reserve         *struct{}       // Reserve details IF this horse is a reserve
	//Ballot          *struct{}       // Ballot order details
	//LongHandicap    *struct{}       // The long handicap details for this horse (if applicable)
	//Medication      *struct{}       // Medication taken by the horse in the form race
	//Travelled       *struct{}       // Distance travelled by horse to course
	//FormRace        []struct{}      // Previous race form for this horse
	//PinSticker      []struct{}      // Pin sticker comments
	//Analysis        *struct{}       // Analysis of horses chance of winning
	//Message         UNUSED       // Other textual messages associated with horse
}

// CardTrainer holds horse trainer details. This field is sent with racing cards
// and have more information then Trainer object.
type CardTrainer struct {
	ID          int    // Identifier for trainer
	Name        string // The name of the trainer
	Nationality string // The nationality of the trainer eg IRE
	Location    string // Where the trainer is based
	//PersonForm UNUSED // Indicates how well the trainer is currently doing
}

// CardHorseStatus is an enum for horse status values.
type CardHorseStatus string

// CardJockey contains data about the person riding a horse in a race. This
// object is sent only in race cards and contains less detauls than Jockey
// object.
type CardJockey struct {
	ID        int        // Identifier for jockey
	Name      string     // The name of the jockey
	Allowance UnitsValue // Allowance of the jockey units in which allowance value is pecified
	//PersonForm UNUSED  // Indicates how well the jockey is currently doing
}

// Breeding describes a horse from the racing horse direct lineage.
type Breeding struct {
	Relation HorseRelation // Sire (father), Dam (mother), DamSire (maternal grandfather)
	Name     string        // The name of the horse
	Bred     string        // The country of breeding of the horse
	YearBord int           // When the horse was born (if known)
}

// HorseRelation describes a breeding relation between two horses.
type HorseRelation string

// Rating is a single instance of race ratings.
type Rating struct {
	Type  string // Type of rating e.g. Official.
	Value int    // Rating value e.g. 57.
}

// RaceType is an enum for race types - Flat, Hurdle, Chase, National Hunt Flat.
type RaceType string

// TrackType is an enum for horse racing track types.
type TrackType string

// Sex is an enum of horse sex values.
type Sex string

// List of allowed CardMeetingStatus values.
const (
	CardMeetingDormant    CardMeetingStatus = "Dormant"    // the meeting is going ahead as planned
	CardMeetingInspection CardMeetingStatus = "Inspection" // the meeting is subject to an inspection
	CardMeetingAbandoned  CardMeetingStatus = "Abandoned"  // the meeting has been abandoned
)

// List of allowed RaceType values.
const (
	RaceFlat             RaceType = "Flat"
	RaceHurdle           RaceType = "Hurdle"
	RaceChase            RaceType = "Chase"
	RaceNationalHuntFlat RaceType = "N_H_Flat" // National Hunt Flat
)

// List of allowed TrackType values.
const (
	TrackTurf       TrackType = "Turf"
	TrackFibresand  TrackType = "Fibresand"
	TrackPolytrack  TrackType = "Polytrack"
	TrackEquitrack  TrackType = "Equitrack"
	TrackDirt       TrackType = "Dirt"
	TrackSand       TrackType = "Sand"
	TrackAllWeather TrackType = "AllWeather"
)

// List of allowed CardHorseStatus values.
const (
	CardHorseRunner   CardHorseStatus = "Runner"
	CardHorseDoubtful CardHorseStatus = "Doubtful"
)

// List of allowed horse sex values.
const (
	Filly    Sex = "f" // filly (young female horse <= 4 years)
	Colt     Sex = "c" // colt (young male horse <= 4 years)
	Mare     Sex = "m" // mare (female horse)
	Stallion Sex = "h" // horse (male non castrated)
	Gelding  Sex = "g" // gelding (male castrated)
	Ridgling Sex = "r" // ridgling (half castrated naturally?)
)

// List of allowed HorseRelation values.
const (
	Sire    HorseRelation = "Sire"    // father
	Dam     HorseRelation = "Dam"     // mother
	DamSire HorseRelation = "DamSire" // maternal grandfather
)

// UnmarshalXML implements xml.Unmarshaler interface.
func (c *RacingCard) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		Meetings []CardMeeting `xml:"Meeting"` // The meeting(s)
	}{}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	*c = data.Meetings
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (m *CardMeeting) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		ID               int               `xml:"id,attr"`       // Meeting internal database ID
		Country          string            `xml:"country,attr"`  // The country where the meeting is being held
		Course           string            `xml:"course,attr"`   // The course where the meeting is being held
		Date             xmlDate           `xml:"date,attr"`     // Date when the meeting starts (format ISO 8601:1988 yyyymmdd)
		Status           CardMeetingStatus `xml:"status,attr"`   // Meeting status, one of Dormant, Inspection, Abandoned
		DeclarationStage string            `xml:"decStage,attr"` // Declaration stage of races at the meeting (summarised), one of Early, Final, Mixed
		WeatherForecast  struct {
			Data string `xml:",chardata"`
		} `xml:"WeatherForecast"` // Forecasted weather for meeting
		Inspection xmlTimeElement `xml:"Inspection"` // Present if the meeting is subject to an inspection
		Abandoned  struct {
			Data string `xml:",chardata"` // Gives the reason for a meeting being abandoned
		} `xml:"Abandoned"` // Present if the meeting has been abandoned
		DrawAdvantage struct {
			Data string `xml:",chardata"` // Generalised comment about advantage gained from stalls position
		} `xml:"DrawAdvantage"` // Draw Advantage details (where appropriate)
		AdvancedGoing struct {
			Data string `xml:",chardata"` // Indication of expected going at the meeting
		} `xml:"AdvancedGoing"` // The advanced going for the meeting.
		//Messages          UNUSED  `xml:"Message"`         // Other textual messages associated with meeting
		Races []CardRace `xml:"Race"` // The race(s)
		//MultiBets         []struct{ TODO }  `xml:"MultiBet"`        // Multi-race bets available on this meeting
	}{}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	*m = CardMeeting{
		ID:      data.ID,
		Country: data.Country,
		Course:  data.Course,
		Date:    time.Time(data.Date),
		Status:  data.Status,
		//DeclarationStage UNUSED
		WeatherForecast: data.WeatherForecast.Data,
		Inspection:      time.Time(data.Inspection),
		AbandonedReason: data.Abandoned.Data,
		DrawAdvantage:   data.DrawAdvantage.Data,
		AdvancedGoing:   data.AdvancedGoing.Data,
		//Messages UNUSED
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *CardRace) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		ID        int       `xml:"id,attr"`        // The internal identifier for the race
		Date      string    `xml:"date,attr"`      // The date of the race (format ISO 8601:1988 yyyymmdd)
		Time      string    `xml:"time,attr"`      // The time of the race (format ISO 8601:1988 hhmm+/-hhmm)
		RaceType  RaceType  `xml:"raceType,attr"`  // Type of race (Flat, Hurdle, Chase, National Hunt Flat)
		TrackType TrackType `xml:"trackType,attr"` // The type of surface being raced on
		Handicap  xmlYesNo  `xml:"handicap,attr"`  // Whether or not this race is a handicap
		Trifecta  xmlYesNo  `xml:"trifecta,attr"`  // Whether or not this race has a trifecta associated with it
		Showcase  xmlYesNo  `xml:"showcase,attr"`  // Whether or not this is a showcase race
		Class     int       `xml:"class,attr"`     // The class of the race
		//DeclarationStage UNUSED `xml:"decStage,attr"`        //Declaration stage of the race. Early - used for early declarations (fourday etc). Final - used for final declarations (overnight etc)
		MaxRunners int `xml:"maxRunners,attr"` // Maximum field size
		NumFences  int `xml:"numFences,attr"`  // Number of fences to be jumped
		Title      struct {
			Data string `xml:",chardata"`
		} `xml:"Title"` // The title of the race
		AddedMoney   *xmlMoneyValue `xml:"AddedMoney"`   // The money added to the prize fund for the race
		PenaltyValue *xmlMoneyValue `xml:"PenaltyValue"` // The prize money awarded to the winner
		PrizeMoney   struct {
			Currency string `xml:"currency,attr"` // The currency of the prize money
			Prize    []struct {
				Position int `xml:"position,attr"` // Finishing position the prize is for
				Amount   int `xml:"amount,attr"`   // Prize amount (currency specified in PrizeMoney element)
			} `xml:"Prize"` // Prize Element
		} // Prize money awarded for the race
		//Fees UNUSED `xml:"Fees"           // Fees associated with the race
		Eligibility struct {
			Type string `xml:"type,attr"` // The type of horses eligible for the race. Example: 3yo plus.
		} `xml:"Eligibility"` // The horses eligible in the race
		Distance xmlUnitsValueText `xml:"Distance"` // The distance of the race
		//WeightsRaised UNUSED `xml:"WeightsRaised"`  // Amount weights raised (at overnight stage)
		LastWinner *struct {
			Year         int    `xml:"year,attr"`   // The year of the corresponding race
			NoRaceReason string `xml:"noRace,attr"` // Reason if race was not run
			Runners      int    `xml:"ran,attr"`    // The number of horses that raced
			//Horses     []TODO `xml:"Horse"`       // The winner(s) details (if race run)
		} `xml:"LastWinner"` // The winner of corresponding race last year
		//Conditions      UNUSED `xml:"Conditions"` // The conditions for the race (penalty weights etc)
		//Televised       UNUSED `xml:"Televised"`  // Television coverage details
		//RaceFlags       UNUSED `xml:"RaceFlags"`  // Optional extra info breaking down type of race etc.
		//PreviewComments UNUSED `xml:"Preview"`    // Preview text comment(s)
		//Selections      UNUSED `xml:"Selections"` // Selections (tips) for race
		//DrawBias        UNUSED `xml:"DrawBias"`   // The effect of the draw in this race (Flat races only)
		Ratings []Rating `xml:"Rating"` // Race ratings
		//Messages UNUSED `xml:"Message"`        // Other textual messages associated with race
		//Totes  []TODO `xml:"Tote"`  // Tote bets applicable to this race
		Horses []CardHorse `xml:"Horse"` // The horse(s)
	}{}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	startTime, err := time.Parse("20060102T1504-0700", fmt.Sprintf("%sT%s", data.Date, data.Time))
	if err != nil {
		return fmt.Errorf("parsing CardRace.date and CardRace.time: %v", err)
	}
	prizes := make(map[int]decimal.Number)
	for _, prize := range data.PrizeMoney.Prize {
		prizes[prize.Position] = decimal.FromInt(prize.Amount)
	}
	*r = CardRace{
		ID:        data.ID,
		StartTime: startTime,
		RaceType:  data.RaceType,
		TrackType: data.TrackType,
		Handicap:  bool(data.Handicap),
		Trifecta:  bool(data.Trifecta),
		Showcase:  bool(data.Showcase),
		Class:     data.Class,
		//DeclarationStage UNUSED
		MaxRunners:    data.MaxRunners,
		NumFences:     data.NumFences,
		Title:         data.Title.Data,
		AddedMoney:    (*MoneyValue)(data.AddedMoney),
		PenaltyValue:  (*MoneyValue)(data.PenaltyValue),
		PrizeCurrency: data.PrizeMoney.Currency,
		Prizes:        prizes,
		//Fees        UNUSED
		Eligibility: data.Eligibility.Type,
		Distance:    UnitsValueText(data.Distance),
		//WeightsRaised   UNUSED
		//LastWinner      *TODO
		//Conditions      UNUSED
		//Televised       UNUSED
		//RaceFlags       UNUSED
		//PreviewComments UNUSED
		//Selections      UNUSED
		//DrawBias        UNUSED
		Ratings: data.Ratings,
		//Messages UNUSED
		//Totes  []TODO
		Horses: data.Horses,
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (h *CardHorse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		ID     int             `xml:"id,attr"`     // The internal identifier for the horse
		Name   string          `xml:"name,attr"`   // The name of the horse
		Bred   string          `xml:"bred,attr"`   // The country of breeding of the horse
		Status CardHorseStatus `xml:"status,attr"` // Horse status - Runner, Doubtful
		Cloth  struct {
			Number int `xml:"number,attr"` // Saddlecloth or racecard number of horse
			//Coupled UNUSED `xml:"coupled,attr"` // In races where two or more horses have been "coupled" together, these horses share the same "number" but have an additional letter to be able to tell them apart. For example 1 and 1a.
		} `xml:"Cloth"` // The saddlecloth number for the horse
		Drawn struct {
			Stall int `xml:"stall"` // The stall this horse will start from
		} `xml:"Drawn"` // The stall the horse starts from (Flat races only)
		//FormFigures     []TODO `xml:"FormFigures"`     // Recent form (figures) for the horse
		//LastRunDays     []TODO `xml:"LastRunDays"`     // Number of days since the horse last ran
		//RaceHistoryStat []TODO `xml:"RaceHistoryStat"` // The race history for the horse
		Age struct {
			Years int `xml:"years,attr"` // The age of the horse in years.
		} `xml:"Age"` // The age of the horse (in years)
		Weight        xmlUnitsValueText `xml:"Weight"`        // The weight carried by the horse
		WeightPenalty xmlUnitsValue     `xml:"WeightPenalty"` // Extra weight incurred through recent win
		Trainer       CardTrainer       `xml:"Trainer"`       // Details of the trainer of the horse
		Owner         struct {
			Name string `xml:"name,attr"` // The name of the owner
		} `xml:"Owner"` // Details of the owner of the horse
		Breeder struct {
			Name string `xml:"name,attr"` // The name of the breeder
		} `xml:"Breeder"` // Details of the breeder of the horse
		Jockey        CardJockey `xml:"Jockey"` // Details of the jockey of the horse
		JockeyColours struct {
			Filename    string `xml:"filename,attr"`    // The name of the graphics file which represents the colours
			Description string `xml:"description,attr"` // Textual description of jockey colours
		} `xml:"JockeyColours"` // Details of the jockey's colours (silks)
		//Tackle          []TODO `xml:"Tackle"`          // The tackle which the horse will be wearing
		//Career          []TODO `xml:"Career"`          // The career performance for the horse
		Colours []struct {
			Type string `xml:"type,attr"` // Colour of horse (e.g. ch = chestnut)
		} `xml:"Colour"` // The colour(s) of the horse
		Sex struct {
			Type Sex `xml:"type,attr"` // f = filly, c = colt, m = mare, h = horse, g = gelding, r = ridgling
		} `xml:"Sex"` // The sex of the horse
		Breeding []Breeding `xml:"Breeding"` // The lineage of the horse
		//Lineage         *struct{}  `xml:"Lineage"`         // Lineage comment for horse
		//FoalDate        *struct{}  `xml:"FoalDate"`        // Date horse was foaled
		//Comment         *struct{}  `xml:"Comment"`         // Textual comment for the horse
		//ForecastPrice   *struct{}  `xml:"ForecastPrice"`   // The betting forecast price for the horse
		//StartingPrice   *struct{}  `xml:"StartingPrice"`   // Starting price of horse (used in LastWinner context)
		//Rating          []struct{} `xml:"Rating"`          // Ratings associated with this horse
		//Reserve         *struct{}  `xml:"Reserve"`         // Reserve details IF this horse is a reserve
		//Ballot          *struct{}  `xml:"Ballot"`          // Ballot order details
		//LongHandicap    *struct{}  `xml:"LongHandicap"`    // The long handicap details for this horse (if applicable)
		//Medication      *struct{}  `xml:"Medication"`      // Medication taken by the horse in the form race
		//Travelled       *struct{}  `xml:"Travelled"`       // Distance travelled by horse to course
		//FormRace        []struct{} `xml:"FormRace"`        // Previous race form for this horse
		//PinSticker      []struct{} `xml:"PinSticker"`      // Pin sticker comments
		//Analysis        *struct{}  `xml:"Analysis"`        // Analysis of horses chance of winning
		//Message       UNUSED  `xml:"Message"`         // Other textual messages associated with horse
	}{
		Status: CardHorseRunner,
	}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	colours := make([]string, 0, len(data.Colours))
	for _, c := range data.Colours {
		colours = append(colours, c.Type)
	}
	*h = CardHorse{
		ID:                data.ID,
		Name:              data.Name,
		Bred:              data.Bred,
		Status:            data.Status,
		ClothNumber:       data.Cloth.Number,
		DrawnStall:        data.Drawn.Stall,
		AgeInYears:        data.Age.Years,
		Weight:            UnitsValueText(data.Weight),
		WeightPenalty:     UnitsValue(data.WeightPenalty),
		Trainer:           data.Trainer,
		OwnerName:         data.Owner.Name,
		BreederName:       data.Breeder.Name,
		Jockey:            data.Jockey,
		JockeyColours:     data.JockeyColours.Description,
		JockeyColoursFile: data.JockeyColours.Filename,
		Colours:           colours,
		Sex:               data.Sex.Type,
		Breeding:          data.Breeding,
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (r *Rating) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		Type  string `xml:"type,attr"`  // Type of rating e.g. Official.
		Value int    `xml:"value,attr"` // Rating value e.g. 57.
	}{}
	*r = Rating(data)
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (t *CardTrainer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		ID          int    `xml:"id,attr"`          // Identifier for trainer
		Name        string `xml:"name,attr"`        // The name of the trainer
		Nationality string `xml:"nationality,attr"` // The nationality of the trainer eg IRE
		Location    string `xml:"location,attr"`    // Where the trainer is based
		//PersonForm UNUSED `xml:"PersonForm"` // Indicates how well the trainer is currently doing
	}{}

	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	*t = CardTrainer{
		ID:          data.ID,
		Name:        data.Name,
		Nationality: data.Nationality,
		Location:    data.Location,
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (j *CardJockey) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		ID        int           `xml:"id,attr"`   // Identifier for jockey
		Name      string        `xml:"name,attr"` // The name of the jockey
		Allowance xmlUnitsValue `xml:"Allowance"` // The allowance of the jockey
		//PersonForm UNUSED  `xml:"PersonForm"` // Indicates how well the jockey is currently doing
	}{}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	*j = CardJockey{
		ID:        data.ID,
		Name:      data.Name,
		Allowance: UnitsValue(data.Allowance),
		//PersonForm UNUSED
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (b *Breeding) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		Type     HorseRelation `xml:"type,attr"`     // Sire (father), Dam (mother), DamSire (maternal grandfather)
		Name     string        `xml:"name,attr"`     // The name of the horse
		Bred     string        `xml:"bred,attr"`     // The country of breeding of the horse
		YearBorn int           `xml:"yearBorn,attr"` // When the horse was born (if known)
	}{}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	*b = Breeding{
		Relation: data.Type,
		Name:     data.Name,
		Bred:     data.Bred,
		YearBord: data.YearBorn,
	}
	return nil
}
