package main

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

type Broadcast string

var Responses = [11]string{
	"Want to go fishing?",
	"Im gonna go grill!",
	"Gotta go get some milk...",
	"Wanna throw the pigskin around?",
	"Good stuff sport!",
	"Good stuff champ!",
	"Hi $user! I'm dad!",
	"Want a burger?",
	"Gotta go to work...",
	"Get off my lawn!",
	"Hey $user, want to go for a ride in my jeep?"}
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
	Embed    Broadcast = "EMBED"
	Plain    Broadcast = "PLAIN"
	DontSend Broadcast = "DONTSEND"
	Meme     Broadcast = "MEME"
	NumMemes int       = 10
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
	id     string
	joke   string
	status int
}

func HandleMsg(session *discord.Session, msg *discord.MessageCreate) {
	if len(msg.Mentions) > 0 {
		if msg.Mentions[0].ID == "772951738247544882" {
			var str, im, br = dadRequest(&msg.Content)
			if br == Meme && im != nil {
				_, _ = session.ChannelFileSend(msg.ChannelID, "dad-meme.jpg", im)
				return
			}
			var _, err = session.ChannelMessageSend(msg.ChannelID, strings.Replace(str, "$user", msg.Author.Username, -1))
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
			_ = bot.UpdateListeningStatus("the election! I am voting Kanye!")
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

func dadRequest(content *string) (string, io.Reader, Broadcast) {
	var dc = strings.ToLower(*content)
	if contains(dc, "i love you") {
		return "I don't love you", nil, Plain
	} else if contains(dc, "joke") {
		return Jokes[rand.Intn(len(Jokes))], nil, Plain
	} else if contains(dc, "help") {
		return "Say \"joke\" to get a dad joke, or anything else to get other responses", nil, Embed
	} else if contains(dc, "exit") {
		return "You will never leave", nil, Plain
	} else if contains(dc, "meme") {
		r := getMeme()
		return "", r, Meme
	} else if contains(dc, "weeb") {
		return "bruh", nil, Plain
	} else if contains(dc, "kill me") {
		return "Ok!", nil, Plain
	} else {
		return Responses[rand.Intn(len(Responses))], nil, Plain
	}
}
func getMeme() io.Reader {
	var num = rand.Intn(NumMemes)
	f, err := os.Open("assets/m-" + strconv.Itoa(num) + ".jpg")
	HandleErr(err)
	return f
}
