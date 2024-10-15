package main

import (
	types "SEC-miniproject2/src/Patient/Types"
	utilities "SEC-miniproject2/src/Patient/Utilities"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var client *http.Client
var port int
var data int
var recshares []int
var patientstotal int
var hosport int
var maxdata int

func patientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println(port, "/patients POST received")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, "Error reading body: ", err)
			w.WriteHeader(400)
			return
		}

		patients := &types.Patients{}
		err = json.Unmarshal(body, patients)
		if err != nil {
			log.Fatal(port, "Error unmarshalling: ", err)
			w.WriteHeader(400)
			return
		}

		shares := utilities.CreateShares(maxdata, data, patientstotal)

		log.Println(port, "Sending shares")

		for i, share := range shares {
			if i == patientstotal-1 {
				break
			}
			shareowned := types.Share{Share: share}

			b, err := json.Marshal(shareowned)
			if err != nil {
				log.Fatal(port, "Error marshalling: ", err)
				w.WriteHeader(400)
				return
			}
			url := fmt.Sprintf("http://localhost:%d/shares", patients.PortsList[i])
			resp, err := client.Post(url, "string", bytes.NewReader(b))
			if err != nil {
				log.Fatal(port, "Error sending share to: ", patients.PortsList[i], " : ", err)
				w.WriteHeader(400)
				return
			}
			log.Println(port, "Sent share to: ", patients.PortsList[i])
			log.Println(port, "Recieved response: ", resp.StatusCode)
		}

		recshares = append(recshares, shares[len(shares)-1])

		if len(recshares) == patientstotal {
			sendaggregatedshares()
		}
		w.WriteHeader(200)
	}
}

func sendaggregatedshares() {
	log.Println(port, "Sending aggregated shares")

	var aggregatedshares int

	for _, share := range recshares {
		aggregatedshares += share
	}

	log.Println(port, "Aggregated shares: ", aggregatedshares)

	aggregate := types.Share{Share: aggregatedshares}

	b, err := json.Marshal(aggregate)
	if err != nil {
		log.Fatal(port, "Error marshalling: ", err)
		return
	}

	log.Println(port, "Sending aggregated shares: ", aggregatedshares)
	url := fmt.Sprintf("http://localhost:%d/shares", hosport)
	resp, err := client.Post(url, "string", bytes.NewReader(b))
	if err != nil {
		log.Fatal(port, "Error sending aggregated shares: ", err)
		return
	}

	log.Println(port, "Sent aggregated shares, recieved response: ", resp.StatusCode)
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println(port, "/shares POST received")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, "Error reading body: ", err)
			w.WriteHeader(400)
			return
		}

		othershare := &types.Share{}
		err = json.Unmarshal(body, othershare)
		if err != nil {
			log.Fatal(port, "Error unmarshalling: ", err)
			w.WriteHeader(400)
			return
		}

		recshares = append(recshares, othershare.Share)

		if len(recshares) == patientstotal {
			sendaggregatedshares()
		}
		w.WriteHeader(200)
	}
}

func patserver() {
	log.Println(port, "Starting patient server")

	mux := http.NewServeMux()
	mux.HandleFunc("/patients", patientHandler)
	mux.HandleFunc("/shares", shareHandler)

	err := http.ListenAndServeTLS(utilities.PortToString(port), "server.crt", "server.key", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var R int

	flag.IntVar(&port, "port", 81, "Port number for patient")
	flag.IntVar(&hosport, "hospitalport", 80, "Port number for hospital")
	flag.IntVar(&patientstotal, "t", 3, "Total number of patients")
	flag.IntVar(&R, "R", 500, "Threshold value")

	flag.Parse()

	maxdata = R / 3

	rand.Seed(time.Now().UnixNano())
	data = rand.Intn(maxdata)

	log.Println(port, "New patient with data: ", data)

	cert, err := os.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(cert)

	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: certpool}}}

	go patserver()

	log.Println(port, "Patient server started")

	url := fmt.Sprintf("http://localhost:%d/patient", hosport)

	patPort := types.Patient{Port: port}

	b, err := json.Marshal(patPort)
	if err != nil {
		log.Fatal(port, "Error marshalling: ", err)
		return
	}

	resp, err := client.Post(url, "string", bytes.NewReader(b))
	if err != nil {
		log.Fatal(port, "Error sending patient port to hospital: ", err)
		return
	}

	log.Println(port, "Sent patient port to hospital, recieved response: ", resp.StatusCode)

	for {
	}

}
