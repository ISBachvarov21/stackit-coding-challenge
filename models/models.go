package models

import "time"

// User represents a user in the database, only the ID is stored for this challenge.
type User struct {
	Id string `pg:"id, pk, default:gen_random_uuid()" json:"id"`
}

type Event struct {
	Id               string    `pg:"id, pk, default:gen_random_uuid()" json:"id"`
	Name             string    `pg:"name" json:"name"`
	Description      string    `pg:"description" json:"description"`
	AvailableTickets int       `pg:"available_tickets" json:"available_tickets"`
	Date             time.Time `pg:"date, notnull" json:"date"`
	CreatedAt        time.Time `pg:"created_at, notnull, default:now()" json:"created_at"`
}

type Ticket struct {
	Id         string    `pg:"id, pk, default:gen_random_uuid()" json:"id"`
	EventId    string    `pg:"event_id, notnull" json:"event_id"`
	UserId     string    `pg:"user_id, notnull" json:"user_id"`
	DateBooked time.Time `pg:"date_booked, notnull, default:now()" json:"date_booked"`
}
