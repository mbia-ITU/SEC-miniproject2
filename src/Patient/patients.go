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
		log.Println(port, ": Patient received POST /patients")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, ": Error when reading request:", err)
			w.WriteHeader(400)
			return
		}
		patients := &types.Patients{}
		err = json.Unmarshal(body, patients)
		if err != nil {
			log.Fatal(port, ": Error unmarshalling:", err)
			w.WriteHeader(400)
			return
		}

		shares := utilities.CreateShares(maxdata, data, patientstotal)

		log.Println(port, ": Sending shares to patients")
		for i, share := range shares {
			if i == patientstotal-1 {
				break
			}
			ownShare := types.Share{
				Share: share,
			}
			b, err := json.Marshal(ownShare)
			if err != nil {
				log.Fatal(port, ": Error when marshalling:", err)
				w.WriteHeader(400)
				return
			}
			url := fmt.Sprintf("https://localhost:%d/shares", patients.PortsList[i])
			resp, err := client.Post(url, "string", bytes.NewReader(b))
			if err != nil {
				log.Fatal(port, ": Error sending share to", patients.PortsList[i], ":", err)
				w.WriteHeader(400)
				return
			}
			log.Println(port, ": Sent share to, ", patients.PortsList[i], ". Received response code:", resp.StatusCode)
		}

		recshares = append(recshares, shares[len(shares)-1])

		if len(recshares) == patientstotal {
			sendaggregatedshares()
		}

		w.WriteHeader(200)
	}
}

func sendaggregatedshares() {
	log.Println(port, ": Computing aggregate share")

	var aggregateShare int

	for _, share := range recshares {
		aggregateShare = aggregateShare + share
	}

	log.Println(port, ": aggregate share is ", aggregateShare)

	aggregate := types.Share{
		Share: aggregateShare,
	}

	b, err := json.Marshal(aggregate)
	if err != nil {
		log.Fatal(port, ": Error marshalling aggregate share:", err)
		return
	}

	log.Println(port, ": Sending aggregate share", aggregateShare, "to hospital")
	url := fmt.Sprintf("https://localhost:%d/shares", hosport)
	response, err := client.Post(url, "string", bytes.NewReader(b))
	if err != nil {
		log.Fatal(port, ": Error sending aggregate share:", err)
		return
	}
	log.Println(port, ": Sent aggregate share to hospital, received response code", response.StatusCode)
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println(port, ": Patient received POST /shares")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(port, ": Error reading request:", err)
			w.WriteHeader(400)
			return
		}
		foreignShare := &types.Share{}
		err = json.Unmarshal(body, foreignShare)
		if err != nil {
			log.Fatal(port, ": Error unmarshalling share:", err)
			w.WriteHeader(400)
			return
		}

		recshares = append(recshares, foreignShare.Share)

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
	var r int

	flag.IntVar(&port, "port", 8081, "port for patient")
	flag.IntVar(&hosport, "h", 8080, "port of the hospital")
	flag.IntVar(&patientstotal, "t", 3, "the total amount of patients")
	flag.IntVar(&r, "r", 500, "the max value that the final computation can have")

	flag.Parse()

	maxdata = r / 3

	data = rand.Intn(maxdata)

	log.Println(port, ": New patient with data =", data)

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

	go patserver()

	log.Println(port, ": Patient registering with hospital")

	url := fmt.Sprintf("https://localhost:%d/patient", hosport)

	ownPort := types.Patient{
		Port: port,
	}

	b, err := json.Marshal(ownPort)
	if err != nil {
		log.Fatal(port, ": Error when marshalling patient:", err)
	}

	response, err := client.Post(url, "string", bytes.NewReader(b))
	if err != nil {
		log.Fatal(port, ": Error when regisering with hospital:", err)
	}
	log.Println(port, ": Registered with hospital, received response code", response.Status)

	for {
	}

}
