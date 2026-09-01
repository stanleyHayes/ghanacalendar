package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	calendar "github.com/stanleyHayes/ghanacalendar/internal/calendar"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type server struct{ engine *calendar.Engine }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://calendar.digitalghana.dev")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func errorJSON(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func (s server) holidays(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil || year < 2024 || year > 2026 {
		errorJSON(w, 400, "year must be 2024, 2025, or 2026")
		return
	}
	items := s.engine.Holidays(year, r.URL.Query().Get("includePending") == "true")
	format := r.URL.Query().Get("format")
	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		out := csv.NewWriter(w)
		_ = out.Write([]string{"id", "name", "date", "observedDate", "kind", "status", "sourceId"})
		for _, item := range items {
			date, observed := "", ""
			if item.Date != nil {
				date = *item.Date
			}
			if item.ObservedDate != nil {
				observed = *item.ObservedDate
			}
			_ = out.Write([]string{item.ID, item.Name, date, observed, item.Kind, item.Status, item.SourceID})
		}
		out.Flush()
		return
	}
	if format == "ics" {
		w.Header().Set("Content-Type", "text/calendar")
		fmt.Fprintln(w, "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//Digital Ghana//GhanaCalendar//EN")
		for _, item := range items {
			if item.Date == nil {
				continue
			}
			date := strings.ReplaceAll(*item.Date, "-", "")
			fmt.Fprintf(w, "BEGIN:VEVENT\r\nUID:%s@calendar.digitalghana.dev\r\nDTSTART;VALUE=DATE:%s\r\nSUMMARY:%s\r\nEND:VEVENT\r\n", item.ID, date, item.Name)
		}
		fmt.Fprintln(w, "END:VCALENDAR")
		return
	}
	writeJSON(w, 200, map[string]any{"version": s.engine.Dataset().Version, "timezone": "Africa/Accra", "holidays": items})
}
func (s server) workingDay(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	working, holiday, err := s.engine.IsWorkingDay(date)
	if err != nil {
		errorJSON(w, 400, "date must use YYYY-MM-DD")
		return
	}
	writeJSON(w, 200, map[string]any{"date": date, "timezone": "Africa/Accra", "workingDay": working, "holiday": holiday})
}
func (s server) shift(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || count < -3660 || count > 3660 {
		errorJSON(w, 400, "count must be between -3660 and 3660")
		return
	}
	result, err := s.engine.AddWorkingDays(date, count)
	if err != nil {
		errorJSON(w, 400, "date must use YYYY-MM-DD")
		return
	}
	writeJSON(w, 200, map[string]any{"date": date, "count": count, "result": result, "timezone": "Africa/Accra"})
}

func (s server) adjacent(w http.ResponseWriter, r *http.Request, direction int) {
	date := r.URL.Query().Get("date")
	result, err := s.engine.AddWorkingDays(date, direction)
	if err != nil {
		errorJSON(w, 400, "date must use YYYY-MM-DD")
		return
	}
	writeJSON(w, 200, map[string]any{"date": date, "result": result, "timezone": "Africa/Accra"})
}

func (s server) graphql(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10)).Decode(&payload); err != nil {
		errorJSON(w, 400, "invalid GraphQL request")
		return
	}
	date, _ := payload.Variables["date"].(string)
	if strings.Contains(payload.Query, "isWorkingDay") {
		working, holiday, err := s.engine.IsWorkingDay(date)
		if err != nil {
			errorJSON(w, 400, "date must use YYYY-MM-DD")
			return
		}
		writeJSON(w, 200, map[string]any{"data": map[string]any{"isWorkingDay": map[string]any{"date": date, "workingDay": working, "holiday": holiday, "timezone": "Africa/Accra"}}})
		return
	}
	if strings.Contains(payload.Query, "holidays") {
		yearValue, ok := payload.Variables["year"].(float64)
		if !ok || yearValue < 2024 || yearValue > 2026 {
			errorJSON(w, 400, "year must be 2024, 2025, or 2026")
			return
		}
		writeJSON(w, 200, map[string]any{"data": map[string]any{"holidays": s.engine.Holidays(int(yearValue), true)}})
		return
	}
	errorJSON(w, 400, "supported operations are holidays and isWorkingDay")
}
func (s server) between(w http.ResponseWriter, r *http.Request) {
	count, err := s.engine.WorkingDaysBetween(r.URL.Query().Get("start"), r.URL.Query().Get("end"))
	if err != nil {
		errorJSON(w, 400, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"start": r.URL.Query().Get("start"), "end": r.URL.Query().Get("end"), "workingDays": count, "endExclusive": true})
}
func main() {
	engine, err := calendar.Load("data/holidays.json")
	if err != nil {
		log.Fatal(err)
	}
	app := server{engine}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "version": engine.Dataset().Version})
	})
	mux.HandleFunc("GET /v1/holidays", app.holidays)
	mux.HandleFunc("GET /v1/working-day", app.workingDay)
	mux.HandleFunc("GET /v1/add-working-days", app.shift)
	mux.HandleFunc("GET /v1/next-working-day", func(w http.ResponseWriter, r *http.Request) { app.adjacent(w, r, 1) })
	mux.HandleFunc("GET /v1/previous-working-day", func(w http.ResponseWriter, r *http.Request) { app.adjacent(w, r, -1) })
	mux.HandleFunc("GET /v1/working-days-between", app.between)
	mux.HandleFunc("POST /graphql", app.graphql)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
