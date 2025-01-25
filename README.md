# Wolt Backend Engineering Task

This project is a backend service for calculating the delivery order price for Wolt. It includes logic to compute the delivery fee, small order surcharge, total price, and delivery distance using data retrieved from a mock Home Assignment API.

## Features

- **GET /api/v1/delivery-order-price**: Computes the delivery price breakdown, including:
  - Total price
  - Small order surcharge
  - Delivery fee
  - Delivery distance

## Technologies Used

- **Programming Language**: Go (Golang)
- **Router**: Chi
- **HTTP Client**: Built-in `net/http`
- **Dependency Injection**: Manual via constructors
- **Testing**: Go testing framework

## Directory Structure

```
backend-task-wolt-go/
├── cmd/                     # Entry point for the application
│   └── server/              # Main server setup
│       └── main.go          # Application entry point
├── configs/                 # Configuration files
│   └── config.yaml          # Application configuration file
├── internal/                # Core application logic
│   ├── api/                 # HTTP handler logic
│   │   └── handler.go       # Request handling and response generation
│   ├── client/              # External API client
│   │   └── client.go        # HTTP client to interact with external APIs
│   ├── models/              # Data models for static and dynamic API responses
│   │   └── models.go        # Definitions for API data models
│   ├── service/             # Business logic
│   │   ├── calculator.go    # Logic for fee, surcharge, and distance calculations
│   │   └── venue_service.go # Venue-related service logic
│   ├── utils/               # Utility functions
│   │   ├── calculator.go    # Helper functions for calculations
│   │   └── load_config.go   # Configuration loading logic
├── go.mod                   # Module dependencies
├── go.sum                   # Module checksum
└── README.md                # Documentation

```

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo/backend-task-wolt-go.git
   cd backend-task-wolt-go
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

4. **Access the service**:
   The service will be available at `http://localhost:8000`.

## API Usage

### Endpoint

**GET /api/v1/delivery-order-price**

### Query Parameters

| Parameter    | Type    | Description                          | Example                        |
|--------------|---------|--------------------------------------|--------------------------------|
| `venue_slug` | string  | Unique identifier for the venue      | `home-assignment-venue-helsinki` |
| `cart_value` | integer | Total value of items in the cart     | `1000`                         |
| `user_lat`   | float   | Latitude of the user's location      | `60.17094`                     |
| `user_lon`   | float   | Longitude of the user's location     | `24.93087`                     |

### Example Request

```bash
curl "http://localhost:8000/api/v1/delivery-order-price?venue_slug=home-assignment-venue-helsinki&cart_value=1000&user_lat=60.17094&user_lon=24.93087"
```

### Example Response

```json
{
  "total_price": 1190,
  "small_order_surcharge": 0,
  "cart_value": 1000,
  "delivery": {
    "fee": 190,
    "distance": 177
  }
}
```

## Development

### Adding a New Feature
1. Add business logic to the relevant service in `internal/service`.
2. Update the handler in `internal/api`.
3. Write unit tests for new logic in `test/`.

### Testing
Run unit and integration tests:
```bash
go test ./...
```

## Future Improvements

- Add caching for frequently accessed venue data.
- Implement rate limiting to prevent abuse.
- Add more robust error handling and logging.
- Use environment variables for configuration.
