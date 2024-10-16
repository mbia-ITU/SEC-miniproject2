package main

import (
	types "SEC-miniproject2/src/Hospital/Types"
	utilities "SEC-miniproject2/src/Hospital/Utilities"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var client *http.Client
var patients []int
var regpatients int
var patientstotal int
var port int
var data int
var recshares int

func patientHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(port, ": Hospital received /patients")
	if r.Method == "POST" {
		log.Println(port, ": Hospital received POST /patients")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, ": Error reading request body:", err)
			w.WriteHeader(400)
			return
		}
		receivedPort := &types.Patient{}
		err = json.Unmarshal(body, receivedPort)
		if err != nil {
			log.Fatal(port, ": Error unmarshalling port:", err)
			w.WriteHeader(400)
			return
		}
		log.Println(port, ": Registered new patient at port", receivedPort.Port)
		patients = append(patients, receivedPort.Port)
		regpatients++
		if regpatients == patientstotal {
			sendport()
		}
		w.WriteHeader(200)
	}
}

func sendport() {
	log.Println(port, ": Sending ports to patients")
	for i, p := range patients {
		otherPatients := make([]int, len(patients))
		copy(otherPatients, patients)
		otherPatients[i] = otherPatients[len(otherPatients)-1]
		otherPatients = otherPatients[:len(otherPatients)-1]

		log.Println(port, ": Sending ports", otherPatients, " to", p)
		url := fmt.Sprintf("https://localhost:%d/patients", p)
		patientPorts := types.Patients{
			PortsList: otherPatients,
		}

		b, err := json.Marshal(patientPorts)
		if err != nil {
			log.Fatal(port, ": Error marshalling patientPorts:", err)
		}

		response, err := client.Post(url, "application/json", bytes.NewReader(b))
		if err != nil {
			log.Fatal(port, ": Error posting patientPorts to", p, ":", err)
		}
		log.Println(port, ": Sent ports to ", p, ". Received response code", response.Status)

	}
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println(port, ": Hospital received POST /shares")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, ": Error reading request body:", err)
			w.WriteHeader(400)
			return
		}
		share := &types.Share{}
		err = json.Unmarshal(body, share)
		if err != nil {
			log.Fatal(port, ": Error unmarshalling share:", err)
			w.WriteHeader(400)
			return
		}
		data = data + share.Share
		recshares++
		log.Println(port, ": Hospital received share", share.Share, ", total of", recshares)

		if recshares == patientstotal {
			log.Println("Computation finished: The final value is", data)
		}
		w.WriteHeader(200)
	}
}

func hosserver() {
	log.Println(port, ": Creating hospital server")

	mux := http.NewServeMux()
	mux.HandleFunc("/patient", patientHandler)
	mux.HandleFunc("/shares", shareHandler)

	err := http.ListenAndServeTLS(utilities.PortToString(port), "server.crt", "server.key", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.IntVar(&port, "port", 8080, "port of hospital")
	flag.IntVar(&patientstotal, "t", 3, "the total amount of patients")

	flag.Parse()

	data = 0

	cert, err := os.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	go hosserver()

	for {
	}

}
