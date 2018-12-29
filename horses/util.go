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

// ParseRacingFile unmarshals Racing XML file contents to RacingFile object.
// This function should be used for files that passes IsRacingFile() check.
func ParseRacingFile(xmlBlob []byte) (*RacingFile, error) {
	var obj RacingFile
	if err := xml.Unmarshal(xmlBlob, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// ParseRacingCardFile unmarshals RacingCard XML file contents to RacingCardFile
// object. This function should be used for files that passes IsRacingCardFile()
// check.
func ParseRacingCardFile(xmlBlob []byte) (*RacingCardFile, error) {
	var obj RacingCardFile
	if err := xml.Unmarshal(xmlBlob, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
