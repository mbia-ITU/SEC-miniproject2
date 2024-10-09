package main

import (
	types "SEC-miniproject2/src/Hospital/Types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var client *http.Client
var patients []int
var regpatients int
var patientstotal int
var port int
var data int
var recshares int

func patientHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(port, "/patiens received")
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
