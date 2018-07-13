package main

import (
	"net/url"
	"os"

	"github.com/lordofthejars/diferencia/core"
	"github.com/lordofthejars/diferencia/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diferencia",
	Short: "Interact with Diferencia",
}

func areHttpsClientAttributesCorrect(caCert, clientCert, clientKey string) bool {
	return (len(caCert) == 0 && len(clientCert) == 0 && len(clientKey) == 0) || (len(caCert) > 0 && len(clientCert) > 0 && len(clientKey) > 0)
}

func main() {

	var port int
	var serviceName, primaryURL, secondaryURL, candidateURL, difference string
	var allowUnsafeOperations, noiseDetection bool
	var storeResults string
	var prometheus bool
	var prometheusPort int
	var headers bool
	var ignoreHeadersValues []string
	var ignoreValuesOf []string
	var ignoreValuesFile string
	var logLevel string
	var insecureSkipVerify bool
	var caCert, clientCert, clientKey string

	var cmdStart = &cobra.Command{
		Use:   "start",
		Short: "Start Diferencia",
		Long:  `start is used to start Diferencia server to start spreading calls across network`,
		Run: func(cmd *cobra.Command, args []string) {
			config := core.DiferenciaConfiguration{}

			config.Port = port
			config.ServiceName = serviceName
			config.Primary = primaryURL
			config.Secondary = secondaryURL
			config.Candidate = candidateURL
			config.StoreResults = storeResults
			config.NoiseDetection = noiseDetection
			config.AllowUnsafeOperations = allowUnsafeOperations
			config.Headers = headers
			config.IgnoreHeadersValues = ignoreHeadersValues
			config.Prometheus = prometheus
			config.PrometheusPort = prometheusPort
			config.IgnoreValues = ignoreValuesOf
			config.IgnoreValuesFile = ignoreValuesFile
			config.InsecureSkipVerify = insecureSkipVerify
			config.CaCert = caCert
			config.ClientCert = clientCert
			config.ClientKey = clientKey

			differenceMode, err := core.NewDifference(difference)

			if err != nil {
				logrus.Errorf("Error while setting difference mode. %s", err.Error())
				os.Exit(1)
			}
			config.DifferenceMode = differenceMode

			if !areHttpsClientAttributesCorrect(caCert, clientCert, clientKey) {
				logrus.Errorf("Https Client options should either not provided or all of them provided but not only some. caCert: %s, clientCert: %s, clientkey: %s.", caCert, clientCert, clientKey)
				os.Exit(1)
			}

			if noiseDetection && len(secondaryURL) == 0 {
				logrus.Errorf("If Noise Detection is enabled, you need to provide a secondary URL as well")
				os.Exit(1)
			}

			if !noiseDetection && (config.IsIgnoreValuesFileSet() || config.IsIgnoreValuesSet()) {
				logrus.Infof("ignoreValues or ignoreValuesFile attributes are set but noise detection is disabled, so they are going to be ignored.")
			}

			if len(config.ServiceName) == 0 {
				candidateURL, _ := url.Parse(config.Candidate)
				config.ServiceName = candidateURL.Hostname()
			}

			log.Initialize(logLevel)
			core.StartProxy(&config)
		},
	}

	cmdStart.Flags().IntVar(&port, "port", 8080, "Listening port of Diferencia proxy")
	cmdStart.Flags().StringVar(&serviceName, "serviceName", "", "Sets service name under test. By default it takes candidate hostname")
	cmdStart.Flags().StringVarP(&primaryURL, "primary", "p", "", "Primary Service URL")
	cmdStart.Flags().StringVarP(&secondaryURL, "secondary", "s", "", "Secondary Service URL")
	cmdStart.Flags().StringVarP(&candidateURL, "candidate", "c", "", "Candidate Service URL")
	cmdStart.Flags().StringVarP(&difference, "difference", "d", "Strict", "Difference mode to compare JSONs")
	cmdStart.Flags().BoolVarP(&allowUnsafeOperations, "unsafe", "u", false, "Allow none safe operations like PUT, POST, PATCH, ...")
	cmdStart.Flags().BoolVarP(&noiseDetection, "noisedetection", "n", false, "Enable noise detection. Secondary URL must be provided.")
	cmdStart.Flags().StringVar(&storeResults, "storeResults", "", "Directory where output is set. If not specified then nothing is stored. Useful for local development.")

	cmdStart.Flags().StringVarP(&logLevel, "logLevel", "l", "error", "Set log level")

	cmdStart.Flags().BoolVar(&headers, "headers", false, "Enable Http headers comparision")
	cmdStart.Flags().StringSliceVar(&ignoreHeadersValues, "ignoreHeadersValues", nil, "List of headers key where their value must be ignored for comparision purposes.")

	cmdStart.Flags().StringSliceVar(&ignoreValuesOf, "ignoreValues", nil, "List of JSON Pointers of values that must be ignored for comparision purposes.")
	cmdStart.Flags().StringVar(&ignoreValuesFile, "ignoreValuesFile", "", "File location where each line is a JSON pointers definition for ignoring values.")

	cmdStart.Flags().BoolVar(&prometheus, "prometheus", false, "Enable Prometheus endpoint")
	cmdStart.Flags().IntVar(&prometheusPort, "prometheusPort", 8081, "Prometheus port")

	cmdStart.Flags().BoolVar(&insecureSkipVerify, "insecureSkipVerify", false, "Sets Insecure Skip Verify flag in Http Client")
	cmdStart.Flags().StringVar(&caCert, "caCert", "", "Certificate Authority path (PEM)")
	cmdStart.Flags().StringVar(&clientCert, "clientCert", "", "Client Certificate path (X509)")
	cmdStart.Flags().StringVar(&clientKey, "clientKey", "", "Client Key path (X509)")

	cmdStart.MarkFlagRequired("primary")
	cmdStart.MarkFlagRequired("candidate")

	rootCmd.AddCommand(cmdStart)

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

}
