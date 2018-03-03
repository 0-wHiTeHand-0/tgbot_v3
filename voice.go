package main

import (
//	"crypto/md5"
//	"encoding/hex"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
//    "os"
    //   "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CmdConfigVoice struct {
	Reg          *regexp.Regexp
	Enabled      bool   `json:"enable"`
	Espeak_param string `json:"espeak_param"`
}

func NewCmdVoice(config *CmdConfigVoice){
    config.Reg = regexp.MustCompile(`^/voice(?:(@[a-zA-Z0-9_]{1,20}bot)?( [ a-zÁÉÍÓÚáéíóúñÑA-Z0-9\.,?!]{1,500})?$)`)
}

func VoiceRun(txt string) ([]byte ,string) {
    m := strings.SplitN(txt, " ", 2)
    if len(m) < 2 {
        return []byte{}, "Sometimes silence is golden. Now it's not."
    }
    stdout1, err := exec.Command("espeak", Commandssl.Voice.Espeak_param, "--stdout", "-s125", m[1]).Output() //No he encontrado la manera de hacer un pipe multiple en go
    if err != nil {
        log.Println("Espeak error. If you want to speak, you must install it.")
        return []byte{}, "Espeak error"
    }
    ex := exec.Command("opusenc", "-", "-")
    stdin2, _ := ex.StdinPipe()
    stdout2, _ := ex.StdoutPipe()
    err = ex.Start()
    if err != nil {
        log.Println("Opusenc error. If you want to speak, you must install it.")
        return []byte{}, "Opusenc error"
    }
    stdin2.Write(stdout1)
    stdin2.Close()
    grepbytes, _ := ioutil.ReadAll(stdout2)
    ex.Wait()
   // f, _ := os.Create("/tmp/dat2.ogg")
   // defer f.Close()
   // f.Write(grepbytes)
   // f.Sync()
    return grepbytes, ""
}
