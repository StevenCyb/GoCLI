package cli

type Context struct {
	commands  []string
	arguments map[string]string
	options   map[string]string
}

func NewContext() *Context {
	return &Context{
		commands:  make([]string, 0),
		arguments: make(map[string]string),
		options:   make(map[string]string),
	}
}

func (c *Context) VisitedCommand(command string) bool {
	for _, cmd := range c.commands {
		if cmd == command {
			return true
		}
	}
	return false
}

func (c *Context) VisitedArgument(argument string) bool {
	_, exists := c.arguments[argument]
	return exists
}

func (c *Context) UsedOption(option string) bool {
	_, exists := c.options[option]
	return exists
}

func (c *Context) GetCommand(command string) *string {
	if c.VisitedCommand(command) {
		return &command
	}
	return nil
}

func (c *Context) GetArgument(argument string) *string {
	if value, exists := c.arguments[argument]; exists {
		return &value
	}
	return nil
}

func (c *Context) GetOption(option string) *string {
	if value, exists := c.options[option]; exists {
		return &value
	}
	return nil
}
