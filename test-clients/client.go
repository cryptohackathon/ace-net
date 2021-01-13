package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	ClientSequenceId   int    `json:"clientSequenceId"`
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
	registerAPIUrl := host + "/ace/register"

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponseRegistrationInfo
	response, err := http.Post(registerAPIUrl, "application/json", nil)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	} else {
		apiResponseData, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(apiResponseData, &apiResponse)
		if err != nil {
			return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
		}
		*regInfo = apiResponse.Data
		return nil
	}
}

func status(regInfo RegistrationInfo, poolDataPayload *PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/status/" + regInfo.PoolLabel

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponsePoolDataPayload
	response, err := http.Get(APIUrl)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	} else {
		apiResponseData, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(apiResponseData, &apiResponse)
		if err != nil {
			return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
		}
		*poolDataPayload = apiResponse.Data
		return nil
	}
}

func postPublicKeyShare(regInfo RegistrationInfo, poolDataPayload *PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/status/" + regInfo.PoolLabel

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponsePoolDataPayload
	response, err := http.Get(APIUrl)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	} else {
		apiResponseData, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(apiResponseData, &apiResponse)
		if err != nil {
			return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
		}
		*poolDataPayload = apiResponse.Data
		return nil
	}
}

func main() {
	fmt.Println("Starting the client...")

	var regInfo RegistrationInfo
	var statusData PoolDataPayload
	err := register(&regInfo)
	if err != nil {
		fmt.Printf("Error while registring")
	}
	fmt.Print(regInfo)
	err = status(regInfo, &statusData)
	if err != nil {
		fmt.Printf("Error while checking status")
	}
	fmt.Print(statusData)
	fmt.Print(*statusData.CreationTime)

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
	fmt.Println("Terminating the client...")
}
