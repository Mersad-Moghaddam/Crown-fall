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
	if Commitment(seed) == Commitment([]byte("wrong-seed")) {
		t.Fatal("wrong seed matched commitment")
	}
}

func TestCryptoSourceProducesIndependentSeeds(t *testing.T) {
	first, err := (CryptoSource{}).Seed()
	if err != nil {
		t.Fatal(err)
	}
	second, err := (CryptoSource{}).Seed()
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != SeedSize || len(second) != SeedSize || string(first) == string(second) {
		t.Fatal("cryptographic source did not produce independent fixed-size seeds")
	}
}
