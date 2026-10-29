## API Routes

All routes (except `/health`) require an `Authorization: Bearer <uuid>` header.
- **Owner** routes expect the token to match the `OWNER_UUID` environment variable.
- **User** routes expect the token to be a valid user ID.

| Method | Path | Access | Description |
|--------|------|--------|-------------|
| GET | `/health` | Public | Health check, returns `200 OK` |
| POST | `/event` | Owner | Create a new event (JSON body: name, description, available tickets, date) |
| DELETE | `/event/:id` | Owner | Delete an event and all of its booked tickets |
| GET | `/event/:id/tickets/amount` | Owner | Get the number of booked tickets for an event (`{"amount": n}`) |
| GET | `/event/:id/tickets` | Owner | List all tickets booked for an event |
| GET | `/event/:id` | User | Get details of a single event |
| GET | `/events` | User | List all events |
| POST | `/event/:id/tickets` | User | Book a ticket for an event (fails if sold out or already booked by the user) |
| GET | `/user/tickets` | User | List all tickets booked by the authenticated user |

## Example requests

Assumes the service is running locally on port 8080. Replace `$OWNER_UUID` with your owner UUID, `$USER_UUID` with an existing user's ID, and `$EVENT_ID` with an event's ID.

```bash
# Health check
curl http://localhost:8080/health

# Create an event (owner only)
curl -X POST http://localhost:8080/event \
  -H "Authorization: Bearer $OWNER_UUID" \
  -H "Content-Type: application/json" \
  -d '{"name":"Concert","description":"Live music","available_tickets":100,"date":"2026-12-01T20:00:00Z"}'

# List all events
curl http://localhost:8080/events \
  -H "Authorization: Bearer $USER_UUID"

# Get a single event
curl http://localhost:8080/event/$EVENT_ID \
  -H "Authorization: Bearer $USER_UUID"

# Book a ticket for an event
curl -X POST http://localhost:8080/event/$EVENT_ID/tickets \
  -H "Authorization: Bearer $USER_UUID"

# List your booked tickets
curl http://localhost:8080/user/tickets \
  -H "Authorization: Bearer $USER_UUID"

# Get the amount of booked tickets for an event (owner only)
curl http://localhost:8080/event/$EVENT_ID/tickets/amount \
  -H "Authorization: Bearer $OWNER_UUID"

# List all booked tickets for an event (owner only)
curl http://localhost:8080/event/$EVENT_ID/tickets \
  -H "Authorization: Bearer $OWNER_UUID"

# Delete an event (owner only)
curl -X DELETE http://localhost:8080/event/$EVENT_ID \
  -H "Authorization: Bearer $OWNER_UUID"
```

## How to run the service
```bash
export OWNER_UUID=your_owner_uuid; go run main.go
```
This command will start the service on port 8080. It will also set the `OWNER_UUID` environment variable, which is required for the service to function properly. Make sure to replace `your_owner_uuid` with a valid UUID.

## Assumptions made for this challenge
- The service is stateless and does not persist data between restarts.
- The `OWNER_UUID` is a single value, stored in an environment variable that represents the owner of the service.
- The service does not implement user registration or authentication. It does check each passed UUID against the users table in the database to see if the user exists.
- The service does not implement any rate limiting or throttling.
- The service assumes what the models for events and tickets look like, as they were not provided in the challenge description. The models are defined in the `models` package.
- The service assumes concurrent access to the same event is possible, and handles it by using database transactions to ensure that ticket booking is atomic and consistent.

## Future improvements
- Implement user authentication.
- Implement rate limiting and throttling.
- Implement pagination for listing events and tickets.
- Implement more detailed error handling and logging.
- Use a more robust HTTP server framework that supports middleware and better routing, f.e. Gin or Echo.