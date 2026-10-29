package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/julienschmidt/httprouter"
)

func createEventHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	event := new(models.Event)

	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if event.Name == "" {
		writeError(w, http.StatusBadRequest, "Event name is required")
		return
	} else if event.Description == "" {
		writeError(w, http.StatusBadRequest, "Event description is required")
		return
	} else if event.AvailableTickets <= 0 {
		writeError(w, http.StatusBadRequest, "Available tickets must be greater than 0")
		return
	} else if event.Date.IsZero() {
		writeError(w, http.StatusBadRequest, "Event date is required")
		return
	}

	if _, err := db.Model(event).Insert(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Println("Failed to encode event:", err)
	}
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	eventId := p.ByName("id")

	// get event by id
	event := &models.Event{Id: eventId}
	if err := db.Model(event).WherePK().Select(); err != nil {
		writeError(w, http.StatusNotFound, "Event with id "+eventId+" not found")
		return
	}

	// delete event by id, cascade delete tickets
	if _, err := db.Model(&models.Ticket{}).Where("event_id = ?", eventId).Delete(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete tickets for event")
		return
	}

	if _, err := db.Model(event).WherePK().Delete(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Event deleted successfully"}); err != nil {
		log.Println("Failed to encode response:", err)
	}
}

func getEventHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	eventId := p.ByName("id")
	event := &models.Event{Id: eventId}

	if err := db.Model(event).WherePK().Select(); err != nil {
		writeError(w, http.StatusNotFound, "Event with id "+eventId+" not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Println("Failed to encode event:", err)
	}
}

func getAllEventsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events := []models.Event{}
	if err := db.Model(&events).Select(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get events")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Println("Failed to encode events:", err)
	}
}
