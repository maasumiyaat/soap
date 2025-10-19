#!/usr/bin/env python3
"""
Python UDP SOAP Client Example
Demonstrates how to interact with the Go UDP SOAP server from Python
"""

import socket
import time

def send_udp_soap_request(server_host, server_port, soap_xml):
    """Send a SOAP request over UDP and return the response"""
    try:
        # Create UDP socket
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.settimeout(10)  # 10 second timeout
        
        # Send SOAP request
        sock.sendto(soap_xml.encode('utf-8'), (server_host, server_port))
        
        # Receive response
        response, addr = sock.recvfrom(4096)
        sock.close()
        
        return response.decode('utf-8')
    
    except socket.timeout:
        return "Error: Request timed out"
    except Exception as e:
        return f"Error: {str(e)}"

def main():
    server_host = "localhost"
    server_port = 8181
    
    print("=== Python UDP SOAP Client Test ===")
    
    # Test 1: GetUserByID
    print("\n1. Testing GetUserByID...")
    get_user_request = '''<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserByID xmlns="urn:user-service">
      <id>1</id>
    </GetUserByID>
  </soap:Body>
</soap:Envelope>'''
    
    response = send_udp_soap_request(server_host, server_port, get_user_request)
    print(f"Response:\n{response}\n")
    
    # Test 2: CreateUser
    print("2. Testing CreateUser...")
    create_user_request = '''<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateUser xmlns="urn:user-service">
      <name>Charlie Brown</name>
      <email>charlie@example.com</email>
    </CreateUser>
  </soap:Body>
</soap:Envelope>'''
    
    response = send_udp_soap_request(server_host, server_port, create_user_request)
    print(f"Response:\n{response}\n")
    
    # Test 3: UpdateUser (assuming the created user has ID 2)
    print("3. Testing UpdateUser...")
    update_user_request = '''<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UpdateUser xmlns="urn:user-service">
      <id>2</id>
      <name>Charlie Brown Jr.</name>
      <email>charlie.jr@example.com</email>
    </UpdateUser>
  </soap:Body>
</soap:Envelope>'''
    
    response = send_udp_soap_request(server_host, server_port, update_user_request)
    print(f"Response:\n{response}\n")
    
    # Test 4: Invalid Operation (should return fault)
    print("4. Testing Invalid Operation...")
    invalid_request = '''<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <InvalidOperation xmlns="urn:user-service">
      <param>test</param>
    </InvalidOperation>
  </soap:Body>
</soap:Envelope>'''
    
    response = send_udp_soap_request(server_host, server_port, invalid_request)
    print(f"Response:\n{response}\n")
    
    print("=== Python UDP SOAP Client Test Completed ===")

if __name__ == "__main__":
    main()
