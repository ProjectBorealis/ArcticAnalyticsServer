package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	// AuthorizationHeader is the header containing the token.
	AuthorizationHeader = "Authorization"

	// ContentTypeHeader is the header containing the content type.
	ContentTypeHeader = "Content-Type"
)

// PerformanceResultsAttribute defines key/value attribute for the results.
type PerformanceResultsAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PerformanceResultsEvent defines an event with attributes.
type PerformanceResultsEvent struct {
	EventName  string                        `json:"eventName"`
	Attributes []PerformanceResultsAttribute `json:"attributes"`
}

// PerformanceResults are the results obtained from the performance test.
type PerformanceResults struct {
	DateTime  time.Time `json:"datetime"`
	BuildInfo string    `json:"buildInfo"`
	SessionID string    `json:"sessionId"`
	UserID    string    `json:"userId"`

	Events []PerformanceResultsEvent `json:"events"`
}

// Server with endpoints for registering performance results.
type Server struct {
	h  *mux.Router
	rw *ResultWriter

	sharedSecret string
}

// New returns a new Server.
func New(rw *ResultWriter, sharedSecret string) *Server {
	s := &Server{rw: rw, sharedSecret: sharedSecret}

	s.h = mux.NewRouter()
	s.h.HandleFunc("/", s.resultsHandler).Methods("POST").Headers(ContentTypeHeader, "application/json")
	s.h.HandleFunc("/example", s.exampleHandler).Methods("GET")

	return s
}

// Handler returns the http.Handler
func (s *Server) Handler() http.Handler {
	return s.h
}

func (s *Server) resultsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// decode hmac
	messageMAC, err := base64.StdEncoding.DecodeString(r.Header.Get(AuthorizationHeader))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read results, limited to 1mb
	// io.LimitReader(r.Body, 1<<20)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshal results
	var results PerformanceResults
	if err = json.Unmarshal(buf, &results); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check date within valid range (-/+ a minute)
	min := time.Now().Add(-time.Minute)
	max := time.Now().Add(time.Minute)
	if !results.DateTime.After(min) || !results.DateTime.Before(max) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// compute and compare hmac
	mac := hmac.New(sha256.New, []byte(s.sharedSecret))
	mac.Write(buf)
	if !hmac.Equal(mac.Sum(nil), messageMAC) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.rw.AppendResults(&results); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) exampleHandler(w http.ResponseWriter, r *http.Request) {
	pr := &PerformanceResults{}
	pr.DateTime = time.Now()
	pr.SessionID = "d1ac887243389d94544e4d9cc5524ab5-2018.10.03-00.56.26"
	pr.UserID = "d1ac887243389d94544e4d9cc5524ab5"
	pr.BuildInfo = "1.0.0.7"

	pr.Events = append(pr.Events, PerformanceResultsEvent{
		EventName: "Score",
		Attributes: []PerformanceResultsAttribute{
			{
				Name:  "Score.Num",
				Value: "0",
			},
		},
	})

	body, _ := json.Marshal(pr)

	mac := hmac.New(sha256.New, []byte(s.sharedSecret))
	mac.Write(body)

	fmt.Fprintf(w, "curl -v -H 'Content-Type: application/json' -H 'Authorization: %s' http://%s -d '%s'\n", base64.StdEncoding.EncodeToString(mac.Sum(nil)), r.Host, string(body))
}
