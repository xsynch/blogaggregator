package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/user"

	"github.com/xsynch/blogaggregator/internal/config"
)

type state struct{
	*config.Config
}

func main(){
	cfgFile, err := config.Read()
	if err != nil {
		log.Printf("Error occurred: %s",err)
		
	}
	user, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cfgFile.Setuser(user.Username)
	if err != nil {
		log.Fatalf("Error setting username: %s",err)
	}

	cfgFile, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading the config file: %s",err)
	}
	cfgData, err := json.Marshal(&cfgFile)
	if err != nil {
		log.Fatalf("Error reading the config file: %s",err)
	}
	fmt.Printf("%s",cfgData)
	
	
}
