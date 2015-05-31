package main

import (
	"fmt"
	"encoding/json"
	"github.com/m4rw3r/uuid"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)


func members(c web.C, w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(getMembers())
}

func member(c web.C, w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	member, err := getMember(uuid)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(member)
}


func main() {
	fmt.Printf("%+v\n", getMembers())
	goji.Get("/members", members)
	goji.Get("/member/:uuid", member)
	goji.Serve()
}
