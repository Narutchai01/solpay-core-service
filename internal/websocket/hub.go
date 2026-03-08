package websocket

import "log"

type TargetMessage struct {
	TxID    string
	Payload []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Notify     chan TargetMessage
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Notify:     make(chan TargetMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client registered for TxID: %s (total clients: %d)", client.TxID, len(h.Clients))

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client unregistered for TxID: %s (total clients: %d)", client.TxID, len(h.Clients))
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}

		case msg := <-h.Notify:
			log.Printf("Notifying transaction status for TxID: %s (total clients: %d)", msg.TxID, len(h.Clients))
			for client := range h.Clients {
				if client.TxID == msg.TxID {
					select {
					case client.Send <- msg.Payload:
						log.Printf("Sent notification to client for TxID: %s", msg.TxID)
					default:
						close(client.Send)
						delete(h.Clients, client)
						log.Printf("Client buffer full, removed client for TxID: %s", msg.TxID)
					}
				}
			}
		}
	}
}

func (h *Hub) NotifyTransactionStatus(txID string, payload []byte) {
	h.Notify <- TargetMessage{TxID: txID, Payload: payload}
}
