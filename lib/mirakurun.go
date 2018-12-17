package mpmirakurun

import (
	"flag"
	"fmt"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type MirakurunPlugin struct {
	Prefix   string
	Target   string
	Tempfile string
}

var graphdef = map[string]mp.Graphs{}

// FetchMetrics interface for mackerelplugin
func (m MirakurunPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

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
