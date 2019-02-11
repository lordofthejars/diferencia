package core

import (
	"bufio"
	"bytes"
	jsonenc "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lordofthejars/diferencia/difference/header"
	"github.com/lordofthejars/diferencia/difference/plain"

	"github.com/lordofthejars/diferencia/difference/json"
	"github.com/lordofthejars/diferencia/exporter"
	"github.com/lordofthejars/diferencia/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sirupsen/logrus"
)

// Difference algorithm
type Difference int

func (difference Difference) String() string {
	names := [...]string{
		"Strict",
		"Subset",
		"Schema"}

	if difference < Strict || difference > Schema {
		return "Unknown"
	}
	return names[difference]
}

//NewDifference creator from String
func NewDifference(difference string) (Difference, error) {

	switch difference {
	case "Strict":
		return Strict, nil
	case "Subset":
		return Subset, nil
	case "Schema":
		return Schema, nil
	}

	return -1, fmt.Errorf("Cannot find %s difference mode", difference)
}

// Config object
var Config *DiferenciaConfiguration

// HttpClient interface to make requests with changed URL
var HttpClient Client = &HTTPClient{
	config: Config,
}

var prometheusCounter *prometheus.CounterVec

const (
	// Strict mode everything should be exactly the same
	Strict Difference = 0
	// Subset mode that the candidate must be a subset of primary
	Subset Difference = 1
	// Schema mode where the schema must be equal but not the values
	Schema Difference = 2
)

// DiferenciaConfiguration object
type DiferenciaConfiguration struct {
	Port                  int        `json:"port,omitempty"`
	ServiceName           string     `json:"serviceName,omitempty"`
	Primary               string     `json:"primary,omitempty"`
	Secondary             string     `json:"secondary,omitempty"`
	Candidate             string     `json:"candidate,omitempty"`
	StoreResults          string     `json:"storeResults,omitempty"`
	DifferenceMode        Difference `json:"-"`
	NoiseDetection        bool       `json:"noiseDetection,omitempty"`
	AllowUnsafeOperations bool       `json:"allowUnsafeOperartions,omitempty"`
	Prometheus            bool       `json:"prometheus,omitempty"`
	PrometheusPort        int        `json:"prometheusPort,omitempty"`
	Headers               bool       `json:"headers,omitempty"`
	IgnoreHeadersValues   []string   `json:"ignoreHeadersValues,omitempty"`
	IgnoreValues          []string   `json:"ignoreValues,omitempty"`
	IgnoreValuesFile      string     `json:"ignoreValuesFile,omitempty"`
	InsecureSkipVerify    bool       `json:"insecureSkipVerify,omitempty"`
	CaCert                string     `json:"caCert,omitempty"`
	ClientCert            string     `json:"clientCert,omitempty"`
	ClientKey             string     `json:"clientKey,omitempty"`
	AdminPort             int        `json:"adminPort,omitempty"`
	ForcePlainText        bool       `json:"forcePlainText,omitempty"`
	LevenshteinPercentage int        `json:"levenshteinPercentage,omitempty"`
	Mirroring             bool       `json:"mirroring,omitempty"`
	ReturnResult          bool       `json:"returnResult,omitempty"`
}

// UpdateConfiguration with configured params
func (conf *DiferenciaConfiguration) UpdateConfiguration(updateConfig DiferenciaConfigurationUpdate) error {

	if updateConfig.isReturnResultSet() {
		returnResult, err := updateConfig.getReturnResult()
		if err != nil {
			return err
		}
		conf.ReturnResult = returnResult
	}

	if updateConfig.isServiceNameSet() {
		conf.SetServiceName(updateConfig.ServiceName)
		prometheusCounter = metrics.RegisterNumberOfRegressions(Config.ServiceName)
	}

	if updateConfig.isPrimarySet() {
		conf.Primary = updateConfig.Primary
	}

	if updateConfig.isSecondarySet() {
		conf.Secondary = updateConfig.Secondary
	}

	if updateConfig.isCandidateSet() {
		conf.Candidate = updateConfig.Candidate
		// Updates service name for new candidate in case of service name not set
		if !updateConfig.isServiceNameSet() {
			conf.SetServiceName(updateConfig.ServiceName)
			prometheusCounter = metrics.RegisterNumberOfRegressions(Config.ServiceName)
		}
	}

	if updateConfig.isModeSet() {
		mode, err := updateConfig.getMode()

		if err != nil {
			return err
		}

		conf.DifferenceMode = mode
	}

	if updateConfig.isNoiseDetectionSet() {
		noise, err := updateConfig.getNoiseDetection()

		if err != nil {
			return err
		}

		conf.NoiseDetection = noise

	}

	return nil
}

