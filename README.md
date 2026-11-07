# RFC 15 — Terminal routing lab

**Covered RFC 15. Next RFC 16.** The RFC Editor index lists number 14 as
**Not Issued**. RFC 13's published source was checked before continuing.

C. Stephen Carr's *Network subsystem for time sharing hosts* is dated
25 September 1969. It proposes an early ARPANET Telnet subsystem that lets
a local terminal act as a terminal on a remote time-sharing host. The
subsystem wraps host system calls for opening, transmitting over and
closing connections. A primary connection carries terminal interaction.
An auxiliary connection can carry file data. A declared escape character
routes the following command to the local subsystem.

The RFC Editor lists status **Unknown**, stream **Legacy**, and empty
updates, updated-by, obsoletes and obsoleted-by fields, checked 7 November
2026. RFC 11 supplies historical background through a textual reference.

- [RFC Editor index](https://www.rfc-editor.org/rfc-index.txt)
- [RFC 15 official text](https://www.rfc-editor.org/rfc/rfc15.txt)
- [RFC 15 catalog metadata](https://www.rfc-editor.org/rfc/rfc15.json)

## Exercise

1. **Cloud operator console training.** Replay a transcript to distinguish
   local session commands from remote input. This is the exercise built here.
2. **Banking support tool tests.** Exercise the boundary between session
   management and customer-system input using synthetic transcripts.
3. **Gateway audit replay.** Label which route handled an input and show
   where a disconnected session rejected it.

The educational simulator uses complete UTF-8 lines and a fictional session
state machine. The original proposal describes terminal character forwarding,
host permissions and auxiliary file copying. This lab focuses on routing
semantics through a local CLI and deliberately uses synthetic host labels.
Its scope is teaching, with no network access or command execution.
A production operator console needs independently enforced authorization,
an authenticated encrypted transport, and audit controls.

Start with `#` as the escape prefix. Use `#connect SRI` or
`#connect UTAH` to open a simulated session, `#logout` to close it,
and `#escape !` to change the active prefix. A doubled active prefix
sends one literal prefix as text. The doubled-prefix rule and line-start
recognition are lab design choices. After `#escape !`, `#logout`
becomes ordinary remote text and `!logout` closes the session.
Unknown or malformed local commands are rejected without changing
session state. Text outside an active session is rejected.

## Run

Requires Go 1.22 or later. Replay the included transcript from the repo root.

```sh
go run ./cmd/terminallab < examples/session.txt
go run ./cmd/terminallab -html < examples/session.txt > docs/index.html
go test -race ./...
go vet ./...
go build ./...
```

The sample has 13 input lines. Its JSON and HTML traces show five local
actions, four simulated remote lines and four rejections. A replay can
contain rejected commands and still produce a report. Invalid UTF-8, a line
over 4,096 bytes, more than 1,000 lines, or an input/output failure causes
the CLI to fail. The HTML page is a saved trace rendered by Go.

Tests check local command isolation, literal escape handling, escape changes,
session state after rejected commands, input limits and HTML escaping.
MIT licensed.
