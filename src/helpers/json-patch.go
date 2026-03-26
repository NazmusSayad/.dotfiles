package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
)

func MergeJSONObject(prev string, next string) (string, error) {
	prevEntries, err := parseObjectEntries([]byte(prev))
	if err != nil {
		return "", err
	}

	nextEntries, err := parseObjectEntries([]byte(next))
	if err != nil {
		return "", err
	}

	patch, err := jsonpatch.CreateMergePatch([]byte(prev), []byte(next))
	if err != nil {
		return "", err
	}

	merged, err := jsonpatch.MergePatch([]byte(prev), patch)
	if err != nil {
		return "", err
	}

	mergedEntries, err := parseObjectEntries(merged)
	if err != nil {
		return "", err
	}

	mergedByKey := map[string]objectEntry{}
	for _, entry := range mergedEntries {
		mergedByKey[entry.Key] = entry
	}

	prevByKey := map[string]objectEntry{}
	for _, entry := range prevEntries {
		prevByKey[entry.Key] = entry
	}

	ordered := []objectEntry{}
	seen := map[string]bool{}

	for _, entry := range prevEntries {
		mergedEntry, ok := mergedByKey[entry.Key]
		if !ok {
			continue
		}

		if jsonEqual(entry.Val, mergedEntry.Val) {
			ordered = append(ordered, entry)
			seen[entry.Key] = true
			continue
		}

		ordered = append(ordered, mergedEntry)
		seen[entry.Key] = true
	}

	for _, entry := range nextEntries {
		if seen[entry.Key] {
			continue
		}

		mergedEntry, ok := mergedByKey[entry.Key]
		if !ok {
			continue
		}

		ordered = append(ordered, mergedEntry)
		seen[entry.Key] = true
	}

	indent := detectIndent(prev)
	closingIndent := detectClosingIndent(prev)
	indentUnit := detectIndentUnit(indent, closingIndent)
	separator := detectSeparator(prev)

	if len(ordered) == 0 {
		if stringsHasNewline(prev) {
			return "{\n" + closingIndent + "}", nil
		}

		return "{}", nil
	}

	var out bytes.Buffer
	out.WriteByte('{')
	out.WriteByte('\n')

	for i, entry := range ordered {
		formatted, err := indentValue(entry.Val, indent, indentUnit, prevByKey[entry.Key].Val)
		if err != nil {
			return "", err
		}

		if i > 0 {
			out.WriteString(separator)
		}

		out.WriteString(indent)
		out.WriteString(entry.KeyRaw)
		out.WriteString(": ")
		out.WriteString(formatted)
	}

	out.WriteByte('\n')
	out.WriteString(closingIndent)
	out.WriteByte('}')

	return out.String(), nil
}

type objectEntry struct {
	Key    string
	KeyRaw string
	Val    []byte
}

func parseObjectEntries(raw []byte) ([]objectEntry, error) {
	i := skipSpace(raw, 0)
	if i >= len(raw) || raw[i] != '{' {
		return nil, fmt.Errorf("expected object")
	}

	i++
	entries := []objectEntry{}

	for {
		i = skipSpace(raw, i)
		if i >= len(raw) {
			return nil, fmt.Errorf("unexpected end of object")
		}

		if raw[i] == '}' {
			return entries, nil
		}

		keyStart := i
		keyEnd, err := consumeString(raw, i)
		if err != nil {
			return nil, err
		}

		var key string
		if err := json.Unmarshal(raw[keyStart:keyEnd], &key); err != nil {
			return nil, err
		}

		i = skipSpace(raw, keyEnd)
		if i >= len(raw) || raw[i] != ':' {
			return nil, fmt.Errorf("expected colon after key")
		}

		i++
		i = skipSpace(raw, i)
		valueStart := i
		valueEnd, err := consumeValue(raw, i)
		if err != nil {
			return nil, err
		}

		entries = append(entries, objectEntry{
			Key:    key,
			KeyRaw: string(raw[keyStart:keyEnd]),
			Val:    append([]byte{}, raw[valueStart:valueEnd]...),
		})

		i = skipSpace(raw, valueEnd)
		if i >= len(raw) {
			return nil, fmt.Errorf("unexpected end of object")
		}

		if raw[i] == ',' {
			i++
			continue
		}

		if raw[i] == '}' {
			return entries, nil
		}

		return nil, fmt.Errorf("expected comma or closing brace")
	}
}

