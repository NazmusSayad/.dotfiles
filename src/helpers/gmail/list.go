package gmail

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ListOptions struct {
	Limit        int
	From         string
	Subject      string
	Text         string
	After        string
	Before       string
	Unread       bool
	IncludeTrash bool
}

func ListMails(domain string, options ListOptions) error {
	if options.Limit < 1 || options.Limit > 100 {
		return fmt.Errorf("--limit must be between 1 and 100")
	}

	query := []string{"to:" + domain}
	for _, filter := range []struct{ name, operator, value string }{
		{"from", "from:", options.From},
		{"subject", "subject:", options.Subject},
		{"text", "", options.Text},
	} {
		if filter.value == "" {
			continue
		}
		if strings.Contains(filter.value, `"`) {
			return fmt.Errorf(`--%s must not contain "`, filter.name)
		}
		query = append(query, filter.operator+`"`+filter.value+`"`)
	}
	for _, filter := range []struct{ name, operator, value string }{
		{"after", "after:", options.After},
		{"before", "before:", options.Before},
	} {
		if filter.value == "" {
			continue
		}
		date, err := time.Parse("2006-01-02", filter.value)
		if err != nil {
			return fmt.Errorf("--%s must be a date like 2026-10-07", filter.name)
		}
		query = append(query, filter.operator+date.Format("2006/01/02"))
	}
	if options.Unread {
		query = append(query, "is:unread")
	}
	if options.IncludeTrash {
		query = append(query, "in:anywhere")
	}

	var result struct {
		Messages []struct {
			ID      string `json:"id"`
			Date    string `json:"date"`
			DateISO string `json:"internalDateIso"`
			From    string `json:"from"`
			Subject string `json:"subject"`
		} `json:"messages"`
	}
	if err := runGog(&result, "gmail", "messages", "search", strings.Join(query, " "), "--max", strconv.Itoa(options.Limit)); err != nil {
		return err
	}

	if len(result.Messages) == 0 {
		fmt.Print("\nNo mails found.\n")
		return nil
	}

	var markdown strings.Builder
	for _, message := range result.Messages {
		sentAt, err := time.Parse(time.RFC3339, message.DateISO)
		if err != nil {
			return fmt.Errorf("failed to parse mail date: %w", err)
		}
		markdown.WriteString("\n### " + message.Subject + "\n")
		markdown.WriteString("**ID:** " + message.ID + "\n")
		markdown.WriteString("**Date:** " + message.Date + " | " + timeAgo(sentAt) + "\n")
		markdown.WriteString("**From:** " + message.From + "\n")
	}
	fmt.Print(markdown.String())
	return nil
}