func (conf DiferenciaConfiguration) GetServiceName() string {
	return conf.ServiceName
}

// SetServiceName and case of empty it calculates from candidate
func (conf *DiferenciaConfiguration) SetServiceName(serviceName string) {

	if len(serviceName) == 0 {
		candidateURL, _ := url.Parse(conf.Candidate)
		conf.ServiceName = candidateURL.Hostname()
	}

}

// IsStoreResultsSet in configuration object
func (conf DiferenciaConfiguration) IsStoreResultsSet() bool {
	return len(conf.StoreResults) > 0
}

// IsIgnoreValuesSet in configuration object
func (conf DiferenciaConfiguration) IsIgnoreValuesSet() bool {
	return conf.IgnoreValues != nil && len(conf.IgnoreValues) > 0
}

// IsIgnoreValuesFileSet in configuration object
func (conf DiferenciaConfiguration) IsIgnoreValuesFileSet() bool {
	return len(conf.IgnoreValuesFile) > 0
}

// AreHttpsClientParamsSet checking if https is enabled
func (conf DiferenciaConfiguration) AreHttpsClientParamsSet() bool {
	return (len(conf.CaCert) > 0 && len(conf.ClientCert) > 0 && len(conf.ClientKey) > 0)
}

// Print configuration
func (conf DiferenciaConfiguration) Print() {
	fmt.Printf("Port: %d\n", conf.Port)
	fmt.Printf("Prometheus Port: %d\n", conf.PrometheusPort)
	fmt.Printf("Admin Port %d\n", conf.AdminPort)
	fmt.Printf("Service Name: %s\n", conf.ServiceName)
	fmt.Printf("Primary: %s\n", conf.Primary)
	fmt.Printf("Secondary: %s\n", conf.Secondary)
	fmt.Printf("Candidate: %s\n", conf.Candidate)
	fmt.Printf("Difference Mode: %s\n", conf.DifferenceMode.String())
	fmt.Printf("Noise Detection: %t\n", conf.NoiseDetection)
	fmt.Printf("Store Results: %s\n", conf.StoreResults)
	fmt.Printf("Ignore Values of: %v\n", conf.IgnoreValues)
	fmt.Printf("Ignore Values File: %s\n", conf.IgnoreValuesFile)
	fmt.Printf("Headers: %t\n", conf.Headers)
	fmt.Printf("Ignored Headers Values of: %v\n", conf.IgnoreHeadersValues)
	fmt.Printf("Allow Unsafe Operations: %t\n", conf.AllowUnsafeOperations)
	fmt.Printf("Insecure Skip Verify Port: %t\n", conf.InsecureSkipVerify)
	fmt.Printf("Ca Cert Path: %s\n", conf.CaCert)
	fmt.Printf("Client Cert Path: %s\n", conf.ClientCert)
	fmt.Printf("Client Key Path: %s\n", conf.ClientKey)
	fmt.Printf("Prometheus Enabled: %t\n", conf.Prometheus)
	fmt.Printf("Levenshtein Percentage: %d\n", conf.LevenshteinPercentage)
	fmt.Printf("Force Plain Text: %t\n", conf.ForcePlainText)
	fmt.Printf("Mirroring: %t\n", conf.Mirroring)
	fmt.Printf("Return Result: %t\n", conf.ReturnResult)
}

type DiferenciaError struct {
	code    int
	message string
}

func (e *DiferenciaError) Error() string {
	return fmt.Sprintf("with message: %s, and code %d", e.message, e.code)
}

type Communicationcontent struct {
	Content    []byte
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
}

