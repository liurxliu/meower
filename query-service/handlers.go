package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/liurxliu/meower/db"
	"github.com/liurxliu/meower/event"
	"github.com/liurxliu/meower/schema"
	"github.com/liurxliu/meower/search"
	"github.com/liurxliu/meower/util"
)

func onMeowCreated(m event.MeowCreateMessage) {
	meow := schema.Meow{
		ID:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := search.InsertMeow(context.Background(), meow); err != nil {
		log.Println(err)
	}
}

func searchMeowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")

	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	meows, err := search.SearchMeows(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []schema.Meow{})
		return
	}
	util.ResponseOk(w, meows)
}

func listMeowsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}

	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	meows, err := db.ListMeows(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch meows")
		return
	}
	util.ResponseOk(w, meows)
}
