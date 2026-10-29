package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/go-pg/pg/v10"
	"github.com/julienschmidt/httprouter"
)

var (
	errNoTicketsAvailable = errors.New("no tickets available")
	errAlreadyBooked      = errors.New("ticket already booked")
)

func getEventBookedTicketsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	eventId := p.ByName("id")

	if err := db.Model(&models.Event{Id: eventId}).WherePK().Select(); err != nil {
		writeError(w, http.StatusNotFound, "Event with id "+eventId+" not found")
		return
	}

	// get booked tickets for event
	tickets := []models.Ticket{}
	if err := db.Model(&tickets).Where("event_id = ?", eventId).Select(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get booked tickets for event")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		log.Println("Failed to encode tickets:", err)
	}
}

func getAmountBookedTicketsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	eventId := p.ByName("id")

	// get event by id
	event := &models.Event{Id: eventId}
	if err := db.Model(event).WherePK().Select(); err != nil {
		writeError(w, http.StatusNotFound, "Event with id "+eventId+" not found")
		return
	}

	count, err := db.Model(&models.Ticket{}).Where("event_id = ?", eventId).Count()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get booked tickets for event")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Amount int `json:"amount"`
	}{
		Amount: count,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to encode ticket amount:", err)
	}
}

func bookTicketHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := userIdFromContext(r)
	eventId := p.ByName("id")

	event := &models.Event{Id: eventId}
	if err := db.Model(event).WherePK().Select(); err != nil {
		writeError(w, http.StatusNotFound, "Event with id "+eventId+" not found")
		return
	}

	ticket := &models.Ticket{}
	err := db.RunInTransaction(r.Context(), func(tx *pg.Tx) error {
		// atomically claim a ticket, guarding against overselling
		res, err := tx.Model((*models.Event)(nil)).
			Set("available_tickets = available_tickets - 1").
			Where("id = ?", eventId).
			Where("available_tickets > 0").
			Update()
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return errNoTicketsAvailable
		}

		// prevent the same user from booking the same event twice
		count, err := tx.Model(ticket).
			Where("event_id = ?", eventId).
			Where("user_id = ?", userId).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return errAlreadyBooked
		}

		ticket.EventId = eventId
		ticket.UserId = userId
		_, err = tx.Model(ticket).Insert()
		return err
	})

	if err != nil {
		switch err {
		case errNoTicketsAvailable:
			writeError(w, http.StatusBadRequest, "No available tickets for event with id "+eventId)
		case errAlreadyBooked:
			writeError(w, http.StatusConflict, "User already booked a ticket for event with id "+eventId)
		default:
			writeError(w, http.StatusInternalServerError, "Failed to book ticket for event with id "+eventId)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		log.Println("Failed to encode ticket:", err)
	}
}

func getUserBookedTicketsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userId := userIdFromContext(r)

	tickets := []models.Ticket{}
	if err := db.Model(&tickets).Where("user_id = ?", userId).Select(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get booked tickets for user with id "+userId)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		log.Println("Failed to encode tickets:", err)
	}
}
