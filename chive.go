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

func NewCmdChive(config *CmdConfigChive) {
    config.Reg = regexp.MustCompile(`^/chive(?:(@[a-zA-Z0-9_]+bot)?( refill)?$)`)
    //    config.Reg = regexp.MustCompile(`^/chive(?:(@[a-zA-Z0-9_]+bot)?( .+)?$)`)
}

func ChiveRun(in string) (string){
    if (in == "/chive refill"){
        log.Println("Refilling the Chive pool because of a request")
        chiveURLs = []string{}
    }
    if len(chiveURLs)==0{
        ChiveRefill()
    }else if (len(chiveURLs)<15){
        log.Println("Refilling the Chive pool...")
        go ChiveRefill()
    }
    if len(chiveURLs)==0{
        return "Chive Error: No pics in the pool!"
    }
    picUrl := chiveURLs[0]
    chiveURLs = chiveURLs[1:len(chiveURLs)]
    //    resp, err := http.Get(picUrl)
    //	if err != nil {
    //		return "", err
    //	}
    //	defer resp.Body.Close()
    //	imgData, err := ioutil.ReadAll(resp.Body)
    //	if err != nil {
    //		return "", err
    //	}
    return "Keep calm and chive on...\n" + picUrl
}

func ChiveRefill(){
    var category struct {
        Post_Count struct {
            Total_Posts int
        }
        Posts []struct {
            Guid int
        }
    }

    var post struct {
        Posts []struct {
            Items []struct {
                Identity struct {
                    URL string
                }
            }
        }
    }

    //Coger una página aleatoria (cada página tiene 40 posts)
    randPage := rand.Intn(90) // Total de posts 3827 ahora mismo. Divido por 40 para obtener las paginas (da 95,675 que redondeo a 90)
    resp, err := http.Get("http://api.thechive.com/api/category/404664888?key=" + Commandssl.Chive.ApiKey + "&page=" + strconv.Itoa(randPage))
    if err != nil {
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
    err = json.Unmarshal(repBody, &category)
    if err != nil {
        log.Println(err)
        return
    }
    //Coger un post aleatorio de la página aleatoria seleccionada
    randPost := rand.Intn(40)
    //fmt.Println("LONGITUD: " + strconv.Itoa(len(category.Posts)))
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
    }
}
