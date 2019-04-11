package greyhounds

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// IsFinalResultsFile given a file name and meeting ID returns true if file
// should contain final results.
func IsFinalResultsFile(name string, meetingID int) bool {
	// The format is: b<date><meetingid><racetime>.xml e.g. b20140601896972052.xml
	return strings.HasPrefix(name, "b") && len(name) == len(fmt.Sprintf("b20140601%d2052.xml", meetingID))
}

// ParseFile unmarshals XML file contents to DogRacing object.
func ParseFile(xmlBlob []byte) (*DogRacing, error) {
	var obj DogRacing
	if err := xml.Unmarshal(xmlBlob, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// ParseResult parses PA position value and returns placement position
// and whether the dog did not finish the race.
func ParseResult(position string) (int, bool) {
	if position == "DN" {
		return 0, true
	}
	if placed, err := strconv.Atoi(position); err == nil {
		return placed, false
	}
	return 0, false
}
