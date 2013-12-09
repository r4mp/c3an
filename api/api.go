package api

import (
	"github.com/r4mp/c3an/core"
	"log"
	"net/http"
	"strconv"
)

type Result struct {
	Result interface{} `json:",omitempty"`
	Error  string      `json:",omitempty"`
	Time   string
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {

	log.Println("Register device...")

	r.ParseForm()
	core.RegisterDevice(r.PostFormValue("token"))
}

func UnregisterDevice(w http.ResponseWriter, r *http.Request) {

	log.Println("Unregister device...")

	r.ParseForm()
	core.UnregisterDevice(r.PostFormValue("token"))
}

func getDeviceList(w http.ResponseWriter, r *http.Request) {

}

func SendNotification(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	iBadge, err := strconv.Atoi(r.PostFormValue("badge"))

	if err != nil {
		iBadge = 0
	}

	var jsonResult []byte
	var e error

	if r.PostFormValue("token") == "" {
		/*jsonResult, e = */ core.SendNotificationToAllRegisteredDevices(r.PostFormValue("message"), iBadge, r.PostFormValue("sound"))
	} else {
		jsonResult, e = core.SendNotificationToSingleDevice(r.PostFormValue("message"), iBadge, r.PostFormValue("sound"), r.PostFormValue("token"))
	}

	if e == nil {
		w.Write(jsonResult)
	} else {
		w.Write([]byte("{ Error: \"" + e.Error() + "\"}"))
	}
}
