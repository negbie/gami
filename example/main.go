package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/negbie/gami"
	"github.com/negbie/gami/event"
)

func main() {

	ami, err := gami.Dial("127.0.0.1:5038")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	ami.Run()
	defer ami.Close()

	//install manager
	go func() {
		for {
			select {
			//handle network errors
			case err := <-ami.NetError:
				log.Println("Network Error:", err)
				//try new connection every second
				<-time.After(time.Second)
				if err := ami.Reconnect(); err == nil {
					//call start actions
					ami.Action("Events", gami.Params{"EventMask": "on"}, time.Second*5)
				}

			case err := <-ami.Error:
				log.Println("error:", err)
			//wait events and process
			case ev := <-ami.Events:
				log.Printf("Event Detect: %v", *ev)
				//if want type of events
				log.Println("EventType:", event.New(ev))
			}
		}
	}()

	if err := ami.Login("admin", "root"); err != nil {
		log.Fatal(err)
	}

	if rs, err := ami.Action("Ping", nil, time.Second*5); err != nil {
		log.Fatal(rs)
	}

	//or with can do async
	pingResp, pingErr := ami.AsyncAction("Ping", gami.Params{"ActionID": "miping"})
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	if _, err := ami.Action("Events", gami.Params{"EventMask": "on"}, time.Second*5); err != nil {
		fmt.Print(err)
	}

	log.Println("future ping:", <-pingResp)
}
