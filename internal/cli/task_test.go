package cli

import (
	"testing"
)

func TestParseCustomFields(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  nil,
		},
		{
			name:  "single field",
			input: []string{"12345=hello"},
			want:  map[string]string{"12345": "hello"},
		},
		{
			name:  "multiple fields",
			input: []string{"12345=hello", "67890=world"},
			want:  map[string]string{"12345": "hello", "67890": "world"},
		},
		{
			name:  "value with equals sign",
			input: []string{"12345=a=b"},
			want:  map[string]string{"12345": "a=b"},
		},
		{
			name:  "empty value",
			input: []string{"12345="},
			want:  map[string]string{"12345": ""},
		},
		{
			name:    "missing equals",
			input:   []string{"12345"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCustomFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseCustomFields() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.want == nil && got != nil {
				t.Fatalf("expected nil, got %v", got)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("got[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
