package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/birudeghi/surfboard-embed/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi/v5"
)

func appInfoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO add session management to all the apps
	// sess := globalSessions.SessionStart(w, r)

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	// Add struct type validation

	appId := chi.URLParam(r, "appId")

	fmt.Println(appId)

	if appId == "" {

		message := "Missing or invalid appId."
		res := ErrorSchema{}
		res.Errors = append(res.Errors, Error{"invalid request", message})
		jsonResponse(w, res, http.StatusBadRequest)
		return
	}

	client, ctx, cancel, err := db.Connect(uri)

	if err != nil {

		log.Fatal(err)
		message := "If problems persist, please contact Surfboard."
		res := ErrorSchema{}
		res.Errors = append(res.Errors, Error{"internal server error", message})
		jsonResponse(w, res, http.StatusInternalServerError)
		return
	}

	defer db.Close(client, ctx, cancel)

	objID, err := primitive.ObjectIDFromHex(appId)

	fmt.Println(objID)

	if err != nil {
		panic(err)
	}

	var filter, field interface{}

	filter = bson.M{
		"_id": objID,
	}

	field = bson.M{
		"compatibleBoards": 1,
		"libs":             1,
		"files":            1,
	}

	result := db.QueryOne(client, ctx, database, collection, filter, field)

	if err != nil {
		log.Fatal(err)
		message := "Missing or invalid appId."
		res := ErrorSchema{}
		res.Errors = append(res.Errors, Error{"invalid request", message})
		jsonResponse(w, res, http.StatusBadRequest)
		return
	}

	res := SurfSchema{}

	result.Decode(&res) // TODO Current: null error message

	jsonResponse(w, res, http.StatusOK)

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// TODO add session management to all the apps
	// sess := globalSessions.SessionStart(w, r)

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	req := SurfSchema{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		log.Fatal(err)
		message := "If problems persist, please contact Surfboard."
		res := ErrorSchema{}
		res.Errors = append(res.Errors, Error{"internal server error", message})
		jsonResponse(w, res, http.StatusInternalServerError)
		return
	}

	client, ctx, cancel, err := db.Connect(uri)

	if err != nil {

		log.Fatal(err)
		message := "If problems persist, please contact Surfboard."
		res := ErrorSchema{}
		res.Errors = append(res.Errors, Error{"internal server error", message})
		jsonResponse(w, res, http.StatusInternalServerError)
		return
	}

	defer db.Close(client, ctx, cancel)

	result, err := db.InsertOne(client, ctx, database, collection, req)

	if err != nil {
		log.Fatal(err)
		message := "If problems persist, please contact Surfboard."
		res := ErrorSchema{}
		res.Errors = append(res.Errors, Error{"invalid request", message})
		jsonResponse(w, res, http.StatusBadRequest)
		return
	}

	res := req

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {

		res.AppId = id.Hex()
	}

	jsonResponse(w, res, http.StatusOK)

}
