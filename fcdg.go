package main

import (
		"path/filepath"
		"math/rand"
		"regexp"
)

type CmdConfigFcdg struct {
    f_slice       []string
		Reg						*regexp.Regexp
    Enabled bool   `json:"enable"`
    Path    string `json:"path_folder"`
}

func NewCmdFcdg(config *CmdConfigFcdg){
    if config.Path == "" {
        config.Path = "4cdg"
    }
		config.Reg = regexp.MustCompile(`^/4cdg(?:@[a-zA-Z0-9_]+bot| rules)?$`)
}

func FcdgRun(text string) (string, error) {
				var	err error
				if text == "/4cdg rules" {
								return "*'4chan drinking card game' rules*\n\n1. The left of the phone owner starts.\n2. Players take a card when it's their turn. They must do what the card says.\n3. You win the game when everyone else pass' out.\n\nCard types:\nAction: This is a standard 'do what it says' card.\nInstant: This card may be kept and used at anytime in the game.\nMandatory: Everyone must play this card.\nStatus: This is constant for the whole game or the timeframe indicated on the card.", nil
				}
				if len(Commandssl.Fcdg.f_slice) == 0 {
								Commandssl.Fcdg.f_slice, err = filepath.Glob(Commandssl.Fcdg.Path + "/*.jpg")
								if err != nil {
												return "", err
								}
								temp, err := filepath.Glob(Commandssl.Fcdg.Path + "/*.png") //Cochinada maxima, pero funciona.
								if err != nil {
												return "", err
								}
								Commandssl.Fcdg.f_slice = append(Commandssl.Fcdg.f_slice, temp...)
				}
				rndInt := rand.Intn(len(Commandssl.Fcdg.f_slice))
				imgName := Commandssl.Fcdg.f_slice[rndInt]
				Commandssl.Fcdg.f_slice = append(Commandssl.Fcdg.f_slice[:rndInt], Commandssl.Fcdg.f_slice[rndInt+1:]...)
				return imgName, nil
}
