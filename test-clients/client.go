package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"github.com/fentec-project/bn256"
	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/fullysec"
	"github.com/pkg/errors"
)

// RegistrationInfo - registration info
type RegistrationInfo struct {
	ClientSequenceID   int      `json:"clientSequenceId"`
	PoolLabel          string   `json:"poolLabel"`
	RegistrationExpiry string   `json:"registrationExpiry"`
	Status             string   `json:"status"`
	SlotLabels         []string `json:"slotLabels"`
	InnerVector        []int    `json:"innerVector"`
}

// APIResponseRegistrationInfo - API response for RegistrationInfo
type APIResponseRegistrationInfo struct {
	Data         RegistrationInfo        `json:"data"`
	ErrorDetails *map[string]interface{} `json:"errorDetails"`
	ErrorMessage *map[string]interface{} `json:"errorMessage"`
	Status       string                  `json:"status"`
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
	DecryptionKeys   *[][]string `json:"decryptionKeys"`
	Histogram        *[]int      `json:"histogram"`
	SlotLabels       []string    `json:"slotLabels"`
	InnerVector      []int       `json:"innerVector"`
}

// APIResponsePoolDataPayload - API response for PoolDataPayload
type APIResponsePoolDataPayload struct {
	Data         PoolDataPayload `json:"data"`
	ErrorDetails *string         `json:"errorDetails"`
	ErrorMessage *string         `json:"errorMessage"`
	Status       string          `json:"status"`
}

// APIResponsePoolDataPayloadArray - API response for PoolDataPayload
type APIResponsePoolDataPayloadArray struct {
	Data         []PoolDataPayload `json:"data"`
	ErrorDetails *string           `json:"errorDetails"`
	ErrorMessage *string           `json:"errorMessage"`
	Status       string            `json:"status"`
}

// PublicKeyShareRequest - public key share
type PublicKeyShareRequest struct {
	ClientSequenceID   int    `json:"clientSequenceId"`
	PoolLabel          string `json:"poolLabel"`
	RegistrationExpiry string `json:"registrationExpiry"`
	KeyShare           string `json:"keyShare"`
}

// CypherAndDKRequest - cyphertext vector and partial decryption key
type CypherAndDKRequest struct {
	ClientSequenceID   int      `json:"clientSequenceId"`
	PoolLabel          string   `json:"poolLabel"`
	CypherText         []string `json:"cypherText"`
	DecryptionKeyShare []string `json:"decryptionKeyShare"`
}

// HistogramPayload - histogram payload
type HistogramPayload struct {
	Secret    string `json:"secret"`
	PoolLabel string `json:"poolLabel"`
	Histogram []int  `json:"histogram"`
}

func register(host string, regInfo *RegistrationInfo) error {
	APIUrl := host + "/ace/register"

	var apiResponse APIResponseRegistrationInfo
	response, err := http.Post(APIUrl, "application/json", nil)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("REGISTER RESPONSE: %v\n", string(apiResponseData))
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*regInfo = apiResponse.Data
	return nil

}

func status(host string, regInfo RegistrationInfo, poolDataPayload *PoolDataPayload) error {
	APIUrl := host + "/ace/status/" + regInfo.PoolLabel

	var apiResponse APIResponsePoolDataPayload
	response, err := http.Get(APIUrl)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("STATUS RESPONSE: %s\n", string(apiResponseData))
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil

}

