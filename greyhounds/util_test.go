package greyhounds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFinalResultsFile(t *testing.T) {
	tests := []struct {
		fileName  string
		meetingID int
		expected  bool
	}{
		{
			fileName:  "b2018041433736119270028.xml",
			meetingID: 337361,
			expected:  false,
		},
		{
			fileName:  "b2018041433735610380001.xml",
			meetingID: 337361,
			expected:  false,
		},
		{
			fileName:  "b201804143373611927.xml",
			meetingID: 337361,
			expected:  true,
		},
		{
			fileName:  "b201804143373561038.xml",
			meetingID: 337361,
			expected:  true,
		},
	}

	for _, test := range tests {
		final := IsFinalResultsFile(test.fileName, test.meetingID)
		assert.Equal(t, test.expected, final)
	}
}

func TestParseResult(t *testing.T) {
	tests := []struct {
		position      string
		expectedPlace int
		expectedDNF   bool
	}{
		{
			position:      "2",
			expectedPlace: 2,
			expectedDNF:   false,
		},
		{
			position:      "1",
			expectedPlace: 1,
			expectedDNF:   false,
		},
		{
			position:      "DN",
			expectedPlace: 0,
			expectedDNF:   true,
		},
		{
			position:      "",
			expectedPlace: 0,
			expectedDNF:   false,
		},
		{
			position:      "0",
			expectedPlace: 0,
			expectedDNF:   false,
		},
	}

	for _, test := range tests {
		place, dnf := ParseResult(test.position)
		assert.Equal(t, test.expectedPlace, place)
		assert.Equal(t, test.expectedDNF, dnf)
	}
}
