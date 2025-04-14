package myjson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func RunJson() {
	// "1" is different from '1'
	dataToEncode := []interface{}{"1", '1', "😎", true, map[string]any{
		"name": "Alice",
		"face": '😎',
	}, '🫟'}
	for _, eachData := range dataToEncode {
		if bytes, err := json.Marshal(eachData); err == nil {
			fmt.Println("bytes: ", bytes)
			fmt.Println("string(bytes): ", string(bytes))
			fmt.Println()
		}
	}

	/*
		Marshal & UnMarshal
	*/

	p1 := person{Name: "Rob"}
	p2 := person{Name: "Pike"}

	persons := []person{p1, p2}

	bytes, err := json.MarshalIndent(persons, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println("persons as json: ", string(bytes))

	// convert back json to go type person{}
	personsList := []person{}
	err = json.Unmarshal(bytes, &personsList)
	if err != nil {
		panic(err)
	}
	fmt.Println("personsList: ", personsList)

	// or you can unmarshal it as golang map
	personsMap := []map[string]interface{}{}
	err = json.Unmarshal(bytes, &personsMap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("personsMap %#v\n", personsMap)

	/*
		Sending Json in HTTP API Responses
	*/

	mux := http.ServeMux{}
	mux.HandleFunc("POST /encode", encodeHandler)
	mux.HandleFunc("POST /decode", decodeHandler)

	log.Println("Server started at :4567")
	http.ListenAndServe(":4567", &mux)
}

type person struct {
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	dataFromUser := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&dataFromUser)
	personName, ok := dataFromUser["name"].(string)
	if !ok {
		w.Write([]byte(`{"error":"'Name' field not found in the body"}`))
		return
	}
	p := person{
		Name:      personName,
		CreatedAt: time.Now().Format(time.RFC1123),
	}
	json.NewEncoder(w).Encode([]person{p, p})
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	var persons []person
	err := json.NewDecoder(r.Body).Decode(&persons)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Your data is decoded. Persons Count : %d", len(persons))})
}
