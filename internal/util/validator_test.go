package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	testCases := map[string]struct {
		url  string
		want bool
	}{
		"valid_https": {
			url:  "https://example.com",
			want: true,
		},
		"valid_http": {
			url:  "http://example.com",
			want: true,
		},
		"another_valid_https": {
			url:  "https://anotherexample.com",
			want: true,
		},
		"missing_host": {
			url:  "https://",
			want: false,
		},
		"wrong_slashes": {
			url:  "http:///example.com",
			want: false,
		},
		"wrong_dots": {
			url:  "http://example..com",
			want: false,
		},
		"without_two_dots": {
			url:  "http//example.com",
			want: false,
		},
		"wrong_two_dots": {
			url:  "http:://example.com",
			want: false,
		},
		"invalid_ftp": {
			url:  "ftp://example.com",
			want: false,
		},
		"empty": {
			url:  "",
			want: false,
		},
		"only_domain_no_scheme": {
			url:  "example.com",
			want: false,
		},
		"https_with_path": {
			url:  "https://example.com/path",
			want: true,
		},
		"https_with_query": {
			url:  "https://example.com/path?query=1",
			want: true,
		},
		"missing_tld": {
			url:  "https://example",
			want: false,
		},
		"too_short_tld": {
			url:  "https://example.a",
			want: false,
		},
		"too_long_tld": {
			url:  "https://example.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want: false,
		},
		"local_file": {
			url:  "file:///Users/example/file.txt",
			want: false,
		},
		"data": {
			url:  "data:text/plain;base64,SGVsbG8sIFdvcmxkIQ==",
			want: false,
		},
		"mailto": {
			url:  "mailto:someone@example.com",
			want: false,
		},
		"tel": {
			url:  "tel:+1234567890",
			want: false,
		},
		"https_with_fragment": {
			url:  "https://example.com/path#section",
			want: true,
		},
		"https_with_subdomain": {
			url:  "https://sub.example.com",
			want: true,
		},
		"https_with_subdomains": {
			url:  "https://sub.sub.example.com",
			want: false,
		},
		"https_with_port": {
			url:  "https://example.com:443",
			want: true,
		},
		"https_with_bad_port": {
			url:  "https://example.com:bad",
			want: false,
		},
		"https_with_negative_port": {
			url:  "https://example.com:-7",
			want: false,
		},
		"https_with_too_big_port": {
			url:  "https://example.com:65536",
			want: false,
		},
		"https_with_two_ports": {
			url:  "https://example.com:80:80",
			want: false,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, IsURL(tc.url))
		})
	}
}
