package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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
	BuildInfo string `json:"buildInfo"`
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`

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
	s.h.HandleFunc("/v1/user/performance", s.resultsHandler).Methods("POST").Headers(ContentTypeHeader, "application/json")
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
	messageMAC, err := hex.DecodeString(r.Header.Get(AuthorizationHeader))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read results, limited to 1mb
	buf, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	// validate
	if err = validate(buf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// unmarshal results
	var results PerformanceResults
	if err = json.Unmarshal(buf, &results); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse datetime
	sessionID := strings.SplitN(results.SessionID, "-", 2)
	if len(sessionID) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	datetime, err := time.Parse("2006.01.02-15.04.05", sessionID[1])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// check date within valid range (-/+ 10 minutes)
	min := time.Now().Add(-time.Minute * 10)
	max := time.Now().Add(time.Minute * 10)
	if !datetime.After(min) || !datetime.Before(max) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// compute and compare hmac
	mac := hmac.New(sha256.New, []byte(s.sharedSecret))
	mac.Write(buf)
	if !hmac.Equal(mac.Sum(nil), messageMAC) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// store additional metadata
	md := &Metadata{
		DateTime: time.Now().UTC(),
		IP:       r.RemoteAddr,
	}

	// write results
	if err = s.rw.AppendResults(md, &results); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) exampleHandler(w http.ResponseWriter, r *http.Request) {
	pr := &PerformanceResults{}
	pr.SessionID = "d1ac887243389d94544e4d9cc5524ab5-" + time.Now().UTC().Format("2006.01.02-15.04.05")
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

	fmt.Fprintf(w, "curl -v -H 'Content-Type: application/json' -H 'Authorization: %s' https://%s/v1/user/performance -d '%s'\n", hex.EncodeToString(mac.Sum(nil)), r.Host, string(body))
}
