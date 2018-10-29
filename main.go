package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/go-env"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
)

const (
	url = "https://discordapp.com/api/webhooks/495073722135478293/q__VqkG-ZSirBEXRHaVpZpJZsbK8cvsY0XMKPIxAgqrNRYcM-B1FRrpYwBVYMnmj_G7j"
)

type discordMsg struct {
	Content  string `json:"content"`
	Username string `json:"username"`
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

func main() {
	if env := os.Getenv("INN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	secret, _ := env.Get("BIT_SECRET")
	hook, _ := bitbucket.New(bitbucket.Options.UUID(secret))

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/discord/:id/:token", func(c *gin.Context) {
		log.Print("Webhooks reached!")

		payload, err := hook.Parse(
			r,
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
		dataJson, _ := json.Marshal(data)
		post(url, dataJson)

		w.Header().Set("Content-Type", "application/json")
		w.Write(dataJson)

		switch payload.(type) {
		case bitbucket.RepoPushPayload:
			repoPush := payload.(bitbucket.RepoPushPayload)
			// Do whatever you want from here...

			url := "https://discordapp.com/api/webhooks/495073722135478293/q__VqkG-ZSirBEXRHaVpZpJZsbK8cvsY0XMKPIxAgqrNRYcM-B1FRrpYwBVYMnmj_G7j"
			data := discordMsg{
				Content:  "test content",
				Username: "bitbucket",
			}
			dataJson, _ := json.Marshal(data)
			post(url, dataJson)

			w.Header().Set("Content-Type", "application/json")
			w.Write(dataJson)
			fmt.Printf("%+v", repoPush)

		case bitbucket.PullRequestCreatedPayload:
			pullRequestCreated := payload.(bitbucket.PullRequestCreatedPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequestCreated)
		}
	})

	r.Run(":9091")
}
