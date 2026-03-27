package impl_test

import (
	"testing"

	music "github.com/xmtlzzz/vMusic/apps/music/impl"
	"github.com/xmtlzzz/vMusic/apps/user/impl"
)

func TestLikeListScan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		{name: "null", input: nil, expected: 0},
		{name: "empty array", input: []byte("[]"), expected: 0},
		{name: "empty object", input: []byte("{}"), expected: 0},
		{name: "single object", input: []byte(`{"music_name":"song","file_name":"song.mp3"}`), expected: 1},
		{name: "array", input: []byte(`[{"music_name":"song","file_name":"song.mp3"}]`), expected: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var list impl.LikeList
			if err := (&list).Scan(test.input); err != nil {
				t.Fatalf("scan failed: %v", err)
			}
			if len(list) != test.expected {
				t.Fatalf("expected %d items, got %d", test.expected, len(list))
			}
		})
	}
}

func TestLikeListValue(t *testing.T) {
	list := impl.LikeList{
		&music.Objector{MusicName: "song", FileName: "song.mp3"},
	}

	value, err := list.Value()
	if err != nil {
		t.Fatalf("value failed: %v", err)
	}

	if value == nil {
		t.Fatal("expected json value")
	}
}
