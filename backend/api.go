package backend

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple/recon"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 30 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  nil,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

func StartAPI(ctx context.Context, wg *sync.WaitGroup, apEmit chan recon.ReconAPExt) {
	defer wg.Done()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET"},
	}))

	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	api.GET("/subscribe", wsHandler(apEmit))

	go r.Run(":8080")

	for {
		select {
		case <-ctx.Done():
			logrus.Warn("stopping api server...")
			return
		}
	}
}

func wsHandler(apEmit chan recon.ReconAPExt) func(*gin.Context) {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logrus.Errorf("failed to upgrade websocket connection: %v", err)
			return
		}
		logrus.Info("websocket connection opened")

		exit := make(chan bool)

		go func(c *gin.Context, exit chan bool) {
			for {
				_, _, err := ws.ReadMessage()
				if err != nil {
					exit <- true
					return
				}
			}
		}(c, exit)

		go func(c *gin.Context, exit chan bool) {
			for {
				select {
				case <-exit:
					// Connection has been closed
					logrus.Info("websocket connection closed")
					ws.Close()
					return
				case newAp := <-apEmit:
					// New AP's have come in, write them to the websocket
					err := ws.WriteJSON(newAp)
					if err != nil {
						logrus.Errorf("failed to send new AP via websocket: %v", err)
					}
				}
			}
		}(c, exit)
	}
}
