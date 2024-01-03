package meeting_place

import (
	"context"
	"encoding/json"
	"log"
	"meetup/internal/repository/postgres/meeting_place"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader       = websocket.Upgrader{}
	placeClients   = make(map[*websocket.Conn]string)
	monitorClients = make(map[*websocket.Conn]string)
)

type dataNotFound struct {
	Message string `json:"message"`
	Err     bool   `json:"error"`
}

var notFoundMessage = dataNotFound{
	Message: "not found",
	Err:     true,
}

var jsonNotFoundMessage, _ = json.Marshal(notFoundMessage)

// ========= FOR PLACE =======
func (mc *Controller) sendDataToPlace(conn *websocket.Conn, placeID string) {
	data, er := mc.meeting_place.GetDetailByPlaceID(context.Background(), placeID)
	if er != nil {

		log.Println("Error getting meeting place data:", er)

		err := conn.WriteMessage(websocket.TextMessage, jsonNotFoundMessage)
		if err != nil {
			log.Println("WebSocket write failed:", err)
			conn.Close()
			delete(placeClients, conn)
		}
		return
	}

	responseJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Println("WebSocket write failed:", err)
		conn.Close()
		delete(placeClients, conn)
	}
}

func (mc *Controller) HandleWebSocketForPlace(c *gin.Context) {
	placeID := c.Query("place_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("WebSocket upgrade failed: ", err)
		return
	}
	defer conn.Close()

	placeClients[conn] = placeID

	mc.sendDataToPlace(conn, placeID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read failed:", err)
			break
		}

	}

	delete(placeClients, conn)
}

func (mc *Controller) broadcastPlace() {

	for conn := range placeClients {
		mc.sendDataToPlace(conn, placeClients[conn])
	}
}

// ========= FOR MONITOR (SHOWS ALL THE PERSON WITH PLACES) =======

func (mc *Controller) sendDataToMonitor(conn *websocket.Conn, meetingID string) {
	data, count, er := mc.meeting_place.MeetingPlaceList(context.Background(), meetingID)
	if er != nil {

		log.Println("Error getting meeting place data:", er)

		err := conn.WriteMessage(websocket.TextMessage, jsonNotFoundMessage)
		if err != nil {
			log.Println("WebSocket write failed:", err)
			conn.Close()
			delete(monitorClients, conn)
		}
		return
	}

	dataRespond := meeting_place.MeetingPlaceResponseForSocket{
		List:  data,
		Count: count,
	}

	responseJSON, err := json.Marshal(dataRespond)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Println("WebSocket write failed:", err)
		conn.Close()
		delete(monitorClients, conn)
	}
}

func (mc *Controller) HandleWebSocketForMonitor(c *gin.Context) {
	placeID := c.Query("meeting_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("WebSocket upgrade failed: ", err)
		return
	}
	defer conn.Close()

	monitorClients[conn] = placeID

	mc.sendDataToMonitor(conn, placeID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read failed:", err)
			break
		}

	}

	delete(monitorClients, conn)
}

func (mc *Controller) broadcastMonitor() {
	for conn := range monitorClients {
		mc.sendDataToMonitor(conn, monitorClients[conn])
	}
}
