#  Hotel Booking & Payment Service

## Technologies Used
- **Go** – Backend programming language  
- **PostgreSQL** – Database  
- **GORM** – ORM for Go  
- **CronJob** – Automatic checkout processing  
- **Testify** – Unit testing framework  
- **Xendit API** – 3rd party payment processing  

## Features
- User registration and login  
- Room management (CRUD)  
- Booking/reservation management  
- Payment creation and simulation (Xendit)  
- Automatic checkout with CronJob  

## Limitations
- Webhook processing must be called manually because the project is running locally and cannot receive callbacks from Xendit.

# Setup and Run Instructions

## Prerequisites
- Docker and Docker Compose installed
- PostgreSQL (or your preferred database)

## Database Setup

1. Create your database
2. Update the `.env` file with your database credentials:
   ```
   DB_HOST=your_host
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   ```

## Running the Services

### Using Docker Compose

Start all services with:
```bash
docker-compose up
```

### Service Ports

The following services will be available on these ports:

- **Auth Service**: `http://localhost:8081`
- **Reservation Service**: `http://localhost:8082`
- **Payment Service**: `http://localhost:8083`

## Main Flow

### Payment Flow (`/payment`)

1. Payment Service receives payment creation request
2. Payment Service contacts Xendit API
3. Xendit returns payment URL
4. Payment Service responds with payment details including the Xendit payment URL

### Webhook Flow (`/webhook`)

1. Xendit calls webhook endpoint
2. Payment Service processes webhook
3. Payment Service calls Reservation Service to validate booking
4. Reservation Service checks:
   - Date validation
   - Room availability
5. Updates statuses:
   - Room status → `unavailable`
   - Booking status → `paid`
6. Reservation Service notifies Payment Service
7. Payment Service responds with webhook success

### Check-in Flow (`/checkin`)

Validates before check-in:
- User authorization (can only check-in own booking)
- Booking not already checked in
- Booking not already completed
- Payment completed (booking status must not be `pending`)
- Check-in date must be today or after
- Room status must be `available`

On successful check-in:
- Room status → `unavailable`
- Booking status → `checked_in`

### Check-out Flow (`/checkout`)

Validates before check-out:
- User authorization (can only check-out own booking)
- Booking must be in `checked_in` status
- Booking not already completed

On successful check-out:
- Room status → `available`
- Booking status → `completed`

---

## API Documentation

Access the Postman collection for API testing:

🔗 [Postman Collection](https://fabuzard-business-8028509.postman.co/workspace/fahreza-alghifary's-Workspace~32dfcd0e-697d-46f0-828a-b3ab4024ffd6/collection/49104373-a5e80818-3677-4c61-a911-9a4861023202?action=share&creator=49104373&active-environment=49104373-d83b1a04-e5b9-4a10-a9b3-7a37119c4816)

## Quick Start

1. Clone the repository
2. Set up your database
3. Configure your `.env` file with database credentials
4. Run `docker-compose up`
5. Access services on their respective ports (8081, 8082, 8083)
6. Use the Postman collection to test the APIs
