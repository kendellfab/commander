commander
=========

A golang tool to help abstract a command line command set.

1. Create new command instance.

  command := commander.NewCommander("Application Name")
  
2. Parse flags

  flag.Parse()
  
3. Register commands

  command.RegisterCommand("title", "invoke", "description", func (args []string){})
  
4. Execute commands

  command.ExecuteCommand(flag.Args())
