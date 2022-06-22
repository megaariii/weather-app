package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

type Cuaca struct {
	Water 			int
	Wind			int
	WaterStatus 	string
	WindStatus 		string
}

var cuaca Cuaca

func main() {
	go statusCuaca()
	mux := http.NewServeMux()
	endPoint := http.HandlerFunc(updateCuaca)
	mux.Handle("/", middleware1(middleware2(endPoint)))
	fmt.Println("Listening to port 8080")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func statusCuaca() {
	for {
		cuaca.Water = rand.Intn(100)
		cuaca.Wind = rand.Intn(100)

		if cuaca.Water < 5 {
			cuaca.WaterStatus = "Aman"
		} else if (cuaca.Water >= 6 && cuaca.Water <= 8)  {
			cuaca.WaterStatus = "Siaga"
		} else {
			cuaca.WaterStatus = "Bahaya"
		}
	
		if cuaca.Wind < 6 {
			cuaca.WindStatus = "Aman"
		} else if (cuaca.Wind >= 7 && cuaca.Wind <= 15)  {
			cuaca.WindStatus = "Siaga"
		} else {
			cuaca.WindStatus = "Bahaya"
		}

		jsonStrin, _ := json.Marshal(&cuaca)
		ioutil.WriteFile("cuaca.json", jsonStrin, os.ModePerm)

		time.Sleep(15 * time.Second)
	}
}

func updateCuaca(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./cuaca.json")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	
	err = json.Unmarshal(content, &cuaca)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	
	tpl, _ := template.ParseFiles("template.html")

	data := Cuaca {
		Water: cuaca.Water,
		Wind: cuaca.Wind,
		WaterStatus: cuaca.WaterStatus,
		WindStatus: cuaca.WindStatus,
	}


	tpl.Execute(w, data)
}

func middleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Pertama")
		next.ServeHTTP(w, r)
	})
}

func middleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Kedua")
		next.ServeHTTP(w, r)
	})
}