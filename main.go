package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/julienschmidt/httprouter"
)

var db *pg.DB
var ownerUUID string

// User represents a user in the database, only the ID is stored for this challenge.
type User struct {
	Id string `pg:"id, pk, default:gen_random_uuid()"`
}

type Event struct {
	Id               string    `pg:"id, pk, default:gen_random_uuid()"`
	Name             string    `pg:"name" json:"name"`
	Description      string    `pg:"description" json:"description"`
	AvailableTickets int       `pg:"available_tickets" json:"available_tickets"`
	Date             time.Time `pg:"date, notnull" json:"date"`
	CreatedAt        time.Time `pg:"created_at, notnull, default:now()"`
}

type Ticket struct {
	Id         string    `pg:"id, pk, default:gen_random_uuid()"`
	EventId    string    `pg:"event_id, notnull"`
	UserId     string    `pg:"user_id, notnull"`
	DateBooked time.Time `pg:"date_booked, notnull, default:now()"`
}

func initDatabase() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "stackit_coding_challenge",
	})

	err := createSchema(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Event)(nil),
		(*Ticket)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

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

	event := new(Event)

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
	} else if event.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Event description is required"))
		if err != nil {
			return
		}
	} else if event.AvailableTickets <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Available tickets must be greater than 0"))
		if err != nil {
			return
		}
	} else if event.Date.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Event date is required"))
		if err != nil {
			return
		}
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
		return
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
	event := &Event{Id: eventId}
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
	_, err = db.Model(&Ticket{}).Where("event_id = ?", eventId).Delete()
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
	tickets := []Ticket{}
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

	err := db.Model(&User{Id: userId}).WherePK().Select()
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
	event := &Event{Id: eventId}
	err = db.Model(event).WherePK().Select()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Event with id " + eventId + " not found"))
		if err != nil {
			return
		}
		return
	}

	var tickets []Ticket
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

func getEventHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	if err := db.Model(&User{Id: userId}).WherePK().Select(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("User with id " + userId + " not found"))
		if err != nil {
			return
		}
		return
	}

	eventId := p.ByName("id")
	event := &Event{Id: eventId}

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
	userId := r.Header.Get("Authorization")

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	events := []Event{}
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
	event := &Event{Id: eventId}
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

	ticket := &Ticket{
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
	if err := db.Model(&User{Id: userId}).WherePK().Select(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("User with id " + userId + " not found"))
		if err != nil {
			return
		}
		return
	}

	tickets := []Ticket{}
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

func main() {
	var err error
	ownerUUID = os.Getenv("OWNER_UUID")
	fmt.Println("OWNER_UUID:", ownerUUID)

	db, err = initDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := httprouter.New()
	router.GET("/health", healthCheckHandler)
	router.POST("/event", createEventHandler)
	router.DELETE("/event/:id", deleteEventHandler)
	router.GET("/event/:id/tickets/amount", getAmountBookedTicketsHandler)
	router.GET("/event/:id/tickets", getEventBookedTicketsHandler)
	router.GET("/event/:id", getEventHandler)
	router.GET("/events", getAllEventsHandler)
	router.POST("/event/:id/tickets", bookTicketHandler)
	router.GET("/user/tickets", getUserBookedTicketsHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
