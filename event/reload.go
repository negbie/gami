// Package event for AMI
package event

type SystemReload struct {
	Privilege []string
	Module    string `AMI:"Module"`
	Message   string `AMI:"Message"`
}

func init() {
	eventTrap["Reload"] = SystemReload{}
}
