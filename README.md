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
├── main.go                      # Application entry point (HTTP + UDP servers)
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums  
├── user.db                      # BoltDB database file (created at runtime)
├── database/
│   ├── bbolt.go                # Database initialization
│   └── user_repository.go      # User data access layer (CRUD operations)
├── handler/
│   ├── soap_handler.go         # HTTP SOAP request handlers
│   └── udp_soap_handler.go     # UDP SOAP request handlers
├── model/
│   ├── user.go                 # User domain models and all SOAP operations
│   └── soap.go                 # SOAP envelope structures
├── service/
│   └── user_service.go         # Business logic layer (all operations)
└── examples/
    └── udp_client.go           # UDP SOAP client example
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

4. **Servers will start on:**
   ```
   HTTP SOAP: http://localhost:8180/soap/user
   UDP SOAP:  localhost:8181
   ```

## API Usage

### HTTP SOAP Endpoint
- **URL**: `http://localhost:8180/soap/user`
- **Method**: `POST`
- **Content-Type**: `text/xml; charset=utf-8`
- **SOAPAction**: `""` (empty)

### UDP SOAP Endpoint
- **Address**: `localhost:8181`
- **Protocol**: `UDP`
- **Message Format**: `XML SOAP Envelope`
- **Max Message Size**: `4KB`

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

#### 2. CreateUser

Creates a new user with the provided name and email.

**SOAP Request:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateUser xmlns="urn:user-service">
      <name>Bob Smith</name>
      <email>bob@example.com</email>
    </CreateUser>
  </soap:Body>
</soap:Envelope>
```

**SOAP Response (Success):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateUserResponse xmlns="urn:user-service">
      <User>
        <id>2</id>
        <name>Bob Smith</name>
        <email>bob@example.com</email>
      </User>
    </CreateUserResponse>
  </soap:Body>
</soap:Envelope>
```

#### 3. UpdateUser

Updates an existing user's information.

**SOAP Request:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UpdateUser xmlns="urn:user-service">
      <id>2</id>
      <name>Bob Johnson</name>
      <email>bob.johnson@example.com</email>
    </UpdateUser>
  </soap:Body>
</soap:Envelope>
```

#### 4. DeleteUser

Deletes a user by ID.

**SOAP Request:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeleteUser xmlns="urn:user-service">
      <id>2</id>
    </DeleteUser>
  </soap:Body>
</soap:Envelope>
```

**SOAP Response (Success):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeleteUserResponse xmlns="urn:user-service">
      <success>true</success>
      <message>User with ID 2 deleted successfully</message>
    </DeleteUserResponse>
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

### Using cURL (HTTP SOAP)

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

### Using UDP SOAP (Go Client Example)

The repository includes a UDP client example in `examples/udp_client.go`. To test UDP SOAP:

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Run the UDP client:**
   ```bash
   go run examples/udp_client.go
   ```

**UDP Client Code Example:**
```go
package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    // Connect to UDP SOAP server
    conn, err := net.Dial("udp", "localhost:8181")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Send SOAP request
    soapRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>`

    _, err = conn.Write([]byte(soapRequest))
    if err != nil {
        log.Fatal(err)
    }

    // Read response
    buffer := make([]byte, 4096)
    conn.SetReadDeadline(time.Now().Add(10 * time.Second))
    n, err := conn.Read(buffer)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %s\n", buffer[:n])
}
```

### Using netcat (UDP SOAP)

You can also test UDP SOAP using netcat:

```bash
# Save SOAP request to file
cat > request.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>
EOF

# Send UDP request
nc -u localhost 8181 < request.xml
```

### UDP vs HTTP SOAP Comparison

| Feature | HTTP SOAP | UDP SOAP |
|---------|-----------|----------|
| **Transport** | TCP (reliable) | UDP (unreliable) |
| **Connection** | Connection-oriented | Connectionless |
| **Message Size** | No practical limit | Limited to 4KB |
| **Delivery** | Guaranteed delivery | Best-effort delivery |
| **Order** | Maintains order | No order guarantee |
| **Overhead** | Higher (TCP + HTTP headers) | Lower (UDP only) |
| **Performance** | Slower, more overhead | Faster, less overhead |
| **Use Cases** | Web services, enterprise | Real-time, low-latency apps |
| **Error Handling** | HTTP status codes + SOAP faults | SOAP faults only |
| **Stateful** | Can maintain sessions | Stateless |

### When to Use UDP SOAP

**Use UDP SOAP when:**
- **Low latency** is critical
- **High throughput** is required
- Message loss is acceptable
- **Simple request-response** patterns
- **Real-time** applications
- **IoT devices** with limited resources
- **Internal microservices** communication

**Avoid UDP SOAP when:**
- **Message delivery** must be guaranteed
- **Large messages** (> 4KB) are common
- **Complex transactions** requiring ACID properties
- **Security** is paramount (harder to secure UDP)
- **Order of operations** matters
- Working with **unreliable networks**

## Examples and Testing

The `examples/` directory contains various client implementations to test the UDP SOAP service:

### 1. Go UDP Client (`examples/udp_client.go`)
A comprehensive Go client that tests all SOAP operations:
```bash
# Start the server first
go run main.go

# In another terminal, run the client
go run examples/udp_client.go
```

### 2. Python UDP Client (`examples/udp_client.py`)
Cross-platform Python client demonstrating interoperability:
```bash
# Make sure Python 3 is installed
python3 examples/udp_client.py
```

### 3. Shell Script with netcat (`examples/test_udp_soap.sh`)
Simple shell script using netcat for quick testing:
```bash
# Make the script executable (already done)
chmod +x examples/test_udp_soap.sh

# Run the test script
./examples/test_udp_soap.sh
```

### UDP SOAP Implementation Features

✅ **Concurrent Request Handling**: Uses goroutines for concurrent UDP request processing  
✅ **Error Handling**: Proper SOAP fault responses for errors  
✅ **Operation Routing**: Dynamic operation detection and routing  
✅ **Cross-Platform**: Works with clients written in any language  
✅ **Logging**: Comprehensive logging for debugging  
✅ **Graceful Shutdown**: Proper resource cleanup  
✅ **Message Size Limits**: 4KB buffer for UDP messages  
✅ **Timeout Handling**: Configurable timeouts for reliability
