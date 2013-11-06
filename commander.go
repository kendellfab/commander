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

type Commander struct {
	Application string
	Commands    map[string]*Command
}

type Command struct {
	Name, Invoke, Description string
	Run                       func(args []string)
}

func NewCommander(application string) *Commander {
	commander := &Commander{
		Application: application,
		Commands:    make(map[string]*Command),
	}
	return commander
}

func (c *Commander) RegisterCommand(name, invoke, description string, f func(args []string)) {
	command := &Command{
		name, invoke, description, f,
	}

	c.Commands[invoke] = command
}

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

func (c *Commander) usage() {
	renderTpl(os.Stderr, usageTemplate, c)
	os.Exit(2)
}

func (c *Commander) help(invoke string) {
	if command, ok := c.Commands[invoke]; ok {
		data := make(map[string]interface{})
		data["app"] = c
		data["command"] = command
		renderTpl(os.Stderr, helpTemplate, data)
	}
}

func renderTpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
