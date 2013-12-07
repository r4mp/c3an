package core

import (
	"encoding/json"
	"github.com/anachronistic/apns"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

var server = "localhost"
var db = "c3an"
var collection = "users"

type ApnResult struct {
	Alert   interface{} `json:"alert,omitempty"`
	Success bool        `json:"success,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type User struct {
	Name        string
	DeviceToken string
}

func SendNotificationToAllRegisteredDevices(message string, badge int, sound string) /*([]byte, error)*/ {

	session, err := mgo.Dial(server)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	result := []User{}
	c := session.DB(db).C(collection)
	err = c.Find(nil).All(&result) // TODO: Checken, ob 'nil' hier zu Problemen fuehren kann!!!

	if err != nil {
		//panic(err)
		return // Wenn Collection oder Abfragefeld nicht vorhanden
	} else {
		for _, element := range result {
			log.Println(element.DeviceToken)
			SendNotificationToSingleDevice(message, badge, sound, element.DeviceToken)
		}
	}
}

func SendNotificationToSingleDevice(message string, badge int, sound string, token string) ([]byte, error) {

	payload := apns.NewPayload()

	payload.Alert = message
	payload.Badge = badge
	payload.Sound = sound

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "certs/cert.pem", "certs/key.pem")
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()

	res := ApnResult{alert, resp.Success, resp.Error}

	return json.Marshal(res)
}

func RegisterDevice(token string) {
	session, err := mgo.Dial(server)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	if checkIfDeviceIsAlreadyRegisted(token) || token == "" {
		return
	} else {
		c := session.DB(db).C(collection)
		err = c.Insert(&User{"Name", token})

		if err != nil {
			panic(err)
		}
	}
}

func checkIfDeviceIsAlreadyRegisted(token string) bool {

	session, err := mgo.Dial(server)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	result := User{}
	c := session.DB(db).C(collection)
	err = c.Find(bson.M{"devicetoken": token}).One(&result)

	if err != nil {
		//panic(err)
		return false // Wenn Collection oder Abfragefeld nicht vorhanden
	} else {
		if result.DeviceToken == "" { // TODO: if(Diese Abfrage sicher)
			return false
		} else {
			return true
		}

	}
}
