package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/audstanley/david/app"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bcpt",
	Short: "Password hash generator",
	Long:  `bcpt is a CLI tool to generate password hashes (bcrypt, argon2, scrypt) for david configuration.`,
}

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark password hashing algorithms",
	Long:  `Measure the performance of different password hashing algorithms on this system.`,
	Run: func(cmd *cobra.Command, args []string) {
		benchmarkHashing()
	},
}

func init() {
	RootCmd.AddCommand(benchmarkCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return RootCmd
}

func GetPasswdCmd() *cobra.Command {
	var passwdCmd *cobra.Command
	for _, c := range RootCmd.Commands() {
		if c.Use == "passwd" {
			passwdCmd = c
			break
		}
	}
	return passwdCmd
}

func benchmarkHashing() {
	password := "testpassword123"
	iterations := 100

	fmt.Println("Password Hashing Benchmark")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Printf("Testing with password: %s\n", password)
	fmt.Printf("Iterations per algorithm: %d\n", iterations)
	fmt.Println()

	// Benchmark bcrypt
	fmt.Println("Bcrypt (cost 10):")
	bcryptHash, _ := app.GeneratePasswordHash(password, "bcrypt", app.HashParams{BcryptCost: 10})
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = app.GeneratePasswordHash(password, "bcrypt", app.HashParams{BcryptCost: 10})
	}
	bcryptGenTime := time.Since(start)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = app.VerifyPasswordHash(password, bcryptHash)
	}
	bcryptVerifyTime := time.Since(start)

	fmt.Printf("  Generate: %s per hash\n", formatTime(bcryptGenTime/time.Duration(iterations)))
	fmt.Printf("  Verify:   %s per hash\n", formatTime(bcryptVerifyTime/time.Duration(iterations)))
	fmt.Println()

	// Benchmark argon2
	fmt.Println("Argon2id:")
	argon2Hash, _ := app.GeneratePasswordHash(password, "argon2", app.HashParams{BcryptCost: 10})
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = app.GeneratePasswordHash(password, "argon2", app.HashParams{BcryptCost: 10})
	}
	argon2GenTime := time.Since(start)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = app.VerifyPasswordHash(password, argon2Hash)
	}
	argon2VerifyTime := time.Since(start)

	fmt.Printf("  Generate: %s per hash\n", formatTime(argon2GenTime/time.Duration(iterations)))
	fmt.Printf("  Verify:   %s per hash\n", formatTime(argon2VerifyTime/time.Duration(iterations)))
	fmt.Println()

	// Benchmark scrypt
	fmt.Println("Scrypt:")
	scryptHash, _ := app.GeneratePasswordHash(password, "scrypt", app.HashParams{BcryptCost: 10})
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = app.GeneratePasswordHash(password, "scrypt", app.HashParams{BcryptCost: 10})
	}
	scryptGenTime := time.Since(start)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = app.VerifyPasswordHash(password, scryptHash)
	}
	scryptVerifyTime := time.Since(start)

	fmt.Printf("  Generate: %s per hash\n", formatTime(scryptGenTime/time.Duration(iterations)))
	fmt.Printf("  Verify:   %s per hash\n", formatTime(scryptVerifyTime/time.Duration(iterations)))
	fmt.Println()

	// Summary table
	fmt.Println("Summary Table:")
	fmt.Println("---------------------------------------------------------------------")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n", "Algorithm", "Gen Total", "Gen/Hash", "Ver Total", "Ver/Hash")
	fmt.Println("---------------------------------------------------------------------")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"bcrypt",
		formatTime(bcryptGenTime),
		formatTime(bcryptGenTime/time.Duration(iterations)),
		formatTime(bcryptVerifyTime),
		formatTime(bcryptVerifyTime/time.Duration(iterations)))
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"argon2",
		formatTime(argon2GenTime),
		formatTime(argon2GenTime/time.Duration(iterations)),
		formatTime(argon2VerifyTime),
		formatTime(argon2VerifyTime/time.Duration(iterations)))
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"scrypt",
		formatTime(scryptGenTime),
		formatTime(scryptGenTime/time.Duration(iterations)),
		formatTime(scryptVerifyTime),
		formatTime(scryptVerifyTime/time.Duration(iterations)))
	fmt.Println("---------------------------------------------------------------------")
}

func formatTime(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d µs", d.Microseconds())
	}
	return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000.0)
}
