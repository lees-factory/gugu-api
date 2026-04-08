package auth

import "testing"

func TestNormalizedClientIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		want       string
	}{
		{name: "ipv4 with port", remoteAddr: "203.0.113.10:54321", want: "203.0.113.10"},
		{name: "ipv6 with port", remoteAddr: "[2001:db8::1]:443", want: "2001:db8::1"},
		{name: "plain ipv4", remoteAddr: "198.51.100.7", want: "198.51.100.7"},
		{name: "plain ipv6", remoteAddr: "2001:db8::2", want: "2001:db8::2"},
		{name: "hostname with port", remoteAddr: "example.com:8080", want: "example.com"},
		{name: "empty", remoteAddr: "   ", want: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := normalizedClientIP(tc.remoteAddr)
			if got != tc.want {
				t.Fatalf("normalizedClientIP(%q) = %q, want %q", tc.remoteAddr, got, tc.want)
			}
		})
	}
}

func TestNormalizedDeviceName(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "trimmed value", value: "  iPhone 15  ", want: "iPhone 15"},
		{name: "empty to default", value: "", want: defaultDeviceName},
		{name: "whitespace to default", value: "   ", want: defaultDeviceName},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := normalizedDeviceName(tc.value)
			if got != tc.want {
				t.Fatalf("normalizedDeviceName(%q) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}
