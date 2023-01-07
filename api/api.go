package api

import(
	"fmt"
	"net/http"
	"io"
)

// Informs load balancer that one user has been processed completely.
// Meaning user has left the server
// loadBalancer (format = https://url:port, http://url:port)
func LoadReleased(loadBalancer string){
	resp, err := http.Post(fmt.Sprintf("%v/v1/server/dequeue", loadBalancer), "text/plain", nil)
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
func AddServer(loadBalancer string){
	resp, err := http.Post(fmt.Sprintf("%v/v1/server/add", loadBalancer), "text/plain", nil)
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
func RemoveServer(loadBalancer string){
	resp, err := http.Post(fmt.Sprintf("%v/v1/server/remove", loadBalancer), "text/plain", nil)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	responseBody,_ := io.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
}