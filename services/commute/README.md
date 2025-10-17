# Duration Service

A microservice that calculates travel duration between two points using Google Maps API.

## Features

- Calculate travel duration and distance between origin and destination
- Support for multiple travel modes (driving, walking, bicycling, transit)
- Traffic-aware calculations for driving mode
- Transit routing preferences
- RESTful API with JSON responses
- Health check endpoint
- Graceful shutdown
- Docker support

## API Endpoints

### POST /duration

Calculate travel duration between two points.

**Request Body:**
```json
{
  "origin": "New York, NY",
  "destination": "Los Angeles, CA",
  "mode": "driving",
  "departure_time": "2024-01-15T09:00:00Z",
  "traffic_model": "best_guess"
}
```

**Response:**
```json
{
  "origin": "New York, NY",
  "destination": "Los Angeles, CA",
  "mode": "driving",
  "duration_seconds": 14400,
  "distance_meters": 3944000,
  "status": "OK"
}
```

### GET /health

Health check endpoint.

**Response:**
```
OK
```

## Configuration

Set the following environment variables:

- `ENVIRONMENT`: Environment type (development/production)
- `API_HOST`: Host to bind the server to (default: "")
- `API_PORT`: Port to bind the server to (default: 8080)
- `GOOGLE_MAPS_API_KEY`: Your Google Maps API key
- `LOG_LEVEL`: Log level (default: info)

## Travel Modes

- `driving`: Car travel (default)
- `walking`: Walking
- `bicycling`: Bicycle travel
- `transit`: Public transportation

## Traffic Models (for driving mode)

- `best_guess`: Default traffic model
- `pessimistic`: Assumes heavy traffic
- `optimistic`: Assumes light traffic

## Transit Routing Preferences

- `less_walking`: Prefer routes with less walking
- `fewer_transfers`: Prefer routes with fewer transfers

## Running the Service

### Local Development

1. Copy `.env.example` to `.env` and configure your settings
2. Install dependencies: `go mod download`
3. Run the service: `go run main.go`

### Docker

1. Build the image: `docker build -t duration-service .`
2. Run the container: `docker run -p 8080:8080 --env-file .env duration-service`

## Example Usage

```bash
# Calculate driving duration
curl -X POST http://localhost:8080/duration \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "New York, NY",
    "destination": "Boston, MA",
    "mode": "driving"
  }'

# Health check
curl http://localhost:8080/health
```

## Error Handling

The service returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid request parameters
- `500 Internal Server Error`: Service errors
- `503 Service Unavailable`: Google Maps API errors

## Dependencies

- Google Maps Distance Matrix API
- Go 1.24+
- Docker (optional)
