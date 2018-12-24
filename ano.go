// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
//    "strings"
//    "strconv"
)

const picsURL = "http://ano.lolcathost.org/pics/"

type CmdConfigAno struct {
    Reg1         *regexp.Regexp
    Reg2        *regexp.Regexp
    Enabled     bool `json:"enable"`
}

func NewCmdAno(config *CmdConfigAno) {
    config.Reg1 = regexp.MustCompile(`^(?:.+$)`)
    config.Reg2 = regexp.MustCompile(`^/ano(?:(@[a-zA-Z0-9_]{1,20}bot)?$)`)
}

func AnoRunRandom(setnum int) ([]string, []int, []int, error) {

    type sPic struct{
        ID string `json:"id"`
        Sizew int `json:"sizew"`
        Sizeh int `json:"sizeh"`
    }
    var respData struct {
        Pics []sPic `json:"pics"`
        Pic sPic `json:"pic"`
    }
    reqData := struct {
        Method  string  `json:"method"`
        Num     int     `json:"num"`
    }{
        "random",
        setnum,
    }
    reqBody, err := json.Marshal(reqData)
    if err != nil {
        return []string{}, []int{},[]int{}, err
    }
    resp, err := http.Post("http://ano.lolcathost.org/json/pic.json",
    "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return []string{}, []int{}, []int{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []string{}, []int{}, []int{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
    }
    repBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []string{}, []int{}, []int{}, err
    }
    err = json.Unmarshal(repBody, &respData)
    if err != nil {
        return []string{}, []int{}, []int{}, err
    }
    var IDs []string
    var width,height []int
    if setnum == 1{
        IDs = append(IDs, picsURL+respData.Pic.ID)
        height = append(height, respData.Pic.Sizeh)
        width = append(width, respData.Pic.Sizew)
    }else{
        for _, tmp := range respData.Pics{
            IDs = append(IDs, picsURL + tmp.ID)
            height = append(height, tmp.Sizeh)
            width = append(width, tmp.Sizew)
        }
    }
    return IDs, height, width, nil
}
/*
func AnoRunTags(in string, offset string) ([]string, string, []int, []int, error){

    var respData struct {
        Pics []struct {
            ID string `json:"id"`
            Sizew int `json:"sizew"`
            Sizeh int `json:"sizeh"`
        }
        Total   string  `json:"total"`
    }
    var IDs []string
    var newoffset string

    reqData := struct {
        Method string   `json:"method"`
        Offset int      `json:"offset"`
        Tags   []string `json:"tags"`
        Limit  int      `json:"limit"`
    }{
        "searchRelated",
        -1,
        strings.SplitN(in, ",", -1),
        25,
    }

    ioffset, err := strconv.Atoi(offset)
    if (ioffset < 0) || (err != nil){
        reqData.Offset = 0
    }else{
        reqData.Offset = ioffset
    }
    reqBody, err := json.Marshal(reqData)
    if err != nil {
        return []string{}, "", []int{}, []int{}, err
    }
    resp, err := http.Post("http://ano.lolcathost.org/json/tag.json",
    "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return []string{}, "", []int{}, []int{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []string{}, "", []int{}, []int{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
    }
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []string{}, "", []int{}, []int{}, err
    }
    err = json.Unmarshal(respBody, &respData)
    if err != nil {
        return []string{}, "", []int{}, []int{}, err
    }

    var width,height []int
    for _, tmp := range respData.Pics{
        IDs = append(IDs, picsURL + tmp.ID)
        height = append(height, tmp.Sizeh)
        width = append(width, tmp.Sizew)
    }
    totalint , err := strconv.Atoi(respData.Total)
    newoffset = ""
    if ioffset + reqData.Limit < totalint{
        newoffset = strconv.Itoa(ioffset + reqData.Limit)
    }
    return IDs, newoffset, height, width, err
}*/
