package main

import (
	"net/url"
	"os"

	"github.com/lordofthejars/diferencia/core"
	"github.com/lordofthejars/diferencia/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diferencia",
	Short: "Interact with Diferencia",
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

			differenceMode, err := core.NewDifference(difference)

			if err != nil {
				log.Error("Error while setting difference mode. %s", err.Error())
				os.Exit(1)
			}
			config.DifferenceMode = differenceMode

			if noiseDetection && len(secondaryURL) == 0 {

				log.Error("If Noise Detection is enabled, you need to provide a secondary URL as well")
				os.Exit(1)

			}

			if len(config.ServiceName) == 0 {
				candidateURL, _ := url.Parse(config.Candidate)
				config.ServiceName = candidateURL.Hostname()
			}

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

	cmdStart.Flags().BoolVar(&headers, "headers", false, "Enable Http headers comparision")
	cmdStart.Flags().StringSliceVar(&ignoreHeadersValues, "ignoreHeadersValues", nil, "List of headers key where its value should be ignored for comparision purposes")

	cmdStart.Flags().BoolVar(&prometheus, "prometheus", false, "Enable Prometheus endpoint")
	cmdStart.Flags().IntVar(&prometheusPort, "prometheusPort", 8081, "Prometheus port")

	cmdStart.MarkFlagRequired("primary")
	cmdStart.MarkFlagRequired("candidate")

	rootCmd.AddCommand(cmdStart)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

}
