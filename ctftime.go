package main

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "strconv"
)

type CmdConfigCtftime struct {
    Enabled bool   `json:"enable"`
    Channel_ids    []int `json:"channel_ids"`
}

func Ctftime_apireq() (error, string) {
    type Event struct {
        Onsite bool `json:"onsite"`
        Title string `json:"title"`
        Url	string `json:"url"`
        Participants int `json:"participants"`
        Format string `json:"format"`
        Start string `json:"start"`
        Finish string `json:"finish"`
    }
    res, err := http.Get("https://ctftime.org/api/v1/events/?limit=20")
    if err != nil {
        return err, ""
    }
    pjson, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        return err, ""
    }
    var events []Event
    err = json.Unmarshal(pjson, &events)
    if err != nil {
        return err, ""
    }
    txt := "*<--- CTFTime upcoming events --->*\n"
    for _, i := range events{
        if (!i.Onsite) && (i.Participants > 25){
            txt += "\nTitle: " + i.Title
            txt += "\nURL: " + i.Url
            txt += "\nParticipants: " + strconv.Itoa(i.Participants)
            txt += "\nFormat: " + i.Format
            txt += "\nStart: "
            loc, _ := time.LoadLocation("Europe/Madrid")
            t1, e := time.Parse(time.RFC3339,i.Start)
            if e != nil{
                txt += "Parser error"
            }else{
                txt += t1.In(loc).Format(time.UnixDate)
            }
            txt += "\nFinish: "

            t1, e = time.Parse(time.RFC3339,i.Finish)
            if e != nil{
                txt += "Parser error"
            }else{
                txt += t1.In(loc).Format(time.UnixDate)
            }
            txt += "\n"
        }
    }
    txt = txt[:len(txt)-1]
    return nil, txt
}
