package api

import(
	"fmt"
	"net/http"
	"io"
	"bytes"
	"encoding/json"
)

type HostRequestBody struct{
	Host string `json:"host"`
}

// Informs load balancer that one user has been processed completely.
// Meaning user has left the server
// loadBalancer (format = https://url:port, http://url:port)
func LoadReleased(loadBalancer, serverHostUrl string){
	reqBody, err := json.Marshal(HostRequestBody{Host: serverHostUrl})
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	resp, err := http.Post(fmt.Sprintf("%v/v1/server/dequeue", loadBalancer), "application/json", bytes.NewReader(reqBody))
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	responseBody,_ := io.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
}


// Adds the current server to the load balancer
// loadBalancer (format = https://url:port, http://url:port)
func AddServer(loadBalancer, serverHostUrl string){
	reqBody, err := json.Marshal(HostRequestBody{Host: serverHostUrl})
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	resp, err := http.Post(fmt.Sprintf("%v/v1/server/add", loadBalancer), "application/json", bytes.NewReader(reqBody))

	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
}


// removes the current server from the load balancer
// loadBalancer (format = https://url:port, http://url:port)
func RemoveServer(loadBalancer, serverHostUrl string){
	reqBody, err := json.Marshal(HostRequestBody{Host: serverHostUrl})
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	resp, err := http.Post(fmt.Sprintf("%v/v1/server/remove", loadBalancer), "application/json", bytes.NewReader(reqBody))
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	responseBody,_ := io.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
}