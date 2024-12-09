Product Management API
Overview
This is a Product Management API built with Go (Golang). The API allows you to create, retrieve, update, and delete product details. It supports asynchronous image processing, caching with Redis, and uses PostgreSQL as the database.

Features
Product Management:
Create a product with details and images.
Retrieve all products or specific products by ID.
Update and delete product details.
Asynchronous Image Processing:
(Optional) Uses RabbitMQ for message queuing.
Downloads and compresses product images in a separate microservice.
Caching:
Uses Redis to cache GET /products/:id responses to reduce database load.
Implements cache invalidation on updates and deletions.
Pagination:
Supports pagination for listing products.
Tech Stack
Backend: Go (Golang)
Database: PostgreSQL
Caching: Redis
Message Queue: RabbitMQ (Optional)
Frameworks: Gorilla Mux
Setup Instructions
Prerequisites
Go (v1.20 or higher)
PostgreSQL
Redis
RabbitMQ 

Installation
Clone the repository:

bash
git clone https://github.com/khushi2411/zocket.git
cd zocket

Install dependencies:
Copy code
go mod tidy
Set up the PostgreSQL database:

Create a database:
CREATE DATABASE products_db;
Create the products and users tables:

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_images TEXT[],
    compressed_product_images TEXT[],
    product_price NUMERIC NOT NULL
);

Set up Redis:
Start Redis:
redis-server
(Optional) Set up RabbitMQ:

Install RabbitMQ:
bash
sudo apt install rabbitmq-server
Start RabbitMQ:
bash
sudo systemctl start rabbitmq-server
Configure environment variables (create a .env file):

env
DATABASE_URL=postgres://<username>:<password>@localhost:5432/products_db?sslmode=disable
REDIS_ADDR=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

API Endpoints
Product Endpoints
Create a Product
POST /create-product
Request Body (JSON):
json
{
    "user_id": 1,
    "product_name": "Camera",
    "product_description": "High-quality camera",
    "product_images": ["image1.jpg", "image2.jpg"],
    "product_price": 499.99
}


Get All Products
GET /get-products?page=1&limit=10
Response:
json
[
    {
        "id": 1,
        "user_id": 1,
        "product_name": "Camera",
        "product_description": "High-quality camera",
        "product_images": ["image1.jpg", "image2.jpg"],
        "product_price": 499.99
    }
]
Get Product by ID
GET /products/:id


PUT /update-product/:id
Request Body (JSON):
json
{
    "user_id": 1,
    "product_name": "Updated Camera",
    "product_description": "Updated high-quality camera",
    "product_images": ["image1.jpg"],
    "product_price": 599.99
}
Response:
Product Updated successfully

Delete a Product
DELETE /delete-product/:id
Response:
Product deleted successfully

Testing
Use Postman to test the endpoints:
Import the Postman collection for faster testing.
Configure the base URL as http://localhost:8080.

Future Enhancements:
Asynchronous Image Processing:
Implement RabbitMQ for image compression and storage.
Store compressed images in S3 and update the database.
Advanced Caching:
Implement fine-grained cache invalidation strategies.
Authentication:
Add user authentication and role-based authorization.
Testing:
Add unit and integration tests for all endpoints.
