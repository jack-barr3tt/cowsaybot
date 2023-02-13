package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	if os.Getenv("TOKEN") == "" {
		println("TOKEN not set in .env file")
		return
	}

	bot, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		println(err)
	}

	err = bot.Open()

	if err != nil {
		println("Unable to log in. Check your token is correct.")
		return
	}

	println("Bot started!")

	// Add the handler for when messages are sent
	bot.AddHandler(messageCreate)

	// Close the bot when the program is closed
	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	println("Press Ctrl+C to exit")
	<-stop

	println("Shutting down...")
}

func cowsay(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	// Execute the cowsay command
	cmd := exec.Command("cowsay", msg)
	stdout, err := cmd.Output()

	if err != nil {
		// If cowsay is not found, send an error message
		if err.Error() == "exec: \"cowsay\": executable file not found in $PATH" {
			s.ChannelMessageSend(m.ChannelID, "`cowsay` not found on your system. Please install it.")
			return
		}
		println(err.Error())
	}

	s.ChannelMessageSend(m.ChannelID, "```"+string(stdout)+"```")
}

func moo(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Split on spaces, grab the tail, then join on spaces
	msg := strings.Join(strings.Split(m.Content, " ")[1:], " ")

	if msg == "ping" {
		// Replace the message with the ping
		msg = strconv.Itoa(int(m.Timestamp.Unix()-time.Now().Unix())) + " ms"
	}

	cowsay(s, m, msg)
}

func mooquote(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Create a request to the zenquotes API
	req, err := http.NewRequest("GET", "https://zenquotes.io/api/random", nil)

	if err != nil {
		println(err)
	}

	// Send the request
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		println(err)
	}

	defer res.Body.Close()

	var result []map[string]interface{}

	// Get the response body
	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		println(err)
	}

	// Parse the response body
	json.Unmarshal(bytes, &result)

	// The quote is in the q key in the first item of the array
	msg := result[0]["q"].(string)

	cowsay(s, m, msg)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// moo command
	if strings.HasPrefix(m.Content, "moo ") {
		moo(s, m)
		// mooquote command
	} else if strings.HasPrefix(m.Content, "mooquote") {
		mooquote(s, m)
	}
}
