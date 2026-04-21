package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func MergeJSONObject(prev string, next string) (string, error) {
	prevLayout, err := parseObjectLayout([]byte(prev))
	if err != nil {
		return "", err
	}
	prevEntries := prevLayout.Entries

	nextLayout, err := parseObjectLayout([]byte(next))
	if err != nil {
		return "", err
	}
	nextEntries := nextLayout.Entries

	if jsonEqual([]byte(prev), []byte(next)) {
		return prev, nil
	}

	patch, err := jsonpatch.CreateMergePatch([]byte(prev), []byte(next))
	if err != nil {
		return "", err
	}

	merged, err := jsonpatch.MergePatch([]byte(prev), patch)
	if err != nil {
		return "", err
	}

	mergedLayout, err := parseObjectLayout(merged)
	if err != nil {
		return "", err
	}
	mergedEntries := mergedLayout.Entries

	mergedByKey := map[string]objectEntry{}
	for _, entry := range mergedEntries {
		mergedByKey[entry.Key] = entry
	}

	replaceOps := []replaceOp{}
	for _, entry := range prevEntries {
		mergedEntry, ok := mergedByKey[entry.Key]
		if !ok {
			continue
		}

		if jsonEqual(entry.Val, mergedEntry.Val) {
			continue
		}

		value := mergedEntry.Val
		if isJSONObject(entry.Val) && isJSONObject(value) {
			mergedObject, err := MergeJSONObject(string(entry.Val), string(value))
			if err != nil {
				return "", err
			}

			value = []byte(mergedObject)
		}

		replaceOps = append(replaceOps, replaceOp{
			start: entry.ValueStart,
			end:   entry.ValueEnd,
			text:  value,
		})
	}

	out := applyReplaceOps([]byte(prev), replaceOps)

	for {
		layout, err := parseObjectLayout(out)
		if err != nil {
			return "", err
		}

		removed := false
		for i := len(layout.Entries) - 1; i >= 0; i-- {
			if _, ok := mergedByKey[layout.Entries[i].Key]; ok {
				continue
			}

			start, end := removalRange(layout, i)
			out = append(append([]byte{}, out[:start]...), out[end:]...)
			removed = true
			break
		}

		if !removed {
			break
		}
	}

	presentLayout, err := parseObjectLayout(out)
	if err != nil {
		return "", err
	}

	presentKeys := map[string]bool{}
	for _, entry := range presentLayout.Entries {
		presentKeys[entry.Key] = true
	}

	for _, nextEntry := range nextEntries {
		mergedEntry, ok := mergedByKey[nextEntry.Key]
		if !ok || presentKeys[nextEntry.Key] {
			continue
		}

		field := buildField([]byte(next), nextEntry, mergedEntry)
		out, err = appendField(out, []byte(next), nextLayout, field)
		if err != nil {
			return "", err
		}

		presentLayout, err = parseObjectLayout(out)
		if err != nil {
			return "", err
		}

		presentKeys[nextEntry.Key] = true
	}

	return string(out), nil
}

type objectEntry struct {
	Key        string
	KeyRaw     string
	Val        []byte
	PreStart   int
	KeyStart   int
	KeyEnd     int
	ValueStart int
	ValueEnd   int
	CommaPos   int
}

type objectLayout struct {
	OpenBrace  int
	CloseBrace int
	Entries    []objectEntry
}

type replaceOp struct {
	start int
	end   int
	text  []byte
}

func parseObjectEntries(raw []byte) ([]objectEntry, error) {
	layout, err := parseObjectLayout(raw)
	if err != nil {
		return nil, err
	}

	return layout.Entries, nil
}

