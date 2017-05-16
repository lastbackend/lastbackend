//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package sockets

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Rooms map[string]*Room
}

func (h *Hub) NewConnection(id string, conn *websocket.Conn) *Client {
	var room *Room
	room = h.GetRoom(id)
	if room == nil {
		room = h.AddRoom(id)
	}

	client := &Client{
		Room: room,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	room.AddClient(client)
	return client
}

func (h *Hub) DelConection(id string, client *Client) {
	if room, ok := h.Rooms[id]; ok {
		room.DelClient(client)
	}
}

func (h *Hub) AddRoom(id string) *Room {
	h.Rooms[id] = NewRoom()
	go func() {
		h.Rooms[id].Listen()
	}()
	return h.Rooms[id]
}

func (h *Hub) GetRoom(id string) *Room {
	if room, ok := h.Rooms[id]; ok {
		return room
	}
	return nil
}

func (h *Hub) DelRoom(id string) {
	if len(h.Rooms[id].Clients) == 0 {
		delete(h.Rooms, id)
	}
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}
