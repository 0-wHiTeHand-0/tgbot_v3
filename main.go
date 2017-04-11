package main

import (
    "log"
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "github.com/go-telegram-bot-api/telegram-bot-api"
		"github.com/robfig/cron"
		"time"
    "math/rand"
   	"strings"
	 	"strconv"
    "crypto/md5"
    "encoding/hex"
)

func parseConfig(file string) (config, error) {
    b, err := ioutil.ReadFile(file)
    if err != nil {
        return config{}, err
    }
    var cfg config
    if err := json.Unmarshal(b, &cfg); err != nil {
        return config{}, err
    }
    log.Println(cfg)
    return cfg, nil
}

type config struct {
    Token          string     `json:"token"`
    AllowedIDs     []int      `json:"allowed_ids"`
    Commands       cmdConfigs `json:"commands"`
}
type cmdConfigs struct {
    Ban     CmdConfigBan        `json:"ban"`
    Quotes  CmdConfigQuote      `json:"quotes"`
    Bitchx  CmdConfigBitchx     `json:"bitchx"`
    Ano     CmdConfigAno        `json:"ano"`
    Chive   CmdConfigChive      `json:"chive"`
    Voice   CmdConfigVoice      `json:"voice"`
		Fcdg		CmdConfigFcdg				`json:"4cdg"`
		Ctftime	CmdConfigCtftime		`json:"ctftime"`
}

func send_pic(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, path string, f bool){
    pic := tgbotapi.NewPhotoUpload(msg.Chat.ID, path)
    if (f){pic.ReplyToMessageID = msg.MessageID}
    bot.Send(pic)
    return
}

func bytes_to_filebytes(a []byte) tgbotapi.FileBytes{
    hasher := md5.New()
    hasher.Write(a)
    tmpid := hex.EncodeToString(hasher.Sum(nil))
    file := tgbotapi.FileBytes{
        Name: tmpid,
        Bytes: a,
    }
    return file
}

