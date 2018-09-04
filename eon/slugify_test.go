package eon

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := map[string]string{
		"HTTPServer":    "http-server",
		"HTTPServerID":  "http-server-id",
		"IP":            "ip",
		"IPAddress":     "ip-address",
		"LogLevel":      "log-level",
		"Name":          "name",
		"NodeID":        "node-id",
		"NodeIPAddress": "node-ip-address",
	}
	for v, expect := range tests {
		got := string(slugify(v))
		if got != expect {
			t.Errorf("mismatching slug output for %q: expected %q, got %q", v, expect, got)
		}
	}
}
