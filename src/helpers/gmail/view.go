package gmail

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
)

type messagePart struct {
	MimeType string `json:"mimeType"`
	Body     struct {
		Data string `json:"data"`
	} `json:"body"`
	Parts []messagePart `json:"parts"`
}

func ViewMail(id string, format string) error {
	if format != "text" && format != "html" {
		return fmt.Errorf("--format must be text or html")
	}

	var result struct {
		Message struct {
			ID           string      `json:"id"`
			InternalDate int64       `json:"internalDate,string"`
			Payload      messagePart `json:"payload"`
		} `json:"message"`
		Headers struct {
			Date    string `json:"date"`
			From    string `json:"from"`
			To      string `json:"to"`
			Subject string `json:"subject"`
		} `json:"headers"`
	}
	if err := runGog(&result, "gmail", "get", id); err != nil {
		return err
	}

	body, err := readBody(result.Message.Payload, format)
	if err != nil {
		return err
	}

	var markdown strings.Builder
	markdown.WriteString("# " + result.Headers.Subject + "\n")
	markdown.WriteString("**ID:** " + result.Message.ID + "\n")
	markdown.WriteString("**From:** " + result.Headers.From + "\n")
	markdown.WriteString("**To:** " + result.Headers.To + "\n")
	markdown.WriteString("**Date:** " + result.Headers.Date + " | " + timeAgo(time.UnixMilli(result.Message.InternalDate)) + "\n")
	markdown.WriteString("\n---\n")
	markdown.WriteString(body + "\n")
	fmt.Print(markdown.String())
	return nil
}

func readBody(payload messagePart, format string) (string, error) {
	textData := findPartData(payload, "text/plain")
	htmlData := findPartData(payload, "text/html")

	if format == "html" {
		if htmlData == "" {
			return "", fmt.Errorf("mail has no html body")
		}
		html, err := decodePartData(htmlData)
		if err != nil {
			return "", err
		}
		return strings.TrimRight(html, " \r\n"), nil
	}

	var text string
	if textData != "" {
		decoded, err := decodePartData(textData)
		if err != nil {
			return "", err
		}
		text = decoded
	} else if htmlData != "" {
		html, err := decodePartData(htmlData)
		if err != nil {
			return "", err
		}
		converted, err := html2text.FromString(html)
		if err != nil {
			return "", fmt.Errorf("failed to convert html body to text: %w", err)
		}
		text = converted
	} else {
		return "", fmt.Errorf("mail has no text body")
	}

	text = strings.TrimRight(text, " \r\n")
	lastLineStart := strings.LastIndex(text, "\n") + 1
	if strings.Trim(text[lastLineStart:], "- ") == "" {
		text = strings.TrimRight(text[:lastLineStart], " \r\n")
	}
	return text, nil
}

func findPartData(part messagePart, mimeType string) string {
	if part.MimeType == mimeType && part.Body.Data != "" {
		return part.Body.Data
	}
	for _, child := range part.Parts {
		if data := findPartData(child, mimeType); data != "" {
			return data
		}
	}
	return ""
}

func decodePartData(data string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(data, "="))
	if err != nil {
		return "", fmt.Errorf("failed to decode mail body: %w", err)
	}
	return string(decoded), nil
}
