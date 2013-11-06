// Commander provides an interface to register, parse, and execute command line commands.
package commander

import (
	"io"
	"os"
	"text/template"
)

const (
	usageTemplate = `usage: {{.Application}} command [arguments]
The commands are:
{{range .Commands}}
  {{.Invoke | printf "%-11s"}} {{.Name}}
{{end}}
Use {{.Application}} help [command] for more information.`

	helpTemplate = `usage: {{.app.Application}} {{.command.Invoke}}
{{.command.Description}}`
)

// The struct that holds the application instance and command.
type Commander struct {
	Application string
	Commands    map[string]*Command
}

// The struct that holds the command name, invocation, description and the
// function to execute.
type Command struct {
	Name, Invoke, Description string
	Run                       func(args []string)
}

// Returns an new commander struct, that knows its name.
func NewCommander(application string) *Commander {
	commander := &Commander{
		Application: application,
		Commands:    make(map[string]*Command),
	}
	return commander
}

// Use this on the returned commander to register a command.
func (c *Commander) RegisterCommand(name, invoke, description string, f func(args []string)) {
	command := &Command{
		name, invoke, description, f,
	}

	c.Commands[invoke] = command
}

// Pass in the arg string to find the command and execute it.
func (c *Commander) ExecuteCommand(args []string) bool {
	if len(args) < 1 {
		c.usage()
		return false
	}
	invoke := args[0]
	if invoke == "help" {
		c.help(args[1])
		return true
	}
	if command, ok := c.Commands[invoke]; ok {
		command.Run(args[1:])
		return true
	}
	return false
}

// Internal func to show the usage of registerd commands.
func (c *Commander) usage() {
	renderTpl(os.Stderr, usageTemplate, c)
	os.Exit(2)
}

// Internal func to show how to use a particular command.
func (c *Commander) help(invoke string) {
	if command, ok := c.Commands[invoke]; ok {
		data := make(map[string]interface{})
		data["app"] = c
		data["command"] = command
		renderTpl(os.Stderr, helpTemplate, data)
	}
}

// Renders the desired usage or help template with the data required.
func renderTpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
