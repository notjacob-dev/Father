package main

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

type Broadcast string

var Responses = [6]string{
	"Want to go fishing?",
	"Im gonna go grill!",
	"Gotta go get some milk...",
	"Wanna throw the pigskin around?",
	"Good stuff sport!",
	"Good stuff champ!"}
var Jokes = [15]string{
	"\"Dad, did you get a haircut?\" \"No, I got them all cut!\"",
	"My wife is really mad at the fact that I have no sense of direction. So I packed up my stuff and right!",
	"How do you get a squirrel to like you? Act like a nut.",
	"Why don't eggs tell jokes? They'd crack each other up.",
	"\"I don't trust stairs. They're always up to something.\"",
	"What do you call someone with no body and no nose? Nobody knows.",
	"Did you hear the rumor about butter? Well, I'm not going to spread it!",
	"Why couldn't the bicycle stand up by itself? It was two tired.",
	"\"Dad, can you put my shoes on?\" \"No, I don't think they'll fit me.\"",
	"Why can't a nose be 12 inches long? Because then it would be a foot.",
	"This graveyard looks overcrowded. People must be dying to get in.",
	"What time did the man go to the dentist? Tooth hurt-y.",
	"How many tickles does it take to make an octopus laugh? Ten tickles.",
	"What concert costs just 45 cents? 50 Cent featuring Nickelback!",
	"How do you make a tissue dance? You put a little boogie in it."}


const (
	Embed Broadcast = "EMBED"
	Plain Broadcast = "PLAIN"
	DontSend Broadcast = "DONTSEND"
)

func contains(str string, substr string) bool {
	return strings.Contains(str, substr)
}

func readLine(file os.File) string {
	var bytes = make([]byte, 1024)
	var _, err = file.Read(bytes)
	if err == io.EOF {
		return ""
	}
	return string(bytes)
}

const dadJoke string = "https://icanhazdadjoke.com/"

type Joke struct {
	id string
	joke string
	status int
}


func HandleMsg(session *discord.Session, msg *discord.MessageCreate) {
	if len(msg.Mentions) > 0 {
		if msg.Mentions[0].ID == "772951738247544882" {
			var str, _ = dadRequest(&msg.Content)
			var _, err = session.ChannelMessageSend(msg.ChannelID, str)
			HandleErr(err)
		}
	}
}

func main() {
	start()
}

func start() {
	var file, ex = checkTokenFile()
	if ex {
		var tk = readLine(*file)
		if tk != "" {
			var bot, err = discord.New(fmt.Sprintf("Bot %s", strings.Replace(tk, "\x00", "", -1)))
			HandleErr(err)
			bot.AddHandler(HandleMsg)
			err = bot.Open()
			HandleErr(err)
			fmt.Println("Bot started!")
			sc := make(chan os.Signal, 1)
			signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
			<-sc
			_ = bot.Close()
			fmt.Println("Bot shut down!")
		}
	}
}

func checkTokenFile() (*os.File, bool) {
	var _, err = os.Stat("bot.token")
	if os.IsNotExist(err) {
		var file, err = os.Create("bot.token")
		HandleErr(err)
		return file, false
	} else {
		var file, err = os.OpenFile("bot.token", os.O_RDWR, 0644)
		HandleErr(err)
		return file, true
	}
}

func dadRequest(content *string) (string, Broadcast) {
	var dc = strings.ToLower(*content)
	if contains(dc, "i love you") {
		return "I don't love you", Plain
	} else if contains(dc, "joke") {
		return Jokes[rand.Intn(len(Jokes))], Plain
	} else if contains(dc, "help") {
		return "Say \"joke\" to get a dad joke, or anything else to get other responses", Embed
	} else if contains(dc, "exit") {
		return "You will never leave", Plain
	} else {
		return Responses[rand.Intn(len(Responses))], Plain
	}
}
