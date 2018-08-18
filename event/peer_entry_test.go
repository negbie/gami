package event

import (
	"testing"

	"github.com/negbie/gami"
)

func TestPeerEntry(t *testing.T) {
	fixture := map[string]string{
		"Channeltype":    "ChannelType",
		"Objectname":     "ObjectName",
		"Chanobjecttype": "ChannelObjectType",
		"Ipaddress":      "IPAddress",
		"Ipport":         "IPPort",
		"Dynamic":        "Dynamic",
		"Natsupport":     "NatSupport",
		"Videosupport":   "VideoSupport",
		"Textsupport":    "TextSupport",
		"Acl":            "ACL",
		"Status":         "Status",
		"Realtimedevice": "RealtimeDevice",
	}
	ev := gami.AMIEvent{
		ID:        "PeerEntry",
		Privilege: []string{"all"},
		Params:    fixture,
	}
	evtype := New(&ev)
	if _, ok := evtype.(PeerEntry); !ok {
		t.Log("PeerEntry type assertion")
		t.Fail()
	}

	testEvent(t, fixture, evtype)
}
