package exporter

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

type Interaction struct {
	URL        string `json:"url"`
	Content    string `json:"content"`
	StatusCode int    `json:"status"`
}

type Interactions struct {
	Primary        Interaction  `json:"primary"`
	Secondary      *Interaction `json:"secondary,omitempty"`
	Candidate      Interaction  `json:"candidate"`
	DifferenceMode string       `json:"differenceMode"`
	Result         bool         `json:"result"`
	Processed      time.Time    `json:"processedDate"`
}

func CreateInteraction(url string, content []byte, statusCode int) Interaction {
	return Interaction{
		URL:        url,
		Content:    string(content),
		StatusCode: statusCode,
	}
}

func CreateInteractions(primary Interaction, secondary *Interaction, candidate Interaction, differenceMode string, result bool) Interactions {
	interactions := Interactions{
		Primary:        primary,
		Candidate:      candidate,
		DifferenceMode: differenceMode,
		Result:         result,
		Processed:      time.Now(),
	}

	if secondary != nil {
		interactions.Secondary = secondary
	}

	return interactions
}

func ExportToFile(file string, interactions Interactions) error {
	f, err := os.Create(file)

	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = json.NewEncoder(w).Encode(interactions)
	w.Flush()
	return err
}
