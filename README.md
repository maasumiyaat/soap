# SOAP API with BoltDB

A Go-based SOAP web service for user management using BoltDB as the embedded database.

## Table of Contents
- [Overview](#overview)
- [REST vs SOAP: Key Differences](#rest-vs-soap-key-differences)
- [SOAP Envelope Structure](#soap-envelope-structure)
- [Project Structure](#project-structure)
- [Installation & Setup](#installation--setup)
- [API Usage](#api-usage)
- [SOAP Request/Response Examples](#soap-requestresponse-examples)
- [Architecture](#architecture)

## Overview

This project demonstrates a SOAP web service implementation in Go that provides user management functionality. It uses:
- **BoltDB (bbolt)**: Embedded key-value database for persistence
- **XML/SOAP**: Standard SOAP protocol for web service communication
- **Clean Architecture**: Layered approach with handlers, services, and repositories

## REST vs SOAP: Key Differences

| Aspect | REST | SOAP |
|--------|------|------|
| **Protocol** | Architectural style over HTTP | Protocol with strict standards |
| **Message Format** | JSON, XML, HTML, plain text | XML only |
| **Transport** | HTTP/HTTPS primarily | HTTP, SMTP, TCP, UDP |
| **Operation Style** | HTTP verbs (GET, POST, PUT, DELETE) | Function-based operations |
| **State Management** | Stateless | Can be stateful or stateless |
| **Error Handling** | HTTP status codes | SOAP fault elements |
| **Standards** | Loose guidelines | Strict WS-* standards |
| **Envelope** | No envelope structure | Mandatory SOAP envelope |
| **Discovery** | No standard (OpenAPI/Swagger) | WSDL (Web Services Description Language) |
| **Security** | HTTPS, OAuth, JWT | WS-Security, built-in security |
| **Performance** | Lightweight, faster | More overhead due to XML |
| **Caching** | HTTP caching supported | Limited caching |

### When to Use SOAP vs REST

**Use SOAP when:**
- Enterprise integration with legacy systems
- Strict security requirements
- ACID transaction support needed
- Formal contract-based communication
- Complex operations requiring stateful interactions

**Use REST when:**
- Web and mobile applications
- Microservices architecture
- Public APIs
- Performance is critical
- Simpler, more flexible communication needed

## SOAP Envelope Structure

SOAP messages are wrapped in a standardized XML envelope structure:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <!-- Optional: Authentication, routing, transaction info -->
  </soap:Header>
  <soap:Body>
    <!-- Required: The actual message payload -->
    <YourOperation xmlns="urn:your-service">
      <parameter>value</parameter>
    </YourOperation>
  </soap:Body>
</soap:Envelope>
```

### Components of SOAP Envelope

1. **Envelope**: Root element that defines the XML document as a SOAP message
2. **Header** (Optional): Contains meta-information like authentication, routing, transaction details
3. **Body** (Required): Contains the actual message payload - your operation and parameters
4. **Fault** (Optional): Used for error reporting within the Body

## Project Structure

```
soap-bbolt-api/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums  
├── user.db                 # BoltDB database file (created at runtime)
├── database/
│   ├── bbolt.go           # Database initialization
│   └── user_repository.go  # User data access layer
├── handler/
│   └── soap_handler.go     # SOAP HTTP request handlers
├── model/
│   ├── user.go            # User domain models and SOAP operations
│   └── soap.go            # SOAP envelope structures
└── service/
    └── user_service.go     # Business logic layer
```

## Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd soap-bbolt-api
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

4. **Server will start on:**
   ```
   http://localhost:8180/soap/user
   ```

## API Usage

### Endpoint
- **URL**: `http://localhost:8180/soap/user`
- **Method**: `POST`
- **Content-Type**: `text/xml; charset=utf-8`
- **SOAPAction**: `""` (empty)

### Available Operations

#### 1. GetUserByID

Retrieves a user by their unique ID.

**SOAP Request:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>
```

**SOAP Response (Success):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByIDResponse xmlns="urn:user-service">
      <User>
        <id>1</id>
        <name>Alice Johnson</name>
        <email>alice@example.com</email>
      </User>
    </GetUserByIDResponse>
  </soap:Body>
</soap:Envelope>
```

**SOAP Response (Error):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <soap:Fault>
      <faultcode>Server</faultcode>
      <faultstring>User with ID 999 not found</faultstring>
    </soap:Fault>
  </soap:Body>
</soap:Envelope>
```

## SOAP Request/Response Examples

### Using cURL

```bash
curl -X POST http://localhost:8180/soap/user \
  -H "Content-Type: text/xml; charset=utf-8" \
  -H "SOAPAction: " \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>'
```

### Using Postman

1. **Method**: POST
2. **URL**: `http://localhost:8180/soap/user`
3. **Headers**:
   - `Content-Type: text/xml; charset=utf-8`
   - `SOAPAction: `
4. **Body** (raw XML):
   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
   <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
     <soap:Body>
       <GetUserByID xmlns="urn:user-service">
         <id>1</id>
       </GetUserByID>
     </soap:Body>
   </soap:Envelope>
   ```

### Using SoapUI

1. Create new SOAP project
2. Set endpoint: `http://localhost:8180/soap/user`
3. Import the request structure or manually create the XML request
4. Execute the request

## Architecture

The application follows a clean architecture pattern:

### Layers

1. **Handler Layer** (`handler/soap_handler.go`):
   - Handles HTTP requests
   - Parses SOAP envelopes
   - Routes operations to appropriate services
   - Formats SOAP responses

2. **Service Layer** (`service/user_service.go`):
   - Contains business logic
   - Validates input
   - Coordinates with repository layer

3. **Repository Layer** (`database/user_repository.go`):
   - Data access abstraction
   - BoltDB operations
   - Data persistence

4. **Model Layer** (`model/`):
   - Domain entities
   - SOAP operation structures
   - Request/Response models

### Data Flow

```
HTTP Request → SOAP Handler → User Service → User Repository → BoltDB
                     ↓
HTTP Response ← SOAP Handler ← User Service ← User Repository ← BoltDB
```

### Error Handling

The application implements SOAP fault handling:
- **Client faults**: Invalid requests, malformed XML
- **Server faults**: Database errors, service failures
- **MustUnderstand faults**: Unknown operations

### Database

- **BoltDB**: Embedded key-value store
- **Bucket**: "Users" bucket for user storage
- **Key**: User ID (string)
- **Value**: JSON-serialized user object

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

[Add your license information here]
