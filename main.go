package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status `json:"status"`
}

func updateData() {
	for {

		var data = Data{Status: Status{}}
		waterMin := 1
		waterMax := 30

		data.Status.Water = rand.Intn(waterMax-waterMin+1) + waterMax

		data.Status.Wind = rand.Intn(waterMax-waterMin+1) + waterMax

		b, err := json.MarshalIndent(&data, "", " ")

		if err != nil {
			log.Fatalln("error while marshalling json data  =>", err.Error())
		}

		err = ioutil.WriteFile("data.json", b, 0644)

		if err != nil {
			log.Fatalln("error while writing value to data.json file  =>", err.Error())
		}
		fmt.Println("menggungu 5 detik")
		time.Sleep(time.Second * 5)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	type DataStatus struct {
		Data
		WaterStatus string `json:"waterStatus"`
		WindStatus  string `json:"windStatus"`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, _ := template.ParseFiles("index.html")

		var data = Data{Status: Status{}}

		b, err := ioutil.ReadFile("data.json")

		if err != nil {
			fmt.Fprint(w, "error braderku")
			return
		}

		err = json.Unmarshal(b, &data)
		if err != nil {
			fmt.Fprint(w, "error braderku")
			return
		}

		// check status wind
		windStatus := ""
		if data.Wind > 15 {
			windStatus = "bahaya"
		} else if data.Wind > 6 {
			windStatus = "siaga"
		} else {
			windStatus = "aman"
		}

		// check status water
		waterStatus := ""
		if data.Wind > 8 {
			waterStatus = "bahaya"
		} else if data.Wind > 6 {
			waterStatus = "siaga"
		} else {
			waterStatus = "aman"
		}

		var dataStatus = DataStatus{
			Data:        data,
			WindStatus:  windStatus,
			WaterStatus: waterStatus,
		}

		err = tpl.ExecuteTemplate(w, "index.html", dataStatus)
		if err != nil {
			fmt.Fprint(w, "error braderku")
			return
		}

	})

	http.ListenAndServe(":8080", nil)
}
