package util

import (
	"strings"
	"testing"
)

func TestParseTime(t *testing.T) {
	creationTimestamp := "2024-11-28 17:22:12 -0300 -03"
	expected := "2024-11-28 17:22"

	resp, err := ParseTime(creationTimestamp)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.EqualFold(resp, expected) {
		t.Fatal(err)
	}
}
