package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pusher "github.com/pusher/pusher-http-go"
)

var client = pusher.Client{
	AppID: "963722",
	Key: "8fb6e8dc4e3a6884c1a4",
	Secret: "fce5d56f30ba8dd1e7ea",
	Cluster: "us2",
	Secure: true,
}

type user struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func registerNewUser(rw http.ResponseWriter, req *http.Request){
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var newUser user

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		panic(err)
	}

	client.Trigger("update", "new-user", newUser)

	json.NewEncoder(rw).Encode(newUser)
}

func pusherAuth(res http.ResponseWriter, req *http.Request){
	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePrivateChannel(params)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, string(response))
}

func testEndpoint(res http.ResponseWriter, req *http.Request){
	fmt.Fprintf(res, "API works for ECT\n")
}

func main(){
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/new/user", registerNewUser)
	http.HandleFunc("/pusher/auth",pusherAuth)
	http.HandleFunc("/testEndpoint", testEndpoint)

	log.Fatal(http.ListenAndServe(":8090",nil))
}
