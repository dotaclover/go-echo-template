package main

import (
	"myapp/cmd/commands"
	"os"
)

var BuildTime string

func main() {
	if BuildTime == "" {
		BuildTime = "dev"
	}
	commands.SetBuildTime(BuildTime)
	commands.SetVersion("0.1.0")

	registry := commands.NewCommandRegistry()
	commands.RegisterCommands(registry)

	args := os.Args[1:]

	if len(args) > 0 && (args[0] == "--version" || args[0] == "-v" || args[0] == "version") {
		commands.PrintVersion()
		return
	}

	if len(args) == 0 {
		args = []string{"serve"}
	}

	if err := registry.Execute(args); err != nil {
		os.Stderr.WriteString("Error: " + err.Error() + "\n")
		os.Exit(1)
	}
}
