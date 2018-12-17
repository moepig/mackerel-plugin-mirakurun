package mpmirakurun

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var statusServer *httptest.Server
var stub map[string]string

var jsonStr = map[string]string{
}

func TestGraphDefinition(t *testing.T) {
	var mirakurun MirakurunPlugin

	graphdef := mirakurun.GraphDefinition()
}

func TestMain(m *testing.M) {
	os.Exit(mainTest(m))
}

func mainTest(m *testing.M) int {
	flag.Parse()

	router := mux.NewRouter()

	return m.Run()
}

func TestFetchMetrics(t *testing.T) {
	// response a valid stats json
	stub = jsonStr

	// get metrics
	p := MirakurunPlugin {
		Target: strings.Replace(statusServer.URL, "http://", "", 1),
		Prefix: "mirakurun",
	}
	metrics, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("Failed to FetchMetrics: %s", err)
		return
	}

	// check the metrics1
	expected := map[string]float64{

	}

	for k, v := range expected {
		value, ok := metrics[k]
		if !ok {
			t.Errorf("metric of %s cannot be fetched", k)
			continue
		}
		if v != value {
			t.Errorf("metric of %s should be %v, but %v", k, v, value)
		}
	}
}

func TestFetchMetricsFail(t *testing.T) {
	p := MirakurunPlugin{
		Target: strings.Replace(statusServer.URL, "http://", "", 1),
		Prefix: "redash",
	}

	// return error against an invalid stats json
	stub = map[string]string{
		"status":   "{feature: [],}",
		"recorded": "[]",
	}
	_, err := p.FetchMetrics()
	if err == nil {
		t.Errorf("FetchMetrics should return error: stub=%v", stub)
	}
}
