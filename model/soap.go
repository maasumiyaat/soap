package model

import "encoding/xml"

type SoapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SoapBody
}

type SoapBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	// Payload will hold the specific request or response struct (e.g., GetUserByIDRequest)
	Payload interface{} `xml:",innerxml"`
}

type SoapFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	Code    string   `xml:"faultcode"`
	String  string   `xml:"faultstring"`
}

func NewSoapEnvelope(payload interface{}) SoapEnvelope {
	return SoapEnvelope{
		Body: SoapBody{
			Payload: payload,
		},
	}
}

func NewSoapFault(code, message string) SoapEnvelope {
	fault := SoapFault{
		Code:   code,
		String: message,
	}
	return NewSoapEnvelope(fault)
}
