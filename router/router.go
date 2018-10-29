package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-hook/router/middleware/header"
	"github.com/go-hook/router/middleware/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-ggz/ggz/config"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
)

func Load() http.Handler {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(logger.SetLogger())
	r.User(header.Options)

	root := r.Group("/")
	{
		root.Post("/discord/:id/:token", handleDiscord)
	}

	return r
}

func handleDiscord(c *gin.Context) {
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

}
