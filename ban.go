package main

import (
    "regexp"
    "strconv"
    "strings"
    "log"
    "time"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

var BannedIDs map[int]time.Time
var BannedIDs_sg map[int]int

type CmdConfigBan struct {
    Enabled		bool	`json:"enable"`
    Allowed		[]int	`json:"allow"`
    Pre_Ban_ids	[]int	`json:"pre_banned_ids"`
    Default_time	int	`json:"default_time"`
    Reg         *regexp.Regexp
}

func NewCmdBan(config *CmdConfigBan) {
    config.Reg = regexp.MustCompile(`^/ban(?: ([0-9]+|-)$)`)
    BannedIDs = make(map[int]time.Time)
    BannedIDs_sg = make(map[int]int)
    for _, i := range config.Pre_Ban_ids{
        BannedIDs[i] = time.Now()
        BannedIDs_sg[i] = config.Default_time
    }
}

func banned_user(usu *tgbotapi.User) (bool) {
    if BannedIDs_sg[usu.ID] == 0{
        return false
    }
    tmp := time.Since(BannedIDs[usu.ID])
    if int64(tmp.Seconds()) < int64(BannedIDs_sg[usu.ID]) {
        log.Println("Warning: user " + usu.FirstName + " blocked")
        return true
    }else{
        BannedIDs[usu.ID] = time.Now()
        return false
    }
}

func BanRun(inmsg *tgbotapi.Message) string {
    flag := false
    for _, i:=range Commandssl.Ban.Allowed{
        if i==inmsg.From.ID{
            flag=true
            break
        }
    }
    if (!flag){
        return ""
    }

    m := strings.SplitN(inmsg.Text, " ", 2)
    var message string
    if m[1] == "-" {
        BannedIDs = make(map[int]time.Time)
        BannedIDs_sg = make(map[int]int)
        message = "Ban list cleared! Trolls are free again!"
    }else if ((inmsg.ReplyToMessage!=nil) && (len(m)==2)){
        BannedIDs[inmsg.ReplyToMessage.From.ID] = time.Now()
        BannedIDs_sg[inmsg.ReplyToMessage.From.ID], _ = strconv.Atoi(m[1])
        message = "Ban set to " + inmsg.ReplyToMessage.From.FirstName + " for " + m[1] + " seconds. Keep calm and relax your boobies, little troll."
    }else{
        message = "You have to reply a BotCommand message."
    }
    return message
}
