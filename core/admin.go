package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DiferenciaConfigurationUpdate struct {
	ServiceName    string `json:"serviceName,omitempty"`
	Primary        string `json:"primary,omitempty"`
	Secondary      string `json:"secondary,omitempty"`
	Candidate      string `json:"candidate,omitempty"`
	NoiseDetection string `json:"noiseDetection,omitempty"`
	Mode           string `json:"mode,omitempty"`
}

func (config DiferenciaConfigurationUpdate) isServiceNameSet() bool {
	return len(config.ServiceName) > 0
}

func (config DiferenciaConfigurationUpdate) isPrimarySet() bool {
	return len(config.Primary) > 0
}

func (config DiferenciaConfigurationUpdate) isSecondarySet() bool {
	return len(config.Secondary) > 0
}

func (config DiferenciaConfigurationUpdate) isCandidateSet() bool {
	return len(config.Candidate) > 0
}

func (config DiferenciaConfigurationUpdate) isNoiseDetectionSet() bool {
	return len(config.NoiseDetection) > 0
}

func (config DiferenciaConfigurationUpdate) isModeSet() bool {
	return len(config.Mode) > 0
}

func (config DiferenciaConfigurationUpdate) getMode() (Difference, error) {
	return NewDifference(config.Mode)
}

func (config DiferenciaConfigurationUpdate) getNoiseDetection() (bool, error) {
	return strconv.ParseBool(config.NoiseDetection)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {

	mutex.Lock()
	defer mutex.Unlock()

	if r.Method == http.MethodPut {
		var updateConfig DiferenciaConfigurationUpdate

		if err := json.NewDecoder(r.Body).Decode(&updateConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		if err := Config.UpdateConfiguration(updateConfig); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			type Alias DiferenciaConfiguration
			json.NewEncoder(w).Encode(&struct {
				Mode string `json:"differenceMode,omitempty"`
				*Alias
			}{
				Mode:  Config.DifferenceMode.String(),
				Alias: (*Alias)(Config),
			})
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

}
