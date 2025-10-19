package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// Simple UDP SOAP client for testing
type UDPSOAPClient struct {
	serverAddr string
	conn       *net.UDPConn
}

func NewUDPSOAPClient(serverAddr string) *UDPSOAPClient {
	return &UDPSOAPClient{
		serverAddr: serverAddr,
	}
}

func (c *UDPSOAPClient) Connect() error {
	udpAddr, err := net.ResolveUDPAddr("udp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to UDP server: %v", err)
	}

	c.conn = conn
	return nil
}

func (c *UDPSOAPClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *UDPSOAPClient) SendSOAPRequest(soapXML string) (string, error) {
	// Send SOAP request
	_, err := c.conn.Write([]byte(soapXML))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	// Set read timeout
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// Read response
	buffer := make([]byte, 4096)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(buffer[:n]), nil
}

func main() {
	client := NewUDPSOAPClient("localhost:8181")

	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	fmt.Println("=== UDP SOAP Client Test ===")

	// Test 1: GetUserByID
	fmt.Println("\n1. Testing GetUserByID...")
	getUserRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>`

	response, err := client.SendSOAPRequest(getUserRequest)
	if err != nil {
		log.Printf("GetUserByID failed: %v", err)
	} else {
		fmt.Printf("Response:\n%s\n", response)
	}

	// Test 2: CreateUser
	fmt.Println("\n2. Testing CreateUser...")
	createUserRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateUser xmlns="urn:user-service">
      <name>Bob Smith</name>
      <email>bob@example.com</email>
    </CreateUser>
  </soap:Body>
</soap:Envelope>`

	response, err = client.SendSOAPRequest(createUserRequest)
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
	} else {
		fmt.Printf("Response:\n%s\n", response)
	}

	// Test 3: UpdateUser
	fmt.Println("\n3. Testing UpdateUser...")
	updateUserRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UpdateUser xmlns="urn:user-service">
      <id>2</id>
      <name>Bob Johnson</name>
      <email>bob.johnson@example.com</email>
    </UpdateUser>
  </soap:Body>
</soap:Envelope>`

	response, err = client.SendSOAPRequest(updateUserRequest)
	if err != nil {
		log.Printf("UpdateUser failed: %v", err)
	} else {
		fmt.Printf("Response:\n%s\n", response)
	}

	// Test 4: DeleteUser
	fmt.Println("\n4. Testing DeleteUser...")
	deleteUserRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeleteUser xmlns="urn:user-service">
      <id>2</id>
    </DeleteUser>
  </soap:Body>
</soap:Envelope>`

	response, err = client.SendSOAPRequest(deleteUserRequest)
	if err != nil {
		log.Printf("DeleteUser failed: %v", err)
	} else {
		fmt.Printf("Response:\n%s\n", response)
	}

	// Test 5: Invalid Operation
	fmt.Println("\n5. Testing Invalid Operation (should return fault)...")
	invalidRequest := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <InvalidOperation xmlns="urn:user-service">
      <param>test</param>
    </InvalidOperation>
  </soap:Body>
</soap:Envelope>`

	response, err = client.SendSOAPRequest(invalidRequest)
	if err != nil {
		log.Printf("Invalid operation test failed: %v", err)
	} else {
		fmt.Printf("Response:\n%s\n", response)
	}

	fmt.Println("\n=== UDP SOAP Client Test Completed ===")
}
