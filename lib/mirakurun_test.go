package mpmirakurun

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var statusServer *httptest.Server
var stub string

var jsonStr = `{
  "version": "2.7.0",
  "process": {
    "arch": "x64",
    "platform": "linux",
    "versions": {
      "http_parser": "2.8.0",
      "node": "8.11.3",
      "v8": "6.2.414.54",
      "uv": "1.19.1",
      "zlib": "1.2.11",
      "ares": "1.10.1-DEV",
      "modules": "57",
      "nghttp2": "1.32.0",
      "napi": "3",
      "openssl": "1.0.2o",
      "icu": "60.1",
      "unicode": "10.0",
      "cldr": "32.0",
      "tz": "2017c"
    },
    "pid": 5893,
    "memoryUsage": {
      "rss": 171012096,
      "heapTotal": 84508672,
      "heapUsed": 54182232,
      "external": 35296339
    }
  },
  "epg": {
    "gatheringNetworks": [],
    "storedEvents": 8416
  },
  "streamCount": {
    "tunerDevice": 0,
    "tsFilter": 0,
    "decoder": 0
  },
  "errorCount": {
    "uncaughtException": 1,
    "bufferOverflow": 0,
    "tunerDeviceRespawn": 18985
  },
  "timerAccuracy": {
    "last": 313.08,
    "m1": {
      "avg": 806.59415,
      "min": 270.528,
      "max": 1365.336
    },
    "m5": {
      "avg": 662.05513,
      "min": -876.187,
      "max": 1376.414
    },
    "m15": {
      "avg": 586.0672522222222,
      "min": -1207.258,
      "max": 2405.415
    }
  }
}`

func TestGraphDefinition(t *testing.T) {
	var mirakurun MirakurunPlugin

	graphdef := mirakurun.GraphDefinition()
}

func TestMain(m *testing.M) {
	os.Exit(mainTest(m))
}

func mainTest(m *testing.M) int {
	flag.Parse()

	statusServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, stub)
		}))

	return m.Run()
}

func TestFetchMetrics(t *testing.T) {
	// response a valid stats json
	stub = jsonStr

	// get metrics
	p := MirakurunPlugin{
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
	stub = "{}"

	_, err := p.FetchMetrics()
	if err == nil {
		t.Errorf("FetchMetrics should return error: stub=%v", stub)
	}
}
