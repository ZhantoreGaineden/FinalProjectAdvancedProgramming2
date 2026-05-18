# Pet Store Microservices Final Project

**Student:** Zhantore Gaineden  
**Project:** Final Project Advanced Programming 2  
**Topic:** 19. Pet Store  

## Project Description

Pet Store is a microservice-based application for managing pets, users, orders, and notifications. The system is implemented in Go using Clean Architecture, gRPC communication, REST API Gateways, PostgreSQL, Redis cache, NATS message broker, Docker Compose, migrations, transactions, and email notifications.

The project is implemented individually, but structured as a three-member microservice project. Each simulated participant owns one domain, one microservice, and one API Gateway.

## Team Structure

### Participant 1 — Pet Domain
Responsibilities:
- Pet Service
- Pet Gateway
- Pet PostgreSQL database
- Redis cache for pet data

### Participant 2 — User Domain
Responsibilities:
- User Service
- User Gateway
- User PostgreSQL database
- User registration and login
- `user.registered` NATS event

### Participant 3 — Order Domain
Responsibilities:
- Order Service
- Order Gateway
- Order PostgreSQL database
- Transactional order creation
- `order.created` and `order.status_updated` NATS events

### Shared Service
- Notification Service
- NATS message broker
- SMTP/email log sender
- Docker Compose
- Migrations

## Architecture

```txt
Frontend / Postman / curl
        |
        | REST API
        v
+----------------+   +----------------+   +----------------+
|  Pet Gateway   |   |  User Gateway  |   | Order Gateway  |
|    :8081       |   |     :8082      |   |     :8083      |
+----------------+   +----------------+   +----------------+
        |                    |                    |
        | gRPC               | gRPC               | gRPC
        v                    v                    v
+----------------+   +----------------+   +----------------+
|  Pet Service   |   |  User Service  |   | Order Service  |
|    :50051      |   |     :50052     |   |     :50053     |
+----------------+   +----------------+   +----------------+
        |                    |                    |
        v                    v                    v
   PostgreSQL           PostgreSQL           PostgreSQL
     pet_db              user_db              order_db
        |
        v
      Redis

User Service and Order Service publish events to NATS.
Notification Service subscribes to NATS and sends email notifications.
Microservices
1. Pet Service

gRPC endpoints:

CreatePet
GetPet
ListPets
UpdatePet
DeletePet

REST endpoints through Pet Gateway:

POST /api/pets
GET /api/pets
GET /api/pets/:id
PUT /api/pets/:id
DELETE /api/pets/:id

Features:

CRUD for pets
PostgreSQL storage
Redis cache for pet list and pet by ID
2. User Service

gRPC endpoints:

RegisterUser
LoginUser
GetUser
UpdateUser
DeleteUser

REST endpoints through User Gateway:

POST /api/users/register
POST /api/users/login
GET /api/users/:id
PUT /api/users/:id
DELETE /api/users/:id

Features:

User registration
User login
Password hashing with bcrypt
PostgreSQL storage
Publishes user.registered event to NATS
3. Order Service

gRPC endpoints:

CreateOrder
GetOrder
ListUserOrders
UpdateOrderStatus
CancelOrder

REST endpoints through Order Gateway:

POST /api/orders
GET /api/orders/:id
GET /api/users/:id/orders
PATCH /api/orders/:id/status
POST /api/orders/:id/cancel

Features:

Order creation
Order status update
Order cancellation
Transactional creation of order and order items
PostgreSQL storage
Publishes order.created and order.status_updated events to NATS
4. Notification Service

Features:

Subscribes to NATS events:
user.registered
order.created
order.status_updated
Sends email notifications
Uses email log mode by default
Supports SMTP configuration through environment variables
Technologies Used
Go
gRPC
Protocol Buffers
Gin
PostgreSQL
Redis
NATS
Docker
Docker Compose
golang-migrate
bcrypt
Requirements Coverage
Requirement	Implementation
Clean Architecture	Each service has entity, repository, usecase, delivery layers
3+ Microservices	Pet, User, Order, Notification
API Gateway by each member	Pet Gateway, User Gateway, Order Gateway
12+ gRPC endpoints	15 gRPC endpoints total
Message Queue	NATS
Database	PostgreSQL
Cache	Redis in Pet Service
Migrations	SQL migrations with migrate runner
Transactions	Order creation uses database transaction
Emails	Notification Service sends email logs / SMTP-supported emails
Docker	Dockerfiles and Docker Compose
Testing	go test ./...
Project Structure
.
├── proto
│   ├── pet.proto
│   ├── user.proto
│   ├── order.proto
│   └── gen
├── pet-service
├── user-service
├── order-service
├── notification-service
├── pet-gateway
├── user-gateway
├── order-gateway
├── migrations-runner
├── docker
│   └── init-db.sql
├── docker-compose.yml
├── go.mod
└── README.md
How to Run
1. Clone repository
git clone https://github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2.git
cd FinalProjectAdvancedProgramming2
2. Start all services
docker compose up --build -d
3. Check containers
docker compose ps
4. Check health endpoints
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

Expected result:

{"status":"pet-gateway is running"}
{"status":"user-gateway is running"}
{"status":"order-gateway is running"}
Demo Commands
Register user
curl -X POST http://localhost:8082/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Zhantore Gaineden","email":"zhantore@example.com","password":"123456"}'
Login user
curl -X POST http://localhost:8082/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"zhantore@example.com","password":"123456"}'
Create pet
curl -X POST http://localhost:8081/api/pets \
  -H "Content-Type: application/json" \
  -d '{"name":"Buddy","category":"dog","breed":"Golden Retriever","age":2,"price":500,"status":"available"}'
List pets
curl http://localhost:8081/api/pets
Create order

Replace USER_ID_HERE and PET_ID_HERE with real IDs from previous responses.

curl -X POST http://localhost:8083/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id":"USER_ID_HERE",
    "user_email":"zhantore@example.com",
    "items":[
      {
        "pet_id":"PET_ID_HERE",
        "price":500
      }
    ]
  }'
Update order status

Replace ORDER_ID_HERE with a real order ID.

curl -X PATCH http://localhost:8083/api/orders/ORDER_ID_HERE/status \
  -H "Content-Type: application/json" \
  -d '{"status":"paid","user_email":"zhantore@example.com"}'
View notification logs
docker compose logs notification-service

Expected logs:

Welcome email after user registration
Order created email after order creation
Order status updated email after status update
Database Tables

Pet database:

docker exec -it petstore-postgres psql -U postgres -d pet_db -c "\dt"

User database:

docker exec -it petstore-postgres psql -U postgres -d user_db -c "\dt"

Order database:

docker exec -it petstore-postgres psql -U postgres -d order_db -c "\dt"
Clean Start

To remove all containers and database volume:

docker compose down -v
docker compose up --build -d

The migrations-runner container automatically applies all migrations after PostgreSQL becomes healthy.

Ports
Service	Port
Pet Gateway	8081
User Gateway	8082
Order Gateway	8083
Pet Service gRPC	50051
User Service gRPC	50052
Order Service gRPC	50053
PostgreSQL	5432
Redis	6379
NATS	4222
Notes

This project was implemented individually, but the architecture is organized as a three-participant project. Each participant domain has its own microservice and API Gateway.
