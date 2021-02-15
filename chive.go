// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
    "log"
)

type CmdConfigChive struct {
    Reg     *regexp.Regexp
    Enabled bool	`json:"enable"`
    ApiKey	string	`json:"api_key"`
}

var chiveURLs []string
var chiveCaptions []string

func NewCmdChive(config *CmdConfigChive) {
    config.Reg = regexp.MustCompile(`^/chive(?:(@[a-zA-Z0-9_]{1,20}bot)?( refill)?$)`)
    //    config.Reg = regexp.MustCompile(`^/chive(?:(@[a-zA-Z0-9_]+bot)?( .+)?$)`)
}

func ChiveRun(in string) (string){
    if (in == "/chive refill"){
        log.Println("Refilling the Chive pool because of a request")
        chiveURLs = []string{}
	chiveCaptions = []string{}
    }
    if len(chiveURLs)==0{
        ChiveRefill()
    }else if (len(chiveURLs)<6){
        log.Println("Refilling the Chive pool...")
        go ChiveRefill()
    }
    if len(chiveURLs)==0{
        return "Chive Error: No pics in the pool!"
    }
    picUrl := chiveURLs[0]
    chiveURLs = chiveURLs[1:len(chiveURLs)]
    picCaption := chiveCaptions[0]
    chiveCaptions = chiveCaptions[1:len(chiveCaptions)]
    //    resp, err := http.Get(picUrl)
    //	if err != nil {
    //		return "", err
    //	}
    //	defer resp.Body.Close()
    //	imgData, err := ioutil.ReadAll(resp.Body)
    //	if err != nil {
    //		return "", err
    //	}
    return "Keep calm and chive on... " + picCaption + "\n" + picUrl
}

func ChiveRefill(){
    type T_img struct{
	URL string
    }
    type Attachment struct {
	Caption string
	Image T_img
    }
    type Item struct {
        Attachments []Attachment
    }
    type t_posts struct {
        Items []Item
    }

    resp, err := http.Get("https://api4.thechive.com/api4/category?category_id=10665&page=1&key=" + Commandssl.Chive.ApiKey)

    if err != nil {
        log.Println(err)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        log.Println("HTTP error: " + resp.Status + " " + strconv.Itoa(resp.StatusCode))
        return
    }
    repBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        return
    }
    posts := t_posts{}
    err = json.Unmarshal(repBody, &posts)
    if err != nil {
        log.Println(err)
        return
    }
    item_num := rand.Intn(len(posts.Items))
 //   log.Println("len(post.Items) = ", len(posts.Items), "item_num = ", item_num)
    for i := len(posts.Items[item_num].Attachments); i>0; i-- {
	attach_num := rand.Intn(i)
//	log.Println("attach_num = ", attach_num, " ; i = ", i, " len Attachments = ", len(posts.Items[item_num].Attachments))
	item := posts.Items[item_num].Attachments[attach_num]
        chiveURLs = append(chiveURLs, item.Image.URL)
	chiveCaptions = append(chiveCaptions, item.Caption)
	posts.Items[item_num].Attachments[attach_num] = posts.Items[item_num].Attachments[len(posts.Items[item_num].Attachments)-1]
	posts.Items[item_num].Attachments = posts.Items[item_num].Attachments[:len(posts.Items[item_num].Attachments)-1]
    }
 /*   //Coger un post aleatorio de la página aleatoria seleccionada
    randPost := rand.Intn(40)
    if len(category.Posts) < randPost+1 {
        log.Println("Posts argument empty!")
        return
    }
    postNum := category.Posts[randPost].Guid

    resp, err = http.Get("http://api.thechive.com/api/post/" + strconv.Itoa(postNum) + "?key=" + Commandssl.Chive.ApiKey)
    if err != nil {
        log.Println(err)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        log.Println("HTTP error: " + resp.Status + " " + strconv.Itoa(resp.StatusCode))
        return
    }
    repBody, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        return
    }
    err = json.Unmarshal(repBody, &post)
    if err != nil {
        log.Println(err)
        return
    }

    //Cojo todas las fotos
    for _, url := range post.Posts[0].Items{
        chiveURLs = append(chiveURLs, url.Identity.URL)
    }*/
}
