package main

import (
	"context"
	"database/sql"
	"fmt"
	"html"
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

func handleFetchFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%s takes a time duration for an argument: %s <time_between_requests>",cmd.name,cmd.name)
	}
	
	timeBetweenRequests,err  := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err 
	}
	fmt.Printf("Collecting feeds every %v seconds\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
	url,err := scrapeFeeds(s)
	if err != nil {
		return err 
	}
	rssFeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err 
	}
	fmt.Printf("%s\n",returnStringFromHTML(rssFeed.Channel.Title))
	// fmt.Printf("%s\n",returnStringFromHTML(rssFeed.Channel.Link))
	// fmt.Printf("%s\n",returnStringFromHTML(rssFeed.Channel.Description))
	// for _,val := range rssFeed.Channel.Item{
	// 	fmt.Printf("%s\n%s\n%s\n%s\n",returnStringFromHTML(val.Title), returnStringFromHTML(val.Link),returnStringFromHTML(val.Description),returnStringFromHTML(val.PubDate))
	// }
	
}
}

func returnStringFromHTML(htmlText string) string {
	return html.UnescapeString(htmlText)
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("%s takes two arguments: %s <name> <url>", cmd.name, cmd.name)
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	// current_user,err  := s.db.GetUser(context.Background(),s.cfg.CURRENTUSER)
	// if err != nil {
	// 	return err 
	// }
	// user_id := current_user.ID
	cfParams :=database.CreateFeedParams{ID: uuid.New(),CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: feedName, Url: feedURL, UserID: user.ID }
	_,err := s.db.CreateFeed(context.Background(), cfParams)
	if err != nil {
		return err 
	}
	fmt.Printf("%s Feed added successfully\n",feedName)
	feedFollowParams := database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: cfParams.CreatedAt, UpdatedAt: cfParams.UpdatedAt, UserID: cfParams.UserID, FeedID: cfParams.ID}
	_,err = s.db.CreateFeedFollow(context.Background(),feedFollowParams)
	if err != nil {
		return err 
	}
	return nil 


}

func handleGetAllFeeds(s *state, cmd command) error {
	feeds,err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err 
	}
	
	for _,feed := range feeds {
		dbUser, err := s.db.GetUserByid(context.Background(),feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n",feed.Name, feed.Url, dbUser.Name)
	}
	return nil 
}

func handleFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%s only takes one argument: %s <url>",cmd.name,cmd.name)
	}
	url := cmd.args[0]
	// user_id, err := s.db.GetUser(context.Background(), s.cfg.CURRENTUSER)
	// if err != nil {
	// 	return err 
	// }
	feed_info,err := s.db.GetFeedIDByName(context.Background(), url)
	if err != nil {
		return err 
	}
	params := database.CreateFeedFollowParams{ID: uuid.New(),CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID:user.ID, FeedID: feed_info.ID }
	follow,err := s.db.CreateFeedFollow(context.Background(),params)
	if err != nil {
		return err 		
	}
	fmt.Printf("%s %s",follow.FeedName,follow.UserName)
	return nil 
}

func handleFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%s takes no arguments",cmd.name)
	}
	// userName := s.cfg.CURRENTUSER
	// dbUser,err := s.db.GetUser(context.Background(), userName)
	// if err != nil {
	// 	return err 
	// }
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err 
	}
	for _,value := range follows {
		// u,err := s.db.GetUserByid(context.Background(), value.UserID)
		// if err != nil {
		// 	return err 
		// }
		feed_name, err := s.db.GetFeedsByID(context.Background(), value.FeedID)
		if err != nil {
			return err 
		}
		fmt.Printf("- %s\n",feed_name.Name)
	}
	return nil 
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1{
		return fmt.Errorf("%s accepts one argument: %s <url>",cmd.name, cmd.name)
	}
	url := cmd.args[0]
	feedID, err := s.db.GetFeedIDByName(context.Background(),url)
	if err != nil {
		return err 
	}
	params := database.UnfollowParams{FeedID: feedID.ID, UserID: user.ID}
	err = s.db.Unfollow(context.Background(), params)
	if err != nil {
		return err 
	}
	fmt.Printf("%s has been removed from %s", feedID.Name, user.Name)
	return nil 

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
		



	
	return  func(s *state, cmd command) error {
		
		userName, err := s.db.GetUser(context.Background(), s.cfg.CURRENTUSER)
		if err != nil {
			return err
		}
		
		return  handler(s,cmd,userName)

	}



}

func scrapeFeeds(s *state) (string, error) {
	feed,err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return "", err 
	}
	params := database.MarkFeedFetchedParams{LastFetchedAt: sql.NullTime{time.Now(),true}, ID: feed.ID}
	err = s.db.MarkFeedFetched(context.Background(),params)
	if err != nil {
		return "",err 
	}
	return feed.Url,nil  

	
}