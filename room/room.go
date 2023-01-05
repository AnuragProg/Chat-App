package room

import (
	"sync"
	"strings"
	"github.com/gorilla/websocket"
)

type Room struct{
	Name string
	Users []*websocket.Conn
	LastMessages []string
}

var userMut = sync.Mutex{}
var msgMut = sync.Mutex{}


func NewRoom(roomname string) *Room{
	return &Room{
		Name: strings.Trim(roomname, " "),
		Users: []*websocket.Conn{},
		LastMessages: []string{},
	}
}

func (r *Room)JoinUser(ws *websocket.Conn){
	userMut.Lock()
	msgMut.Lock()
	defer func(){
		msgMut.Unlock()
		userMut.Unlock()
	}()
	r.Users = append(r.Users, ws)
	for _, msg := range r.LastMessages{
		ws.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func (r *Room)RemoveUser(user *websocket.Conn){
	userMut.Lock()
	defer userMut.Unlock()

	for i, ws := range r.Users{
		if ws == user{
			r.Users[i] = r.Users[len(r.Users)-1]
			r.Users = r.Users[:len(r.Users)-1]
		}
	}
}

// Msg should be of format "SenderName:Msg"
func (r *Room)SendMsg(sender *websocket.Conn, mt int, msg string){
	msgMut.Lock()
	defer msgMut.Unlock()
	r.LastMessages = append(r.LastMessages, msg)
	if len(r.LastMessages) > 5 {
		r.LastMessages = r.LastMessages[1:]
	}
	for _, ws := range r.Users{
		if ws == sender {
			continue
		}	
		go ws.WriteMessage(mt, []byte(msg))
	}
}