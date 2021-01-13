package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Starting the client...")

	// register
	registerAPIUrl := "http://localhost:9500/ace/register"

	response, err := http.Post(registerAPIUrl, "application/json", nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

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
