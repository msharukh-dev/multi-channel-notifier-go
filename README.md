# Webhook Notification API

A robust, production-ready API service for sending email, SMS, and webhook notifications. Perfect for integrating notifications into your application with reliability tracking and usage analytics.

## Features

✅ **Multi-channel notifications** - Email, SMS, and webhooks
✅ **API Key authentication** - Secure client access
✅ **Usage tracking** - Daily and monthly quota management  
✅ **Status tracking** - Real-time notification delivery status
✅ **PostgreSQL backend** - Reliable data persistence
✅ **Rate limiting** - Built-in quota enforcement
✅ **REST API** - Simple, intuitive endpoints
✅ **Production-ready** - Error handling, logging, CORS support

## Quick Start

### Prerequisites
- Go 1.20+
- PostgreSQL 12+

### Installation

1. **Clone and setup**
```bash
cd webhook-api
go mod tidy
```

2. **Configure environment**
```bash
cp .env.example .env
```

Edit `.env` with your settings:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=webhook_api

MAILTRAP_API_TOKEN=your_token
MAILTRAP_FROM_EMAIL=noreply@yourdomain.com
```

3. **Start the server**
```bash
go run main.go
```

Server runs on `http://localhost:8080`

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```
GET /health
```

### 1. Register Client

Create a new client account and receive API key.

**Endpoint:** `POST /register`

**Request Body:**
```json
{
  "client_name": "My Company",
  "email": "contact@mycompany.com",
  "website": "https://mycompany.com",
  "webhook_url": "https://mycompany.com/webhooks/notifications",
  "daily_limit": 1000,
  "monthly_limit": 30000
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Client registered successfully",
  "data": {
    "client_id": 1,
    "client_name": "My Company",
    "email": "contact@mycompany.com",
    "api_key": "550e8400-e29b-41d4-a716-446655440000",
    "daily_limit": 1000,
    "monthly_limit": 30000
  }
}
```

**Save the API key** - You'll need it for all other requests.

### 2. Send Notification

Send an email, SMS, or webhook notification.

**Endpoint:** `POST /send`

**Headers:**
```
X-API-Key: your_api_key
Content-Type: application/json
```

**Request Body:**
```json
{
  "type": "email",
  "to": "user@example.com",
  "subject": "Welcome to Our Service",
  "message": "Hello! This is your notification."
}
```

**Supported Types:**
- `email` - Send email via Mailtrap
- `sms` - Send SMS via Twilio
- `webhook` - POST to webhook URL

**Response (202 Accepted):**
```json
{
  "status": "success",
  "message": "Notification queued for delivery",
  "data": {
    "notification_id": 42,
    "type": "email",
    "to": "user@example.com",
    "status": "pending",
    "created_at": "2024-01-19T10:30:45Z"
  }
}
```

### 3. Check Notification Status

Get the delivery status of a sent notification.

**Endpoint:** `GET /status/:id`

**Headers:**
```
X-API-Key: your_api_key
```

**Example Request:**
```
GET /api/v1/status/42
X-API-Key: your_api_key
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Notification status retrieved",
  "data": {
    "id": 42,
    "type": "email",
    "to": "user@example.com",
    "subject": "Welcome to Our Service",
    "status": "sent",
    "sent_at": "2024-01-19T10:30:46Z",
    "retry_count": 0,
    "created_at": "2024-01-19T10:30:45Z",
    "updated_at": "2024-01-19T10:30:46Z"
  }
}
```

**Status Values:**
- `pending` - Queued for delivery
- `sent` - Successfully delivered
- `failed` - Delivery failed

### 4. Get Usage Statistics

Check your account's current usage and remaining quota.

**Endpoint:** `GET /usage`

**Headers:**
```
X-API-Key: your_api_key
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Usage retrieved successfully",
  "data": {
    "today_usage": 234,
    "monthly_usage": 5678,
    "daily_limit": 1000,
    "monthly_limit": 30000,
    "remaining_today": 766,
    "remaining_this_month": 24322,
    "percentage_today": 23.4,
    "percentage_month": 18.9,
    "last_reset": "2024-01-19T00:00:00Z"
  }
}
```

