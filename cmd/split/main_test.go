package main

import (
	"bytes"
	"testing"
)

func TestRunPrintsEachShareInOrder(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"100", "ana=3", "ben", "cem=2"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	if want := "ana\t50.00\nben\t16.67\ncem\t33.33\n"; stdout.String() != want {
		t.Errorf("printed\n%s\nwant\n%s", stdout.String(), want)
	}
}

func TestRunRefusesAWeightThatIsNotAWholeNumber(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"100", "ana=1.5"}, &stdout, &stderr); code != 2 || stdout.Len() != 0 {
		t.Errorf("exit %d, printed %q", code, stdout.String())
	}
}
