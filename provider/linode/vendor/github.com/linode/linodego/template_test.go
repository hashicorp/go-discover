// +build ignore

package linodego_test

import (
	"context"

	. "github.com/linode/linodego"

	"testing"
)

func TestGetTemplate_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetTemplate_missing")
	defer teardown()

	i, err := client.GetTemplate(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing template, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing template, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing template, got %v", e.Code)
	}
}

func TestGetTemplate_found(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetTemplate_found")
	defer teardown()

	i, err := client.GetTemplate(context.Background(), "linode/ubuntu16.04lts")
	if err != nil {
		t.Errorf("Error getting template, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "linode/ubuntu16.04lts" {
		t.Errorf("Expected a specific template, but got a different one %v", i)
	}
}
func TestListTemplates(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListTemplates")
	defer teardown()

	i, err := client.ListTemplates(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing templates, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of templates, but got none %v", i)
	}
}
