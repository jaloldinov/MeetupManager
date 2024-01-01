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
	Message: "not found user",
	Err:     true,
}

var jsonNotFoundMessage, _ = json.Marshal(notFoundMessage)

// ========= FOR PLACE =======
func (mc *Controller) sendDataToPlace(conn *websocket.Conn, placeID string) {
	data, er := mc.meeting_place.GetDetailByPlaceID(context.Background(), placeID)
	if er != nil {
		// log it and send an error response
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

	// Send the meeting place data to the connected client
	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Println("WebSocket write failed:", err)
		conn.Close()
		delete(placeClients, conn)
	}
}

func (mc *Controller) HandleWebSocketForPlace(c *gin.Context) {
	placeID := c.Query("place_id")

	// Upgrade the HTTP server connection to the WebSocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("WebSocket upgrade failed: ", err)
		return
	}
	defer conn.Close()

	// Add the new client to the clients map
	placeClients[conn] = placeID

	// Send the user data to the connected client when they connect
	mc.sendDataToPlace(conn, placeID)

	// Continuously read and write messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read failed:", err)
			break
		}

	}

	// Remove the client from the clients map when the connection is closed
	delete(placeClients, conn)
}

func (mc *Controller) broadcastPlace() {
	// Iterate over all connected clients
	for conn := range placeClients {
		// Send the updated  to client
		mc.sendDataToPlace(conn, placeClients[conn])
	}
}

// ========= FOR MONITOR (SHOWS ALL THE PERSON WITH PLACES) =======

func (mc *Controller) sendDataToMonitor(conn *websocket.Conn, meetingID string) {
	data, count, er := mc.meeting_place.MeetingPlaceList(context.Background(), meetingID)
	if er != nil {
		// log it and send an error response
		log.Println("Error getting meeting place data:", er)

		err := conn.WriteMessage(websocket.TextMessage, jsonNotFoundMessage)
		if err != nil {
			log.Println("WebSocket write failed:", err)
			conn.Close()
			delete(monitorClients, conn)
		}
		return
	}

	// Create a MeetingPlaceResponseForSocket structure
	dataRespond := meeting_place.MeetingPlaceResponseForSocket{
		List:  data,
		Count: count,
	}

	// Convert the entire structure to JSON
	responseJSON, err := json.Marshal(dataRespond)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}

	// Send the entire structure as a single JSON message to the connected client
	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Println("WebSocket write failed:", err)
		conn.Close()
		delete(monitorClients, conn)
	}
}

func (mc *Controller) HandleWebSocketForMonitor(c *gin.Context) {
	placeID := c.Query("meeting_id")

	// Upgrade the HTTP server connection to the WebSocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("WebSocket upgrade failed: ", err)
		return
	}
	defer conn.Close()

	// Add the new client to the clients map
	monitorClients[conn] = placeID

	// Send the user data to the connected client when they connect
	mc.sendDataToMonitor(conn, placeID)

	// Continuously read and write messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read failed:", err)
			break
		}

	}

	// Remove the client from the clients map when the connection is closed
	delete(monitorClients, conn)
}

func (mc *Controller) broadcastMonitor() {
	// Iterate over all connected clients
	for conn := range monitorClients {
		// Send the updated  to client
		mc.sendDataToMonitor(conn, monitorClients[conn])
	}
}
