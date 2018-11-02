package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tungsheng/go-hook/config"
	"github.com/tungsheng/go-hook/router/middleware/header"
	"github.com/tungsheng/go-hook/router/middleware/logger"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
)

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
	c.JSON(200, gin.H{
		"message": "pong test",
	})
}

func handleDiscordGet(c *gin.Context) {
	accessToken := c.Query("access_token")
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("access_token = %s", accessToken),
	})
}

func handleDiscord(c *gin.Context) {

	hook, _ := bitbucket.New(bitbucket.Options.UUID("k3jhvK38dvBAkOk482P"))

	payload, err := hook.Parse(
		c.Request,
		bitbucket.RepoPushEvent,
		bitbucket.PullRequestCreatedEvent,
		bitbucket.PullRequestUpdatedEvent,
		bitbucket.PullRequestApprovedEvent,
		bitbucket.PullRequestUnapprovedEvent,
		bitbucket.PullRequestMergedEvent,
		bitbucket.PullRequestDeclinedEvent,
		bitbucket.PullRequestCommentCreatedEvent,
		bitbucket.PullRequestCommentUpdatedEvent,
		bitbucket.PullRequestCommentDeletedEvent,
	)

	if err != nil {
		if err == bitbucket.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}

	id := c.Param("id")
	token := c.Param("token")
	url := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", id, token)
	data := discordMsg{
		Content:  "test content",
		Username: "bitbucket",
	}

	dataJSON, _ := json.Marshal(data)
	post(url, dataJSON)
	pJSON, _ := json.Marshal(payload)

	c.JSON(200, gin.H{
		"message": "pong",
		"data":    pJSON,
	})

}
