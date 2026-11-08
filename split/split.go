// Package split divides a bill among people in whole cents, so the shares
// always add up to the bill.
package split

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Cents is an amount of money in hundredths of the currency's unit.
type Cents int64

// String writes the amount with two decimal places, such as 12.05 or -0.40.
func (c Cents) String() string {
	sign := ""
	if c < 0 {
		sign = "-"
		c = -c
	}
	return fmt.Sprintf("%s%d.%02d", sign, c/100, c%100)
}

// ParseCents reads an amount such as 12, 12.5 or 12.05. It refuses more than
// two decimal places, because a share of a cent cannot be paid.
func ParseCents(text string) (Cents, error) {
	whole, fraction, found := strings.Cut(strings.TrimSpace(text), ".")
	if whole == "" || found && fraction == "" || len(fraction) > 2 {
		return 0, fmt.Errorf("split: %q is not an amount with at most two decimal places", text)
	}
	units, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("split: %q is not an amount: %w", text, err)
	}
	cents := int64(0)
	if fraction != "" {
		if cents, err = strconv.ParseInt(fraction+strings.Repeat("0", 2-len(fraction)), 10, 64); err != nil || cents < 0 {
			return 0, fmt.Errorf("split: %q is not an amount", text)
		}
	}
	if strings.HasPrefix(whole, "-") {
		return Cents(units*100 - cents), nil
	}
	return Cents(units*100 + cents), nil
}

// Share is what one person owes.
type Share struct {
	Name   string
	Amount Cents
}

// ErrNoPeople is returned when a bill is split among nobody.
var ErrNoPeople = errors.New("split: a bill needs at least one person")

// Even splits total among the people in equal shares. The cents that do not
// divide evenly go one each to the first people, in the order given.
func Even(total Cents, people []string) ([]Share, error) {
	weights := make(map[string]int, len(people))
	for _, name := range people {
		weights[name] = 1
	}
	return Weighted(total, people, weights)
}

// Weighted splits total in proportion to each person's weight, such as the
// nights each stayed. Each share is rounded down to the cent, and the cents
// left over go one each to the people whose shares lost the most in rounding,
// ties to the earlier person in the order given. Every weight must be at
// least 1, and every person must have one.
func Weighted(total Cents, people []string, weights map[string]int) ([]Share, error) {
	if len(people) == 0 {
		return nil, ErrNoPeople
	}
	if total < 0 {
		return nil, fmt.Errorf("split: the total %s is negative", total)
	}
	sum := 0
	seen := make(map[string]bool, len(people))
	for _, name := range people {
		if seen[name] {
			return nil, fmt.Errorf("split: %s is named twice", name)
		}
		seen[name] = true
		weight, ok := weights[name]
		if !ok || weight < 1 {
			return nil, fmt.Errorf("split: %s has no weight of at least 1", name)
		}
		sum += weight
	}

	shares := make([]Share, len(people))
	remainders := make([]int64, len(people))
	allotted := Cents(0)
	for index, name := range people {
		exact := int64(total) * int64(weights[name])
		shares[index] = Share{Name: name, Amount: Cents(exact / int64(sum))}
		remainders[index] = exact % int64(sum)
		allotted += shares[index].Amount
	}

	order := make([]int, len(people))
	for index := range order {
		order[index] = index
	}
	sort.SliceStable(order, func(a, b int) bool {
		return remainders[order[a]] > remainders[order[b]]
	})
	for _, index := range order[:total-allotted] {
		shares[index].Amount++
	}
	return shares, nil
}
