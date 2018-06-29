package main

import (
	"os"
	"log"
	"github.com/alounce/clirescue/trackerapi"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "clirescue"
	app.Usage = "CLI tool to talk to the Pivotal Tracker's API"

	app.Commands = []cli.Command{
		{
			Name:  "me",
			Usage: "prints out Tracker's representation of your account",
			Action: me,
		},
	}

	app.Run(os.Args)
}

func me (c *cli.Context) {
	user, err := trackerapi.Me()
	if err != nil {
		log.Fatalf("HTTP Request failed: %v", err)
	}
	// print out user info
	log.Println("----------------------")
	log.Println("User Info")
	log.Println("----------------------")
	log.Printf("User Login: %v", user.Username)
	log.Printf("User Initials: %v", user.Initials)
	log.Printf("User Email: %v", user.Email)
	log.Printf("TimeZone offser: %v", user.Timezone.Offset)

}
