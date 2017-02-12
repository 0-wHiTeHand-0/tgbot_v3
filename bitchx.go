package main

import (
    "io/ioutil"
    "log"
    "math/rand"
    "strings"
)

type CmdConfigBitchx struct {
    f_quotes       []string
    Enabled bool   `json:"enable"`
    Path    string `json:"path"`
}

func NewCmdBitchx(config *CmdConfigBitchx){
    if config.Path == "" {
        config.Path = "bitchx.txt"
    }
    f, err := ioutil.ReadFile(config.Path) //Si es un fichero grande, se puede liar. En este caso es de mas o menos 3k. No problemo.
    if err != nil {
        log.Fatalln(err)
    }
    config.f_quotes = strings.Split(string(f), "\n")
}

func BitchxRun(uname string) string {
        rndInt := rand.Intn(len(Commandssl.Bitchx.f_quotes))
        return strings.Replace(Commandssl.Bitchx.f_quotes[rndInt], "$0", uname, -1)
}
