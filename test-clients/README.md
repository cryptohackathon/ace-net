# ACE test-clients

## Prerequisites

- [go](https://golang.org/),
- [GoFE](https://github.com/fentec-project/gofe)

## About

`client.go` implements routines for data encryption and decryption based on the decentralized multi-client functional encryption scheme for inner product ([DMCFE](https://eprint.iacr.org/2017/989.pdf)) implemented in [GoFE](https://github.com/fentec-project/gofe). It can be used to simulate the ACE* framework process by initializing clients that submit encrypted data to pools and an analytics server that serves the pool management system to extract hisograms from the gathered data.

## Simulation

- `./run-clients.sh` starts the simulation of clients. This is the `CLIENT` operation mode.
- `./run-analytics-server.sh` starts the simultation of the analytics server. This is the `ANALYTICS` operation mode.

## Functionality

`client.go` contains
- API protocol for interaction with the backend system implementation ([api-back](../api-back))
- wrapper for the GoFE DMCFE scheme
- randomized simulation of clients and analytics server