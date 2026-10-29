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
	router.POST("/event", withOwner(createEventHandler))
	router.DELETE("/event/:id", withOwner(deleteEventHandler))
	router.GET("/event/:id/tickets/amount", withOwner(getAmountBookedTicketsHandler))
	router.GET("/event/:id/tickets", withOwner(getEventBookedTicketsHandler))
	router.GET("/event/:id", withUser(getEventHandler))
	router.GET("/events", withUser(getAllEventsHandler))
	router.POST("/event/:id/tickets", withUser(bookTicketHandler))
	router.GET("/user/tickets", withUser(getUserBookedTicketsHandler))
	router.DELETE("/ticket/:id", withUser(cancelTicketHandler))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
