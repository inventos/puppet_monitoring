// Describes json report model comes from puppet server

package impl

import (
	"encoding/json"
	"log"
)

// Unmarshal json string into structured data
func (p ReportItem) FromJson(data []byte) (*ReportItem, error) {
	var result ReportItem
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
	}
	return &result, err
}

// Top level info
type ReportItem struct {
	Host          string     `json:"host"`
	Environment   string     `json:"environment"`
	PuppetVersion string     `json:"puppet_version"`
	Status        string     `json:"status"`
	Metrics       MetricItem `json:"metrics"`
	Logs          []LogItem  `json:"logs"`
}

// Metrics - container of metrics failures
type MetricItem struct {
	Resources MetricResourcesItem `json:"resources"`
	Events    MetricEventsItem    `json:"events"`
}

// Describes how many resources / restarts was failed to implement
type MetricResourcesItem struct {
	Failed          int `json:"failed"`
	FailedToRestart int `json:"failed_to_restart"`
}

// Describes how many events was failed to implement
type MetricEventsItem struct {
	Failed int `json:"failure"`
}

// Describes log item
type LogItem struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}
