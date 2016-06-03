package strftime

import (
	"bytes"
	"testing"
	"time"
)

type TestCase struct {
	format, value string
}

var testTime = time.Date(2009, time.November, 8, 23, 1, 2, 3, time.UTC)
var testCases = []*TestCase{
	&TestCase{"%a", "Sun"},
	&TestCase{"%A", "Sunday"},
	&TestCase{"%b", "Nov"},
	&TestCase{"%B", "November"},
	&TestCase{"%c", "Sun, 08 Nov 2009 23:01:02 UTC"},
	&TestCase{"%d", "08"},
	&TestCase{"%H", "23"},
	&TestCase{"%I", "11"},
	&TestCase{"%j", "312"},
	&TestCase{"%m", "11"},
	&TestCase{"%M", "01"},
	&TestCase{"%p", "PM"},
	&TestCase{"%S", "02"},
	&TestCase{"%U", "45"},
	&TestCase{"%w", "0"},
	&TestCase{"%W", "44"},
	&TestCase{"%x", "11/08/09"},
	&TestCase{"%X", "23:01:02"},
	&TestCase{"%y", "09"},
	&TestCase{"%Y", "2009"},
	&TestCase{"%Z", "UTC"},
	&TestCase{"%3n", "000"},
	&TestCase{"%6n", "000000"},
	&TestCase{"%9n", "000000003"},

	// Escape
	&TestCase{"%%%Y", "%2009"},
	&TestCase{"%3%%", "%3%"},
	&TestCase{"%3%3n", "%3000"},
	&TestCase{"%3xy%3n", "%3xy000"},
	// Embedded
	&TestCase{"/path/%Y/%m/report", "/path/2009/11/report"},
	//Empty
	&TestCase{"", ""},
}

func TestFormats(t *testing.T) {
	for _, tc := range testCases {
		value := Format(tc.format, testTime)
		if value != tc.value {
			t.Fatalf("error in %s: got %s instead of %s", tc.format, value, tc.value)
		}
	}
}

func TestUnknown(t *testing.T) {
	unknownFormat := "%g"
	value := Format(unknownFormat, testTime)
	if unknownFormat != value {
		t.Fatalf("error to in %s: got %s instead of %s", unknownFormat, value, unknownFormat)
	}
}

func TestFormatter_ValidFormats(t *testing.T) {
	for _, tc := range testCases {
		formatter := NewFormatter(tc.format)
		value := formatter.Format(testTime)
		if value != tc.value {
			t.Fatalf("error in %s: got %s instead of %s", tc.format, value, tc.value)
		}
		buf := bytes.NewBuffer(make([]byte, 0, 0))
		formatter.FormatTo(buf, testTime)
		if string(buf.Bytes()) != tc.value {
			t.Fatalf("error in %s: got %s instead of %s", tc.format, value, tc.value)
		}
	}
}

func TestFormatter_InvalidFormats(t *testing.T) {
	unknownFormat := "%g"
	formatter := NewFormatter(unknownFormat)
	value := formatter.Format(testTime)
	if unknownFormat != value {
		t.Fatalf("error to in %s: get %s instead of %s", unknownFormat, value, unknownFormat)
	}
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	formatter.FormatTo(buf, testTime)
	if unknownFormat != value {
		t.Fatalf("error to in %s: get %s instead of %s", unknownFormat, value, unknownFormat)
	}
}

type EdgeCase struct {
	t             time.Time
	format, value string
}

var (
	firstDay    = time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	firstSunday = time.Date(2009, time.January, 4, 0, 0, 0, 0, time.UTC)
	firstMonday = time.Date(2009, time.January, 5, 0, 0, 0, 0, time.UTC)
	jan2016     = time.Date(2016, time.January, 23, 00, 11, 10, 0, time.UTC)
)
var edgeCases = []*EdgeCase{
	&EdgeCase{firstDay, "%W", "00"},
	&EdgeCase{firstDay, "%U", "00"},
	&EdgeCase{firstSunday, "%W", "00"},
	&EdgeCase{firstSunday, "%U", "01"},
	&EdgeCase{firstMonday, "%W", "01"},
	&EdgeCase{firstMonday, "%U", "01"},
	&EdgeCase{jan2016, "%Yw%U", "2016w03"},
	&EdgeCase{time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2004-W53-6"},
	&EdgeCase{time.Date(2005, time.January, 2, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2004-W53-0"},
	&EdgeCase{time.Date(2005, time.December, 31, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2005-W52-6"},
	&EdgeCase{time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2007-W01-1"},
	&EdgeCase{time.Date(2007, time.December, 30, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2007-W52-0"},
	&EdgeCase{time.Date(2007, time.December, 31, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2008-W01-1"},
	&EdgeCase{time.Date(2008, time.January, 1, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2008-W01-2"},
	&EdgeCase{time.Date(2008, time.December, 28, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2008-W52-0"},
	&EdgeCase{time.Date(2008, time.December, 29, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W01-1"},
	&EdgeCase{time.Date(2008, time.December, 30, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W01-2"},
	&EdgeCase{time.Date(2008, time.December, 31, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W01-3"},
	&EdgeCase{time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W01-4"},
	&EdgeCase{time.Date(2009, time.December, 31, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W53-4"},
	&EdgeCase{time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W53-5"},
	&EdgeCase{time.Date(2010, time.January, 2, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W53-6"},
	&EdgeCase{time.Date(2010, time.January, 3, 0, 0, 0, 0, time.UTC), "%G-W%V-%w", "2009-W53-0"},
}

func TestEdgeCases(t *testing.T) {
	for _, tc := range edgeCases {
		value := Format(tc.format, tc.t)
		if value != tc.value {
			t.Fatalf("error in %s (%v): got %s instead of %s", tc.format, tc.t, value, tc.value)
		}
	}
}
