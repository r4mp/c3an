package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/r4mp/c3an/core"
	"log"
)

type Engel struct {
	Text   string
	Time   string
	Id     int
	Issuer string
}

func receiveMessages() {

	var engel Engel

	client, err := socketio.Dial("http://ara.uberspace.de:62961/")

	if err != nil {
		log.Fatal(err)
	}

	client.On("engel", func(ns *socketio.NameSpace, message string) {
		log.Println(message)

		err := json.Unmarshal([]byte(message), &engel)

		if err == nil {
			core.SendNotificationToAllRegisteredDevices(engel.Text, 0, "bingbong.aiff")
		} else {
			log.Println("Can't unmarshal message: " + message)
		}

	})

	client.On("reload", func(ns *socketio.NameSpace, message string) {
		log.Println(message)
	})

	client.Run()
}
