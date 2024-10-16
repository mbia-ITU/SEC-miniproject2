# Mandatory Hand-in 2: Secure storage and processing of medical experiment data

This is a Mandatory Hand-in for the course Security 1, BSc at ITU.
The code and report has been written by Mikkel Bistrup Andersen (mbia).

## Structure of repo
The code in this repo is a solution to Mandatory Hand-in 2 for the course Security 1 at ITU. The code is structured as follows:

### Hospital
contains the code for the hospital client/server along with types and utilities. I choose to use this approach because it makes it easier to manage the code and it is the same way I did it during the Distributed Systems course.

The hospital has two endpoints. /patients for handling POST requests with patient ports and /shares for handling POST requests with aggregate shares.

### Patient
contains the code for the patient client/server along with types and utilities.

The patient has two endpoints. /patients for handling POST requests with the list of patients ports and /shares for handling POST requests with the shares of other patients.

## Compile and run instructions
First a certificate and key must be generated for the server. This can be done with the following commands:
- `openssl genrsa -out server.key 2048`
- `openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -addext "subjectAltName = DNS:localhost"`

The server can then be started with the following command in different terminals:
- `go run src/Hospital/hospital.go`
- `go run src/Patient/patient.go -port=8081`
- `go run src/Patient/patient.go -port=8082`
- `go run src/Patient/patient.go -port=8083`
once three instances of the patient has been started along with an instance of the hospital, then the program shoul run as intended.