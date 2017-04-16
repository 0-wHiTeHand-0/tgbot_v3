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
        go handle_updates(bot, update)
    }
}

func handle_updates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    if (update.Message!=nil){
        if banned_user(update.Message.From){
            log.Println("Banned user " + update.Message.From.FirstName + " blocked!")
            return
        }
        log.Printf("Message -> [%s] %s (id: %d, alias: %s)", update.Message.From.FirstName, update.Message.Text, update.Message.Chat.ID, update.Message.From.UserName)
        var msg tgbotapi.MessageConfig
        if (Commandssl.Ano.Enabled && Commandssl.Ano.Reg2.MatchString(update.Message.Text)){
            txt := "What have been seen cannot be unseen!\n"
            pic, _, _, err := AnoRunRandom(1)
            if err != nil || len(pic) != 1{
                txt += "Error requesting an ANO random pic."
                log.Println(err)
            }else{
                txt += pic[0]
            }
            msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
        }else if (Commandssl.Ban.Enabled && Commandssl.Ban.Reg.MatchString(update.Message.Text)){
            txt := BanRun(update.Message)
            if (txt!=""){
                msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
                msg.ReplyToMessageID = update.Message.MessageID
            }else{
                send_pic(bot, update.Message, "a12.jpg" , true)
                return
            }
        }else if (Commandssl.Quotes.Enabled && Commandssl.Quotes.Reg.MatchString(update.Message.Text)){
            txt := QuotesRun(update.Message)
            if (txt!=""){
                msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
                msg.ParseMode = "Markdown"
            }else{
                send_pic(bot, update.Message, "a12.jpg", true)
                return
            }
        }else if (Commandssl.Chive.Enabled && Commandssl.Chive.Reg.MatchString(update.Message.Text)){
            msg = tgbotapi.NewMessage(update.Message.Chat.ID, ChiveRun(update. Message.Text))
        }else if (Commandssl.Voice.Enabled && Commandssl.Voice.Reg.MatchString(update.Message.Text)){
            b, st := VoiceRun(update.Message.Text)
            if st == ""{
                bot.Send(tgbotapi.NewVoiceUpload(update.Message.Chat.ID, bytes_to_filebytes(b)))
                return
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
                return
            }
        }else if((update.Message.Text=="")||(update.Message.ReplyToMessage!=nil)){
            return
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
            return
        }
        log.Printf("Inline -> [%s] %s (alias: %s)", update.InlineQuery.From.FirstName, update.InlineQuery.Query, update.InlineQuery.From.UserName)
        var resu []interface{}
        var ano_urls []string
        var height,width []int
        var offset, newoffset string
        var err error
        if update.InlineQuery.Query == ""{
            newoffset = "1"
            ano_urls, height, width, err = AnoRunRandom(25)
            if err != nil{
                log.Println(err)
                return
            }
        }else if Commandssl.Ano.Reg1.MatchString(update.InlineQuery.Query){
            if update.InlineQuery.Offset == ""{
                offset = "0"
            }else{
                offset = update.InlineQuery.Offset
            }
            ano_urls, newoffset, height, width, err = AnoRunTags(update.InlineQuery.Query, offset)
            if err != nil{
                log.Println(err)
                return
            }
        }
        for i, ano_url := range ano_urls{
            hasher := md5.New()
            hasher.Write([]byte(ano_url))
            tmpid := hex.EncodeToString(hasher.Sum(nil))
            if (strings.Contains(ano_url, ".jpg")||strings.Contains(ano_url, ".png")||strings.Contains(ano_url, ".jpeg")){
                tmp1:=tgbotapi.NewInlineQueryResultPhoto(tmpid, ano_url)
                tmp1.ThumbURL = ano_url
                tmp1.Height = height[i]
                tmp1.Width = width[i]
                resu = append(resu, tmp1)
            }else if strings.Contains(ano_url, ".gif"){
                tmp1:=tgbotapi.NewInlineQueryResultGIF(tmpid, ano_url)
                tmp1.ThumbURL = ano_url
                tmp1.Height = height[i]
                tmp1.Width = width[i]
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
            NextOffset: newoffset,
        }
        if _, err := bot.AnswerInlineQuery(inline); err != nil {
            log.Println(err)
        }
    }
}
