package horses

import (
	"encoding/xml"
	"strings"
)

// IsRacingFile given a file name returns true if file should should contain
// Racing message.
func IsRacingFile(name string) bool {
	return strings.HasPrefix(name, "b")
}

// IsRacingCardFile given a file name returns true if file should should contain
// RacingCard message.
func IsRacingCardFile(name string) bool {
	return strings.HasPrefix(name, "c")
}

// ParseRacing unmarshals Racing XML file to Racing object. This function should
// be used for files that passes IsRacingFile() check.
func ParseRacing(xmlBlob []byte) (*Racing, error) {
	var obj Racing
	if err := xml.Unmarshal(xmlBlob, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// ParseRacingCard unmarshals Racing XML file to Racing object. This function
// should be used for files that passes IsRacingCardFile() check.
func ParseRacingCard(xmlBlob []byte) (*RacingCard, error) {
	var obj RacingCard
	if err := xml.Unmarshal(xmlBlob, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
