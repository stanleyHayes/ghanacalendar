package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

type Occurrence struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Date         *string `json:"date"`
	ObservedDate *string `json:"observedDate,omitempty"`
	Kind         string  `json:"kind"`
	Status       string  `json:"status"`
	SourceID     string  `json:"sourceId"`
}
type Dataset struct {
	Version     string       `json:"version"`
	Timezone    string       `json:"timezone"`
	Occurrences []Occurrence `json:"occurrences"`
}
type Engine struct {
	dataset  Dataset
	holidays map[string]Occurrence
}

func Load(path string) (*Engine, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dataset Dataset
	if err := json.Unmarshal(raw, &dataset); err != nil {
		return nil, err
	}
	if dataset.Timezone != "Africa/Accra" {
		return nil, errors.New("timezone must be Africa/Accra")
	}
	holidays := map[string]Occurrence{}
	for _, occurrence := range dataset.Occurrences {
		if occurrence.Status == "confirmed" && occurrence.Date != nil {
			holidays[*occurrence.Date] = occurrence
		}
		if occurrence.Status == "confirmed" && occurrence.ObservedDate != nil {
			holidays[*occurrence.ObservedDate] = occurrence
		}
	}
	return &Engine{dataset: dataset, holidays: holidays}, nil
}
func parseDate(value string) (time.Time, error) { return time.Parse(DateLayout, value) }
func formatDate(value time.Time) string         { return value.Format(DateLayout) }
func (engine *Engine) Dataset() Dataset         { return engine.dataset }
func (engine *Engine) Holidays(year int, includePending bool) []Occurrence {
	result := []Occurrence{}
	for _, item := range engine.dataset.Occurrences {
		if item.Date == nil {
			if includePending && item.Status == "pending" && strings.HasSuffix(item.ID, fmt.Sprint(year)) {
				result = append(result, item)
			}
			continue
		}
		date, err := parseDate(*item.Date)
		if err == nil && date.Year() == year {
			result = append(result, item)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Date == nil {
			return false
		}
		if result[j].Date == nil {
			return true
		}
		return *result[i].Date < *result[j].Date
	})
	return result
}
func (engine *Engine) IsWorkingDay(value string) (bool, *Occurrence, error) {
	date, err := parseDate(value)
	if err != nil {
		return false, nil, err
	}
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return false, nil, nil
	}
	if item, ok := engine.holidays[value]; ok {
		return false, &item, nil
	}
	return true, nil, nil
}
func (engine *Engine) AddWorkingDays(value string, count int) (string, error) {
	date, err := parseDate(value)
	if err != nil {
		return "", err
	}
	step := 1
	if count < 0 {
		step = -1
		count = -count
	}
	for added := 0; added < count; {
		date = date.AddDate(0, 0, step)
		working, _, _ := engine.IsWorkingDay(formatDate(date))
		if working {
			added++
		}
	}
	return formatDate(date), nil
}
func (engine *Engine) NextWorkingDay(value string) (string, error) {
	return engine.AddWorkingDays(value, 1)
}
func (engine *Engine) PreviousWorkingDay(value string) (string, error) {
	return engine.AddWorkingDays(value, -1)
}
func (engine *Engine) WorkingDaysBetween(startValue, endValue string) (int, error) {
	start, err := parseDate(startValue)
	if err != nil {
		return 0, err
	}
	end, err := parseDate(endValue)
	if err != nil {
		return 0, err
	}
	if end.Before(start) {
		return 0, errors.New("end must not be before start")
	}
	count := 0
	for cursor := start; cursor.Before(end); cursor = cursor.AddDate(0, 0, 1) {
		working, _, _ := engine.IsWorkingDay(formatDate(cursor))
		if working {
			count++
		}
	}
	return count, nil
}
