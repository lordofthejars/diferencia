package exporter

import (
	"encoding/json"
	"net/http"
	"sync"
)

// URLCall contains the tuple Http Method Path
type URLCall struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// URLCounterMap is a type-safe and concurrent map storing for each URL the number of errors encountered
type URLCounterMap struct {
	sync.RWMutex
	internal map[URLCall]int
}

// Entry tuple for endpoint and number of errors
type Entry struct {
	Endpoint URLCall `json:"endpoint"`
	Errors   int     `json:"errors"`
}

// NewURLCounterMap creates a new instance of the map
func NewURLCounterMap() *URLCounterMap {
	return &URLCounterMap{
		internal: make(map[URLCall]int),
	}
}

// Inc by 1 given key or initialize to 1
func (m *URLCounterMap) Inc(method, path string) int {

	m.Lock()
	defer m.Unlock()
	call := URLCall{method, path}

	counter, ok := m.internal[call]
	newCounter := counter
	if ok {
		newCounter = counter + 1
		m.internal[call] = newCounter
	} else {
		newCounter = 1
		m.internal[call] = newCounter
	}

	return newCounter
}

// Get count for given method, path
func (m *URLCounterMap) Get(method, path string) (int, bool) {
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

// Entries returns a list of tuple key, value contained inside map
func (m *URLCounterMap) Entries() []Entry {
	m.RLock()
	defer m.RUnlock()

	var entries []Entry

	for key, value := range m.internal {
		entries = append(entries, Entry{Endpoint: key, Errors: value})
	}

	// To avoid null representation on JSON conversion
	if entries == nil {
		entries = make([]Entry, 0)
	}

	return entries
}

var stats = NewURLCounterMap()

// Reset Removes all
func (m *URLCounterMap) Reset() {
	for key, _ := range m.internal {
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

// IncrementError stats with a new error
func IncrementError(method, path string) int {
	return stats.Inc(method, path)
}

// StatsHandler to return JSON with stats
func StatsHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats.Entries())

}
