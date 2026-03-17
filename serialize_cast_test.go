package protocache

import "testing"

func TestBytesWordsShareBackingStore(t *testing.T) {
	words := []uint32{1, 2, 3, 4}
	raw := WordsToBytes(words)
	if len(raw) != len(words)*4 {
		t.Fatalf("unexpected raw size: %d", len(raw))
	}

	roundtrip := BytesToWords(raw)
	if len(roundtrip) != len(words) {
		t.Fatalf("unexpected roundtrip size: %d", len(roundtrip))
	}

	roundtrip[1] = 99
	if words[1] != 99 {
		t.Fatalf("expected zero-copy roundtrip, got %d", words[1])
	}

	putUint32(raw[8:], 123)
	if roundtrip[2] != 123 {
		t.Fatalf("expected shared backing store, got %d", roundtrip[2])
	}
}
