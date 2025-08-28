package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DF = "20060102"

func getLastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	date, err := time.Parse(DF, dstart)
	if err != nil {
		return "", fmt.Errorf("Error dayStart: %v", err)
	}
	originalDate := date
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("Repeat empty string")
	}
	parts := strings.Split(repeat, " ")
	switch parts[0] {

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("Invalid formal of d rule")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("Invalid interval(convertation): %v", err)
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("Invalid interval(out of range): %d", interval)
		}
		for !date.After(now) {
			date = date.AddDate(0, 0, interval)
		}
		if date.Equal(originalDate) && originalDate.After(now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(DF), nil

	case "y":
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		if date.Equal(originalDate) && originalDate.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DF), nil

	case "w":
		var weekDays [7]bool
		for _, day := range strings.Split(parts[1], ",") {
			nday, err := strconv.Atoi(day)
			if err != nil {
				return "", fmt.Errorf("Invalid formal of weekday: %v", err)
			}
			if nday < 1 || nday > 7 {
				return "", fmt.Errorf("Invalid formal of weekday: %d", nday)
			}
			w := nday % 7
			weekDays[w] = true
		}
		d := date
		for {
			if d.After(now) && weekDays[int(d.Weekday())] {
				return d.Format(DF), nil
			}
			d = d.AddDate(0, 0, 1)
			if d.Sub(date) > 730*24*time.Hour {
				return "", fmt.Errorf("no date found for w rule")
			}
		}

	case "m":
		var days [32]bool
		var months [13]bool
		last, prelast := false, false
		specialMonths := len(parts) == 3

		//Processing the days of the month
		for _, s := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				return "", fmt.Errorf("Invalid formal of month: %v", err)
			}
			switch n {
			case -1:
				last = true
			case -2:
				prelast = true
			default:
				if n < 1 || n > 31 {
					return "", fmt.Errorf("Day out of range: %d", n)
				}
				days[n] = true
			}
		}

		//Processing months if exist
		if specialMonths {
			for _, s := range strings.Split(parts[2], ",") {
				n, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("Months out of range: %d", n)
				}
				months[n] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				months[i] = true
			}
		}

		d := date
		for i := 0; i < 731; i++ { //two years scanning
			monthNum := int(d.Month())
			if d.After(now) && months[monthNum] {
				lastDay := getLastDayOfMonth(d)
				if days[d.Day()] || (last && d.Day() == lastDay) || (prelast && d.Day() == lastDay-1) {
					return d.Format(DF), nil
				}
			}
			d = d.AddDate(0, 0, 1)
		}
		return "", fmt.Errorf("No months found")

	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}
}
