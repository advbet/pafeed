package greyhounds

import (
	"encoding/xml"
)

// ParseFile unmarshals XML file contents to DogRacing object.
func ParseFile(xmlBlob []byte) (*DogRacing, error) {
	var obj DogRacing
	if err := xml.Unmarshal(xmlBlob, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
