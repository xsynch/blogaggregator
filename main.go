package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/xsynch/blogaggregator/internal/config"
	"github.com/xsynch/blogaggregator/internal/database"
)

func main() {
	
	st := &state{}
	cmds := commands{allCommands: make(map[string]func(*state, command) error)}
	cmds.register("login",handlerLogin)
	cmds.register("register",handlerRegister)
	cmds.register("reset",handleReset)
	cmds.register("users", handleGetUsers)
	cmds.register("agg",handleFetchFeed)
	cmds.register("addfeed",handleAddFeed)
	cmds.register("feeds",handleGetAllFeeds)
	cmds.register("follow", handleFeedFollow)
	cmds.register("following",handleFollowing)

	cfgFile, err := config.Read()
	if err != nil {
		log.Printf("Error occurred: %s", err)

	}
	st.cfg = cfgFile
	if len(os.Args) < 2{
		log.Fatalf("At least two arguments are needed\n")
	}
	dbURL := st.cfg.DBURL
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening the database: %s",err)
	}
	dbQueries := database.New(db)
	st.db = dbQueries

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
