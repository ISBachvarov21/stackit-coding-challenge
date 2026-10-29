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

