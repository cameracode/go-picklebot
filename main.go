package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for the command line params
var (
	Token string
)

const KuteGoAPIURL = "https://kutego-api-xxxxx-ew.a.run.app"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// In n this example, we only care about receiving message events
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the discord session.
	dg.Close()
}

type Pickle struct {
	Name string `json: "name"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!pickle" {
		// Call the KuteGo API and retrieve our cell-shaded 3d pickle rick
		response, err := http.Get(KuteGoAPIURL + "/pickle/" + "arakaki-picklerick")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "arakaki-picklerick.png", response.Request.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get arakaki-picklerick Pickle! (ﾉ◕ヮ◕)ﾉ*:･ﾟ✧")
		}
	}

	if m.Content == "!pickles" {
		// Call the KuteGo API and display the list of available Pickles
		response, err := http.Get(KuteGoAPIURL + "/pickles/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			// Transform our response to a []byte
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			// Put only needed information of the JSON document in our array of Pickle
			var data []Pickle
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}

			// Create a string with all of the Pickle's name and a blank line as seperator
			var pickles strings.Builder
			for _, pickle := range data {
				pickles.WriteString(pickle.Name + "\n")
			}

			// Send a text message with the list of Pickles
			_, err = s.ChannelMessageSend(m.ChannelID, pickles.String())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get list of Pickles! (ﾉ◕ヮ◕)ﾉ*:･ﾟ✧")
		}
	}
}
