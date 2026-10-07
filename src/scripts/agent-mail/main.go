package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
	"github.com/spf13/cobra"
)

func main() {
	domain := os.Getenv("AGENT_MAIL_DOMAIN")
	if domain == "" {
		fmt.Fprintln(os.Stderr, "AGENT_MAIL_DOMAIN is not set")
		os.Exit(1)
	}

	command := &cobra.Command{
		Use:   "agent-mail",
		Short: "Read mails",
		Long:  "Read all mails sent to *@" + domain + ".",
	}

	var limit int
	var from, subject, text, after, before string
	var unread, includeTrash bool
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List mails, newest first",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if limit < 1 || limit > 100 {
				fmt.Fprintln(os.Stderr, "--limit must be between 1 and 100")
				os.Exit(1)
			}

			query := []string{"to:" + domain}
			for _, filter := range []struct{ name, operator, value string }{
				{"from", "from:", from},
				{"subject", "subject:", subject},
				{"text", "", text},
			} {
				if filter.value == "" {
					continue
				}
				if strings.Contains(filter.value, `"`) {
					fmt.Fprintln(os.Stderr, "--"+filter.name+` must not contain "`)
					os.Exit(1)
				}
				query = append(query, filter.operator+`"`+filter.value+`"`)
			}
			for _, filter := range []struct{ name, operator, value string }{
				{"after", "after:", after},
				{"before", "before:", before},
			} {
				if filter.value == "" {
					continue
				}
				date, err := time.Parse("2006-01-02", filter.value)
				if err != nil {
					fmt.Fprintln(os.Stderr, "--"+filter.name+" must be a date like 2026-10-07")
					os.Exit(1)
				}
				query = append(query, filter.operator+date.Format("2006/01/02"))
			}
			if unread {
				query = append(query, "is:unread")
			}
			if includeTrash {
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
			gog(&result, "gmail", "messages", "search", strings.Join(query, " "), "--max", strconv.Itoa(limit))

			if len(result.Messages) == 0 {
				fmt.Println()
				fmt.Println("No mails found.")
				return
			}

			for _, message := range result.Messages {
				sent, err := time.Parse(time.RFC3339, message.DateISO)
				if err != nil {
					fmt.Fprintln(os.Stderr, "failed to parse mail date:", err)
					os.Exit(1)
				}
				fmt.Println()
				fmt.Println("### " + message.Subject)
				fmt.Println("**ID:** " + message.ID)
				fmt.Println("**Date:** " + message.Date + " | " + timeAgo(sent))
				fmt.Println("**From:** " + message.From)
			}
		},
	}
	listCommand.Flags().IntVar(&limit, "limit", 5, "Number of mails to list (1-100)")
	listCommand.Flags().StringVar(&from, "from", "", "Only mails from this sender (name or address)")
	listCommand.Flags().StringVar(&subject, "subject", "", "Only mails whose subject contains this phrase")
	listCommand.Flags().StringVar(&text, "text", "", "Only mails containing this phrase anywhere")
	listCommand.Flags().StringVar(&after, "after", "", "Only mails after this date (YYYY-MM-DD)")
	listCommand.Flags().StringVar(&before, "before", "", "Only mails before this date (YYYY-MM-DD)")
	listCommand.Flags().BoolVar(&unread, "unread", false, "Only unread mails")
	listCommand.Flags().BoolVar(&includeTrash, "include-trash", false, "Also include mails in trash and spam")
	command.AddCommand(listCommand)

	var format string
	viewCommand := &cobra.Command{
		Use:   "view <id>",
		Short: "Show a mail",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if format != "text" && format != "html" {
				fmt.Fprintln(os.Stderr, "--format must be text or html")
				os.Exit(1)
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
			gog(&result, "gmail", "get", args[0])

			textData := findPartData(result.Message.Payload, "text/plain")
			htmlData := findPartData(result.Message.Payload, "text/html")
			var decoded string
			if format == "text" && textData != "" {
				decoded = decodePartData(textData)
			} else if htmlData == "" {
				fmt.Fprintln(os.Stderr, "mail has no "+format+" body")
				os.Exit(1)
			} else if format == "html" {
				decoded = decodePartData(htmlData)
			} else {
				text, err := html2text.FromString(decodePartData(htmlData))
				if err != nil {
					fmt.Fprintln(os.Stderr, "failed to convert html body to text:", err)
					os.Exit(1)
				}
				decoded = text
			}

			body := strings.TrimRight(decoded, " \r\n")
			if format == "text" {
				lastLineStart := strings.LastIndex(body, "\n") + 1
				if strings.Trim(body[lastLineStart:], "- ") == "" {
					body = strings.TrimRight(body[:lastLineStart], " \r\n")
				}
			}

			fmt.Println("# " + result.Headers.Subject)
			fmt.Println("**ID:** " + result.Message.ID)
			fmt.Println("**Date:** " + result.Headers.Date + " | " + timeAgo(time.UnixMilli(result.Message.InternalDate)))
			fmt.Println("**From:** " + result.Headers.From)
			fmt.Println("**To:** " + result.Headers.To)
			fmt.Println()
			fmt.Println("---")
			fmt.Println(body)
		},
	}
	viewCommand.Flags().StringVar(&format, "format", "text", "Body format: text or html")
	command.AddCommand(viewCommand)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

type messagePart struct {
	MimeType string `json:"mimeType"`
	Body     struct {
		Data string `json:"data"`
	} `json:"body"`
	Parts []messagePart `json:"parts"`
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

func decodePartData(data string) string {
	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(data, "="))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to decode mail body:", err)
		os.Exit(1)
	}
	return string(decoded)
}

func gog(result any, args ...string) {
	cmd := exec.Command("gog", append(args, "--json")...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}
	if err := json.Unmarshal(output, result); err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse gog output:", err)
		os.Exit(1)
	}
}

func timeAgo(t time.Time) string {
	elapsed := time.Since(t)
	if elapsed < time.Minute {
		return "just now"
	}
	if elapsed < time.Hour {
		return fmt.Sprintf("%dm ago", int(elapsed.Minutes()))
	}
	if elapsed < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(elapsed.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(elapsed.Hours()/24))
}
