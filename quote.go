package main

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "io/ioutil"
    "log"
    "math/rand"
    "os"
    "regexp"
    "strings"
)

type CmdConfigQuote struct {
    f_quotes       []string
    Reg            *regexp.Regexp
    Enabled bool   `json:"enable"`
    Path    string `json:"path"`
    Allowed []int  `json:"allow"`
}

func NewCmdQuote(config *CmdConfigQuote){
    if config.Path == "" {
        config.Path = "quotes.txt"
    }
    f, err := ioutil.ReadFile(config.Path) //Si es un fichero grande, se puede liar. En este caso es de mas o menos 3k. No problemo.
    if err != nil {
        log.Fatalln(err)
    }

		config.Reg = regexp.MustCompile(`^/quote(?:(@[a-zA-Z0-9_]{1,20}bot)?( [<>] (.|\n){1,500})?$)`)
    config.f_quotes = strings.Split(string(f), "\n")
}

func QuotesRun(inmsg *tgbotapi.Message) (string) {
    //Compruebo que chatID este permitido
    flag := false
    for _, i:=range Commandssl.Quotes.Allowed{
        if ((i == inmsg.From.ID)||(int64(i) == inmsg.Chat.ID)) {
            flag = true
            break
        }
    }
    if flag == false {
        return ""
    }

    m := strings.SplitN(inmsg.Text, " ", 3)
		str := ""
    if len(m) == 1 {
        rndInt := rand.Intn(len(Commandssl.Quotes.f_quotes))
        str = "*<-- Random quote -->*\n"+strings.Replace(Commandssl.Quotes.f_quotes[rndInt], "   ", "\n", -1)
    } else if len(m) == 3 && m[1] == ">" {
        quote := strings.Replace(m[2], "\n", "   ", -1)
        linesFiltered := make([]string, 0)
        for _, line := range Commandssl.Quotes.f_quotes {
            if strings.Contains(strings.ToLower(line), strings.ToLower(quote)) {
                linesFiltered = append(linesFiltered, line)
            }
        }
        if len(linesFiltered) == 0 {
            str = "No quote found :("
        } else {
            rndInt := rand.Intn(len(linesFiltered))
            str = "*<-- Match! -->*\n"+strings.Replace(linesFiltered[rndInt], "   ", "\n", -1)
        }
    } else if len(m) == 3 && m[1] == "<" {
        quote := strings.Replace(m[2], "\n", "   ", -1)
        f, err := os.OpenFile(Commandssl.Quotes.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
        if _, err = f.WriteString(quote + "\n"); err != nil {
            log.Println(err)
            return "mmm...It was an error writing to disk :/"
        }else{f.Close()}

        Commandssl.Quotes.f_quotes = append(Commandssl.Quotes.f_quotes, quote)
        str = "*<-- We have a new quote! -->*\n"+strings.Replace(quote, "   ", "\n", -1)
    }
		return str
}
