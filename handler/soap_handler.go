package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/maasumiyaat/soap/model"
	"github.com/maasumiyaat/soap/service"
)

type UserSOAPHandler struct {
	UserService *service.UserService
}

func (h *UserSOAPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeSOAPFault(w, "Client", "Failed to read request body")
		return
	}

	var initialEnv model.SoapEnvelope
	if err := xml.Unmarshal(body, &initialEnv); err != nil {
		log.Printf("Error un-marshalling SOAP envelope: %v", err)
		h.writeSOAPFault(w, "Client", "Invalid SOAP message")
		return
	}

	if initialEnv.Body.Payload == nil {
		h.writeSOAPFault(w, "Client", "SOAP Body is empty")
		return
	}

	xmlName := initialEnv.Body.Payload.(xml.StartElement).Name
	opName := xmlName.Local

	var responseEnv model.SoapEnvelope
	switch opName {
	case "GetUserByID":
		responseEnv = h.handleGetUserByID(body)
	default:
		h.writeSOAPFault(w, "MustUnderstand", fmt.Sprintf("Unknown operation: %s", opName))
		return
	}

	h.writeSOAPResponse(w, responseEnv)
}

func (h *UserSOAPHandler) writeSOAPFault(w http.ResponseWriter, code, message string) {
	faultEnv := model.NewSoapFault(code, message)
	h.writeSOAPResponse(w, faultEnv)
}

func (h *UserSOAPHandler) writeSOAPResponse(w http.ResponseWriter, env model.SoapEnvelope) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("SOAPAction", "")

	output, err := xml.MarshalIndent(env, "", "  ")
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(output)
}

func (h *UserSOAPHandler) handleGetUserByID(body []byte) model.SoapEnvelope {
	var requestEnv struct {
		XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    struct {
			XMLName xml.Name                 `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
			Request model.GetUserByIDRequest `xml:"urn:user-service GetUserByID"`
		}
	}

	if err := xml.Unmarshal(body, &requestEnv); err != nil {
		log.Printf("Error unmarshalling GetUserByID request: %v", err)
		return model.NewSoapFault("Client", "Invalid GetUserByID Request Structure")
	}

	response, err := h.UserService.HandleGetUserByID(requestEnv.Body.Request)
	if err != nil {
		log.Printf("Service error for GetUserByID: %v", err)
		return model.NewSoapFault("Server", err.Error())
	}

	return model.NewSoapEnvelope(response)
}
