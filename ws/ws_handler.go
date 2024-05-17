package ws

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandlerOpts struct {
	Hub                 *Hub
	ConsultationUsecase usecases.ConsultationUsecase
}

type WebSocketHandler struct {
	hub                 *Hub
	ConsultationUsecase usecases.ConsultationUsecase
}

func NewWebSocketHandler(wshOpts *WebSocketHandlerOpts) *WebSocketHandler {
	return &WebSocketHandler{
		hub:                 wshOpts.Hub,
		ConsultationUsecase: wshOpts.ConsultationUsecase,
	}
}

type CreateRoomRequest struct {
	Id       int64 `json:"id"`
	DoctorId int64 `json:"doctor_id"`
}

func (h *WebSocketHandler) CreateRoom(ctx *gin.Context) {
	var payload CreateRoomRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	id, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}
	userId := int64(id)

	roomId := fmt.Sprintf("consult-%d", payload.Id)

	_, exists := h.hub.Rooms[roomId]
	if exists {
		ctx.Error(custom_errors.BadRequest(custom_errors.ErrRoomAlreadyExists, constants.RoomNotUniqueErrMsg))
		return
	}

	room := &Room{
		Id:       roomId,
		Clients:  make(map[string]*Client),
		UserId:   userId,
		DoctorId: payload.DoctorId,
	}

	h.hub.Rooms[roomId] = room

	ctx.JSON(http.StatusOK, RoomResponse{
		Id:       room.Id,
		UserId:   room.UserId,
		DoctorId: room.DoctorId,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) JoinRoomAsUser(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(err)
		return
	}

	userIdParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		conn.Close()
		return
	}

	userId := int64(userIdParam)

	consultationIdParam, err := strconv.Atoi(ctx.Param("consultationId"))
	if err != nil {
		ctx.Error(err)
		conn.Close()
		return
	}

	consultationId := int64(consultationIdParam)

	consultation, err := h.ConsultationUsecase.GetConsultationById(ctx, consultationId)
	if err != nil {
		ctx.Error(err)
		conn.Close()
		return
	}

	roomId := fmt.Sprintf("consult-%d", consultationId)

	room, exists := h.hub.Rooms[roomId]
	if !exists {
		newRoom := &Room{
			Id:       roomId,
			Clients:  make(map[string]*Client),
			UserId:   userId,
			DoctorId: consultation.Doctor.Id,
		}

		h.hub.Rooms[roomId] = newRoom
		room = newRoom
	}

	if room.UserId != userId {
		ctx.Error(custom_errors.InvalidAuthToken())
		conn.Close()
		return
	}

	userRole := constants.UserRole
	clientId := fmt.Sprintf("%d-%s", userId, userRole)

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		Id:       clientId,
		UserId:   userId,
		UserRole: userRole,
		RoomId:   roomId,
	}

	m := &Message{
		Content:  "user has joined the room",
		RoomId:   roomId,
		Id:       clientId,
		UserId:   userId,
		UserRole: userRole,
		Type:     "read",
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *WebSocketHandler) JoinRoomAsDoctor(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(err)
		return
	}

	doctorIdParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	doctorId := int64(doctorIdParam)

	consultationIdParam, err := strconv.Atoi(ctx.Param("consultationId"))
	if err != nil {
		ctx.Error(err)
		return
	}

	consultationId := int64(consultationIdParam)

	consultation, err := h.ConsultationUsecase.GetConsultationById(ctx, consultationId)
	if err != nil {
		ctx.Error(err)
		conn.Close()
		return
	}

	roomId := fmt.Sprintf("consult-%d", consultationId)

	room, exists := h.hub.Rooms[roomId]
	if !exists {
		newRoom := &Room{
			Id:       roomId,
			Clients:  make(map[string]*Client),
			UserId:   consultation.User.Id,
			DoctorId: doctorId,
		}

		h.hub.Rooms[roomId] = newRoom
		room = newRoom
	}

	if room.DoctorId != int64(doctorId) {
		ctx.Error(custom_errors.InvalidAuthToken())
		conn.Close()
		return
	}

	userRole := constants.DoctorRole
	clientId := fmt.Sprintf("%d-%s", doctorId, userRole)

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		Id:       clientId,
		UserId:   doctorId,
		UserRole: userRole,
		RoomId:   roomId,
	}

	m := &Message{
		Content:  "doctor has joined the room",
		RoomId:   roomId,
		Id:       clientId,
		UserId:   doctorId,
		UserRole: userRole,
		Type:     "read",
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *WebSocketHandler) JoinDoctorRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(err)
		return
	}

	roomId := constants.DoctorRole

	cl := &Client{
		Conn:    conn,
		Message: make(chan *Message, 10),
		Id:      uuid.New().String(),
		RoomId:  roomId,
	}

	m := &Message{
		Content: "user has joined the room",
		RoomId:  roomId,
		Type:    "read",
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomResponse struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	DoctorId int64  `json:"doctor_id"`
}

func (h *WebSocketHandler) GetRooms(ctx *gin.Context) {
	rooms := make([]RoomResponse, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomResponse{
			Id:       r.Id,
			UserId:   r.UserId,
			DoctorId: r.DoctorId,
		})
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    rooms,
	})
}

type ClientResponse struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id,omitempty"`
	UserRole string `json:"user_role,omitempty"`
	RoomId   string `json:"room_id"`
}

func (h *WebSocketHandler) GetClients(ctx *gin.Context) {
	var clients []ClientResponse
	roomId := ctx.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientResponse, 0)
		ctx.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientResponse{
			Id:       c.Id,
			UserId:   c.UserId,
			UserRole: c.UserRole,
			RoomId:   c.RoomId,
		})
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    clients,
	})
}
