package commands

import (
	"fmt"
	"strings"
)

// Command 命令接口
type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

// CommandRegistry 命令注册器
type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]Command)}
}

func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Execute(args []string) error {
	if len(args) == 0 {
		r.ShowHelp()
		return nil
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	if cmdName == "help" || cmdName == "--help" || cmdName == "-h" {
		r.ShowHelp()
		return nil
	}

	cmd, exists := r.commands[cmdName]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmdName)
	}
	return cmd.Execute(cmdArgs)
}

func (r *CommandRegistry) ShowHelp() {
	fmt.Printf("MyApp %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  myapp [command] [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	for name, cmd := range r.commands {
		fmt.Printf("  %-15s %s\n", name, cmd.Description())
	}
	fmt.Println()
	fmt.Println("Global options:")
	fmt.Println("  -v, --version  Show version")
	fmt.Println("  -h, --help     Show help")
}

// ParseFlags 解析命令行参数
func ParseFlags(args []string) (map[string]string, []string) {
	flags := make(map[string]string)
	var remaining []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				flags[parts[0]] = parts[1]
			} else {
				key := arg[2:]
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					flags[key] = args[i+1]
					i++
				} else {
					flags[key] = "true"
				}
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			key := arg[1:]
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags[key] = args[i+1]
				i++
			} else {
				flags[key] = "true"
			}
		} else {
			remaining = append(remaining, arg)
		}
	}
	return flags, remaining
}
