package discover

import (
	"errors"
	"reflect"
	"testing"
)

func TestConfigParse(t *testing.T) {
	tests := []struct {
		s   string
		c   Config
		err error
	}{
		{"", nil, nil},
		{"  ", nil, nil},
		{"provider=aws foo", nil, errors.New(`invalid format: foo`)},
		{"project_name=Test zone_pattern=us-(?west|east).%2b tag_value=consul+server credentials_file=xxx",
			map[string]string{
				"project_name":     "Test",
				"zone_pattern":     "us-(?west|east).+",
				"tag_value":        "consul server",
				"credentials_file": "xxx",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			c, err := Parse(tt.s)
			if got, want := err, tt.err; !reflect.DeepEqual(got, want) {
				t.Fatalf("got error %v want %v", got, want)
			}
			if got, want := c, tt.c; !reflect.DeepEqual(got, want) {
				t.Fatalf("got config %#v want %#v", got, want)
			}
		})
	}
}

func TestConfigString(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"   ", ""},
		{"b=c a=b", "a=b b=c"},
		{"a=b provider=foo x=y", "provider=foo a=b x=y"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			c, err := Parse(tt.in)
			if err != nil {
				t.Fatal("Parse failed: ", err)
			}
			if got, want := c.String(), tt.out; got != want {
				t.Fatalf("got %q want %q", got, want)
			}
		})
	}
}
