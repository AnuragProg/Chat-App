package server

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"chat-app/room"
	"github.com/gorilla/websocket"
)

var mut sync.Mutex

type Server struct{
	UserCount uint64
	Rooms map[string]*room.Room
	SocketVsName map[*websocket.Conn]string
	SocketVsRoom map[*websocket.Conn]*room.Room
}

func NewServer() *Server{
	return &Server{
		UserCount: 0,
		Rooms : make(map[string]*room.Room),
		SocketVsName: make(map[*websocket.Conn]string),
		SocketVsRoom:make( map[*websocket.Conn]*room.Room),
	}
}

func (s *Server)AddUser(ws *websocket.Conn, roomname string){
	mut.Lock()
	defer mut.Unlock()

	// Adding user socket with their name for later access to send message
	// with proper username designation
	atomic.AddUint64(&s.UserCount, 1)
	username := fmt.Sprintf("User %v", s.UserCount)

	s.SocketVsName[ws] = username
	
	// Looking for room corresponding to the roomname
	roomname = strings.Trim(strings.ToLower(roomname), " ")
	if room, exists := s.Rooms[roomname]; exists{
		go room.JoinUser(ws)
		s.SocketVsRoom[ws] = room
	}else{
		// Creating new room and updating the server with new data
		s.createNewRoomAndJoinUser(ws, roomname)
	}

	go s.SendMsg(websocket.TextMessage, "Joined", roomname, ws)
}


func (s *Server)RemoveUser(ws *websocket.Conn){
	mut.Lock()
	defer mut.Unlock()

	room := s.SocketVsRoom[ws]

	delete(s.SocketVsName, ws)
	delete(s.SocketVsRoom, ws)

	go room.RemoveUser(ws)
}

func (s *Server)createNewRoomAndJoinUser(ws *websocket.Conn, roomname string){

	// Adding the newly created room to the rooms list
	// Attaching the user to that room
	room := room.NewRoom(roomname)
	s.Rooms[room.Name] = room
	s.SocketVsRoom[ws] = room

	// Adding the user to that room
	go room.JoinUser(ws)
}

// Send msg to the room
func (s *Server)SendMsg(mt int, msg string, roomname string, sender *websocket.Conn){
	roomname = strings.Trim(roomname, " ")
	room := s.Rooms[roomname]
	username := s.SocketVsName[sender]
	go room.SendMsg(sender, mt, fmt.Sprintf("%v: %v", username, msg))
}