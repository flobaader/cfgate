package cloudflare

import (
	"errors"
	"testing"
)

func TestValidateHostnameDepth(t *testing.T) {
	tests := []struct {
		name        string
		hostname    string
		wantErr     bool
		wantErrType error
	}{
		// Valid cases - apex domains
		{
			name:     "apex domain",
			hostname: "example.com",
			wantErr:  false,
		},
		{
			name:     "apex domain with complex TLD",
			hostname: "example.co.uk",
			wantErr:  false,
		},
		{
			name:     "apex domain .io",
			hostname: "selectcode.io",
			wantErr:  false,
		},

		// Valid cases - single-level subdomains
		{
			name:     "single-level subdomain",
			hostname: "app.example.com",
			wantErr:  false,
		},
		{
			name:     "single-level subdomain with complex TLD",
			hostname: "app.example.co.uk",
			wantErr:  false,
		},
		{
			name:     "www subdomain",
			hostname: "www.example.com",
			wantErr:  false,
		},
		{
			name:     "api subdomain",
			hostname: "api.selectcode.io",
			wantErr:  false,
		},
		{
			name:     "cfgate-demo subdomain",
			hostname: "cfgate-demo.selectcode.io",
			wantErr:  false,
		},

		// Invalid cases - multi-level subdomains
		{
			name:        "two-level subdomain",
			hostname:    "api.staging.example.com",
			wantErr:     true,
			wantErrType: ErrMultiLevelSubdomain,
		},
		{
			name:        "two-level subdomain with testing prefix",
			hostname:    "tunnel-demo.testing.selectcode.io",
			wantErr:     true,
			wantErrType: ErrMultiLevelSubdomain,
		},
		{
			name:        "two-level subdomain with complex TLD",
			hostname:    "api.staging.example.co.uk",
			wantErr:     true,
			wantErrType: ErrMultiLevelSubdomain,
		},
		{
			name:        "three-level subdomain",
			hostname:    "a.b.c.example.com",
			wantErr:     true,
			wantErrType: ErrMultiLevelSubdomain,
		},
		{
			name:        "deep subdomain",
			hostname:    "very.deep.nested.subdomain.example.com",
			wantErr:     true,
			wantErrType: ErrMultiLevelSubdomain,
		},

		// Edge cases
		{
			name:     "empty string",
			hostname: "",
			wantErr:  false, // ExtractZoneFromHostname returns empty, so no validation
		},
		{
			name:     "single word (localhost-like)",
			hostname: "localhost",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHostnameDepth(tt.hostname)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateHostnameDepth(%q) = nil, want error", tt.hostname)
					return
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("ValidateHostnameDepth(%q) error = %v, want error wrapping %v", tt.hostname, err, tt.wantErrType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateHostnameDepth(%q) = %v, want nil", tt.hostname, err)
				}
			}
		})
	}
}

func TestExtractZoneFromHostname(t *testing.T) {
	tests := []struct {
		hostname string
		want     string
	}{
		{"example.com", "example.com"},
		{"app.example.com", "example.com"},
		{"api.staging.example.com", "example.com"},
		{"example.co.uk", "example.co.uk"},
		{"app.example.co.uk", "example.co.uk"},
		{"selectcode.io", "selectcode.io"},
		{"tunnel-demo.testing.selectcode.io", "selectcode.io"},
		{"cfgate-demo.selectcode.io", "selectcode.io"},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			got := ExtractZoneFromHostname(tt.hostname)
			if got != tt.want {
				t.Errorf("ExtractZoneFromHostname(%q) = %q, want %q", tt.hostname, got, tt.want)
			}
		})
	}
}

func TestIsOwnedByCfgate(t *testing.T) {
	tests := []struct {
		name    string
		record  *DNSRecord
		ownerID string
		want    bool
	}{
		{
			name:   "nil record",
			record: nil,
			want:   false,
		},
		{
			name: "heritage=cfgate without ownerID check",
			record: &DNSRecord{
				Content: "heritage=cfgate,cfgate/owner=cluster-a,cfgate/resource=ns/route",
			},
			ownerID: "",
			want:    true,
		},
		{
			name: "heritage=cfgate with matching ownerID",
			record: &DNSRecord{
				Content: "heritage=cfgate,cfgate/owner=cluster-a,cfgate/resource=ns/route",
			},
			ownerID: "cluster-a",
			want:    true,
		},
		{
			name: "heritage=cfgate with non-matching ownerID",
			record: &DNSRecord{
				Content: "heritage=cfgate,cfgate/owner=cluster-a,cfgate/resource=ns/route",
			},
			ownerID: "cluster-b",
			want:    false,
		},
		{
			name: "comment-based ownership",
			record: &DNSRecord{
				Comment: "managed by cfgate",
			},
			ownerID: "",
			want:    true,
		},
		{
			name: "comment-based ownership ignores ownerID",
			record: &DNSRecord{
				Comment: "managed by cfgate",
			},
			ownerID: "any-owner",
			want:    true,
		},
		{
			name: "not owned",
			record: &DNSRecord{
				Content: "some other content",
				Comment: "some other comment",
			},
			want: false,
		},
		{
			name: "prevent substring false-positive",
			record: &DNSRecord{
				Content: "heritage=cfgate,cfgate/owner=ns/foobar,cfgate/resource=test",
			},
			ownerID: "ns/foo",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsOwnedByCfgate(tt.record, tt.ownerID)
			if got != tt.want {
				t.Errorf("IsOwnedByCfgate() = %v, want %v", got, tt.want)
			}
		})
	}
}
