package main

import (
	"fmt"
	"os"

	"dotfiles/src/helpers/gmail"

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

	var listOptions gmail.ListOptions
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List mails, newest first",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if err := gmail.ListMails(domain, listOptions); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	listCommand.Flags().IntVar(&listOptions.Limit, "limit", 5, "Number of mails to list (1-100)")
	listCommand.Flags().StringVar(&listOptions.From, "from", "", "Only mails from this sender (name or address)")
	listCommand.Flags().StringVar(&listOptions.Subject, "subject", "", "Only mails whose subject contains this phrase")
	listCommand.Flags().StringVar(&listOptions.Text, "text", "", "Only mails containing this phrase anywhere")
	listCommand.Flags().StringVar(&listOptions.After, "after", "", "Only mails after this date (YYYY-MM-DD)")
	listCommand.Flags().StringVar(&listOptions.Before, "before", "", "Only mails before this date (YYYY-MM-DD)")
	listCommand.Flags().BoolVar(&listOptions.Unread, "unread", false, "Only unread mails")
	listCommand.Flags().BoolVar(&listOptions.IncludeTrash, "include-trash", false, "Also include mails in trash and spam")
	command.AddCommand(listCommand)

	var format string
	viewCommand := &cobra.Command{
		Use:   "view <id>",
		Short: "Show a mail",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := gmail.ViewMail(args[0], format); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	viewCommand.Flags().StringVar(&format, "format", "text", "Body format: text or html")
	command.AddCommand(viewCommand)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
