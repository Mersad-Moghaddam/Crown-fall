package content

import "testing"

func TestEmbeddedContentValidates(t *testing.T) {
	if err := Validate(); err != nil {
		t.Fatal(err)
	}
}
