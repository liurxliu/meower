package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/liurxliu/meower/db"
	"github.com/liurxliu/meower/event"
	"github.com/liurxliu/meower/schema"
	"github.com/liurxliu/meower/util"
	"github.com/segmentio/ksuid"
)

func createMeowHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create meow")
		return
	}
	meow := schema.Meow{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}

	if err := db.InsertMeow(ctx, meow); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create meow")
		return
	}

	if err := event.PublishMeowCreated(meow); err != nil {
		log.Println(err)
	}

	util.ResponseOk(w, response{ID: meow.ID})
}
