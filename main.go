package main

import (
	"log"
	"os"

	"github.com/xsynch/blogaggregator/internal/config"
)

func main() {
	st := &state{}
	cmds := commands{allCommands: make(map[string]func(*state, command) error)}
	cmds.register("login",handlerLogin)

	cfgFile, err := config.Read()
	if err != nil {
		log.Printf("Error occurred: %s", err)

	}
	st.Config = cfgFile
	if len(os.Args) < 2{
		log.Fatalf("At least two arguments are needed\n")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmdS := command{name: cmdName, args: cmdArgs}

	err = cmds.run(st, cmdS)
	if err != nil {
		log.Fatal(err)
	}
	

	

	// user, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// err = cfgFile.Setuser(user.Username)
	// if err != nil {
	// 	log.Fatalf("Error setting username: %s", err)
	// }

	// cfgFile, err = config.Read()
	// if err != nil {
	// 	log.Fatalf("Error reading the config file: %s", err)
	// }
	// cfgData, err := json.Marshal(&cfgFile)
	// if err != nil {
	// 	log.Fatalf("Error reading the config file: %s", err)
	// }
	// fmt.Printf("%s", cfgData)

}
