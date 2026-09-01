package calendar

import "testing"

func testEngine(t *testing.T) *Engine {
	t.Helper()
	engine, err := Load("../../data/holidays.json")
	if err != nil {
		t.Fatal(err)
	}
	return engine
}
func TestOfficial2026Fixtures(t *testing.T) {
	engine := testEngine(t)
	for _, date := range []string{"2026-01-01", "2026-04-03", "2026-07-03", "2026-12-28"} {
		working, _, err := engine.IsWorkingDay(date)
		if err != nil || working {
			t.Fatalf("expected %s to be non-working", date)
		}
	}
}
func TestWeekendAndHolidayDeterminism(t *testing.T) {
	got, err := testEngine(t).NextWorkingDay("2026-12-25")
	if err != nil {
		t.Fatal(err)
	}
	if got != "2026-12-29" {
		t.Fatalf("got %s", got)
	}
}
func TestHistoricalRulesAreNotRewritten(t *testing.T) {
	engine := testEngine(t)
	if len(engine.Holidays(2024, false)) == 0 || len(engine.Holidays(2025, false)) == 0 {
		t.Fatal("historical versions missing")
	}
	working, _, _ := engine.IsWorkingDay("2024-08-05")
	if working {
		t.Fatal("observed historical holiday missing")
	}
}
func TestPendingMovableHolidayIsNotGuessed(t *testing.T) {
	pending := 0
	for _, item := range testEngine(t).Holidays(2026, true) {
		if item.Status == "pending" {
			pending++
			if item.Date != nil {
				t.Fatal("invented pending date")
			}
		}
	}
	if pending != 3 {
		t.Fatalf("got %d pending", pending)
	}
}
func TestWorkingDayOperations(t *testing.T) {
	engine := testEngine(t)
	got, _ := engine.AddWorkingDays("2026-04-02", 2)
	if got != "2026-04-08" {
		t.Fatalf("got %s", got)
	}
	count, _ := engine.WorkingDaysBetween("2026-04-01", "2026-04-08")
	if count != 3 {
		t.Fatalf("got %d", count)
	}
}