var Commandssl cmdConfigs

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "Usage: tgbot config")
        os.Exit(1)
    }
    rand.Seed(time.Now().UTC().UnixNano())
    configFile := os.Args[1]
    cfg, err := parseConfig(configFile)
    if err != nil {
        log.Fatalln(err)
    }
    Commandssl = cfg.Commands

    if (Commandssl.Ban.Enabled){NewCmdBan(&Commandssl.Ban)}
    if (Commandssl.Quotes.Enabled){NewCmdQuote(&Commandssl.Quotes)}
    if (Commandssl.Bitchx.Enabled){NewCmdBitchx(&Commandssl.Bitchx)}
    if (Commandssl.Ano.Enabled){NewCmdAno(&Commandssl.Ano)}
    if (Commandssl.Chive.Enabled){NewCmdChive(&Commandssl.Chive)}
    if (Commandssl.Voice.Enabled){NewCmdVoice(&Commandssl.Voice)}
    if (Commandssl.Fcdg.Enabled){NewCmdFcdg(&Commandssl.Fcdg)}

    bot, err := tgbotapi.NewBotAPI(cfg.Token)
    if err != nil {
        log.Panic(err)
    }
    bot.Debug = false
    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

		if Commandssl.Ctftime.Enabled{
						c := cron.New()
						c.AddFunc("0 0 18 * * 3", func() {//Every thursday at 18.00
										err, txt := Ctftime_apireq()
										if err != nil{
														log.Println(err)
														return
										}
										for _, i := range Commandssl.Ctftime.Channel_ids{
														msg := tgbotapi.NewMessage(int64(i), txt)
														msg.ParseMode = "Markdown"
														bot.Send(msg)
										}
										log.Println("Sent CTFtime upcoming events!")
						})
						c.Start()
		}

		updates, err := bot.GetUpdatesChan(u)
		for update := range updates {
						if (update.Message!=nil){
										if banned_user(update.Message.From){
														log.Println("Banned user " + update.Message.From.FirstName + " blocked!")
														continue
										}
										log.Printf("Message -> [%s] %s (id: %s)", update.Message.From.FirstName, update.Message.Text, strconv.FormatInt(update.Message.Chat.ID, 10))
										var msg tgbotapi.MessageConfig
										if (Commandssl.Ban.Enabled && Commandssl.Ban.Reg.MatchString(update.Message.Text)){
														txt := BanRun(update.Message)
														if (txt!=""){
																		msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
																		msg.ReplyToMessageID = update.Message.MessageID
														}else{
																		send_pic(bot, update.Message, "a12.jpg" , true)
																		continue
														}
										}else if (Commandssl.Quotes.Enabled && Commandssl.Quotes.Reg.MatchString(update.Message.Text)){
														txt := QuotesRun(update.Message)
														if (txt!=""){
																		msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
																		msg.ParseMode = "Markdown"
														}else{
																		send_pic(bot, update.Message, "a12.jpg", true)
																		continue
														}
										}else if (Commandssl.Chive.Enabled && Commandssl.Chive.Reg.MatchString(update.Message.Text)){
														msg = tgbotapi.NewMessage(update.Message.Chat.ID, ChiveRun(update. Message.Text))
										}else if (Commandssl.Voice.Enabled && Commandssl.Voice.Reg.MatchString(update.Message.Text)){
														b, st := VoiceRun(update.Message.Text)
														if st == ""{
																		bot.Send(tgbotapi.NewVoiceUpload(update.Message.Chat.ID, bytes_to_filebytes(b)))
																		continue
														}else{
																		msg = tgbotapi.NewMessage(update.Message.Chat.ID, st)
														}
										}else if (Commandssl.Fcdg.Enabled && Commandssl.Fcdg.Reg.MatchString(update.Message.Text)){
														txt, err := FcdgRun(update.Message.Text)
														if err != nil{
																		log.Fatalln(err)
														}
														if (update.Message.Text == "/4cdg rules") {
																		msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
																		msg.ParseMode = "Markdown"
														}else{
																		send_pic(bot, update.Message, txt, false)
																		continue
														}
										}else if((update.Message.Text=="")||(update.Message.ReplyToMessage!=nil)){
														continue
										}else{
														m := "Command not found!"
														if (Commandssl.Bitchx.Enabled){m = BitchxRun(update.Message.From.FirstName)}
														msg = tgbotapi.NewMessage(update.Message.Chat.ID, m)
														msg.ReplyToMessageID = update.Message.MessageID
										}
										bot.Send(msg)
						}else if (update.InlineQuery!=nil){
										if banned_user(update.InlineQuery.From){
														log.Println("Banned user " + update.InlineQuery.From.FirstName + " blocked!")
														continue
										}
										log.Printf("Inline -> [%s] %s (alias: %s)", update.InlineQuery.From.FirstName, update.InlineQuery.Query, update.InlineQuery.From.UserName)
										var resu []interface{}
										var ano_urls []string
										var err error
										random := false
										if update.InlineQuery.Query == ""{
														ano_urls, err = AnoRunRandom()
														if err != nil{
																		log.Println(err)
																		continue
														}
														random = true
										}else if Commandssl.Ano.Reg.MatchString(update.InlineQuery.Query){
														ano_urls, err = AnoRunTags(update.InlineQuery.Query)
														if err != nil{
																		log.Println(err)
																		continue
														}
										}
										for _, ano_url := range ano_urls{
														hasher := md5.New()
														hasher.Write([]byte(ano_url))
														tmpid := hex.EncodeToString(hasher.Sum(nil))
														if (strings.Contains(ano_url, ".jpg")||strings.Contains(ano_url, ".png")||strings.Contains(ano_url, ".jpeg")){
																		tmp1:=tgbotapi.NewInlineQueryResultPhoto(tmpid, ano_url)
																		tmp1.ThumbURL = ano_url
																		resu = append(resu, tmp1)
														}else if strings.Contains(ano_url, ".gif"){
																		tmp1:=tgbotapi.NewInlineQueryResultGIF(tmpid, ano_url)
																		tmp1.ThumbURL = ano_url
																		resu = append(resu, tmp1)
														}else{
																		log.Println("ANO is returning strange replies: " + ano_url)
														}
										}
										if len(ano_urls)==0{
														log.Println("Bad syntax, or No pics found!")
										}
										inline := tgbotapi.InlineConfig{
														InlineQueryID: update.InlineQuery.ID,
														IsPersonal: false,
														CacheTime: 0,
														Results: resu,
										}
										if random{
														inline.NextOffset = "1"
										}
										if _, err := bot.AnswerInlineQuery(inline); err != nil {
														log.Println(err)
										}
						}
		}
}