func postPublicKeyShare(host string, publicKeyShareReq PublicKeyShareRequest, poolDataPayload *PoolDataPayload) error {
	APIUrl := host + "/ace/public-key-share"

	var apiResponse APIResponsePoolDataPayload
	jsonValue, _ := json.Marshal(publicKeyShareReq)
	fmt.Printf("PUBLIC: %s\n", jsonValue)
	response, err := http.Post(APIUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("POST PUBLIC KEY SHARE RESPONSE: %s\n", string(apiResponseData))
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil
}

func postCypherAndDecryptionKey(host string, cypherAndDKReq CypherAndDKRequest, poolDataPayload *PoolDataPayload) error {
	APIUrl := host + "/ace/cyphertext-and-dk"

	var apiResponse APIResponsePoolDataPayload
	jsonValue, _ := json.Marshal(cypherAndDKReq)
	response, err := http.Post(APIUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("POST CYPHER AND DECRYPT KEY REQUEST: %s\n", string(jsonValue))
	fmt.Printf("POST CYPHER AND DECRYPT KEY RESPONSE: %s\n", string(apiResponseData))
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil
}

func postHistogram(host string, histogramPayload HistogramPayload, poolDataPayload *PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/histogram"

	var apiResponse APIResponsePoolDataPayload
	jsonValue, _ := json.Marshal(histogramPayload)
	response, err := http.Post(APIUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("HIST REQUEST: %s\n", string(jsonValue))
	fmt.Printf("HIST RESPONSE: %s\n", string(apiResponseData))
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayload = apiResponse.Data
	return nil
}

func listFinalized(host string, poolDataPayloadArray *[]PoolDataPayload) error {
	// register
	APIUrl := host + "/ace/pools?status=FINALIZED"

	// var registrationInfo RegistrationInfo
	var apiResponse APIResponsePoolDataPayloadArray
	response, err := http.Get(APIUrl)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}
	apiResponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("POLL RESPONSE: %s\n", string(apiResponseData))
	err = json.Unmarshal(apiResponseData, &apiResponse)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal APIResponse: %s", err)
	}
	*poolDataPayloadArray = apiResponse.Data
	return nil
}

type Client struct {
	Encryptor *fullysec.DMCFEClient
}

func (c *Client) init(sequence int) error {
	encryptor, err := fullysec.NewDMCFEClient(sequence)
	if err != nil {
		return errors.Wrap(err, "could not instantiate fullysec")
	}
	c.Encryptor = encryptor
	return nil
}

func (c *Client) getPublicKeyEncoding() string {
	return c.g1Base64Encoding(c.Encryptor.ClientPubKey)
}

func (c *Client) setShare(encodedPublicKeys []string) error {
	var err error
	publicKeys := make([]*bn256.G1, len(encodedPublicKeys))
	for i := 0; i < len(encodedPublicKeys); i++ {
		publicKeys[i], err = g1Base64Decoding(encodedPublicKeys[i])
		if err != nil {
			return err
		}
	}
	err = c.Encryptor.SetShare(publicKeys)
	return err
}

func (c *Client) encryptData(data []int, labels []string) ([]string, error) {
	ciphers := make([]string, len(data))
	for i := 0; i < len(data); i++ {
		g1, err := c.Encryptor.Encrypt(big.NewInt(int64(data[i])), labels[i])
		if err != nil {
			return nil, err
		}
		ciphers[i] = c.g1Base64Encoding(g1)
	}
	return ciphers, nil
}

func (c *Client) deriveKeyShare(vector []int) ([]string, error) {

	// TEST
	// oneVec := data.NewConstantVector(3, big.NewInt(1))
	// keyShare, err := c.Encryptor.DeriveKeyShare(oneVec)

	keyShare, err := c.Encryptor.DeriveKeyShare(toVector(vector))
	if err != nil {
		return nil, err
	}
	fmt.Printf("MAKE %d %v LL: %d", len(vector), vector, len(keyShare))
	// LOOKS LIKE KEYSHARES ARE OF DIFFERENT LENGTHS ...
	// result := make([]string, len(vector))
	// for i := 0; i < len(vector); i++ {
	// 	result[i] = g2Base64Encoding(keyShare[i])
	// }

	result := make([]string, len(keyShare))
	for i := 0; i < len(keyShare); i++ {
		result[i] = g2Base64Encoding(keyShare[i])
	}

	return result, nil
}

func (c *Client) g1Base64Encoding(g1 *bn256.G1) string {
	return base64.StdEncoding.EncodeToString(g1.Marshal())
}

func decryptHistogram(ciphers [][]string, keyShares [][]string, labels []string, vector []int) ([]int, error) {
	var err error
	y := toVector(vector)
	b := big.NewInt(int64(len(vector)))
	keys := make([]data.VectorG2, len(keyShares))
	for i := 0; i < len(keyShares); i++ {
		keys[i] = make(data.VectorG2, len(keyShares[i]))
		for j := 0; j < len(keyShares[i]); j++ {
			keys[i][j], err = g2Base64Decoding(keyShares[i][j])
			if err != nil {
				return nil, err
			}
		}
	}

	fmt.Printf("LABS: %v\n", labels)
	histogram := make([]int, len(labels))
	for i := 0; i < len(labels); i++ {
		ciphersi := make([]*bn256.G1, len(vector))
		for j := 0; j < len(vector); j++ {
			ciphersi[j], err = g1Base64Decoding(ciphers[j][i])
			if err != nil {
				return nil, err
			}
		}
		value, err := fullysec.DMCFEDecrypt(ciphersi, keys, y, labels[i], b)
		if err != nil {
			return nil, err
		}
		histogram[i] = int(value.Int64())
	}
	return histogram, nil
}

func g2Base64Encoding(g2 *bn256.G2) string {
	return base64.StdEncoding.EncodeToString(g2.Marshal())
}

func g1Base64Decoding(b64 string) (*bn256.G1, error) {
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	g1 := new(bn256.G1)
	_, err = g1.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}
	return g1, nil
}