func parseObjectLayout(raw []byte) (objectLayout, error) {
	i := skipSpace(raw, 0)
	if i >= len(raw) || raw[i] != '{' {
		return objectLayout{}, fmt.Errorf("expected object")
	}

	openBrace := i
	i++
	entries := []objectEntry{}

	for {
		preStart := i
		i = skipSpace(raw, i)
		if i >= len(raw) {
			return objectLayout{}, fmt.Errorf("unexpected end of object")
		}

		if raw[i] == '}' {
			return objectLayout{OpenBrace: openBrace, CloseBrace: i, Entries: entries}, nil
		}

		keyStart := i
		keyEnd, err := consumeString(raw, i)
		if err != nil {
			return objectLayout{}, err
		}

		var key string
		if err := json.Unmarshal(raw[keyStart:keyEnd], &key); err != nil {
			return objectLayout{}, err
		}

		i = skipSpace(raw, keyEnd)
		if i >= len(raw) || raw[i] != ':' {
			return objectLayout{}, fmt.Errorf("expected colon after key")
		}

		i++
		i = skipSpace(raw, i)
		valueStart := i
		valueEnd, err := consumeValue(raw, i)
		if err != nil {
			return objectLayout{}, err
		}

		j := skipSpace(raw, valueEnd)
		if j >= len(raw) {
			return objectLayout{}, fmt.Errorf("unexpected end of object")
		}

		commaPos := -1
		if raw[j] == ',' {
			commaPos = j
			j++
		} else if raw[j] != '}' {
			return objectLayout{}, fmt.Errorf("expected comma or closing brace")
		}

		entries = append(entries, objectEntry{
			Key:        key,
			KeyRaw:     string(raw[keyStart:keyEnd]),
			Val:        append([]byte{}, raw[valueStart:valueEnd]...),
			PreStart:   preStart,
			KeyStart:   keyStart,
			KeyEnd:     keyEnd,
			ValueStart: valueStart,
			ValueEnd:   valueEnd,
			CommaPos:   commaPos,
		})

		i = j
	}
}

func applyReplaceOps(raw []byte, ops []replaceOp) []byte {
	out := append([]byte{}, raw...)
	for i := len(ops) - 1; i >= 0; i-- {
		op := ops[i]
		out = append(append(append([]byte{}, out[:op.start]...), op.text...), out[op.end:]...)
	}

	return out
}

func removalRange(layout objectLayout, idx int) (int, int) {
	entries := layout.Entries
	entry := entries[idx]

	if len(entries) == 1 {
		return entry.PreStart, entry.ValueEnd
	}

	if idx < len(entries)-1 {
		return entry.PreStart, entries[idx+1].PreStart
	}

	prev := entries[idx-1]
	if prev.CommaPos >= 0 {
		return prev.CommaPos, entry.ValueEnd
	}

	return entry.PreStart, entry.ValueEnd
}

func buildField(nextRaw []byte, nextEntry objectEntry, mergedEntry objectEntry) []byte {
	colonAndSpacing := []byte(":")
	if nextEntry.KeyEnd >= 0 && nextEntry.ValueStart >= nextEntry.KeyEnd && nextEntry.ValueStart <= len(nextRaw) {
		colonAndSpacing = append([]byte{}, nextRaw[nextEntry.KeyEnd:nextEntry.ValueStart]...)
	}

	value := mergedEntry.Val
	if len(nextEntry.Val) > 0 && jsonEqual(nextEntry.Val, mergedEntry.Val) {
		value = nextEntry.Val
	}

	field := append([]byte{}, nextEntry.KeyRaw...)
	field = append(field, colonAndSpacing...)
	field = append(field, value...)
	return field
}

func appendField(raw []byte, nextRaw []byte, nextLayout objectLayout, field []byte) ([]byte, error) {
	layout, err := parseObjectLayout(raw)
	if err != nil {
		return nil, err
	}

	entries := layout.Entries
	if len(entries) == 0 {
		leading := []byte{}
		trailing := []byte{}
		if len(nextLayout.Entries) > 0 {
			first := nextLayout.Entries[0]
			leading = append([]byte{}, nextRaw[nextLayout.OpenBrace+1:first.KeyStart]...)
			if len(nextLayout.Entries) == 1 {
				trailing = append([]byte{}, nextRaw[first.ValueEnd:nextLayout.CloseBrace]...)
			} else {
				last := nextLayout.Entries[len(nextLayout.Entries)-1]
				trailing = append([]byte{}, nextRaw[last.ValueEnd:nextLayout.CloseBrace]...)
			}
		}

		out := append([]byte{}, raw[:layout.OpenBrace+1]...)
		out = append(out, leading...)
		out = append(out, field...)
		out = append(out, trailing...)
		out = append(out, raw[layout.CloseBrace:]...)
		return out, nil
	}

	last := entries[len(entries)-1]
	separator := []byte{','}
	if len(entries) >= 2 {
		prev := entries[len(entries)-2]
		separator = append([]byte{}, raw[prev.ValueEnd:last.KeyStart]...)
	} else {
		leading := raw[last.PreStart:last.KeyStart]
		separator = append([]byte{','}, leading...)
	}

	tail := append([]byte{}, raw[last.ValueEnd:layout.CloseBrace]...)
	out := append([]byte{}, raw[:last.ValueEnd]...)
	out = append(out, separator...)
	out = append(out, field...)
	out = append(out, tail...)
	out = append(out, raw[layout.CloseBrace:]...)
	return out, nil
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

func isJSONObject(raw []byte) bool {
	var value map[string]json.RawMessage
	return json.Unmarshal(raw, &value) == nil
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
