package main

import (
	"os"
	"log"
	"fmt"
	"chat-app/server"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{}

var s *server.Server = server.NewServer()

func setupRouter() *gin.Engine{
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		room := c.Query("room")
		if room == ""{
			c.AbortWithStatus(404)
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil{
			c.AbortWithStatus(404)
			return
		}
		defer func(){
			s.RemoveUser(ws)
			ws.Close()
		}()

		s.AddUser(ws, room)

		for{
			mt, message, err := ws.ReadMessage()
			if err != nil{
				break;
			}
			go s.SendMsg(mt, string(message), room, ws)
		}
	})
	return router
}

func main(){
	if err:= godotenv.Load(); err != nil{
		log.Fatal("couldn't load env file")		
	}	
	router := setupRouter()
	router.Run(fmt.Sprintf(":%v", os.Getenv("PORT")))
}