package handler

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"

	"github.com/maasumiyaat/soap/model"
	"github.com/maasumiyaat/soap/service"
)

type UDPSOAPHandler struct {
	UserService *service.UserService
	conn        *net.UDPConn
}

// NewUDPSOAPHandler creates a new UDP SOAP handler
func NewUDPSOAPHandler(userService *service.UserService) *UDPSOAPHandler {
	return &UDPSOAPHandler{
		UserService: userService,
	}
}

// StartUDPServer starts the UDP SOAP server
func (h *UDPSOAPHandler) StartUDPServer(address string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to start UDP server: %v", err)
	}

	h.conn = conn
	log.Printf("UDP SOAP Server listening on %s", address)

	// Start handling UDP requests
	go h.handleUDPRequests()

	return nil
}

// Stop closes the UDP connection
func (h *UDPSOAPHandler) Stop() {
	if h.conn != nil {
		h.conn.Close()
	}
}

// handleUDPRequests processes incoming UDP SOAP requests
func (h *UDPSOAPHandler) handleUDPRequests() {
	buffer := make([]byte, 4096) // 4KB buffer for UDP messages

	for {
		n, clientAddr, err := h.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP message: %v", err)
			continue
		}

		// Process the SOAP request in a goroutine to handle concurrent requests
		go h.processUDPSOAPRequest(buffer[:n], clientAddr)
	}
}

// processUDPSOAPRequest processes a single UDP SOAP request
func (h *UDPSOAPHandler) processUDPSOAPRequest(data []byte, clientAddr *net.UDPAddr) {
	log.Printf("Received UDP SOAP request from %s, size: %d bytes", clientAddr, len(data))

	// Parse SOAP envelope
	var initialEnv model.SoapEnvelope
	if err := xml.Unmarshal(data, &initialEnv); err != nil {
		log.Printf("Error unmarshalling UDP SOAP envelope: %v", err)
		h.sendUDPSOAPFault(clientAddr, "Client", "Invalid SOAP message")
		return
	}

	if initialEnv.Body.Payload == nil {
		h.sendUDPSOAPFault(clientAddr, "Client", "SOAP Body is empty")
		return
	}

	// Extract operation name
	xmlName := initialEnv.Body.Payload.(xml.StartElement).Name
	opName := xmlName.Local

	var responseEnv model.SoapEnvelope
	switch opName {
	case "GetUserByID":
		responseEnv = h.handleUDPGetUserByID(data)
	case "CreateUser":
		responseEnv = h.handleUDPCreateUser(data)
	case "UpdateUser":
		responseEnv = h.handleUDPUpdateUser(data)
	case "DeleteUser":
		responseEnv = h.handleUDPDeleteUser(data)
	default:
		h.sendUDPSOAPFault(clientAddr, "MustUnderstand", fmt.Sprintf("Unknown operation: %s", opName))
		return
	}

	h.sendUDPSOAPResponse(clientAddr, responseEnv)
}

// sendUDPSOAPFault sends a SOAP fault response via UDP
func (h *UDPSOAPHandler) sendUDPSOAPFault(clientAddr *net.UDPAddr, code, message string) {
	faultEnv := model.NewSoapFault(code, message)
	h.sendUDPSOAPResponse(clientAddr, faultEnv)
}

// sendUDPSOAPResponse sends a SOAP response via UDP
func (h *UDPSOAPHandler) sendUDPSOAPResponse(clientAddr *net.UDPAddr, env model.SoapEnvelope) {
	output, err := xml.MarshalIndent(env, "", "  ")
	if err != nil {
		log.Printf("Error marshalling UDP SOAP response: %v", err)
		return
	}

	// Prepare XML response with header
	response := []byte(xml.Header)
	response = append(response, output...)

	// Send response back to client
	_, err = h.conn.WriteToUDP(response, clientAddr)
	if err != nil {
		log.Printf("Error sending UDP SOAP response: %v", err)
	} else {
		log.Printf("Sent UDP SOAP response to %s, size: %d bytes", clientAddr, len(response))
	}
}

// handleUDPGetUserByID handles GetUserByID operation over UDP
func (h *UDPSOAPHandler) handleUDPGetUserByID(body []byte) model.SoapEnvelope {
	var requestEnv struct {
		XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    struct {
			XMLName xml.Name                 `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
			Request model.GetUserByIDRequest `xml:"urn:user-service GetUserByID"`
		}
	}

	if err := xml.Unmarshal(body, &requestEnv); err != nil {
		log.Printf("Error unmarshalling UDP GetUserByID request: %v", err)
		return model.NewSoapFault("Client", "Invalid GetUserByID Request Structure")
	}

	response, err := h.UserService.HandleGetUserByID(requestEnv.Body.Request)
	if err != nil {
		log.Printf("Service error for UDP GetUserByID: %v", err)
		return model.NewSoapFault("Server", err.Error())
	}

	return model.NewSoapEnvelope(response)
}

// handleUDPCreateUser handles CreateUser operation over UDP
func (h *UDPSOAPHandler) handleUDPCreateUser(body []byte) model.SoapEnvelope {
	var requestEnv struct {
		XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    struct {
			XMLName xml.Name                `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
			Request model.CreateUserRequest `xml:"urn:user-service CreateUser"`
		}
	}

	if err := xml.Unmarshal(body, &requestEnv); err != nil {
		log.Printf("Error unmarshalling UDP CreateUser request: %v", err)
		return model.NewSoapFault("Client", "Invalid CreateUser Request Structure")
	}

	response, err := h.UserService.HandleCreateUser(requestEnv.Body.Request)
	if err != nil {
		log.Printf("Service error for UDP CreateUser: %v", err)
		return model.NewSoapFault("Server", err.Error())
	}

	return model.NewSoapEnvelope(response)
}

// handleUDPUpdateUser handles UpdateUser operation over UDP
func (h *UDPSOAPHandler) handleUDPUpdateUser(body []byte) model.SoapEnvelope {
	var requestEnv struct {
		XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    struct {
			XMLName xml.Name                `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
			Request model.UpdateUserRequest `xml:"urn:user-service UpdateUser"`
		}
	}

	if err := xml.Unmarshal(body, &requestEnv); err != nil {
		log.Printf("Error unmarshalling UDP UpdateUser request: %v", err)
		return model.NewSoapFault("Client", "Invalid UpdateUser Request Structure")
	}

	response, err := h.UserService.HandleUpdateUser(requestEnv.Body.Request)
	if err != nil {
		log.Printf("Service error for UDP UpdateUser: %v", err)
		return model.NewSoapFault("Server", err.Error())
	}

	return model.NewSoapEnvelope(response)
}

// handleUDPDeleteUser handles DeleteUser operation over UDP
func (h *UDPSOAPHandler) handleUDPDeleteUser(body []byte) model.SoapEnvelope {
	var requestEnv struct {
		XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    struct {
			XMLName xml.Name                `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
			Request model.DeleteUserRequest `xml:"urn:user-service DeleteUser"`
		}
	}

	if err := xml.Unmarshal(body, &requestEnv); err != nil {
		log.Printf("Error unmarshalling UDP DeleteUser request: %v", err)
		return model.NewSoapFault("Client", "Invalid DeleteUser Request Structure")
	}

	response, err := h.UserService.HandleDeleteUser(requestEnv.Body.Request)
	if err != nil {
		log.Printf("Service error for UDP DeleteUser: %v", err)
		return model.NewSoapFault("Server", err.Error())
	}

	return model.NewSoapEnvelope(response)
}