func indentValue(raw []byte, prefix string, indentUnit string, prev []byte) (string, error) {
	if len(prev) > 0 && jsonEqual(prev, raw) {
		return string(prev), nil
	}

	var formatted bytes.Buffer
	if err := json.Indent(&formatted, raw, "", indentUnit); err == nil {
		return strings.ReplaceAll(formatted.String(), "\n", "\n"+prefix), nil
	}

	return string(raw), nil
}

func detectIndentUnit(indent string, closingIndent string) string {
	if strings.HasPrefix(indent, closingIndent) && len(indent) > len(closingIndent) {
		return indent[len(closingIndent):]
	}

	if indent != "" {
		return indent
	}

	return "  "
}

func detectSeparator(input string) string {
	if strings.Contains(input, ",\n\n") {
		return ",\n\n"
	}

	return ",\n"
}

func jsonEqual(left []byte, right []byte) bool {
	var leftValue any
	if err := json.Unmarshal(left, &leftValue); err != nil {
		return bytes.Equal(left, right)
	}

	var rightValue any
	if err := json.Unmarshal(right, &rightValue); err != nil {
		return bytes.Equal(left, right)
	}

	leftJSON, err := json.Marshal(leftValue)
	if err != nil {
		return bytes.Equal(left, right)
	}

	rightJSON, err := json.Marshal(rightValue)
	if err != nil {
		return bytes.Equal(left, right)
	}

	return bytes.Equal(leftJSON, rightJSON)
}

func detectIndent(input string) string {
	raw := []byte(input)
	newline := bytes.IndexByte(raw, '\n')
	if newline == -1 {
		return "  "
	}

	i := newline + 1
	start := i
	for i < len(raw) && raw[i] == ' ' {
		i++
	}

	if i == start {
		return "  "
	}

	return string(raw[start:i])
}

func detectClosingIndent(input string) string {
	raw := []byte(input)
	newline := bytes.LastIndexByte(raw, '\n')
	if newline == -1 {
		return ""
	}

	i := newline + 1
	start := i
	for i < len(raw) && raw[i] == ' ' {
		i++
	}

	return string(raw[start:i])
}

func stringsHasNewline(input string) bool {
	return bytes.IndexByte([]byte(input), '\n') != -1
}

func skipSpace(raw []byte, i int) int {
	for i < len(raw) {
		switch raw[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			return i
		}
	}

	return i
}

func consumeString(raw []byte, i int) (int, error) {
	if i >= len(raw) || raw[i] != '"' {
		return 0, fmt.Errorf("expected string")
	}

	i++
	for i < len(raw) {
		switch raw[i] {
		case '\\':
			i += 2
		case '"':
			return i + 1, nil
		default:
			i++
		}
	}

	return 0, fmt.Errorf("unterminated string")
}

func consumeValue(raw []byte, i int) (int, error) {
	if i >= len(raw) {
		return 0, fmt.Errorf("expected value")
	}

	switch raw[i] {
	case '{', '[':
		stack := []byte{raw[i]}
		i++
		for i < len(raw) {
			switch raw[i] {
			case '"':
				end, err := consumeString(raw, i)
				if err != nil {
					return 0, err
				}
				i = end
			case '{', '[':
				stack = append(stack, raw[i])
				i++
			case '}':
				if len(stack) == 0 || stack[len(stack)-1] != '{' {
					return 0, fmt.Errorf("unexpected closing brace")
				}
				stack = stack[:len(stack)-1]
				i++
				if len(stack) == 0 {
					return i, nil
				}
			case ']':
				if len(stack) == 0 || stack[len(stack)-1] != '[' {
					return 0, fmt.Errorf("unexpected closing bracket")
				}
				stack = stack[:len(stack)-1]
				i++
				if len(stack) == 0 {
					return i, nil
				}
			default:
				i++
			}
		}

		return 0, fmt.Errorf("unterminated composite value")
	case '"':
		return consumeString(raw, i)
	default:
		for i < len(raw) {
			switch raw[i] {
			case ',', '}', ']', ' ', '\n', '\r', '\t':
				return i, nil
			default:
				i++
			}
		}

		return i, nil
	}
}
