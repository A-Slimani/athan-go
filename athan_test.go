package main

import (
	"testing"
)

func TestBuildAthanString(t *testing.T) {
	// hours and minutes plural
	result := buildAthanString(3, 12, "Fajr")
	expected := "Fajr in 3 hours and 12 minutes"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}

	// hours plural, minutes singular
	result = buildAthanString(3, 1, "Fajr")
	expected = "Fajr in 3 hours and 1 minute"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}

	// hours singular, minutes plural
	result = buildAthanString(1, 2, "Fajr")
	expected = "Fajr in 1 hour and 2 minutes"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}

	// hours and minutes singular
	result = buildAthanString(1, 1, "Fajr")
	expected = "Fajr in 1 hour and 1 minute"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}

	// minutes zero
	result = buildAthanString(1, 0, "Fajr")
	expected = "Fajr in 1 hour"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}

	// hours zero
	result = buildAthanString(0, 1, "Fajr")
	expected = "Fajr in 1 minute"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}

	// hours and minutes zero
	result = buildAthanString(0, 0, "Fajr")
	expected = "Fajr is now \n"
	if result != expected {
		t.Errorf("Expected: %s, returned: %s ", expected, result)
	}
}

func TestCacheAthanTimes(t *testing.T) {

}