func (c Communicationcontent) isEmpty() bool {
	return len(c.Content) == 0 && c.StatusCode == 0 && c.Header == nil && len(c.Cookies) == 0
}

// Result struct
type Result struct {
	EqualContent         bool
	PrimaryElapsedTime   time.Duration
	CandidateElapsedTime time.Duration
	Diff                 DifferenceDescription
}

// DifferenceDescription offers the description of the differences
type DifferenceDescription struct {
	HeadersDiff string `json:"headersDiff,omitempty"`
	BodyDiff    string `json:"bodyDiff,omitempty"`
	StatusDiff  string `json:"statusDiff,omitempty"`
}

// MarshallJson translate object to byte[]
func (r Result) MarshallJson() ([]byte, error) {
	return jsonenc.Marshal(struct {
		Result                   bool
		PrimaryElapsedTimeNano   int64
		CandidateElapsedTimeNano int64
		Description              *DifferenceDescription `json:"description,omitempty"`
	}{
		Result:                   r.EqualContent,
		PrimaryElapsedTimeNano:   r.PrimaryElapsedTime.Nanoseconds(),
		CandidateElapsedTimeNano: r.CandidateElapsedTime.Nanoseconds(),
		Description:              &r.Diff,
	})
}

func Diferencia(r *http.Request) (Result, Communicationcontent, error) {

	if !Config.AllowUnsafeOperations && !isSafeOperation(r.Method) {
		if !Config.Mirroring {
			logrus.Debugf("Unsafe operations are not allowed and %s method has been received", r.Method)
			return Result{EqualContent: false}, Communicationcontent{}, &DiferenciaError{http.StatusMethodNotAllowed, fmt.Sprintf("Unsafe operations are not allowed and %s method has been received", r.Method)}
		} else {
			// Do the request and return as ok.
		}
	}

	logrus.Debugf("URL %s is going to be processed", r.URL.String())

	// TODO it can be parallelized
	// Get request from primary
	primaryFullURL := CreateUrl(*r.URL, Config.Primary)
	logrus.Debugf("Forwarding call to %s", primaryFullURL)
	primaryStartTime := time.Now()
	primaryBodyContent, primaryStatus, primaryHeader, cookies, err := getContent(r, primaryFullURL)
	primaryElapsedDuration := time.Now().Sub(primaryStartTime)
	if err != nil {
		logrus.Errorf("Error while connecting to Primary site (%s) with %s", primaryFullURL, err.Error())
		return Result{EqualContent: false}, Communicationcontent{Content: primaryBodyContent, StatusCode: primaryStatus, Header: primaryHeader, Cookies: cookies}, &DiferenciaError{http.StatusServiceUnavailable, fmt.Sprintf("Error while connecting to Primary site (%s) with %s", primaryFullURL, err.Error())}
	}

	// Get candidate
	candidateFullURL := CreateUrl(*r.URL, Config.Candidate)
	logrus.Debugf("Forwarding call to %s", candidateFullURL)
	candidateStartTime := time.Now()
	candidateBodyContent, candidateStatus, candidateHeader, _, err := getContent(r, candidateFullURL)
	candidateElapsedDuration := time.Now().Sub(candidateStartTime)
	if err != nil {
		logrus.Errorf("Error while connecting to Candidate site (%s) with %s", candidateFullURL, err.Error())
		return Result{EqualContent: false}, Communicationcontent{Content: primaryBodyContent, StatusCode: primaryStatus, Header: primaryHeader, Cookies: cookies}, &DiferenciaError{http.StatusServiceUnavailable, fmt.Sprintf("Error while connecting to Candidate site (%s) with %s", candidateFullURL, err.Error())}
	}

	var result bool

	var secondaryFullURL string
	var secondaryBodyContent []byte
	var secondaryStatus int
	if Config.NoiseDetection {
		// Get secondary to do the noise cancellation
		secondaryFullURL := CreateUrl(*r.URL, Config.Secondary)
		logrus.Debugf("Forwarding call to %s", secondaryFullURL)
		secondaryBodyContent, secondaryStatus, _, _, err := getContent(r, secondaryFullURL)
		if err != nil {
			logrus.Errorf("Error while connecting to Secondary site (%s) with error %s", candidateFullURL, err.Error())
			return Result{EqualContent: false}, Communicationcontent{Content: primaryBodyContent, StatusCode: primaryStatus, Header: primaryHeader, Cookies: cookies}, &DiferenciaError{http.StatusServiceUnavailable, fmt.Sprintf("Error while connecting to Secondary site (%s) with error %s", candidateFullURL, err.Error())}
		}
		// If status code is equal then we detect noise and and remove from primary and candidate
		// What to do in case of two identical status code but no body content (404) might be still valid since you are testing that nothing is there
		if primaryStatus == secondaryStatus {

			contentType := primaryHeader.Get("Content-Type")
			var err error
			switch {
			case strings.HasPrefix(contentType, "application/json"):
				primaryBodyContent, candidateBodyContent, err = noiseCancellationJson(primaryBodyContent, secondaryBodyContent, candidateBodyContent)
			case strings.HasPrefix(contentType, "text/plain"):
				primaryBodyContent, candidateBodyContent = noiseCancellationText(primaryBodyContent, secondaryBodyContent, candidateBodyContent)
			default:
				{
					if Config.ForcePlainText {
						primaryBodyContent, candidateBodyContent = noiseCancellationText(primaryBodyContent, secondaryBodyContent, candidateBodyContent)
					} else {
						primaryBodyContent, candidateBodyContent, err = noiseCancellationJson(primaryBodyContent, secondaryBodyContent, candidateBodyContent)
					}
				}
			}

			if err != nil {
				logrus.WithError(err).Errorf("Error detecting noise between %s and %s.", primaryFullURL, secondaryFullURL)
				return Result{EqualContent: false}, Communicationcontent{Content: primaryBodyContent, StatusCode: primaryStatus, Header: primaryHeader, Cookies: cookies}, &DiferenciaError{http.StatusBadRequest, fmt.Sprintf("Error detecting noise between %s and %s. (%s)", primaryFullURL, secondaryFullURL, err.Error())}
			}

		} else {
			logrus.Errorf("Status code between %s(%d) and %s(%d) are different", primaryFullURL, primaryStatus, secondaryFullURL, secondaryStatus)
			return Result{EqualContent: false}, Communicationcontent{Content: primaryBodyContent, StatusCode: primaryStatus, Header: primaryHeader, Cookies: cookies}, &DiferenciaError{http.StatusBadRequest, fmt.Sprintf("Status code between %s(%d) and %s(%d) are different", primaryFullURL, primaryStatus, secondaryFullURL, secondaryStatus)}
		}
	}

	result, output := compareResult(candidateBodyContent, primaryBodyContent, candidateStatus, primaryStatus, candidateHeader, primaryHeader)

	if Config.IsStoreResultsSet() {
		primary := exporter.CreateInteraction(primaryFullURL, primaryBodyContent, primaryStatus)
		candidate := exporter.CreateInteraction(candidateFullURL, candidateBodyContent, candidateStatus)
		var secondary exporter.Interaction

		if Config.NoiseDetection {
			secondary = exporter.CreateInteraction(secondaryFullURL, secondaryBodyContent, secondaryStatus)
		}

		interactions := exporter.CreateInteractions(primary, &secondary, candidate, Config.DifferenceMode.String(), result)

		exporter.ExportToFile(Config.StoreResults, interactions)
	}

	logrus.Debugf("Result of comparing %s and %s is %t", primaryFullURL, candidateFullURL, result)

	// If it is a failure let's print the contents.
	if !result {
		logrus.Debugf("************************")
		logrus.Debugf("Explanation of Failure:")
		logrus.Debugf("Primary Status Code: %d Candidate StatusCode: %d", primaryStatus, candidateStatus)
		logrus.Debugf("Primary time: %s Candidate Time: %s", primaryElapsedDuration, candidateElapsedDuration)
		logrus.Debugf("Primary Content:")
		logrus.Debugf(string(primaryBodyContent[:]))
		logrus.Debugf("Candidate Content:")
		logrus.Debugf(string(candidateBodyContent[:]))
		if Config.Headers {
			logrus.Debugf("Primary Headers:")
			logrus.Debugf(createKeyValuePairs(primaryHeader))
			logrus.Debugf("Candidate Headers:")
			logrus.Debugf(createKeyValuePairs(candidateHeader))
		}
		logrus.Debugf("************************")
	}

	return Result{EqualContent: result, PrimaryElapsedTime: primaryElapsedDuration, CandidateElapsedTime: candidateElapsedDuration, Diff: output}, Communicationcontent{Content: primaryBodyContent, StatusCode: primaryStatus, Header: primaryHeader, Cookies: cookies}, nil

}