func g2Base64Decoding(b64 string) (*bn256.G2, error) {
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	g2 := new(bn256.G2)
	_, err = g2.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}
	return g2, nil
}

func toVector(vector []int) data.Vector {
	components := make([]*big.Int, len(vector))
	for i := 0; i < len(vector); i++ {
		components[i] = big.NewInt(int64(vector[i]))
	}
	return data.NewVector(components)
}

func simulateExposure(size int) []int {
	exposure := make([]int, size)
	index := rand.Intn(size)
	exposure[index] = 1
	return exposure
}

func simulateClient(host string) {
	fmt.Println("Starting the client...")

	// Initial delay
	startDelay := time.Duration(rand.Intn(500))
	time.Sleep(startDelay * time.Millisecond)

	var regInfo RegistrationInfo
	var statusData PoolDataPayload
	var publicKeyShareReq PublicKeyShareRequest
	var cypherAndDKReq CypherAndDKRequest

	// Registration
	err := register(host, &regInfo)
	if err != nil {
		fmt.Printf("Error while registring: %v. Terminating client.", err)
		return
	}

	// Client initialization
	client := new(Client)
	err = client.init(regInfo.ClientSequenceID)
	if err != nil {
		fmt.Printf("Error in client initialization: %v\n", err)
		return
	}

	fmt.Println("REGISTRED")
	fmt.Println(regInfo)

	publicKeyShareReq.ClientSequenceID = regInfo.ClientSequenceID
	publicKeyShareReq.PoolLabel = regInfo.PoolLabel
	publicKeyShareReq.RegistrationExpiry = regInfo.RegistrationExpiry

	// Set public key share
	publicKeyShareReq.KeyShare = client.getPublicKeyEncoding() // fmt.Sprintf("<KEY-SHARE-%d>", regInfo.ClientSequenceID)

	// Sending public key share to central server
	err = postPublicKeyShare(host, publicKeyShareReq, &statusData)
	if err != nil {
		fmt.Println("Error while posting public key share status. Terminating client.")
		return
	}
	fmt.Println("PUBLIC KEY SHARE SUBMITTED")
	fmt.Println(statusData.Status)

	// Waiting for collection of all key shares on server
	cnt := 0
	const pollingIterationsLimit = 100
	for statusData.Status != "ENCRYPTION" {
		cnt = cnt + 1
		if cnt > pollingIterationsLimit {
			fmt.Printf("Too many polling iterrations (%d). Terminating client.", cnt)
			return
		}
		err = status(host, regInfo, &statusData)
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

	client.setShare(*statusData.PublicKeys)
	fmt.Println("PUBLIC KEY OBTAINED")
	fmt.Println(statusData.Status)

	cypherAndDKReq.ClientSequenceID = regInfo.ClientSequenceID
	cypherAndDKReq.PoolLabel = regInfo.PoolLabel

	// ==========
	// TODO: provide real data and labels
	// cypherAndDKReq.CypherText = make([]string, 5)
	// copy(cypherAndDKReq.CypherText, []string{"CY", "PH", "ER", "TE", "XT"})
	exposure := simulateExposure(len(statusData.SlotLabels))
	fmt.Printf("Simulating exposure %v\n", exposure)
	cypherAndDKReq.CypherText, err = client.encryptData(exposure, statusData.SlotLabels)
	// []int{1, 0, 0, 0, 0},
	// []string{"CY", "PH", "ER", "TE", "XT"}

	if err != nil {
		fmt.Printf("Error encrypting client data: %v. Terminating client", err)
		return
	}
	// cypherAndDKReq.DecryptionKeyShare = fmt.Sprintf("<DEC-KEY-%d>", regInfo.ClientSequenceID)
	// TODO: provide the size of pool or vector
	// WARN: DecryptionKeyShare is an array of strings
	fmt.Printf("INN: %v", statusData.InnerVector)

	cypherAndDKReq.DecryptionKeyShare, err = client.deriveKeyShare(statusData.InnerVector)
	// []int{1, 1, 1, 1, 1, 1, 1, 1}
	if err != nil {
		fmt.Printf("Error deriving client key share: %v. Terminating client", err)
		return
	}
	// ==========

	err = postCypherAndDecryptionKey(host, cypherAndDKReq, &statusData)
	if err != nil {
		fmt.Printf("Error while posting cyphertexts and decrypt keys status: %v. Terminating client", err)
		return
	}

	fmt.Println("ENCRYPTION SUBMITTED")
	fmt.Println(statusData.Status)

	fmt.Println("END. Terminating client.")

}

func simulateAnalyticsServer(host string, secret string) error {
	var poolDataPayloadArray []PoolDataPayload
	var statusData PoolDataPayload
	var histogramPayload HistogramPayload
	histogramPayload.Secret = secret
	for {
		err := listFinalized(host, &poolDataPayloadArray)
		if err != nil {
			fmt.Printf("Error while polling status: %v. Terminating client.", err)
			return err
		}
		fmt.Printf("FINALIZED:\n")

		for i := 0; i < len(poolDataPayloadArray); i++ {
			histogramPayload.PoolLabel = poolDataPayloadArray[i].PoolLabel

			// Decrypt histogram
			histogram, err := decryptHistogram(
				*poolDataPayloadArray[i].CypherTexts,
				*poolDataPayloadArray[i].DecryptionKeys,
				poolDataPayloadArray[i].SlotLabels,
				poolDataPayloadArray[i].InnerVector)
			if err != nil {
				fmt.Printf("%v x %v\n", len(*poolDataPayloadArray[i].CypherTexts), len((*poolDataPayloadArray[i].CypherTexts)[0]))
				fmt.Printf("%v\n", len(*poolDataPayloadArray[i].DecryptionKeys))
				// fmt.Println(*poolDataPayloadArray[i].DecryptionKeys)
				fmt.Println(poolDataPayloadArray[i].SlotLabels)
				fmt.Println(poolDataPayloadArray[i].InnerVector)
				fmt.Printf("Error decrypting histogram: %v. Terminating client.\n", err)
				return err
			}

			if err != nil {
				fmt.Printf("Histogram calculation failed: %v\n", err)
				return err
			}
			// ==========
			// TODO: deserialize and decrypt
			// fmt.Printf("%d: \n%v\n%v\n", i, *cypherTextPtr, *decryptionKeysPtr)
			fmt.Printf("Histogram: %v\n", histogram)

			// ==========
			// TODO: post decrypted values and create real histogram

			exampleHistogram := make([]*big.Int, 0)
			exampleHistogram = append(exampleHistogram, new(big.Int).SetInt64(1))
			exampleHistogram = append(exampleHistogram, new(big.Int).SetInt64(2))
			exampleHistogram = append(exampleHistogram, new(big.Int).SetInt64(5))
			exampleHistogram = append(exampleHistogram, new(big.Int).SetInt64(0))
			exampleHistogram = append(exampleHistogram, new(big.Int).SetInt64(2))
			// ==========

			histogramPayload.Histogram = histogram
			err = postHistogram(host, histogramPayload, &statusData)
			if err != nil {
				fmt.Printf("Error while posting histogram status: %v. Terminating client.\n", err)
				return err
			}
			fmt.Printf("HISTOGRAM SUBMITTED: %v\n", histogramPayload)
			fmt.Println(statusData.Status)
			// return nil
		}

		delay := time.Duration(2000)
		time.Sleep(delay * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	modePtr := flag.String("mode", "CLIENT", "Client operation mode: CLIENT or ANALYTICS")
	hostPtr := flag.String("host", "http://localhost:9500", "URL of central server")
	secretPtr := flag.String("secret", "", "URL of central server")
	flag.Parse()
	// 	argsWithProg := os.Args

	// 	const defaultHost = "http://localhost:9500"
	// const defaultMode = "CLIENT"

	mode := *modePtr
	host := *hostPtr
	secret := *secretPtr

	fmt.Printf("Running in %s mode. Connected to %s\n", mode, host)
	if mode == "CLIENT" {
		simulateClient(host)
		return
	}

	if mode == "ANALYTICS" {
		err := simulateAnalyticsServer(host, secret)
		if err != nil {
			fmt.Printf("Error in analytics server: %s. Terminating client.", err)
		}
		return
	}

	fmt.Printf("ERROR: Wrong mode: %s\n", mode)
}
