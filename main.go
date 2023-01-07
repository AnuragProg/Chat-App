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

func setupRouter() *gin.Engine{
	router := gin.Default()


	// loadBalancer (format = https://url:port, http://url:port)
	loadBalancer := os.Getenv("LOADBALANCER")
	fmt.Println("Load balancer => ", loadBalancer)
	if loadBalancer == ""{
		panic("Counldn't find load balancer url")
	}

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
			go api.LoadReleased(loadBalancer)
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

func init(){
	if err:= godotenv.Load(); err != nil{
		// log.Fatal("couldn't load env file")		
		fmt.Println("couldn't load env file")

	}	
	go api.AddServer(os.Getenv("LOADBALANCER"))
}

func main(){
	router := setupRouter()
	defer func(){
		go api.RemoveServer(os.Getenv("LOADBALANCER"))
	}()
	port := os.Getenv("PORT")
	// fmt.Println("Port is ", port)
	router.Run(fmt.Sprintf(":%v", port))
}