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
	Long:  `Generate a password hash for use in david configuration. Supports bcrypt, argon2, and scrypt with configurable parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		bcryptCost, _ := cmd.Flags().GetInt("cost")

		// Argon2 parameters
		argon2Memory, _ := cmd.Flags().GetUint32("argon2-memory")
		argon2Iterations, _ := cmd.Flags().GetUint("argon2-iterations")
		argon2Parallelism, _ := cmd.Flags().GetUint("argon2-parallelism")

		// Scrypt parameters
		scryptN, _ := cmd.Flags().GetUint("scrypt-n")
		scryptR, _ := cmd.Flags().GetUint("scrypt-r")
		scryptP, _ := cmd.Flags().GetUint("scrypt-p")

		if password == "" {
			password = readPassword()
		}

		if algorithm == "" {
			algorithm = "bcrypt"
		}

		var params app.HashParams
		switch algorithm {
		case "bcrypt":
			params = app.HashParams{BcryptCost: bcryptCost}
		case "argon2", "argon2id":
			params = app.HashParams{
				Argon2Memory:      argon2Memory,
				Argon2Iterations:  argon2Iterations,
				Argon2Parallelism: argon2Parallelism,
			}
		case "scrypt":
			params = app.HashParams{
				ScryptN: scryptN,
				ScryptR: scryptR,
				ScryptP: scryptP,
			}
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
	passwdCmd.Flags().IntP("cost", "c", 10, "Hash cost factor for bcrypt (default 10)")

	// Argon2 flags
	passwdCmd.Flags().Uint32("argon2-memory", 65536, "Argon2 memory in KiB (default 65536)")
	passwdCmd.Flags().Uint("argon2-iterations", 3, "Argon2 iterations (default 3)")
	passwdCmd.Flags().Uint("argon2-parallelism", 4, "Argon2 parallelism (default 4)")

	// Scrypt flags
	passwdCmd.Flags().Uint("scrypt-n", 16384, "Scrypt N parameter (power of 2, default 16384)")
	passwdCmd.Flags().Uint("scrypt-r", 8, "Scrypt r parameter (default 8)")
	passwdCmd.Flags().Uint("scrypt-p", 1, "Scrypt p parameter (default 1)")

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
