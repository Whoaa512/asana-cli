package cli

import (
	"reflect"
	"testing"

	"github.com/whoaa512/asana-cli/internal/models"
)

func TestEncodeTypedCustomFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		raw       string
		want      any
		wantErr   bool
	}{
		// clear
		{"empty clears text", "text", "", nil, false},
		{"empty clears number", "number", "", nil, false},
		{"empty clears enum", "enum", "", nil, false},

		// text
		{"text string", "text", "In review", "In review", false},
		{"text digit-only stays string", "text", "42", "42", false},

		// number
		{"number integer", "number", "42", float64(42), false},
		{"number float", "number", "3.14", float64(3.14), false},
		{"number invalid", "number", "notanumber", nil, true},
		{"number GID-like stays float", "number", "1213529490866717", float64(1213529490866717), false},

		// enum — GID must be sent as string
		{"enum GID string", "enum", "1213529490866717", "1213529490866717", false},

		// multi_enum — comma-separated GIDs become []string
		{"multi_enum single", "multi_enum", "111", []string{"111"}, false},
		{"multi_enum multiple", "multi_enum", "111,222,333", []string{"111", "222", "333"}, false},

		// people — comma-separated user GIDs become []string
		{"people single", "people", "111", []string{"111"}, false},
		{"people multiple", "people", "111,222", []string{"111", "222"}, false},

		// date
		{"date", "date", "2024-01-15", "2024-01-15", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &models.CustomField{GID: "123", Name: "Test Field", Type: tt.fieldType}
			got, err := encodeTypedCustomFieldValue(field, tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestSplitCustomFieldPairs(t *testing.T) {
	tests := []struct {
		name    string
		pairs   []string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "nil returns nil",
			pairs: nil,
			want:  map[string]string{},
		},
		{
			name:  "single pair",
			pairs: []string{"111=foo"},
			want:  map[string]string{"111": "foo"},
		},
		{
			name:  "multiple pairs",
			pairs: []string{"111=foo", "222=bar"},
			want:  map[string]string{"111": "foo", "222": "bar"},
		},
		{
			name:  "empty value clears field",
			pairs: []string{"111="},
			want:  map[string]string{"111": ""},
		},
		{
			name:  "value with equals sign",
			pairs: []string{"111=foo=bar"},
			want:  map[string]string{"111": "foo=bar"},
		},
		{
			name:    "missing equals",
			pairs:   []string{"111"},
			wantErr: true,
		},
		{
			name:    "empty GID",
			pairs:   []string{"=value"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitCustomFieldPairs(tt.pairs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.pairs == nil {
				if len(got) != 0 {
					t.Errorf("expected empty map for nil input, got %v", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskCustomFieldCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range taskCmd.Commands() {
		if cmd.Use == "custom-field" {
			found = true
			break
		}
	}
	if !found {
		t.Error("task custom-field command should be registered")
	}
}

func TestTaskCustomFieldSubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list":                          false,
		"set <task> <field-gid> <value>": false,
	}
	for _, cmd := range taskCFCmd.Commands() {
		if _, ok := subcommands[cmd.Use]; ok {
			subcommands[cmd.Use] = true
		}
	}
	for use, found := range subcommands {
		if !found {
			t.Errorf("subcommand %q should be registered", use)
		}
	}
}
