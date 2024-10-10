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
	log.Println(port, "/patients received")
	if r.Method == "POST" {
		log.Println(port, "/patients POST received")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, "Error reading body: ", err)
			w.WriteHeader(400)
			return
		}
		recport := &types.Patient{}
		err = json.Unmarshal(body, recport)
		if err != nil {
			log.Fatal(port, "Error unmarshalling: ", err)
			w.WriteHeader(400)
			return
		}
		log.Println(port, "Received patient, Port: ", recport.Port)
		patients = append(patients, recport.Port)
		regpatients++
		if regpatients == patientstotal {
			sendport()
		}
		w.WriteHeader(200)
	}
}

func sendport() {
	log.Println(port, ": Sending port")
	for i, p := range patients {
		patientsother := make([]int, len(patients))
		copy(patientsother, patients)
		patientsother[i] = patientsother[len(patientsother)-1]
		patientsother = patientsother[:len(patientsother)-1]

		log.Println(port, ": Sending ", patientsother, "to", p)
		url := fmt.Sprintf("https://localhost:%d/patients", p)
		patientport := types.Patients{PortsList: patientsother}

		b, err := json.Marshal(patientport)
		if err != nil {
			log.Fatal(port, "Error marshalling: ", err)
		}

		resp, err := client.Post(url, "application/json", bytes.NewReader(b))
		if err != nil {
			log.Fatal(port, ": Error sending ", p, ":", err)
		}
		log.Println(port, ": Sent port to ", p, ". Received response", resp.Status)

	}
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println(port, "/share received")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, "Error reading body: ", err)
			w.WriteHeader(400)
			return
		}
		recshare := &types.Share{}
		err = json.Unmarshal(body, recshare)
		if err != nil {
			log.Fatal(port, "Error unmarshalling: ", err)
			w.WriteHeader(400)
			return
		}
		data += recshare.Share
		recshares++
		log.Println(port, "Received share: ", recshare.Share)

		if recshares == patientstotal {
			log.Println("Final value is: ", data)
		}
		w.WriteHeader(200)

	}
}

func hosserver() {
	log.Println(port, ": Starting server")

	mux := http.NewServeMux()
	mux.HandleFunc("/patients", patientHandler)
	mux.HandleFunc("/share", shareHandler)

	err := http.ListenAndServeTLS(utilities.PortToString(port), "server.crt", "server.key", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.IntVar(&port, "port", 80, "Port number")
	flag.IntVar(&patientstotal, "t", 3, "Total number of patients")
	flag.Parse()

	data = 0
	certificate, err := os.ReadFile("server.crt")
	if err != nil {
		log.Fatal("Error reading certificate: ", err)
	}

	certificatepool := x509.NewCertPool()
	certificatepool.AppendCertsFromPEM(certificate)

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certificatepool,
			},
		},
	}

	go hosserver()

	for {
	}

}
