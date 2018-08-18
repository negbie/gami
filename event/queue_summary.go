// Package event for AMI
package event

// This event is returned after a QueueSummary manager command
type QueueSummary struct {
	Privilege       []string
	Queue           string `AMI:"Queue"`
	LoggedIn        int    `AMI:"Loggedin"`
	Available       int    `AMI:"Available"`
	Callers         int    `AMI:"Callers"`
	HoldTime        int    `AMI:"Holdtime"`
	TalkTime        int    `AMI:"Talktime"`
	LongestHoldTime int    `AMI:"Longestholdtime"`
}

func init() {
	eventTrap["QueueSummary"] = QueueSummary{}
}
