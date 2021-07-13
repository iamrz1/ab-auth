package main

import (
	"github.com/iamrz1/auth/cmd"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
)

// rootCmd is the root of all sub commands in the binary
// it doesn't have a Run method as it executes other sub commands
var rootCmd = &cobra.Command{
	Use:     "ab-boilerplate",
	Short:   "ab-boilerplate is a web server boilerplate",
	Version: "1.0.0",
}

func init() {
	rootCmd.AddCommand(cmd.SrvCmd)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