func createKeyValuePairs(m http.Header) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func noiseCancellationText(primaryBodyContent, secondaryBodyContent, candidateBodyContent []byte) ([]byte, []byte) {

	noiseOperation := plain.NoiseOperation{}
	noiseOperation.Detect(primaryBodyContent, secondaryBodyContent)

	primaryWithoutNoise, candidateWithoutNoise := noiseOperation.Remove(primaryBodyContent, candidateBodyContent)

	return primaryWithoutNoise, candidateWithoutNoise

}

func noiseCancellationJson(primaryBodyContent, secondaryBodyContent, candidateBodyContent []byte) ([]byte, []byte, error) {
	noiseOperation := json.NoiseOperation{}
	manualNoise := manualNoiseDetection()
	noiseOperation.Initialize(manualNoise)
	err := noiseOperation.Detect(primaryBodyContent, secondaryBodyContent)
	if err != nil {
		return nil, nil, err
	}
	primaryWithoutNoise, candidateWithoutNoise, _ := noiseOperation.Remove(primaryBodyContent, candidateBodyContent)

	return primaryWithoutNoise, candidateWithoutNoise, nil
}

func manualNoiseDetection() []string {
	var pointers []string

	if Config.IsIgnoreValuesSet() {
		for _, v := range Config.IgnoreValues {
			pointers = append(pointers, v)
		}
	}

	if Config.IsIgnoreValuesFileSet() {

		lines, err := readLines(Config.IgnoreValuesFile)

		if err != nil {
			logrus.Errorf("Error reading %s that defines ignoring values. %s. Execution will continue ignoring this file.", Config.IgnoreValuesFile, err)
			return pointers
		}

		for _, v := range lines {
			pointers = append(pointers, string(v))
		}

	}

	return pointers
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func compareResult(candidate, primary []byte, candidateStatus, primaryStatus int, candidateHeader, primaryHeader http.Header) (bool, DifferenceDescription) {

	// TODO This method should be refactored to a chain of responsibility pattern
	if primaryStatus == candidateStatus {
		headersDiff := ""
		headerEqual := true
		if Config.Headers {
			headerEqual, headersDiff = header.CompareHeaders(candidateHeader, primaryHeader, Config.IgnoreHeadersValues...)
		}
		// Comparision between documents
		contentType := primaryHeader.Get("Content-Type")
		switch {
		case strings.HasPrefix(contentType, "application/json"):
			bodyEqual, bodyDiff := json.CompareDocuments(candidate, primary, Config.DifferenceMode.String())

			if headerEqual && bodyEqual {
				return bodyEqual, DifferenceDescription{}
			}

			return bodyEqual && headerEqual, DifferenceDescription{HeadersDiff: headersDiff, BodyDiff: bodyDiff}
		case strings.HasPrefix(contentType, "text/plain"):
			return compareText(candidate, primary, Config.LevenshteinPercentage), DifferenceDescription{}
		default:
			{
				if Config.ForcePlainText {
					return compareText(candidate, primary, Config.LevenshteinPercentage), DifferenceDescription{}
				}
				bodyEqual, bodyDiff := json.CompareDocuments(candidate, primary, Config.DifferenceMode.String())

				if headerEqual && bodyEqual {
					return bodyEqual, DifferenceDescription{}
				}

				return bodyEqual && headerEqual, DifferenceDescription{HeadersDiff: headersDiff, BodyDiff: bodyDiff}
			}
		}
	}

	return false, DifferenceDescription{StatusDiff: fmt.Sprintf(`"status": %d => %d`, primaryStatus, candidateStatus)}
}

func compareText(candidate, primary []byte, levenshtein int) bool {
	if levenshtein < 100 {
		dif := int(plain.CalculateSimilarity(primary, candidate) * 100)
		return dif > levenshtein
	}

	return bytes.Equal(candidate, primary)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// If this handler is up and running means that Proxy can start dealing with requests
	w.WriteHeader(http.StatusOK)
}

func diferenciaHandler(w http.ResponseWriter, r *http.Request) {

	mutex.Lock()
	defer mutex.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Error reading body: %v", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	result, primaryCommunication, err := Diferencia(r)
	if err != nil {
		if de, ok := err.(*DiferenciaError); ok {
			w.WriteHeader(de.code)
			fmt.Fprintf(w, de.message)
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if result.EqualContent {
		if Config.Mirroring {
			MirrorResponse(primaryCommunication, w)
		} else {
			if Config.ReturnResult {
				content, _ := result.MarshallJson()
				w.Write(content)
			}
			w.WriteHeader(http.StatusOK)
		}
		exporter.IncrementSuccess(r.Method, r.URL.Path, result.PrimaryElapsedTime, result.CandidateElapsedTime)
	} else {
		// If there is a regression
		if Config.Mirroring {
			MirrorResponse(primaryCommunication, w)
		} else {
			w.WriteHeader(http.StatusPreconditionFailed)
			if Config.ReturnResult {
				content, _ := result.MarshallJson()
				w.Write(content)
			}
		}
		if Config.Prometheus {
			prometheusCounter.WithLabelValues(r.Method, r.URL.Path).Inc()
		}
		exporter.IncrementError(r.Method, r.URL.Path, string(body[:]), r.URL.RequestURI(), result.Diff.HeadersDiff, result.Diff.BodyDiff, result.Diff.StatusDiff, r.Header)
	}
}

func isSafeOperation(method string) bool {
	return method == http.MethodGet || method == http.MethodOptions || method == http.MethodHead
}

func getContent(r *http.Request, url string) ([]byte, int, http.Header, []*http.Cookie, error) {

	newRequest := duplicate(r)
	resp, err := HttpClient.MakeRequest(newRequest, url)

	if err != nil {
		// In case of error in service we should add as metrics as well or assume that the service itself would communicate to metrics?
		return make([]byte, 0), 0, nil, nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return bodyBytes, resp.StatusCode, resp.Header, resp.Cookies(), err

}

// StartProxy server
func StartProxy(configuration *DiferenciaConfiguration) {

	finish := make(chan bool)

	Config = configuration
	initialize()

	go func() {
		// Initialize Proxy server
		proxyMux := http.NewServeMux()
		// Matches everything
		proxyMux.HandleFunc("/", diferenciaHandler)
		proxyMux.HandleFunc("/healthdif", healthHandler)
		logrus.Errorf("Error starting proxy: %s", http.ListenAndServe(":"+strconv.Itoa(Config.Port), proxyMux))
	}()

	go func() {
		if Config.Prometheus {
			//Initialize Prometheus endpoint
			prometheusMux := http.NewServeMux()
			prometheusMux.Handle("/metrics", prometheus.Handler())
			logrus.Errorf("Error starting prometheus endpoint: %s", http.ListenAndServe(":"+strconv.Itoa(Config.PrometheusPort), prometheusMux))
		}
	}()

	go func() {
		// Initialize Admin server
		adminMux := http.NewServeMux()
		adminMux.HandleFunc("/configuration", adminHandler)
		adminMux.HandleFunc("/stats", exporter.StatsHandler)
		adminMux.HandleFunc("/dashboard/details", dashboardDetailsHandler)
		adminMux.HandleFunc("/dashboard/", dashboardHandler)
		logrus.Errorf("Error starting admin: %s", http.ListenAndServe(":"+strconv.Itoa(Config.AdminPort), adminMux))
	}()

	<-finish
}

func initialize() {

	// Print config object
	Config.Print()

	//Initialize Prometheus if required
	if Config.Prometheus {
		prometheusCounter = metrics.RegisterNumberOfRegressions(Config.ServiceName)
	}

}
