package commands

import "github.com/spf13/cobra"

// includes the set of CLI commanders for using and managment user entity CRUD and
// authentication managment.
var (
	user        cobra.Command
	userAuth    cobra.Command
	userLogin   cobra.Command
	userCreate  cobra.Command
	userUpdate  cobra.Command
	userDestroy cobra.Command
)

// InitUserCommands must be called only once.
// It initializes and sets up the commander for user management using the CLI.
func InitUserCommands() {
	user = cobra.Command{}

	// Initalize the subcommands for this commander.
	initUserAuth()
	initUserCreate()
	initUserDestroy()
	initUserLogin()
	initUserUpdate()

	// Register the sub-commands into the root commander.
	user.AddCommand(&userLogin, &userCreate, &userDestroy, &userUpdate)
}

// UserCommander returns the user command for accessing user-level cli commands.
func UserCommander() *cobra.Command {
	return &user
}

func initUserCreate() {
	userCreate = cobra.Command{
		Use:   "login [subcommands] [....subcommand parameters]",
		Short: "login ",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func initUserAuth() {
	userAuth = cobra.Command{
		Use:   "login [subcommands] [....subcommand parameters]",
		Short: "login ",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
func initUserLogin() {
	userLogin = cobra.Command{
		Use:   "login [subcommands] [....subcommand parameters]",
		Short: "login ",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func initUserUpdate() {
	userLogin = cobra.Command{
		Use:   "login [subcommands] [....subcommand parameters]",
		Short: "login ",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func initUserDestroy() {
	userLogin = cobra.Command{
		Use:   "login [subcommands] [....subcommand parameters]",
		Short: "login ",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
