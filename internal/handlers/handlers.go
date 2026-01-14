package handlers

// Handlers holds all handler dependencies
type Handlers struct {
	*ClientHandler
	*HumanResourceHandler
}

// NewHandlers creates a new Handlers instance with all handler dependencies
func NewHandlers(clientHandler *ClientHandler, hrHandler *HumanResourceHandler) *Handlers {
	return &Handlers{
		ClientHandler:        clientHandler,
		HumanResourceHandler: hrHandler,
	}
}
