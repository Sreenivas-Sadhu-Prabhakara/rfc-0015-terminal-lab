package lab

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const MaxLines = 1000
const MaxLineBytes = 4096

type Event struct {
	Line    int    `json:"line"`
	Route   string `json:"route"`
	Host    string `json:"host,omitempty"`
	Message string `json:"message"`
}

// Replay routes a finite transcript through a simulated terminal session.
func Replay(input io.Reader) ([]Event, error) {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), MaxLineBytes+2)
	session := terminal{escape: '#'}
	events := make([]Event, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(events) == MaxLines {
			return nil, fmt.Errorf("transcript exceeds %d lines", MaxLines)
		}
		if len(line) > MaxLineBytes || !utf8.ValidString(line) {
			return nil, fmt.Errorf("invalid UTF-8 or overlong input at line %d", len(events)+1)
		}
		event := session.route(line)
		event.Line = len(events) + 1
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read transcript: %w", err)
	}
	return events, nil
}

type terminal struct {
	escape byte
	host   string
}

func (s *terminal) route(line string) Event {
	if len(line) == 0 || line[0] != s.escape {
		return s.forward(line)
	}
	if len(line) > 1 && line[1] == s.escape {
		return s.forward(line[1:])
	}
	fields := strings.Fields(line[1:])
	if len(fields) == 0 {
		return s.reject("empty local command")
	}
	switch fields[0] {
	case "connect":
		if len(fields) != 2 {
			return s.reject("usage: connect SITE")
		}
		if s.host != "" {
			return s.reject("log out before connecting again")
		}
		host := strings.ToUpper(fields[1])
		if host != "SRI" && host != "UTAH" {
			return s.reject("site is outside the simulator allowlist")
		}
		s.host = host
		return Event{Route: "local", Host: s.host, Message: "opened simulated session"}
	case "logout":
		if len(fields) != 1 || s.host == "" {
			return s.reject("logout requires an active session and no arguments")
		}
		host := s.host
		s.host = ""
		return Event{Route: "local", Host: host, Message: "closed simulated session"}
	case "escape":
		if len(fields) != 2 || len(fields[1]) != 1 || !strings.Contains("#!~%", fields[1]) {
			return s.reject("escape must be one of # ! ~ %")
		}
		s.escape = fields[1][0]
		return Event{Route: "local", Host: s.host, Message: "escape changed to " + fields[1]}
	default:
		return s.reject("unknown local command")
	}
}

func (s *terminal) forward(line string) Event {
	if s.host == "" {
		return s.reject("no active session")
	}
	return Event{Route: "remote", Host: s.host, Message: line}
}

func (s *terminal) reject(reason string) Event {
	return Event{Route: "rejected", Host: s.host, Message: reason}
}
