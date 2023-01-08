package main

import (
	"os"
	"fmt"
	"chat-app/server"
	"chat-app/api"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{}
var s *server.Server = server.NewServer()

// To be extracted from env file
var load_balancer string
var port string
var server_url string


func setupRouter() *gin.Engine{
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		room := c.Query("room")
		fmt.Println("Room is => ", room)
		if room == ""{
			fmt.Println("couldn't find room in query param")
			c.AbortWithStatus(404)
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil{
			fmt.Println(err.Error())
			c.AbortWithStatus(404)
			return
		}
		defer func(){
			go api.LoadReleased(load_balancer, server_url)
			s.RemoveUser(ws)
			ws.Close()
		}()

		s.AddUser(ws, room)

		for{
			mt, message, err := ws.ReadMessage()
			if err != nil{
				fmt.Println(err.Error())
				break;
			}
			go s.SendMsg(mt, string(message), room, ws)
		}
	})
	return router
}

func init(){
	if err:= godotenv.Load(); err != nil{
		// log.Fatal("couldn't load env file")		
		fmt.Println("couldn't load env file")

	}	
	load_balancer = os.Getenv("LOADBALANCER")
	port = os.Getenv("PORT")
	server_url = os.Getenv("SERVERURL")
	go api.AddServer(load_balancer, server_url)
}

func main(){
	router := setupRouter()
	defer func(){
		go api.RemoveServer(load_balancer, port)
	}()
	// fmt.Println("Port is ", port)
	router.Run(fmt.Sprintf(":%v", port))
}