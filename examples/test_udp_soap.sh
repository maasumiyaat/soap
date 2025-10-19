#!/bin/bash

# UDP SOAP Test Script using netcat
# Make sure the Go server is running before executing this script

echo "=== UDP SOAP Test Script ==="
echo "Make sure the server is running: go run main.go"
echo ""

# Function to send UDP SOAP request and display response
send_udp_soap() {
    local operation="$1"
    local xml_request="$2"
    
    echo "Testing $operation..."
    echo "Request:"
    echo "$xml_request"
    echo ""
    echo "Response:"
    
    # Use netcat to send UDP request
    echo "$xml_request" | nc -u -w 5 localhost 8181
    echo ""
    echo "---"
    echo ""
}

# Test 1: GetUserByID
GET_USER_REQUEST='<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>'

send_udp_soap "GetUserByID" "$GET_USER_REQUEST"

# Test 2: CreateUser
CREATE_USER_REQUEST='<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateUser xmlns="urn:user-service">
      <name>Diana Prince</name>
      <email>diana@example.com</email>
    </CreateUser>
  </soap:Body>
</soap:Envelope>'

send_udp_soap "CreateUser" "$CREATE_USER_REQUEST"

# Test 3: UpdateUser
UPDATE_USER_REQUEST='<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UpdateUser xmlns="urn:user-service">
      <id>2</id>
      <name>Diana Wonder Woman</name>
      <email>wonder.woman@example.com</email>
    </UpdateUser>
  </soap:Body>
</soap:Envelope>'

send_udp_soap "UpdateUser" "$UPDATE_USER_REQUEST"

# Test 4: DeleteUser
DELETE_USER_REQUEST='<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeleteUser xmlns="urn:user-service">
      <id>2</id>
    </DeleteUser>
  </soap:Body>
</soap:Envelope>'

send_udp_soap "DeleteUser" "$DELETE_USER_REQUEST"

# Test 5: Invalid Operation
INVALID_REQUEST='<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <InvalidOperation xmlns="urn:user-service">
      <param>test</param>
    </InvalidOperation>
  </soap:Body>
</soap:Envelope>'

send_udp_soap "Invalid Operation" "$INVALID_REQUEST"

echo "=== UDP SOAP Test Completed ==="
