package chat

type Hub struct {
	Clients      map[*Client]bool
	RegisterCh   chan *Client
	UnregisterCh chan *Client
	BroadcastCh  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]bool),
		RegisterCh:   make(chan *Client),
		UnregisterCh: make(chan *Client),
		BroadcastCh:  make(chan []byte),
	}
}

func (h *Hub) RunLoop() {
	for {
		select {
		case client := <-h.RegisterCh:
			h.Clients[client] = true
		case client := <-h.UnregisterCh:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.sendCh)
			}
		case message := <-h.BroadcastCh:
			for client := range h.Clients {
				select {
				case client.sendCh <- message:
				default:
					close(client.sendCh)
					delete(h.Clients, client)
				}
			}
		}
	}
}
