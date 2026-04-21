package helpers

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMergeJSONObjectCoreExact(t *testing.T) {
	tests := []struct {
		name string
		prev string
		next string
		want string
	}{
		{
			name: "unchanged multiline passthrough",
			prev: "{\n  \"a\": 1,\n  \"b\": [1, 2],\n  \"c\": {\n    \"x\": true\n  }\n}",
			next: "{\n  \"a\": 1,\n  \"b\": [1, 2],\n  \"c\": {\n    \"x\": true\n  }\n}",
			want: "{\n  \"a\": 1,\n  \"b\": [1, 2],\n  \"c\": {\n    \"x\": true\n  }\n}",
		},
		{
			name: "unchanged compact passthrough",
			prev: "{\"a\":1,\"b\":2}",
			next: "{\"a\":1,\"b\":2}",
			want: "{\"a\":1,\"b\":2}",
		},
		{
			name: "add key compact",
			prev: "{\"a\":1}",
			next: "{\"a\":1,\"b\":2}",
			want: "{\"a\":1,\"b\":2}",
		},
		{
			name: "update key compact",
			prev: "{\"a\":1,\"b\":2}",
			next: "{\"a\":1,\"b\":3}",
			want: "{\"a\":1,\"b\":3}",
		},
		{
			name: "remove key compact",
			prev: "{\"a\":1,\"b\":2}",
			next: "{\"a\":1}",
			want: "{\"a\":1}",
		},
		{
			name: "replace array compact",
			prev: "{\"items\":[1,2,3]}",
			next: "{\"items\":[4,5]}",
			want: "{\"items\":[4,5]}",
		},
		{
			name: "single value change keeps surrounding formatting",
			prev: "{\n  \"a\": 1,\n\n  \"b\": 2,\n  \"c\": 3\n}",
			next: "{\n  \"a\": 1,\n\n  \"b\": 9,\n  \"c\": 3\n}",
			want: "{\n  \"a\": 1,\n\n  \"b\": 9,\n  \"c\": 3\n}",
		},
		{
			name: "remove key keeps remaining formatting",
			prev: "{\n  \"a\": 1,\n\n  \"b\": 2,\n  \"c\": 3\n}",
			next: "{\n  \"a\": 1,\n  \"c\": 3\n}",
			want: "{\n  \"a\": 1,\n  \"c\": 3\n}",
		},
		{
			name: "add key keeps existing formatting",
			prev: "{\n  \"a\": 1,\n  \"b\": 2\n}",
			next: "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}",
			want: "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeJSONObject(tt.prev, tt.next)
			if err != nil {
				t.Fatalf("MergeJSONObject() error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("string mismatch\nwant:\n%s\n\ngot:\n%s", tt.want, got)
			}
		})
	}
}

func TestMergeJSONObjectNoChangePassthrough(t *testing.T) {
	for i := 0; i < 120; i++ {
		t.Run(fmt.Sprintf("case-%03d", i), func(t *testing.T) {
			prev := noChangeInput(i)
			next := prev

			got, err := MergeJSONObject(prev, next)
			if err != nil {
				t.Fatalf("MergeJSONObject() error = %v", err)
			}

			if got != next {
				t.Fatalf("only target value should change\nwant:\n%s\n\ngot:\n%s", next, got)
			}

			if got != prev {
				t.Fatalf("unchanged object should be byte-identical\nwant:\n%s\n\ngot:\n%s", prev, got)
			}
		})
	}
}

func TestMergeJSONObjectUnchangedFieldsRemainRaw(t *testing.T) {
	for i := 0; i < 40; i++ {
		t.Run(fmt.Sprintf("case-%03d", i), func(t *testing.T) {
			prev := fmt.Sprintf("{\n  \"keepNum\": %d,\n  \"keepArr\": [1, 2, %d],\n  \"target\": %d,\n  \"keepText\": \"line-%d\"\n}", i, i%5, i, i)
			next := fmt.Sprintf("{\n  \"keepNum\": %d,\n  \"keepArr\": [1, 2, %d],\n  \"target\": %d,\n  \"keepText\": \"line-%d\"\n}", i, i%5, i+1000, i)

			got, err := MergeJSONObject(prev, next)
			if err != nil {
				t.Fatalf("MergeJSONObject() error = %v", err)
			}

			if got != next {
				t.Fatalf("only nested target value should change\nwant:\n%s\n\ngot:\n%s", next, got)
			}

			prevByKey := mustEntriesByKey(t, prev)
			gotByKey := mustEntriesByKey(t, got)
			nextByKey := mustEntriesByKey(t, next)

			for _, key := range []string{"keepNum", "keepArr", "keepText"} {
				if !bytes.Equal(prevByKey[key].Val, gotByKey[key].Val) {
					t.Fatalf("unchanged value mutated for key %s\nwant: %s\ngot: %s", key, string(prevByKey[key].Val), string(gotByKey[key].Val))
				}
				if prevByKey[key].KeyRaw != gotByKey[key].KeyRaw {
					t.Fatalf("unchanged key mutated for key %s\nwant: %s\ngot: %s", key, prevByKey[key].KeyRaw, gotByKey[key].KeyRaw)
				}
			}

			if !jsonEqual(gotByKey["target"].Val, nextByKey["target"].Val) {
				t.Fatalf("target key did not update correctly\nwant: %s\ngot: %s", string(nextByKey["target"].Val), string(gotByKey["target"].Val))
			}
		})
	}
}

