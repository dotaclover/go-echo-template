package commands

import "fmt"

var (
	version   string
	buildTime string
)

func SetVersion(v string)   { version = v }
func SetBuildTime(bt string) { buildTime = bt }

func PrintVersion() {
	fmt.Printf("MyApp %s\n", version)
	fmt.Printf("Build Time: %s\n", buildTime)
}

// RegisterCommands 注册所有命令
func RegisterCommands(registry *CommandRegistry) {
	registry.Register(&MigrateCommand{})
	registry.Register(&ServeCommand{})
	registry.Register(&SeedCommand{})
}
