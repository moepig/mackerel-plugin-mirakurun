package mpmirakurun

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type MirakurunPlugin struct {
	Prefix   string
	Target   string
	Tempfile string
}

type Status struct {
	Process *Process `json:"process"`
}

type Process struct {
	MemoryUsage *MemoryUsage `json:"memoryUsage"`
}

type MemoryUsage struct {
	Rss       *int `json:"rss"`
	HeapTotal *int `json:"heapTotal"`
	HeapUsed  *int `json:"heapUsed"`
	External  *int `json:"external"`
}

var graphdef = map[string]mp.Graphs{
	"memoryUsage": mp.Graphs{
		Label: "Mirakurun Memory Usage",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "rss", Label: "RSS", Diff: false},
			{Name: "heapTotal", Label: "Heap Total", Diff: false},
			{Name: "heapUsed", Label: "Heap Used", Diff: false},
		},
	},
}

// FetchMetrics interface for mackerelplugin
func (m MirakurunPlugin) FetchMetrics() (map[string]float64, error) {
	// call status api
	url := fmt.Sprintf("http://%s/api/status", m.Target)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	byteArray, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var status Status
	err = json.Unmarshal(byteArray, &status)
	if err != nil {
		return nil, err
	}

	// mapping to metrics
	stat := make(map[string]float64)

	if status.Process != nil {
		if status.Process.MemoryUsage != nil {
			stat["rss"] = float64(*status.Process.MemoryUsage.Rss)
			stat["heapTotal"] = float64(*status.Process.MemoryUsage.HeapTotal)
			stat["heapUsed"] = float64(*status.Process.MemoryUsage.HeapUsed)
		}
	}
	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m MirakurunPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// MetricKeyPrefix interface for mackerelplugin
func (m MirakurunPlugin) MetricKeyPrefix() string {
	if m.Prefix == "" {
		m.Prefix = "mirakurun"
	}
	return m.Prefix
}

func Do() {
	optPrefix := flag.String("metric-key-prefix", "mirakurun", "Metric key prefix")
	optHost := flag.String("host", "", "mirakurun hostname")
	optPort := flag.String("port", "", "mirakurun port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin MirakurunPlugin

	plugin.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	plugin.Prefix = *optPrefix

	helper := mp.NewMackerelPlugin(plugin)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/.mackerel-plugin-mirakurun"
	}

	helper.Run()
}
