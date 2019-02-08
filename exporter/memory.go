package exporter

import (
	"encoding/json"
	"math"
	"net/http"
	"sync"
	"time"
)

// URLCall contains the tuple Http Method Path
type URLCall struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// CallData contains the information that we want to store for the given URL
type CallData struct {
	Success                   int           `json:"success"`
	Errors                    int           `json:"errors"`
	ErrorDetails              []ErrorData   `json:"errorDetails"`
	PrimaryDurationAllCalls   time.Duration `json:"-"`
	CandidateDurationAllCalls time.Duration `json:"-"`
}

// ErrorData to hold all info when an error occurs
type ErrorData struct {
	FullURI         string      `json:"fullURI"`
	OriginalBody    string      `json:"originalBody,omitempty"`
	OriginalHeaders http.Header `json:"originalHeaders,omitempty"`
	HeaderDiff      string      `json:"headerDiff,omitempty"`
	BodyDiff        string      `json:"bodyDiff,omitempty"`
	StatusDiff      string      `json:"statusDiff,omitempty"`
}

// IncError increments the error counter
func (c *CallData) IncError() {
	c.Errors++
}

// IncSuccess increments the success counter
func (c *CallData) IncSuccess() {
	c.Success++
}

// IncAveragePrimaryTime increments the duration (elapsed time) of primary
func (c *CallData) IncAveragePrimaryTime(d time.Duration) {
	c.PrimaryDurationAllCalls += d
}

// IncAverageCandidateTime increments the duration (elapsed time) of candidate
func (c *CallData) IncAverageCandidateTime(d time.Duration) {
	c.CandidateDurationAllCalls += d
}

// AppendErrorData adds a new Error Data
func (c *CallData) AppendErrorData(errorData ErrorData) {
	c.ErrorDetails = append(c.ErrorDetails, errorData)
}

// URLCounterMap is a type-safe and concurrent map storing for each URL the number of errors encountered
type URLCounterMap struct {
	sync.RWMutex
	internal map[URLCall]CallData
}

// Entry tuple for endpoint and number of errors
type Entry struct {
	Endpoint                 URLCall     `json:"endpoint"`
	ErrorDetails             []ErrorData `json:"errorDetails"`
	Errors                   int         `json:"errors"`
	Success                  int         `json:"success"`
	AveragePrimaryDuration   float32     `json:"averagePrimaryDuration"`
	AverageCandidateDuration float32     `json:"averageCandidateDuration"`
}

// NewURLCounterMap creates a new instance of the map
func NewURLCounterMap() *URLCounterMap {
	return &URLCounterMap{
		internal: make(map[URLCall]CallData),
	}
}

// IncSuccess by 1 the success field and updates the average time
func (m *URLCounterMap) IncSuccess(method, path string, primaryAverage, candidateAverage time.Duration) int {
	m.Lock()
	defer m.Unlock()
	call := URLCall{method, path}

	counter, ok := m.internal[call]
	newCounter := counter
	if ok {
		counter.IncSuccess()
		counter.IncAveragePrimaryTime(primaryAverage)
		counter.IncAverageCandidateTime(candidateAverage)
		newCounter = counter
		m.internal[call] = newCounter
	} else {
		newCounter = CallData{PrimaryDurationAllCalls: primaryAverage, CandidateDurationAllCalls: candidateAverage, Success: 1}
		m.internal[call] = newCounter
	}

	return newCounter.Success
}

// IncErr by 1 the error field
func (m *URLCounterMap) IncErr(method, path string, errorData ErrorData) int {

	m.Lock()
	defer m.Unlock()
	call := URLCall{method, path}

	counter, ok := m.internal[call]
	newCounter := counter

	if ok {
		counter.IncError()
		counter.AppendErrorData(errorData)
		newCounter = counter
		m.internal[call] = newCounter
	} else {
		newCounter = CallData{Errors: 1, PrimaryDurationAllCalls: 0, CandidateDurationAllCalls: 0}
		newCounter.AppendErrorData(errorData)
		m.internal[call] = newCounter
	}

	return newCounter.Errors
}

// Get count for given method, path
func (m *URLCounterMap) Get(method, path string) (CallData, bool) {
	m.RLock()
	defer m.RUnlock()
	call := URLCall{method, path}
	result, ok := m.internal[call]

	return result, ok
}

// Keys returns the list of keys of map
func (m *URLCounterMap) Keys() []URLCall {
	m.RLock()
	defer m.RUnlock()

	var keys []URLCall
	for key := range m.internal {
		keys = append(keys, key)
	}

	return keys
}

// FindEntry finds an entry by method and path
func (m *URLCounterMap) FindEntry(method, path string) Entry {

	url := URLCall{method, path}
	result, ok := m.internal[url]

	if ok {
		return convert(url, result)
	}

	return Entry{}
}

// Entries returns a list of tuple key, value contained inside map
func (m *URLCounterMap) Entries() []Entry {
	m.RLock()
	defer m.RUnlock()

	var entries []Entry

	for key, value := range m.internal {
		entries = append(entries, convert(key, value))
	}

	// To avoid null representation on JSON conversion
	if entries == nil {
		entries = make([]Entry, 0)
	}

	return entries
}

func convert(key URLCall, value CallData) (e Entry) {

	primaryAverage := 0.0
	candidateAverage := 0.0

	if value.Success > 0 {
		primaryAverage = float64(value.PrimaryDurationAllCalls.Nanoseconds() / int64(value.Success))
		primaryAverage = (float64(primaryAverage) / float64(1000000))
	}

	if value.Success > 0 {
		candidateAverage = float64(value.CandidateDurationAllCalls.Nanoseconds() / int64(value.Success))
		candidateAverage = (float64(candidateAverage) / float64(1000000))
	}

	e = Entry{Endpoint: key, Errors: value.Errors, Success: value.Success,
		AveragePrimaryDuration:   float32(math.Round(primaryAverage*100) / 100),
		AverageCandidateDuration: float32(math.Round(candidateAverage*100) / 100),
		ErrorDetails:             value.ErrorDetails}

	return
}

var stats = NewURLCounterMap()

// Reset Removes all
func (m *URLCounterMap) Reset() {
	for key := range m.internal {
		delete(m.internal, key)
	}
}

// Reset Removes all
func Reset() {
	stats.Reset()
}

// Entries that are stored
func Entries() []Entry {
	return stats.Entries()
}

// FindEntry inside stats
func FindEntry(method, path string) Entry {
	return stats.FindEntry(method, path)
}

// IncrementSuccess stats with new success
func IncrementSuccess(method, path string, primaryAverage, candidateAverage time.Duration) int {
	return stats.IncSuccess(method, path, primaryAverage, candidateAverage)
}

// IncrementError stats with a new error
func IncrementError(method, path, body, uri, headersDiff, bodyDiff, stautsDiff string, headers http.Header) int {
	errorData := ErrorData{FullURI: uri, OriginalBody: body, OriginalHeaders: headers, HeaderDiff: headersDiff, BodyDiff: bodyDiff, StatusDiff: stautsDiff}

	return stats.IncErr(method, path, errorData)
}

// StatsHandler to return JSON with stats
func StatsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats.Entries())

}
