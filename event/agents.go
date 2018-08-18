// Package event for AMI
package event

//Agents trigger for agents
type Agents struct {
	Privilege        []string
	Status           string `AMI:"Status"`
	Agent            string `AMI:"Agent"`
	Name             string `AMI:"Name"`
	Channel          string `AMI:"Channel"`
	LoggedInTime     string `AMI:"Loggedintime"`
	TalkingTo        string `AMI:"Talkingto"`
	TalkingToChannel string `AMI:"Talkingtochannel"`
}

type AgentCalled struct {
	Privilege          []string
	Queue              string `AMI:"Queue"`
	AgentCalled        string `AMI:"Agentcalled"`
	AgentName          string `AMI:"Agentname"`
	ChannelCalling     string `AMI:"Channelcalling"`
	DestinationChannel string `AMI:"Destinationchannel"`
	CallerIDNum        string `AMI:"Calleridnum"`
	CallerIDName       string `AMI:"Calleridname"`
	ConnectedLineNum   string `AMI:"Connectedlinenum"`
	ConnectedLineName  string `AMI:"Connectedlinename"`
	Context            string `AMI:"Context"`
	Extension          string `AMI:"Extension"`
	Priority           int    `AMI:"Priority"`
	Uniqueid           string `AMI:"Uniqueid"`
}

func init() {
	eventTrap["Agents"] = Agents{}
	eventTrap["AgentCalled"] = AgentCalled{}
}
