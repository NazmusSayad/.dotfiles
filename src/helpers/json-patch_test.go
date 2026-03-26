package helpers

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestMergeJSONObjectSemantics(t *testing.T) {
	tests := []struct {
		name string
		prev string
		next string
		want string
	}{
		{
			name: "add and update top level keys",
			prev: "{\n  \"name\": \"old\",\n  \"enabled\": false\n}",
			next: "{\n  \"name\": \"new\",\n  \"enabled\": false,\n  \"count\": 3\n}",
			want: "{\n  \"name\": \"new\",\n  \"enabled\": false,\n  \"count\": 3\n}",
		},
		{
			name: "preserve untouched keys",
			prev: "{\n  \"a\": 1,\n  \"b\": 2\n}",
			next: "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}",
			want: "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}",
		},
		{
			name: "replace nested object fields",
			prev: "{\n  \"service\": {\n    \"host\": \"localhost\",\n    \"port\": 8080\n  }\n}",
			next: "{\n  \"service\": {\n    \"host\": \"example.com\",\n    \"port\": 9090,\n    \"tls\": true\n  }\n}",
			want: "{\n  \"service\": {\n    \"host\": \"example.com\",\n    \"port\": 9090,\n    \"tls\": true\n  }\n}",
		},
		{
			name: "replace arrays",
			prev: "{\n  \"items\": [1, 2, 3]\n}",
			next: "{\n  \"items\": [4, 5]\n}",
			want: "{\n  \"items\": [4, 5]\n}",
		},
		{
			name: "replace scalar with object",
			prev: "{\n  \"config\": true\n}",
			next: "{\n  \"config\": {\n    \"enabled\": true\n  }\n}",
			want: "{\n  \"config\": {\n    \"enabled\": true\n  }\n}",
		},
		{
			name: "replace object with scalar",
			prev: "{\n  \"config\": {\n    \"enabled\": true\n  }\n}",
			next: "{\n  \"config\": false\n}",
			want: "{\n  \"config\": false\n}",
		},
		{
			name: "replace object with array",
			prev: "{\n  \"config\": {\n    \"enabled\": true\n  }\n}",
			next: "{\n  \"config\": [\"a\", \"b\"]\n}",
			want: "{\n  \"config\": [\"a\", \"b\"]\n}",
		},
		{
			name: "replace array with object",
			prev: "{\n  \"config\": [\"a\", \"b\"]\n}",
			next: "{\n  \"config\": {\n    \"enabled\": true\n  }\n}",
			want: "{\n  \"config\": {\n    \"enabled\": true\n  }\n}",
		},
		{
			name: "explicit null removes key",
			prev: "{\n  \"a\": 1,\n  \"b\": 2\n}",
			next: "{\n  \"a\": 1,\n  \"b\": null\n}",
			want: "{\n  \"a\": 1\n}",
		},
		{
			name: "omission removes key",
			prev: "{\n  \"a\": 1,\n  \"b\": 2\n}",
			next: "{\n  \"a\": 1\n}",
			want: "{\n  \"a\": 1\n}",
		},
		{
			name: "escaped strings are preserved",
			prev: "{\n  \"text\": \"hello\\nworld\"\n}",
			next: "{\n  \"text\": \"quote: \\\"hi\\\" and slash: \\\\server\"\n}",
			want: "{\n  \"text\": \"quote: \\\"hi\\\" and slash: \\\\server\"\n}",
		},
		{
			name: "unicode strings are preserved",
			prev: "{\n  \"text\": \"hello\"\n}",
			next: "{\n  \"text\": \"hello 世界 😀\"\n}",
			want: "{\n  \"text\": \"hello 世界 😀\"\n}",
		},
		{
			name: "deeply nested objects",
			prev: "{\n  \"a\": {\n    \"b\": {\n      \"c\": {\n        \"d\": 1\n      }\n    }\n  }\n}",
			next: "{\n  \"a\": {\n    \"b\": {\n      \"c\": {\n        \"d\": 2,\n        \"e\": [1, 2, 3]\n      }\n    }\n  }\n}",
			want: "{\n  \"a\": {\n    \"b\": {\n      \"c\": {\n        \"d\": 2,\n        \"e\": [1, 2, 3]\n      }\n    }\n  }\n}",
		},
		{
			name: "arrays of objects",
			prev: "{\n  \"items\": [{\n    \"id\": 1\n  }]\n}",
			next: "{\n  \"items\": [{\n    \"id\": 2\n  }, {\n    \"id\": 3\n  }]\n}",
			want: "{\n  \"items\": [{\n    \"id\": 2\n  }, {\n    \"id\": 3\n  }]\n}",
		},
		{
			name: "empty previous object",
			prev: "{}",
			next: "{\n  \"name\": \"value\"\n}",
			want: "{\n  \"name\": \"value\"\n}",
		},
		{
			name: "empty next object clears object",
			prev: "{\n  \"name\": \"value\"\n}",
			next: "{}",
			want: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeJSONObject(tt.prev, tt.next)
			if err != nil {
				t.Fatalf("MergeJSONObject() error = %v", err)
			}

			assertJSONObjectEqual(t, got, tt.want)
		})
	}
}

