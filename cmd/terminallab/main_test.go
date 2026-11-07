package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHTMLReportEscapesRemoteText(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"-html"}, strings.NewReader("#connect sri\n<script>alert(1)</script>\n"), &out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "<script>") || !strings.Contains(out.String(), "&lt;script&gt;") {
		t.Fatal("remote text was not escaped")
	}
}
