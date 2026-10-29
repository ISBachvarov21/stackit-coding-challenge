package routes

import (
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/julienschmidt/httprouter"
)

var db *pg.DB
var ownerUUID string

// Register wires up all route handlers on the given router.
func Register(router *httprouter.Router, database *pg.DB, owner string) {
	db = database
	ownerUUID = owner

	router.GET("/health", healthCheckHandler)
	router.POST("/event", createEventHandler)
	router.DELETE("/event/:id", deleteEventHandler)
	router.GET("/event/:id/tickets/amount", getAmountBookedTicketsHandler)
	router.GET("/event/:id/tickets", getEventBookedTicketsHandler)
	router.GET("/event/:id", getEventHandler)
	router.GET("/events", getAllEventsHandler)
	router.POST("/event/:id/tickets", bookTicketHandler)
	router.GET("/user/tickets", getUserBookedTicketsHandler)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
