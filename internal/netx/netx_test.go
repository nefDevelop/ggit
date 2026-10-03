package netx

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLocalHostname(t *testing.T) {
	// The resolver is replaced instead of hitting DNS. Resolving gogs.io or
	// google.com from a test made the verdict depend on the network, and the
	// table below failed at random: one run in twelve saw gogs.io answer with a
	// private address and the test reported it as a local hostname.
	original := lookupIP
	t.Cleanup(func() { lookupIP = original })

	lookupIP = func(host string) ([]net.IP, error) {
		if ip, ok := fakeDNS[host]; ok {
			return []net.IP{ip}, nil
		}
		return nil, fmt.Errorf("no such host: %s", host)
	}

	tests := []struct {
		hostname  string
		allowlist []string
		want      bool
	}{
		{hostname: "localhost", want: true},       // #00
		{hostname: "127.0.0.1", want: true},       // #01
		{hostname: "::1", want: true},             // #02
		{hostname: "0:0:0:0:0:0:0:1", want: true}, // #03
		{hostname: "127.0.0.95", want: true},      // #04
		{hostname: "0.0.0.0", want: true},         // #05
		{hostname: "192.168.123.45", want: true},  // #06

		{hostname: "gogs.io", want: false},         // #07
		{hostname: "google.com", want: false},      // #08
		{hostname: "165.232.140.255", want: false}, // #09

		{hostname: "192.168.123.45", allowlist: []string{"10.0.0.17"}, want: true}, // #10
		{hostname: "gogs.local", allowlist: []string{"gogs.local"}, want: false},   // #11

		{hostname: "192.168.123.45", allowlist: []string{"*"}, want: false}, // #12
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.want, IsBlockedLocalHostname(test.hostname, test.allowlist))
		})
	}
}

// fakeDNS maps the hostnames used above to the address each one is given, so no
// lookup ever leaves the process.
var fakeDNS = map[string]net.IP{
	"gogs.io":         net.ParseIP("8.247.106.126"),
	"google.com":      net.ParseIP("142.250.185.100"),
	"165.232.140.255": net.ParseIP("165.232.140.255"),
	"gogs.local":      net.ParseIP("10.0.0.9"),
}

func TestIsLocalIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"192.168.1.1", true},
		{"172.16.0.1", true},
		{"1.1.1.1", false},
		{"8.8.8.8", false},
		{"::1", true},
		{"fe80::1", true},
	}
	for _, test := range tests {
		t.Run(test.ip, func(t *testing.T) {
			assert.Equal(t, test.want, IsLocalIP(net.ParseIP(test.ip)))
		})
	}
}

func TestNewSafeControl(t *testing.T) {
	tests := []struct {
		hostname  string
		allowlist []string
		address   string
		wantErr   bool
	}{
		{hostname: "localhost", address: "127.0.0.1:80", wantErr: true},
		{hostname: "google.com", address: "1.1.1.1:80", wantErr: false},
		{hostname: "google.com", address: "127.0.0.1:80", wantErr: true}, // Case where hostname is safe but resolved IP is local
		{hostname: "localhost", allowlist: []string{"localhost"}, address: "127.0.0.1:80", wantErr: false},
		{hostname: "any", allowlist: []string{"*"}, address: "127.0.0.1:80", wantErr: false},
		{hostname: "invalid-addr", address: "not-an-ip", wantErr: true},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			control := NewSafeControl(test.hostname, test.allowlist)
			if control == nil {
				assert.False(t, test.wantErr)
				return
			}
			err := control("tcp", test.address, nil)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
