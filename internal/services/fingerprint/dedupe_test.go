package fingerprint

import "testing"

func TestDecideFingerprintMatchPrefersExternalID(t *testing.T) {
	got := decideFingerprintMatch("abc123", "hexid")
	if got.MatchBy != "fingerprint_id" || got.ExistingHexID != "abc123" {
		t.Fatalf("expected fingerprint_id match, got %+v", got)
	}
}

func TestDecideFingerprintMatchFallsBackToName(t *testing.T) {
	got := decideFingerprintMatch("", "hexid")
	if got.MatchBy != "_id" || got.ExistingHexID != "hexid" {
		t.Fatalf("expected _id match for same-name legacy doc, got %+v", got)
	}
}

func TestDecideFingerprintMatchInsertsWhenUnknown(t *testing.T) {
	got := decideFingerprintMatch("", "")
	if got.MatchBy != "" {
		t.Fatalf("expected insert, got %+v", got)
	}
}

func TestShouldReuseFingerprintID(t *testing.T) {
	if shouldReuseFingerprintID("") {
		t.Fatal("empty match must not reuse an id")
	}
	if !shouldReuseFingerprintID("hexid") {
		t.Fatal("existing doc should be reused instead of creating a duplicate")
	}
}
