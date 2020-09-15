package exporter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/percona/exporter_shared/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

func collect(c helpers.Collector) []prometheus.Metric {
	m := make([]prometheus.Metric, 0)
	ch := make(chan prometheus.Metric)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for metric := range ch {
			m = append(m, metric)
		}
		wg.Done()
	}()

	c.Collect(ch)
	close(ch)

	wg.Wait()

	return m
}

func filterMetrics(metrics []*helpers.Metric) []*helpers.Metric {
	res := make([]*helpers.Metric, 0, len(metrics))

	for _, m := range metrics {
		m.Value = 0
	}

	return res
}

func getMetricNames(lines []string) map[string]bool {
	names := map[string]bool{}

	for _, line := range lines {
		if strings.HasPrefix(line, "# TYPE ") {
			m := strings.Split(line, " ")
			names[m[2]] = true
		}
	}

	return names
}

func readTestMetrics(filename string) ([]*helpers.Metric, error) {
	m := []*helpers.Metric{}

	buf, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func writeTestDataJSON(filename string, data interface{}) error {
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, os.ModePerm)
}
