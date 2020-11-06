package main

import (
	"errors"
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

type Broadcast string

func GetState(session *discord.Session, user string) (*discord.VoiceState, error) {
	for _, g := range session.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == user {
				return v, nil
			}
		}
	}
	return nil, errors.New("no voice state for user")
}

func JoinVcSession(session *discord.Session, voice *discord.VoiceState) (*discord.VoiceConnection, error) {
	var v, e = session.ChannelVoiceJoin(voice.GuildID, voice.ChannelID, false, false)
	currentVoice = v
	return v, e
}

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

var currentVoice *discord.VoiceConnection = nil

const (
	Embed    Broadcast = "EMBED"
	Plain    Broadcast = "PLAIN"
	DontSend Broadcast = "DONTSEND"
	Meme     Broadcast = "MEME"
	NumMemes int       = 10
	Voice    Broadcast = "VOICE"
	NumSongs int       = 1
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

func PlayAudio(voice *discord.VoiceConnection) {

	r := rand.Intn(NumSongs)

	f, err0 := os.Open("assets/encoded/e-" + strconv.Itoa(r) + ".dca")
	HandleErr(err0)
	dc := dca.NewDecoder(f)

	_ = voice.Speaking(true)
	for {
		frame, err1 := dc.OpusFrame()
		if err1 != io.EOF {
			HandleErr(err1)
		} else {
			break
		}
		select {
		case voice.OpusSend <- frame:
		case <-time.After(time.Second):
			currentVoice = nil
			return
		}
	}
	_ = voice.Speaking(false)
	time.Sleep(250 * time.Millisecond)
	_ = voice.Disconnect()

}

func HandleMsg(session *discord.Session, msg *discord.MessageCreate) {
	if len(msg.Mentions) > 0 {
		if msg.Mentions[0].ID == "772951738247544882" {
			var str, im, br = dadRequest(&msg.Content)
			if br == Meme && im != nil {
				_, _ = session.ChannelFileSend(msg.ChannelID, "dad-meme.jpg", im)
				return
			} else if br == Voice {
				state, _ := GetState(session, msg.Author.ID)
				if state != nil {
					con, err := JoinVcSession(session, state)
					HandleErr(err)
					go PlayAudio(con)
				} else {
					session.ChannelMessageSend(msg.ChannelID, "You aren't in a voice channel! "+Responses[rand.Intn(len(Responses))])
				}
				return
			} else if br == DontSend {
				return
			}
			var _, err = session.ChannelMessageSend(msg.ChannelID, strings.Replace(str, "$user", msg.Author.Username, -1))
			HandleErr(err)
		}
	}
}

func ToDca(index int) {
	_, s := os.Stat("encoded.dca")
	if os.IsNotExist(s) {
		enc, _ := dca.EncodeFile("assets/sound/s-"+strconv.Itoa(index)+".mp3", dca.StdEncodeOptions)
		defer enc.Cleanup()
		out, err01 := os.Create("assets/encoded/e-" + strconv.Itoa(index) + ".dca")
		HandleErr(err01)
		io.Copy(out, enc)
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
			var files, _ = ioutil.ReadDir("assets/sound/")
			for _, f := range files {
				var str = strings.Split(f.Name(), "-")
				var in, _ = strconv.Atoi(strings.Split(str[1], ".")[0])
				ToDca(in)
			}
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
	} else if contains(dc, "music") {
		return "", nil, Voice
	} else if contains(dc, "voiceleave") {
		currentVoice.Disconnect()
		return "", nil, DontSend
	} else {
		return Responses[rand.Intn(len(Responses))], nil, Plain
	}
}
func getMeme() io.Reader {
	var num = rand.Intn(NumMemes)
	f, err := os.Open("assets/memes/m-" + strconv.Itoa(num) + ".jpg")
	HandleErr(err)
	return f
}
