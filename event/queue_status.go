// Package event for AMI
package event

// This event is returned after a QueueStatus manager command
type QueueParams struct {
	Privilege        []string
	Queue            string  `AMI:"Queue"`
	Max              int     `AMI:"Max"`
	Strategy         string  `AMI:"Strategy"`
	Calls            int     `AMI:"Calls"`
	HoldTime         int     `AMI:"Holdtime"`
	TalkTime         int     `AMI:"Talktime"`
	Completed        int     `AMI:"Completed"`
	Abandoned        int     `AMI:"Abandoned"`
	ServiceLevel     int     `AMI:"Servicelevel"`
	ServicelevelPerf float64 `AMI:"Servicelevelperf"`
	Weight           int     `AMI:"Weight"`
	WhenStatsCleared string  `AMI:"Whenstatscleared"`
}

// This event is returned after a QueueStatus manager command
type QueueMember struct {
	Privilege  []string
	Queue      string `AMI:"Queue"`
	Name       string `AMI:"Name"`
	Location   string `AMI:"Location"`
	Membership string `AMI:"Membership"`
	DynamicAge string `AMI:"Dynamicage"`
	Penalty    int    `AMI:"Penalty"`
	CallsTaken int    `AMI:"Callstaken"`
	LastCall   int    `AMI:"Lastcall"`
	WhenAdded  int    `AMI:"Whenadded"`
	Status     int    `AMI:"Status"`
	Paused     int    `AMI:"Paused"`
}

// This event is returned after a QueueStatus manager command
type QueueEntry struct {
	Privilege         []string
	Queue             string `AMI:"Queue"`
	Position          int    `AMI:"Position"`
	Channel           string `AMI:"Channel"`
	Uniqueid          string `AMI:"Uniqueid"`
	CallerIDNum       string `AMI:"Calleridnum"`
	CallerIDName      string `AMI:"Calleridname"`
	ConnectedLineNum  string `AMI:"Connectedlinenum"`
	ConnectedLineName string `AMI:"Connectedlinename"`
	Wait              int    `AMI:"Wait"`
}

func init() {
	eventTrap["QueueParams"] = QueueParams{}
	eventTrap["QueueMember"] = QueueMember{}
	eventTrap["QueueEntry"] = QueueEntry{}
}
