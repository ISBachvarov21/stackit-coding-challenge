package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/julienschmidt/httprouter"
)

func createEventHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	event := new(models.Event)

	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Invalid request body"))
		if err != nil {
			return
		}
		return
	}

	if event.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Event name is required"))
		if err != nil {
			return
		}
		return
	} else if event.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Event description is required"))
		if err != nil {
			return
		}
		return
	} else if event.AvailableTickets <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Available tickets must be greater than 0"))
		if err != nil {
			return
		}
		return
	} else if event.Date.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Event date is required"))
		if err != nil {
			return
		}
		return
	}

	_, err := db.Model(event).Insert()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to create event"))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Println("Failed to encode event:", err)
	}
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	// get event by id
	event := &models.Event{Id: eventId}
	err := db.Model(event).WherePK().Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Event not found"))
		if err != nil {
			return
		}
		return
	}

	// delete event by id, cascade delete tickets
	_, err = db.Model(&models.Ticket{}).Where("event_id = ?", eventId).Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to delete tickets for event"))
		if err != nil {
			return
		}
		return
	}

	_, err = db.Model(event).WherePK().Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to delete event"))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Event deleted successfully"))
	if err != nil {
		return
	}
}

func getEventHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Println("Failed to encode event:", err)
	}
}

func getAllEventsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	events := []models.Event{}
	if err := db.Model(&events).Select(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to get events"))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Println("Failed to encode events:", err)
	}
}
