# Advertisement Management API

This API allows users to manage advertisements by retrieving existing ads based on various conditions and adding new ads to the system.

## Table of Contents

- [Endpoints](#endpoints)
- [Query Parameters](#query-parameters)
- [Database Schema](#database-schema)
- [Usage](#usage)
- [Testing](#testing)


## Endpoints

### GET /api/v1/ad

Retrieves advertisements based on query parameters such as age, gender, country, and platform.

### POST /api/v1/ad

Adds a new advertisement to the system.

## Query Parameters

- `age`: Filter ads based on age range.
- `gender`: Filter ads based on gender.
- `country`: Filter ads based on country.
- `platform`: Filter ads based on platform.

## Database Schema

The MongoDB database contains a collection named `ads` to store advertisement documents. Each advertisement document includes fields such as title, startAt, endAt, and conditions.

## Design Choices

This API was built with scalability and performance in mind. Here are some key design choices:

### MongoDB

MongoDB was chosen as the database solution due to its flexibility, scalability, and ease of use. As the number of requests increases, MongoDB's horizontal scaling capabilities allow the system to easily scale out by adding more nodes to the cluster. This ensures that the API can handle a large number of requests without sacrificing performance.

### Goroutine

Goroutines were utilized to handle concurrent requests efficiently. Each incoming HTTP request can be processed concurrently in its own goroutine, allowing the API to handle a large number of requests simultaneously. This concurrency model helps improve the responsiveness and throughput of the API, ensuring optimal performance even under high load conditions.

## Usage

1. Clone the repository:
```plaintext
git clone https://github.com/meganchlin/advertisement_management_api.git
```

2. Navigate to the project directory:
```plaintext
cd advertisement_management_api
```

3. Install dependencies:
```plaintext
go mod download
```

4. Run MongoDB via Docker:
```plaintext
docker run -p 127.0.0.1:27017:27017 -d --rm --name mongo mongo:7.0.5
```

5. Run the application:
```plaintext
go run main.go
```

## Testing
Run Tests: 
```plaintext
go test
```