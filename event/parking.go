// Package event for AMI
package event

type ParkedCall struct {
    Privilege         []string
    Exten             string `AMI:"Exten"`
    Channel           string `AMI:"Channel"`
    Parkinglot        string `AMI:"Parkinglot"`
    From              string `AMI:"From"`
    Timeout           int    `AMI:"Timeout"`
    CallerIDNum       string `AMI:"CallerIDNum"`
    CallerIDName      string `AMI:"CallerIDName"`
    ConnectedLineNum  string `AMI:"ConnectedLineNum"`
    ConnectedLineName string `AMI:"ConnectedLineName"`
    Uniqueid          string `AMI:"Uniqueid"`
    CustGroup         string `AMI:"CustGroup"`
}

func init() {
    eventTrap["ParkedCall"] = ParkedCall{}
    eventTrap["UnParkedCall"] = ParkedCall{}
    eventTrap["ParkedCallGiveUp"] = ParkedCall{}
    eventTrap["ParkedCallTimeOut"] = ParkedCall{}
}
