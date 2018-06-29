package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// RegisterNumberOfRegressions cpunter to Prometheus register
func RegisterNumberOfRegressions(namespace string) *prometheus.CounterVec {

	filteredNamespace := strings.Replace(namespace, ".", "_", -1)

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: filteredNamespace,
		Name:      "service_regressions_failures_total",
		Help:      "Number of regressions detected by endpoints.",
	},
		[]string{"method", "path"},
	)

	prometheus.MustRegister(counter)

	return counter
}