func TestMergeJSONObjectNestedUnchangedFieldsRemainRaw(t *testing.T) {
	for i := 0; i < 24; i++ {
		t.Run(fmt.Sprintf("case-%03d", i), func(t *testing.T) {
			prev := fmt.Sprintf("{\n  \"outer\": {\n    \"keepA\": \"v-%d\",\n    \"target\": %d,\n    \"keepB\": [1, 2, %d]\n  },\n  \"untouched\": true\n}", i, i, i%7)
			next := fmt.Sprintf("{\n  \"outer\": {\n    \"keepA\": \"v-%d\",\n    \"target\": %d,\n    \"keepB\": [1, 2, %d]\n  },\n  \"untouched\": true\n}", i, i+500, i%7)

			got, err := MergeJSONObject(prev, next)
			if err != nil {
				t.Fatalf("MergeJSONObject() error = %v", err)
			}

			prevRoot := mustEntriesByKey(t, prev)
			gotRoot := mustEntriesByKey(t, got)
			nextRoot := mustEntriesByKey(t, next)

			if !bytes.Equal(prevRoot["untouched"].Val, gotRoot["untouched"].Val) {
				t.Fatalf("untouched root key mutated\nwant: %s\ngot: %s", string(prevRoot["untouched"].Val), string(gotRoot["untouched"].Val))
			}

			prevOuter := mustEntriesByKey(t, string(prevRoot["outer"].Val))
			gotOuter := mustEntriesByKey(t, string(gotRoot["outer"].Val))
			nextOuter := mustEntriesByKey(t, string(nextRoot["outer"].Val))

			for _, key := range []string{"keepA", "keepB"} {
				if !bytes.Equal(prevOuter[key].Val, gotOuter[key].Val) {
					t.Fatalf("nested unchanged value mutated for key %s\nwant: %s\ngot: %s", key, string(prevOuter[key].Val), string(gotOuter[key].Val))
				}
			}

			if !jsonEqual(gotOuter["target"].Val, nextOuter["target"].Val) {
				t.Fatalf("nested target key did not update correctly\nwant: %s\ngot: %s", string(nextOuter["target"].Val), string(gotOuter["target"].Val))
			}
		})
	}
}

func TestMergeJSONObjectErrors(t *testing.T) {
	tests := []struct {
		name string
		prev string
		next string
	}{
		{name: "invalid previous json", prev: "{", next: "{}"},
		{name: "invalid next json", prev: "{}", next: "{"},
		{name: "previous is array", prev: "[]", next: "{}"},
		{name: "next is array", prev: "{}", next: "[]"},
		{name: "previous is scalar", prev: "1", next: "{}"},
		{name: "next is scalar", prev: "{}", next: "1"},
		{name: "previous empty string", prev: "", next: "{}"},
		{name: "next empty string", prev: "{}", next: ""},
		{name: "previous whitespace only", prev: "   \n\t", next: "{}"},
		{name: "next whitespace only", prev: "{}", next: "   \n\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := MergeJSONObject(tt.prev, tt.next); err == nil {
				t.Fatalf("MergeJSONObject() expected error")
			}
		})
	}
}

func mustEntriesByKey(t *testing.T, raw string) map[string]objectEntry {
	t.Helper()

	entries, err := parseObjectEntries([]byte(raw))
	if err != nil {
		t.Fatalf("parseObjectEntries() error = %v\ninput:\n%s", err, raw)
	}

	byKey := make(map[string]objectEntry, len(entries))
	for _, entry := range entries {
		byKey[entry.Key] = entry
	}

	return byKey
}

func noChangeInput(i int) string {
	switch i % 3 {
	case 0:
		return fmt.Sprintf("{\n  \"id\": %d,\n  \"enabled\": true,\n  \"meta\": {\n    \"name\": \"n-%d\",\n    \"arr\": [1, 2, %d]\n  }\n}", i, i, i%7)
	case 1:
		return fmt.Sprintf("{\"id\":%d,\"enabled\":true,\"meta\":{\"name\":\"n-%d\",\"arr\":[1,2,%d]}}", i, i, i%7)
	default:
		return fmt.Sprintf("{\n\t\"id\": %d,\n\t\"enabled\": false,\n\t\"meta\": {\n\t\t\"name\": \"t-%d\",\n\t\t\"arr\": [3, 4, %d]\n\t}\n}", i, i, i%7)
	}
}
