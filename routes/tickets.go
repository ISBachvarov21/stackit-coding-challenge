package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/julienschmidt/httprouter"
)

func getEventBookedTicketsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bearerKey := r.Header.Get("Authorization")

	if bearerKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	key := bearerKey[len("Bearer "):]

	if key != ownerUUID {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	eventId := p.ByName("id")

	// get booked tickets for event
	tickets := []models.Ticket{}
	err := db.Model(&tickets).Where("event_id = ?", eventId).Select()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to get booked tickets for event"))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		log.Println("Failed to encode tickets:", err)
	}
}

func getAmountBookedTicketsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := r.Header.Get("Authorization")

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	userId = userId[len("Bearer "):]

	err := db.Model(&models.User{Id: userId}).WherePK().Select()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("User with id " + userId + " not found"))
		if err != nil {
			return
		}
		return
	}

	eventId := p.ByName("id")

	// get event by id
	event := &models.Event{Id: eventId}
	err = db.Model(event).WherePK().Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Event with id " + eventId + " not found"))
		if err != nil {
			return
		}
		return
	}

	var tickets []models.Ticket
	err = db.Model(&tickets).Where("event_id = ?", eventId).Select()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to get booked tickets for event"))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Amount int `json:"amount"`
	}{
		Amount: len(tickets),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to encode ticket amount:", err)
	}
}

func bookTicketHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := r.Header.Get("Authorization")

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	eventId := p.ByName("id")
	event := &models.Event{Id: eventId}
	if err := db.Model(event).WherePK().Select(); err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Event with id " + eventId + " not found"))
		if err != nil {
			return
		}
		return
	}

	if event.AvailableTickets <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("No available tickets for event with id " + eventId))
		if err != nil {
			return
		}
		return
	}

	ticket := &models.Ticket{
		EventId: eventId,
		UserId:  userId[len("Bearer "):],
	}

	if _, err := db.Model(ticket).Insert(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to book ticket for event with id " + eventId))
		if err != nil {
			return
		}
		return
	}

	event.AvailableTickets--
	if _, err := db.Model(event).WherePK().Update(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to update available tickets for event with id " + eventId))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getUserBookedTicketsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := r.Header.Get("Authorization")

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	userId = userId[len("Bearer "):]
	if err := db.Model(&models.User{Id: userId}).WherePK().Select(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("User with id " + userId + " not found"))
		if err != nil {
			return
		}
		return
	}

	tickets := []models.Ticket{}
	if err := db.Model(&tickets).Where("user_id = ?", userId).Select(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to get booked tickets for user with id " + userId))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		log.Println("Failed to encode tickets:", err)
	}
}
