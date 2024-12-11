package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/xsynch/blogaggregator/internal/config"
	"github.com/xsynch/blogaggregator/internal/database"
)



type state struct {
	cfg *config.Config
	db *database.Queries
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
	userName := cmd.args[0]
	
	if len(cmd.args) != 1 {
		return fmt.Errorf("login should contain one value: login <name>")
	}
	_,err := s.db.GetUser(context.Background(),userName)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%s not found in the database",userName) 
	}
	
	err = s.cfg.Setuser(userName)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s has been set successfully\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	var emptyDBuser =database.User{}
	if len(cmd.args) != 1 {
		return fmt.Errorf("function %s takes two arguments: %s <name>",cmd.name,cmd.name)
	}
	userName := cmd.args[0]
	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: userName,
	}
	dbUser,err := s.db.GetUser(context.Background(),userName)
	if err != nil && err != sql.ErrNoRows {
		return err 
	}
	if dbUser != emptyDBuser {
		return fmt.Errorf("user %s already exists in the db",userName)
	}
	
	_,err = s.db.CreateUser(context.Background(),userParams)
	if err != nil {
		return err 
	}
	err = s.cfg.Setuser(userName)
	if err != nil {
		return err 
	}
	fmt.Printf("%s created successfully\n",userName)
	log.Printf("%s created successfully\n",userName)
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

func handleReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err 
	}
	fmt.Printf("User database cleared successfully")
	return nil 
}

func handleGetUsers(s *state, cmd command) error {
	dbUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err 
	}
	for _, value := range dbUsers{
		if s.cfg.CURRENTUSER == value.Name{
			fmt.Printf("* %s (current)\n",value.Name)
		} else {
			fmt.Printf("* %s\n",value.Name)
		}
	}
	return nil 
}