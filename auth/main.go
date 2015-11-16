package main

import (
	"github.com/ArdanStudios/aggserver/auth/commands"
	"github.com/ArdanStudios/aggserver/auth/session"
	"github.com/spf13/cobra"
)

func main() {
	// Initialize a the authentication Session.
	session.Init(nil)

	cli := &cobra.Command{
		Use:   "auth [subcommands] [....subcommand parameters]",
		Short: "auth provides a cli management tool for doing entity creation and authentication against the database api",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	// Initialize all sub-commanders
	commands.InitUserCommands()
	cli.AddCommand(commands.UserCommander())
	cli.Execute()

}
