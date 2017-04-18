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

package wss

import "fmt"

type Room struct {
	// Registered clients.
	Clients map[*Client]bool
	// Outbound messages to the clients.
	Broadcast chan []byte
}

func (r *Room) AddClient(client *Client) {
	fmt.Println("Add new websocket client")
	r.Clients[client] = true
}

func (r *Room) DelClient(client *Client) {
	if _, ok := r.Clients[client]; ok {
		fmt.Println("Delete websocket client")
		close(client.Send)
		delete(r.Clients, client)
	}
}

func (r *Room) Listen() {
	for {
		select {
		case message := <-r.Broadcast:

			fmt.Println("receive new message for broadcast")
			fmt.Printf("Total connected clients: %d \n", len(r.Clients))
			for client := range r.Clients {
				fmt.Println("send client message")
				client.Send <- message
				fmt.Println("message update sended")
			}
		}
	}
}

func NewRoom() *Room {
	return &Room{
		Broadcast: make(chan []byte),
		Clients:   make(map[*Client]bool),
	}
}
