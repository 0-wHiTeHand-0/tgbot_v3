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
    "strings"
)

const picsURL = "http://ano.lolcathost.org/pics/"

type CmdConfigAno struct {
    Reg         *regexp.Regexp
    Enabled     bool `json:"enable"`
}

func NewCmdAno(config *CmdConfigAno) {
    config.Reg = regexp.MustCompile(`^(?:.+$)`)
}

func AnoRunRandom() ([]string,error) {

    var respData struct {
        Pics []struct {
            ID string
        }
    }
    reqData := struct {
        Method  string  `json:"method"`
        Num     int     `json:"num"`
    }{
        "random",
        15,
    }
    reqBody, err := json.Marshal(reqData)
    if err != nil {
        return []string{}, err
    }
    resp, err := http.Post("http://ano.lolcathost.org/json/pic.json",
    "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return []string{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []string{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
    }
    repBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []string{}, err
    }
    err = json.Unmarshal(repBody, &respData)
    if err != nil {
        return []string{}, err
    }
    var IDs []string
    for _, tmp := range respData.Pics{
        IDs = append(IDs, picsURL + tmp.ID)
    }
    return IDs, nil
}

func AnoRunTags(in string) ([]string, error){

    var respData struct {
        Pics []struct {
            ID string
        }
    }
    var IDs []string

    reqData := struct {
        Method string   `json:"method"`
        Tags   []string `json:"tags"`
        Limit  int      `json:"limit"`
    }{
        "searchRelated",
        strings.SplitN(in, ",", -1),
        25,
    }

    reqBody, err := json.Marshal(reqData)
    if err != nil {
        return []string{}, err
    }
    resp, err := http.Post("http://ano.lolcathost.org/json/tag.json",
    "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return []string{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []string{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
    }
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []string{}, err
    }
    err = json.Unmarshal(respBody, &respData)
    if err != nil {
        return []string{}, err
    }
    for _, tmp := range respData.Pics{
        IDs = append(IDs, picsURL + tmp.ID)
    }
    return IDs, nil
}
