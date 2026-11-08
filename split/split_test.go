package split

import (
	"testing"
)

func TestParseCentsReadsWholeAndDecimalAmounts(t *testing.T) {
	for text, want := range map[string]Cents{
		"12":    1200,
		"12.5":  1250,
		"12.05": 1205,
		"0.01":  1,
		" 3.10": 310,
		"-4.25": -425,
	} {
		got, err := ParseCents(text)
		if err != nil || got != want {
			t.Errorf("ParseCents(%q) = %d, %v, want %d", text, got, err, want)
		}
	}
}

func TestParseCentsRefusesAShareOfACent(t *testing.T) {
	for _, text := range []string{"", ".5", "1.", "1.234", "one", "1.-5"} {
		if got, err := ParseCents(text); err == nil {
			t.Errorf("ParseCents(%q) = %d, want an error", text, got)
		}
	}
}

func TestCentsStringKeepsTwoDecimalPlaces(t *testing.T) {
	for cents, want := range map[Cents]string{1205: "12.05", 40: "0.40", -40: "-0.40", 0: "0.00"} {
		if got := cents.String(); got != want {
			t.Errorf("Cents(%d).String() = %q, want %q", int64(cents), got, want)
		}
	}
}

func total(shares []Share) Cents {
	var sum Cents
	for _, share := range shares {
		sum += share.Amount
	}
	return sum
}

func TestEvenGivesTheLeftoverCentsToTheFirstPeople(t *testing.T) {
	shares, err := Even(1000, []string{"ana", "ben", "cem"})
	if err != nil {
		t.Fatal(err)
	}
	want := []Share{{"ana", 334}, {"ben", 333}, {"cem", 333}}
	for index, share := range shares {
		if share != want[index] {
			t.Errorf("share %d = %+v, want %+v", index, share, want[index])
		}
	}
	if total(shares) != 1000 {
		t.Errorf("the shares add up to %s, not 10.00", total(shares))
	}
}

func TestWeightedSharesAlwaysAddUpToTheTotal(t *testing.T) {
	people := []string{"ana", "ben", "cem", "dev"}
	weights := map[string]int{"ana": 3, "ben": 1, "cem": 2, "dev": 7}
	for amount := Cents(0); amount < 2000; amount += 7 {
		shares, err := Weighted(amount, people, weights)
		if err != nil {
			t.Fatal(err)
		}
		if total(shares) != amount {
			t.Fatalf("%s splits into shares that add up to %s", amount, total(shares))
		}
	}
}

func TestWeightedRefusesABillItCannotSplit(t *testing.T) {
	cases := map[string]struct {
		total   Cents
		people  []string
		weights map[string]int
	}{
		"nobody":           {100, nil, nil},
		"a negative bill":  {-100, []string{"ana"}, map[string]int{"ana": 1}},
		"a missing weight": {100, []string{"ana", "ben"}, map[string]int{"ana": 1}},
		"a zero weight":    {100, []string{"ana"}, map[string]int{"ana": 0}},
		"a name twice":     {100, []string{"ana", "ana"}, map[string]int{"ana": 1}},
	}
	for name, c := range cases {
		if _, err := Weighted(c.total, c.people, c.weights); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}
