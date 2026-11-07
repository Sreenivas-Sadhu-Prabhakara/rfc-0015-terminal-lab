package lab

import (
	"reflect"
	"strings"
	"testing"
)

func TestReplayKeepsLocalCommandsOutOfRemoteTranscript(t *testing.T) {
	input := "before\n#connect unknown\n#connect sri\n  payload  \n##logout\n#escape !\n#logout\n!bogus\n!logout\nafter\n!connect utah\nnext\n"
	events, err := Replay(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var remote []Event
	for _, event := range events {
		if event.Route == "remote" {
			remote = append(remote, event)
		}
	}
	want := []Event{
		{Line: 4, Route: "remote", Host: "SRI", Message: "  payload  "},
		{Line: 5, Route: "remote", Host: "SRI", Message: "#logout"},
		{Line: 7, Route: "remote", Host: "SRI", Message: "#logout"},
		{Line: 12, Route: "remote", Host: "UTAH", Message: "next"},
	}
	if !reflect.DeepEqual(remote, want) {
		t.Fatalf("got %+v, want %+v", remote, want)
	}
	for _, index := range []int{0, 1, 7, 9} {
		if events[index].Route != "rejected" {
			t.Fatalf("line %d was not rejected", index+1)
		}
	}
}

func TestInvalidCommandsPreserveRouteAndEscape(t *testing.T) {
	events, err := Replay(strings.NewReader("#connect sri\n#connect utah\n#escape é\n#logout extra\nmessage\n#logout\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events[1:4] {
		if event.Route != "rejected" || event.Host != "SRI" {
			t.Fatalf("unexpected event: %+v", event)
		}
	}
	if events[4].Host != "SRI" || events[4].Route != "remote" || events[5].Route != "local" {
		t.Fatalf("session changed: %+v", events)
	}
}

func TestReplayInputBoundaries(t *testing.T) {
	maxLine := strings.Repeat("x", MaxLineBytes)
	if _, err := Replay(strings.NewReader("#connect sri\r\n" + maxLine + "\r\n")); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{maxLine + "x", "\xff", strings.Repeat("x\n", MaxLines+1)} {
		events, err := Replay(strings.NewReader(bad))
		if err == nil || events != nil {
			t.Fatalf("malformed transcript returned %d events and %v", len(events), err)
		}
	}
}
