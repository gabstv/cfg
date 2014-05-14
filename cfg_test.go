package cfg

import (
	"testing"
)

func TestSimpleConfigString(t *testing.T) {
	cfg0s := "Title = Sample App Title\r\nDescription = Line 1\\\nLine2\nRating = 1 \"STAR\""
	cfg0, err := ParseString(cfg0s)
	if err != nil {
		t.Fatal(err)
	}
	if cfg0["Title"] != "Sample App Title" {
		t.Error(`cfg0["Title"]`, cfg0["Title"], " != Sample App Title")
	}
	if cfg0["Description"] != "Line 1\nLine2" {
		t.Error(`cfg0["Description"]`, cfg0["Description"], " != Line 1\nLine2")
	}
	if cfg0["Rating"] != "1 \"STAR\"" {
		t.Error(`cfg0["Rating"]`, cfg0["Rating"], " != 1 \"STAR\"")
	}
}
