// Package event for AMI
package event

// This is an async event when a call enters a queue
type QueueJoin struct {
	Privilege         []string
	Queue             string `AMI:"Queue"`
	Channel           string `AMI:"Channel"`
	Uniqueid          string `AMI:"Uniqueid"`
	CallerIDNum       string `AMI:"Calleridnum"`
	CallerIDName      string `AMI:"Calleridname"`
	ConnectedLineNum  string `AMI:"Connectedlinenum"`
	ConnectedLineName string `AMI:"Connectedlinename"`
	Position          int    `AMI:"Position"`
	Count             int    `AMI:"Count"`
}

// This is an async event when a call leaves a queue
type QueueLeave struct {
	Privilege []string
	Queue     string `AMI:"Queue"`
	Channel   string `AMI:"Channel"`
	Uniqueid  string `AMI:"Uniqueid"`
	Position  int    `AMI:"Position"`
	Count     int    `AMI:"Count"`
}

// This is an async event when a caller hangs up while in queue
type QueueCallerAbandon struct {
	Privilege        []string
	Queue            string `AMI:"Queue"`
	Uniqueid         string `AMI:"Uniqueid"`
	Position         int    `AMI:"Position"`
	OriginalPosition int    `AMI:"Originalposition"`
	HoldTime         int    `AMI:"Holdtime"`
}

// This is an async event when a queue call is hung up by the caller or agent or is transfered
type QueueAgentComplete struct {
	Privilege  []string
	Queue      string `AMI:"Queue"`
	Uniqueid   string `AMI:"Uniqueid"`
	Channel    string `AMI:"Channel"`
	Member     string `AMI:"Member"` // This is called Location in other events
	MemberName string `AMI:"Membername"`
	HoldTime   int    `AMI:"Holdtime"`
	TalkTime   int    `AMI:"Talktime"`
	Reason     string `AMI:"Reason"` // caller, agent, transfer
}

// This is an async event sent whenever a queue member's status changes (ringing, extension state change, etc.)
type QueueMemberStatus struct {
	Privilege  []string
	Queue      string `AMI:"Queue"`
	Location   string `AMI:"Location"`
	Name       string `AMI:"Membername"`
	Membership string `AMI:"Membership"`
	Penalty    int    `AMI:"Penalty"`
	CallsTaken int    `AMI:"Callstaken"`
	LastCall   int    `AMI:"Lastcall"`
	WhenAdded  int    `AMI:"Whenadded"`
	Status     int    `AMI:"Status"`
	Paused     int    `AMI:"Paused"`
}

type QueueMemberAdded struct {
	Privilege  []string
	Queue      string `AMI:"Queue"`
	Location   string `AMI:"Location"`
	MemberName string `AMI:"Membername"`
	Membership string `AMI:"Membership"`
	Penalty    int    `AMI:"Penalty"`
	CallsTaken int    `AMI:"Callstaken"`
	LastCall   int    `AMI:"Lastcall"`
	WhenAdded  int    `AMI:"Whenadded"`
	Status     int    `AMI:"Status"`
	Paused     int    `AMI:"Paused"`
}

type QueueMemberRemoved struct {
	Privilege  []string
	Queue      string `AMI:"Queue"`
	Location   string `AMI:"Location"`
	MemberName string `AMI:"Membername"`
}

type QueueMemberPaused struct {
	Privilege  []string
	Queue      string `AMI:"Queue"`
	Location   string `AMI:"Location"`
	MemberName string `AMI:"Membername"`
	Paused     int    `AMI:"Paused"`
	Reason     int    `AMI:"Reason"`
}

type QueueMemberPenalty struct {
	Privilege []string
	Queue     string `AMI:"Queue"`
	Location  string `AMI:"Location"`
	Penalty   int    `AMI:"Penalty"`
}

type QueueStatsCleared struct {
	Privilege        []string
	Queue            string `AMI:"Queue"`
	WhenStatsCleared string `AMI:"Whenstatscleared"`
}

func init() {
	eventTrap["Join"] = QueueJoin{}
	eventTrap["Leave"] = QueueLeave{}
	eventTrap["QueueCallerAbandon"] = QueueCallerAbandon{}
	eventTrap["AgentComplete"] = QueueAgentComplete{}
	eventTrap["QueueMemberStatus"] = QueueMemberStatus{}
	eventTrap["QueueMemberAdded"] = QueueMemberAdded{}
	eventTrap["QueueMemberRemoved"] = QueueMemberRemoved{}
	eventTrap["QueueMemberPaused"] = QueueMemberPaused{}
	eventTrap["QueueMemberPenalty"] = QueueMemberPenalty{}
	eventTrap["QueueStatsCleared"] = QueueStatsCleared{}
}
