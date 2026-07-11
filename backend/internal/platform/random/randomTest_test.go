package random

import "testing"

func TestCommitmentAndDomains(t *testing.T) {
	seed := []byte("fixed-seed")
	if Commitment(seed) != Commitment(seed) {
		t.Fatal("commitment must be deterministic")
	}
	if string(Derive(seed, "roles")) == string(Derive(seed, "sigils")) {
		t.Fatal("domain streams must be separated")
	}
}