func TestMergeJSONObjectFormatting(t *testing.T) {
	tests := []struct {
		name string
		prev string
		next string
		want string
	}{
		{
			name: "preserves two space indentation",
			prev: "{\n  \"a\": 1\n}",
			next: "{\n  \"a\": 2,\n  \"b\": {\n    \"c\": true\n  }\n}",
			want: "{\n  \"a\": 2,\n  \"b\": {\n    \"c\": true\n  }\n}",
		},
		{
			name: "preserves four space indentation",
			prev: "{\n    \"a\": 1\n}",
			next: "{\n    \"a\": 2,\n    \"b\": {\n        \"c\": true\n    }\n}",
			want: "{\n    \"a\": 2,\n    \"b\": {\n        \"c\": true\n    }\n}",
		},
		{
			name: "single line previous uses default indentation",
			prev: "{\"a\":1}",
			next: "{\"a\":2,\"b\":3}",
			want: "{\n  \"a\": 2,\n  \"b\": 3\n}",
		},
		{
			name: "empty object formatting",
			prev: "{}",
			next: "{}",
			want: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeJSONObject(tt.prev, tt.next)
			if err != nil {
				t.Fatalf("MergeJSONObject() error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("MergeJSONObject() formatting mismatch\nwant:\n%s\n\ngot:\n%s", tt.want, got)
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
			got, err := MergeJSONObject(tt.prev, tt.next)
			if err == nil {
				t.Fatalf("MergeJSONObject() error = nil, got = %s", got)
			}
		})
	}
}

func TestMergeJSONObjectIdempotency(t *testing.T) {
	prev := "{\n  \"alpha\": 1,\n  \"nested\": {\n    \"flag\": true\n  }\n}"
	next := "{\n  \"alpha\": 2,\n  \"nested\": {\n    \"flag\": true,\n    \"extra\": \"x\"\n  },\n  \"items\": [1, 2]\n}"

	first, err := MergeJSONObject(prev, next)
	if err != nil {
		t.Fatalf("first MergeJSONObject() error = %v", err)
	}

	second, err := MergeJSONObject(first, next)
	if err != nil {
		t.Fatalf("second MergeJSONObject() error = %v", err)
	}

	assertJSONObjectEqual(t, first, second)

	if first != second {
		t.Fatalf("MergeJSONObject() is not idempotent\nfirst:\n%s\n\nsecond:\n%s", first, second)
	}
}

func TestMergeJSONObjectLargeObject(t *testing.T) {
	prev := "{\n  \"k1\": 1,\n  \"k2\": 2,\n  \"k3\": 3,\n  \"k4\": 4,\n  \"k5\": 5,\n  \"k6\": 6,\n  \"k7\": 7,\n  \"k8\": 8,\n  \"k9\": 9,\n  \"k10\": 10\n}"
	next := "{\n  \"k1\": 1,\n  \"k2\": 20,\n  \"k3\": 3,\n  \"k4\": 40,\n  \"k5\": 5,\n  \"k6\": 60,\n  \"k7\": 7,\n  \"k8\": 80,\n  \"k9\": 9,\n  \"k10\": 100,\n  \"k11\": 110\n}"

	got, err := MergeJSONObject(prev, next)
	if err != nil {
		t.Fatalf("MergeJSONObject() error = %v", err)
	}

	assertJSONObjectEqual(t, got, next)
}

func assertJSONObjectEqual(t *testing.T, got string, want string) {
	t.Helper()

	var gotValue map[string]any
	if err := json.Unmarshal([]byte(got), &gotValue); err != nil {
		t.Fatalf("failed to unmarshal got JSON: %v\n%s", err, got)
	}

	var wantValue map[string]any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("failed to unmarshal want JSON: %v\n%s", err, want)
	}

	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("JSON mismatch\nwant:\n%s\n\ngot:\n%s", prettyJSON(t, want), prettyJSON(t, got))
	}

	if !json.Valid([]byte(got)) {
		t.Fatalf("got invalid JSON: %s", got)
	}
}

func prettyJSON(t *testing.T, input string) string {
	t.Helper()

	var buffer bytes.Buffer
	if err := json.Indent(&buffer, []byte(input), "", "  "); err != nil {
		return input
	}

	return buffer.String()
}
