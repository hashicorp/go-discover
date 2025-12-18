// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package discover

import (
	"errors"
	"reflect"
	"testing"
)

func TestConfigParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s   string
		c   Config
		err error
	}{
		// happy flows
		{``, nil, nil},
		{`key=a`, Config{"key": "a"}, nil},
		{`key=a key2=b`, Config{"key": "a", "key2": "b"}, nil},
		{`key=a+b key2=c/d`, Config{"key": "a+b", "key2": "c/d"}, nil},
		{` key=a    key2=b `, Config{"key": "a", "key2": "b"}, nil},
		{` key = a   key2 = b `, Config{"key": "a", "key2": "b"}, nil},
		{`  "k e \\\" y" = "a \" b" key2=c`, Config{`k e \" y`: `a " b`, "key2": "c"}, nil},
		{`secret_access_key="fpOfcHQJAQBczjAxiVpeyLmX1M0M0KPBST+GU2GvEN4="`, Config{"secret_access_key": "fpOfcHQJAQBczjAxiVpeyLmX1M0M0KPBST+GU2GvEN4="}, nil},

		{`provider=aws foo`, nil, errors.New(`foo: missing '='`)},
		{
			`project_name=Test zone_pattern=us-(?west|east).+ tag_value="consul server" credentials_file=xxx`,
			Config{
				"project_name":     "Test",
				"zone_pattern":     "us-(?west|east).+",
				"tag_value":        "consul server",
				"credentials_file": "xxx",
			},
			nil,
		},

		// errors
		{`key`, nil, errors.New(`key: missing '='`)},
		{`key=`, nil, errors.New(`key: missing value`)},
		{`key="a`, nil, errors.New(`key: unbalanced quotes`)},
		{`key="\`, nil, errors.New(`key: unterminated escape sequence`)},
		{`key=a key=b`, nil, errors.New(`key: duplicate key`)},
		{`key key2`, nil, errors.New(`key: missing '='`)},
		{`secret_access_key=fpOfcHQJAQBczjAxiVpeyLmX1M0M0KPBST+GU2GvEN4=`, nil, errors.New(`secret_access_key: - equals in key's value, enclosing double-quote needed secret_access_key="value-with-=-symbol"`)},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			t.Parallel()
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

// TestConfigString verifies that String() returns config keys in a stable
// order.
func TestConfigString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, out string
	}{
		{``, ``},
		{`   `, ``},
		{`b=c "a a"="b b"`, `"a a"="b b" b=c`},
		{`a=b provider=foo x=y`, `provider=foo a=b x=y`},
		{`provider=k8s namespace=vault label_selector="app.kubernetes.io/name=vault,component=server"`, `provider=k8s label_selector="app.kubernetes.io/name=vault,component=server" namespace=vault`},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
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

// TestConfigRoundTrip verifies that an input string can be parsed and formatted
// back into a parseable string.
func TestConfigRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []string{
		``,
		`"a a"="b b" b=c`,
		`provider=foo a=b x=y`,
		`key=a`,
		`key=a key2=b`,
		`key=a+b key2=c/d`,
		`key=a key2=b`,
		`"k e \\\" y"="a \" b" key2=c`,
		`secret_access_key="fpOfcHQJAQBczjAxiVpeyLmX1M0M0KPBST+GU2GvEN4="`,
		`credentials_file=xxx project_name=Test tag_value="consul server" zone_pattern=us-(?west|east).+`,
		`provider=k8s label_selector="app.kubernetes.io/name=vault,component=server" namespace=vault`,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			c, err := Parse(tt)
			if err != nil {
				t.Fatal("Initial Parse() failed: ", err)
			}
			got, err := Parse(c.String())
			if err != nil {
				t.Fatal("Round trip Parse() failed: ", err)
			}
			if gs := got.String(); gs != tt {
				t.Fatalf("got %q want %q", gs, tt)
			}
		})
	}
}
