package cli

import (
	"fmt"
	"os"
	"syscall"

	"github.com/audstanley/david/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Generate a password hash",
	Long:  `Generate a password hash for use in david configuration. Supports bcrypt, argon2, and scrypt.`,
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		cost, _ := cmd.Flags().GetInt("cost")

		if password == "" {
			password = readPassword()
		}

		if algorithm == "" {
			algorithm = "bcrypt"
		}

		var params app.HashParams
		switch algorithm {
		case "bcrypt":
			params = app.HashParams{BcryptCost: cost}
		case "argon2", "argon2id":
			// Argon2 uses fixed parameters for simplicity
			params = app.HashParams{ScryptN: 16384} // placeholder, not used for argon2
		case "scrypt":
			params = app.HashParams{ScryptN: 16384, ScryptR: 8, ScryptP: 1}
		default:
			fmt.Printf("Unsupported algorithm: %s. Use bcrypt, argon2, or scrypt.\n", algorithm)
			os.Exit(1)
		}

		hash, err := app.GeneratePasswordHash(password, algorithm, params)
		if err != nil {
			fmt.Printf("Error generating hash: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Hashed Password (%s): %s\n", algorithm, hash)
	},
}

func init() {
	RootCmd.AddCommand(passwdCmd)

	passwdCmd.Flags().StringP("password", "p", "", "Password to hash (required)")
	passwdCmd.Flags().String("algorithm", "", "Hash algorithm: bcrypt, argon2, scrypt (default: bcrypt)")
	passwdCmd.Flags().IntP("cost", "c", 10, "Hash cost factor")
	passwdCmd.MarkFlagRequired("password")

	viper.BindPFlags(passwdCmd.Flags())
}

func readPassword() string {
	fmt.Print("Enter password: ")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("An error occurred reading the password: %s\n", err)
		os.Exit(1)
	}

	fmt.Println()
	return string(pw)
}

func GetReadPassword() func() string {
	return readPassword
}
