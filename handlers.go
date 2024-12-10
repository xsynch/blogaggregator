package main

import (
	"fmt"

	"github.com/xsynch/blogaggregator/internal/config"
)



type state struct {
	*config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	allCommands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	// log.Printf("The arguments supplied are: %v",cmd.args)
	
	if len(cmd.args) != 1 {
		return fmt.Errorf("login should contain one value: login <name>")
	}
	
	err := s.Config.Setuser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User: %s has been set successfully\n", cmd.args[0])
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.allCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler,ok := c.allCommands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s",cmd.name)
	}
	
	return handler(s,cmd)
}
