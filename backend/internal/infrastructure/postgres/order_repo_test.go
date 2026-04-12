package postgres

import (
	"net"
	"testing"
)

func TestParsePostgresInetAcceptsHostAndCIDRForms(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "ipv4 host", raw: "172.19.0.1", want: "172.19.0.1"},
		{name: "ipv4 cidr from inet text", raw: "172.19.0.1/32", want: "172.19.0.1"},
		{name: "ipv6 cidr from inet text", raw: "2001:db8::1/128", want: "2001:db8::1"},
		{name: "trims spaces", raw: "  10.0.0.5/32  ", want: "10.0.0.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePostgresInet(tt.raw)
			if err != nil {
				t.Fatalf("parsePostgresInet(%q) returned error: %v", tt.raw, err)
			}
			if got.String() != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got.String())
			}
		})
	}
}

func TestParsePostgresInetRejectsInvalidValues(t *testing.T) {
	for _, raw := range []string{"", "not-an-ip", "172.19.0.1/999"} {
		if got, err := parsePostgresInet(raw); err == nil {
			t.Fatalf("expected error for %q, got %q", raw, got.String())
		}
	}
}

func TestIPParamRejectsNilIP(t *testing.T) {
	if _, err := ipParam(nil); err == nil {
		t.Fatal("expected nil IP to be rejected")
	}
}

func TestIPParamNormalizesIP(t *testing.T) {
	got, err := ipParam(net.ParseIP("172.19.0.1"))
	if err != nil {
		t.Fatalf("ipParam returned error: %v", err)
	}
	if got != "172.19.0.1" {
		t.Fatalf("expected normalized IPv4, got %q", got)
	}
}
