package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const host = "http://localhost:9500"

// RegistrationInfo - registration info
type RegistrationInfo struct {
	ClientSequenceID   int    `json:"clientSequenceId"`
	PoolLabel          string `json:"poolLabel"`
	RegistrationExpiry string `json:"registrationExpiry"`
	Status             string `json:"status"`
}

// APIResponseRegistrationInfo - API response for RegistrationInfo
type APIResponseRegistrationInfo struct {
	Data         RegistrationInfo `json:"data"`
	ErrorDetails *string          `json:"errorDetails"`
	ErrorMessage *string          `json:"errorMessage"`
	Status       string           `json:"status"`
}

// PoolDataPayload - pool data payload
type PoolDataPayload struct {
	Status           string      `json:"status"`
	CreationTime     *string     `json:"creationTime"`
	RegistrationTime *string     `json:"registrationTime"`
	FinalizationTime *string     `json:"finalizationTime"`
	CalculationTime  *string     `json:"calculationTime"`
	PoolLabel        string      `json:"poolLabel"`
	PoolExpiry       string      `json:"poolExpiry"`
	PublicKeys       *[]string   `json:"publicKeys"`
	CypherTexts      *[][]string `json:"cypherTexts"`
	DecryptionKeys   *[]string   `json:"decryptionKeys"`
	Histogram        *[]int      `json:"histogram"`
}

// APIResponsePoolDataPayload - API response for PoolDataPayload
type APIResponsePoolDataPayload struct {
	Data         PoolDataPayload `json:"data"`
	ErrorDetails *string         `json:"errorDetails"`
	ErrorMessage *string         `json:"errorMessage"`
	Status       string          `json:"status"`
}

// PublicKeyShareRequest - pool data payload
type PublicKeyShareRequest struct {
	ClientSequenceID   int    `json:"clientSequenceId"`
	PoolLabel          string `json:"poolLabel"`
	RegistrationExpiry string `json:"registrationExpiry"`
	KeyShare           string `json:"keyShare"`
}

// CypherAndDKRequest - pool data payload
type CypherAndDKRequest struct {
	ClientSequenceID   int      `json:"clientSequenceId"`
	PoolLabel          string   `json:"poolLabel"`
	CypherText         []string `json:"cypherText"`
	DecryptionKeyShare string   `json:"decryptionKeyShare"`
}

func register(regInfo *RegistrationInfo) error {
	// register
	APIUrl := host + "/ace/register"

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponseRegistrationInfo
	response, err := http.Post(APIUrl, "application/json", nil)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*regInfo = apiResponse.Data
	return nil

}

func status(regInfo RegistrationInfo, poolDataPayload *PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/status/" + regInfo.PoolLabel

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponsePoolDataPayload
	response, err := http.Get(APIUrl)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil

}

func postPublicKeyShare(publicKeyShareReq PublicKeyShareRequest, poolDataPayload *PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/public-key-share"

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponsePoolDataPayload
	jsonValue, _ := json.Marshal(publicKeyShareReq)
	fmt.Printf("PUBLIC: %s\n", jsonValue)
	response, err := http.Post(APIUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil
}

func postCypherAndDecryptionKey(cypherAndDKReq CypherAndDKRequest, poolDataPayload *PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/cyphertext-and-dk"

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponsePoolDataPayload
	jsonValue, _ := json.Marshal(cypherAndDKReq)
	response, err := http.Post(APIUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		fmt.Printf("BB: %s", string(apiResponseData))
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil
}

// TODO: post key share, post encryptions, make robust, random timings, backoff retries on error, test with APIs

func main() {
	fmt.Println("Starting the client...")

	startDelay := time.Duration(rand.Intn(500))
	time.Sleep(startDelay * time.Millisecond)

	var regInfo RegistrationInfo
	var statusData PoolDataPayload
	var publicKeyShareReq PublicKeyShareRequest
	var cypherAndDKReq CypherAndDKRequest

	err := register(&regInfo)
	if err != nil {
		fmt.Println("Error while registring. Terminating client.")
		return
	}
	fmt.Println("REGISTRED")
	fmt.Println(regInfo)

	publicKeyShareReq.ClientSequenceID = regInfo.ClientSequenceID
	publicKeyShareReq.PoolLabel = regInfo.PoolLabel
	publicKeyShareReq.RegistrationExpiry = regInfo.RegistrationExpiry
	publicKeyShareReq.KeyShare = fmt.Sprintf("<KEY-SHARE-%d>", regInfo.ClientSequenceID)

	err = postPublicKeyShare(publicKeyShareReq, &statusData)
	if err != nil {
		fmt.Println("Error while checking status. Terminating client.")
		return
	}
	fmt.Println("PUBLIC KEY SHARE SUBMITTED")
	fmt.Println(statusData.Status)

	cnt := 0
	const pollingIterationsLimit = 10
	for statusData.Status != "ENCRYPTION" {
		cnt = cnt + 1
		if cnt > 10 {
			fmt.Printf("Too many polling iterrations (%d). Terminating client.", cnt)
			return
		}
		err = status(regInfo, &statusData)
		if err != nil {
			fmt.Println("Error while checking status. Terminating client.")
			return
		}
		fmt.Println(statusData)
		if statusData.Status == "EXPIRED" {
			fmt.Println("Pool expired. Terminating client.")
			return
		}
		delay := time.Duration(500 + rand.Intn(1000))
		// fmt.Println("POLL DELAY: %f", delay)
		time.Sleep(delay * time.Millisecond)
	}

	fmt.Println("READY TO SUBMIT ENCRYPTION")
	fmt.Println(statusData.Status)

	cypherAndDKReq.ClientSequenceID = regInfo.ClientSequenceID
	cypherAndDKReq.PoolLabel = regInfo.PoolLabel
	cypherAndDKReq.CypherText = make([]string, 5)
	copy(cypherAndDKReq.CypherText, []string{"CY", "PH", "ER", "TE", "XT"})
	cypherAndDKReq.DecryptionKeyShare = fmt.Sprintf("<DEC-KEY-%d>", regInfo.ClientSequenceID)

	err = postCypherAndDecryptionKey(cypherAndDKReq, &statusData)
	if err != nil {
		fmt.Printf("Error while checking status: %v. Terminating client", err)
		return
	}

	fmt.Println("ENCRYPTION SUBMITTED")
	fmt.Println(statusData.Status)

	// fmt.Print(*statusData.CreationTime)

	// response, err := http.Get("https://httpbin.org/ip")
	// if err != nil {
	//     fmt.Printf("The HTTP request failed with error %s\n", err)
	// } else {
	//     data, _ := ioutil.ReadAll(response.Body)
	//     fmt.Println(string(data))
	// }
	// jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
	// jsonValue, _ := json.Marshal(jsonData)
	// response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
	// if err != nil {
	//     fmt.Printf("The HTTP request failed with error %s\n", err)
	// } else {
	//     data, _ := ioutil.ReadAll(response.Body)
	//     fmt.Println(string(data))
	// }
	fmt.Println("END. Terminating client.")
}
