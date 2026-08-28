package main

import (
	"fmt"
	"os"

	"example.com/childfitness/internal/cli"
)

func main() {
	databasePath := "childfitness.db"
	if value := os.Getenv("CHILDFITNESS_DB"); value != "" {
		databasePath = value
	}
	app, err := cli.NewApplication(databasePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer app.Close()
	if err := app.Execute(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
