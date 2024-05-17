package ws

import (
	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
)

type Room struct {
	Id       string             `json:"id"`
	Clients  map[string]*Client `json:"clients"`
	UserId   int64              `json:"user_id,omitempty"`
	DoctorId int64              `json:"doctor_id,omitempty"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	initialRooms := make(map[string]*Room)
	initialRooms[constants.DoctorRole] = &Room{
		Id:      constants.DoctorRole,
		Clients: make(map[string]*Client),
	}

	return &Hub{
		Rooms:      initialRooms,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomId]; ok {
				r := h.Rooms[cl.RoomId]

				if _, ok := r.Clients[cl.Id]; !ok {
					r.Clients[cl.Id] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomId]; ok {
				if _, ok := h.Rooms[cl.RoomId].Clients[cl.Id]; ok {
					if len(h.Rooms[cl.RoomId].Clients) != 0 {

						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomId:   cl.RoomId,
							Id:       cl.Id,
							UserId:   cl.UserId,
							UserRole: cl.UserRole,
							Type:     "read",
						}
					}

					delete(h.Rooms[cl.RoomId].Clients, cl.Id)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomId]; ok {

				for _, cl := range h.Rooms[m.RoomId].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
