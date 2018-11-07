package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tungsheng/gohook/config"
	"github.com/tungsheng/gohook/model"
	"github.com/tungsheng/gohook/router/middleware/header"
	"github.com/tungsheng/gohook/router/middleware/logger"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
)

const color = 14177041

type discordMsg struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	//log.Info().Msg("Global init")
}

// Load initializes the routing of the application.
func Load() http.Handler {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(logger.SetLogger())
	e.Use(header.Options)

	root := e.Group("/")
	{
		root.GET("/test", handleTest)
		root.GET("/disc", handleDiscordGet)
		root.POST("/bitbucket", handleBitBucket)
		root.POST("/discord/:id/:token", handleDiscord)
	}

	return e
}

func post(url string, jsonData []byte) string {
	var jsonStr = []byte(jsonData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func handleTest(c *gin.Context) {
	log.Info().Msg("Discord test!")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong test",
	})
}

func handleDiscordGet(c *gin.Context) {
	accessToken := c.Query("access_token")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("access_token = %s", accessToken),
	})
}

func handleBitBucket(c *gin.Context) {
	var payload bitbucket.RepoPushPayload
	c.BindJSON(&payload)

	author := model.Author{
		Name:    payload.Actor.DisplayName,
		URL:     payload.Actor.Links.HTML.Href,
		IconURL: payload.Actor.Links.Avatar.Href,
	}
	footer := model.Footer{
		Text:    "Powered by Gohook",
		IconURL: "https://bitbucket.org/account/torchchurch/avatar/",
	}
	emd := model.Embed{
		Author:      author,
		Title:       payload.Push.Changes[0].New.Target.Message,
		URL:         payload.Push.Changes[0].New.Target.Links.HTML.Href,
		Description: fmt.Sprintf("%s pushed to %s", author.Name, payload.Push.Changes[0].New.Name),
		Color:       color,
		Footer:      footer,
	}
	wh := model.Webhook{
		Embeds: []model.Embed{emd},
	}

	c.JSON(http.StatusOK, gin.H{
		"webhook": wh,
	})
}

func handleDiscord(c *gin.Context) {
	var payload bitbucket.RepoPushPayload
	c.BindJSON(&payload)

	author := model.Author{
		Name:    payload.Actor.DisplayName,
		URL:     payload.Actor.Links.HTML.Href,
		IconURL: payload.Actor.Links.Avatar.Href,
	}
	footer := model.Footer{
		Text:    "Powered by Gohook",
		IconURL: "https://bitbucket.org/account/torchchurch/avatar/",
	}
	emd := model.Embed{
		Author:      author,
		Title:       payload.Push.Changes[0].New.Target.Message,
		URL:         payload.Push.Changes[0].New.Target.Links.HTML.Href,
		Description: fmt.Sprintf("%s pushed to %s", author.Name, payload.Push.Changes[0].New.Name),
		Color:       color,
		Footer:      footer,
	}
	wh := model.Webhook{
		Embeds: []model.Embed{emd},
	}

	id := c.Param("id")
	token := c.Param("token")
	url := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", id, token)
	whJSON, _ := json.Marshal(wh)

	var jsonStr = []byte(whJSON)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Gohook-Bitbucket")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful!",
		"wh":      wh,
		"body":    string(body),
	})

}
