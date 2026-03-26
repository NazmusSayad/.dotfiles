package helpers

import "testing"

type testCase struct {
	name string
	prev string
	next string
	want string
}

var tests = map[string][]testCase{
	"semantics": {
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
			want: "{\n  \"items\": [\n    4,\n    5\n  ]\n}",
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
			want: "{\n  \"config\": [\n    \"a\",\n    \"b\"\n  ]\n}",
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
			want: "{\n  \"a\": {\n    \"b\": {\n      \"c\": {\n        \"d\": 2,\n        \"e\": [\n          1,\n          2,\n          3\n        ]\n      }\n    }\n  }\n}",
		},
		{
			name: "arrays of objects",
			prev: "{\n  \"items\": [{\n    \"id\": 1\n  }]\n}",
			next: "{\n  \"items\": [{\n    \"id\": 2\n  }, {\n    \"id\": 3\n  }]\n}",
			want: "{\n  \"items\": [\n    {\n      \"id\": 2\n    },\n    {\n      \"id\": 3\n    }\n  ]\n}",
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
			want: "{\n}",
		},
		{
			name: "large object",
			prev: "{\n  \"k1\": 1,\n  \"k2\": 2,\n  \"k3\": 3,\n  \"k4\": 4,\n  \"k5\": 5,\n  \"k6\": 6,\n  \"k7\": 7,\n  \"k8\": 8,\n  \"k9\": 9,\n  \"k10\": 10\n}",
			next: "{\n  \"k1\": 1,\n  \"k2\": 20,\n  \"k3\": 3,\n  \"k4\": 40,\n  \"k5\": 5,\n  \"k6\": 60,\n  \"k7\": 7,\n  \"k8\": 80,\n  \"k9\": 9,\n  \"k10\": 100,\n  \"k11\": 110\n}",
			want: "{\n  \"k1\": 1,\n  \"k2\": 20,\n  \"k3\": 3,\n  \"k4\": 40,\n  \"k5\": 5,\n  \"k6\": 60,\n  \"k7\": 7,\n  \"k8\": 80,\n  \"k9\": 9,\n  \"k10\": 100,\n  \"k11\": 110\n}",
		},
	},
	"formatting": {
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
		{
			name: "preserves unchanged nested key order",
			prev: "{\n  \"openrouter*\": {\n    \"api\": \"https://openrouter.ai/api/v1\",\n    \"name\": \"Openrouter+\",\n    \"models\": {\n      \"A\": {\n        \"id\": \"a\",\n        \"name\": \"A\"\n      },\n      \"B\": {\n        \"id\": \"b\",\n        \"name\": \"B\"\n      }\n    }\n  }\n}",
			next: "{\n  \"openrouter*\": {\n    \"api\": \"https://openrouter.ai/api/v1\",\n    \"name\": \"Openrouter+\",\n    \"models\": {\n      \"A\": {\n        \"id\": \"a\",\n        \"name\": \"A\"\n      }\n    }\n  }\n}",
			want: "{\n  \"openrouter*\": {\n    \"api\": \"https://openrouter.ai/api/v1\",\n    \"name\": \"Openrouter+\",\n    \"models\": {\n      \"A\": {\n        \"id\": \"a\",\n        \"name\": \"A\"\n      }\n    }\n  }\n}",
		},
	},
}

var errorTests = []struct {
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

func TestMergeJSONObject(t *testing.T) {
	for group, cases := range tests {
		for _, tt := range cases {
			t.Run(group+"/"+tt.name, func(t *testing.T) {
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
}