## Error Responses

**400 Bad Request:**
```json
{
  "status": "error",
  "message": "Invalid request: email is required"
}
```

**401 Unauthorized:**
```json
{
  "status": "error",
  "message": "Invalid API key"
}
```

**429 Too Many Requests:**
```json
{
  "status": "error",
  "message": "Daily limit reached. Please try again tomorrow."
}
```

**500 Internal Server Error:**
```json
{
  "status": "error",
  "message": "Failed to save notification: database error"
}
```

## Configuration

### Environment Variables

```env
# Server
PORT=8080
GIN_MODE=debug  # Use 'release' for production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=webhook_api

# Email (Mailtrap)
MAILTRAP_API_TOKEN=your_api_token
MAILTRAP_FROM_EMAIL=noreply@domain.com

# SMS (Twilio) - Optional
TWILIO_ACCOUNT_SID=your_sid
TWILIO_AUTH_TOKEN=your_token
TWILIO_PHONE_NUMBER=+1234567890
```

## Database Schema

### Tables

**clients** - Store customer information
- id, name, email, website, webhook_url
- daily_limit, monthly_limit
- is_active, created_at, updated_at

**api_keys** - API credentials for clients
- id, key, name, client_id
- is_active, created_at, updated_at

**notifications** - Track all sent notifications
- id, client_id, type, to, subject, message
- status, error_message, sent_at, retry_count
- created_at, updated_at

**usage_logs** - Daily usage tracking
- id, client_id, notification_count, date
- created_at, updated_at

## Development

### Project Structure
```
webhook-api/
├── main.go                 # Entry point
├── go.mod                  # Dependencies
├── .env.example            # Environment template
├── config/
│   └── config.go          # Database initialization
├── models/
│   └── notification.go    # Data models
├── controllers/
│   ├── register.go        # Registration API
│   ├── send.go            # Send notification API
│   ├── status.go          # Status API
│   └── usage.go           # Usage API
├── middleware/
│   └── auth.go            # API key validation
├── routes/
│   └── routes.go          # Route definitions
└── utils/
    └── sender.go          # Email/SMS/Webhook sending
```

### Running Locally

```bash
# Install dependencies
go mod tidy

# Create database
createdb webhook_api

# Start server
go run main.go
```

### Database Migration

The application uses GORM's auto-migration. Tables are created automatically on first run.

## Production Deployment

### Environment Setup
```bash
export PORT=8080
export GIN_MODE=release
export DATABASE_URL=postgres://user:pass@prod-db:5432/webhook_api
```

### Build
```bash
go build -o webhook-api
./webhook-api
```

### Docker (Optional)
```dockerfile
FROM golang:1.20 as builder
WORKDIR /app
COPY . .
RUN go build -o webhook-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/webhook-api .
EXPOSE 8080
CMD ["./webhook-api"]
```

## Examples

### Using with cURL

**Register a client:**
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Test Corp",
    "email": "test@testcorp.com",
    "website": "https://testcorp.com",
    "daily_limit": 1000
  }'
```

**Send notification:**
```bash
curl -X POST http://localhost:8080/api/v1/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "type": "email",
    "to": "user@example.com",
    "message": "Hello!"
  }'
```

**Check usage:**
```bash
curl http://localhost:8080/api/v1/usage \
  -H "X-API-Key: your_api_key"
```

## Security Considerations

✅ **API Key validation** - All requests validated against database
✅ **Client isolation** - Clients can only access their own notifications  
✅ **Rate limiting** - Daily and monthly quota enforcement
✅ **Secure headers** - CORS and security headers configured
✅ **Error handling** - No sensitive data in error messages
✅ **SQL injection protection** - Using parameterized queries via GORM

## Support

For issues, feature requests, or questions, please contact: support@yourcompany.com

## License

MIT License - See LICENSE file for details
