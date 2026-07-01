package config

import (
	"os"
	"testing"
)

func TestProxyFromEnvPrecedence(t *testing.T) {
	all := []string{"HTTPS_PROXY", "https_proxy", "ALL_PROXY", "all_proxy", "HTTP_PROXY", "http_proxy"}
	clear := func() {
		for _, k := range all {
			os.Unsetenv(k)
		}
	}

	cases := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"none", nil, ""},
		{"https wins over all and http", map[string]string{"HTTPS_PROXY": "http://h", "ALL_PROXY": "socks5://a", "HTTP_PROXY": "http://p"}, "http://h"},
		{"all wins, http_proxy ignored", map[string]string{"ALL_PROXY": "socks5://a", "HTTP_PROXY": "http://p"}, "socks5://a"},
		{"http_proxy alone is ignored (http-only, not for HTTPS endpoint)", map[string]string{"HTTP_PROXY": "http://p", "http_proxy": "http://p2"}, ""},
		{"uppercase preferred over lowercase", map[string]string{"HTTPS_PROXY": "http://up", "https_proxy": "http://low"}, "http://up"},
		{"lowercase https used when only lowercase set", map[string]string{"https_proxy": "http://low"}, "http://low"},
		{"lowercase all_proxy socks5h", map[string]string{"all_proxy": "socks5h://a:1"}, "socks5h://a:1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clear()
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
			if got := proxyFromEnv(); got != tc.want {
				t.Fatalf("proxyFromEnv() = %q, want %q", got, tc.want)
			}
		})
	}
	clear()
}
