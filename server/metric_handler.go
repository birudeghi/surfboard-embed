package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handles the POST /submit endpoint logic
//returns 201 Accepted with body {Created: bool}
//400 Bad Request will have body {Created: bool, Errortype: string, Message: string}
func submitHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	sess := globalSessions.SessionStart(w, r)

	req := ReqSchema{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || req.EventType != "submit" || req.TimeTaken == 0 {

		errorType := "invalid_request"
		message := "Request sent is invalid. Please try again with a different request."
		res := generateError(errorType, message)
		jsonResponse(w, res, http.StatusBadRequest)
		return
	}

	sess.Set("FormCompletionTime", req.TimeTaken)
	sess.Set("WebsiteUrl", req.WebsiteUrl)
	sess.Set("SessionId", req.SessionId)

	fmt.Printf("%+v\n", sess.Get("items"))

	res := MetricResSchema{true}
	jsonResponse(w, res, http.StatusAccepted)
}

//Handles POST /screen-resize endpoint logic
//returns 201 Accepted with body {Created: bool}
//400 Bad Request will have body {Created: bool, Errortype: string, Message: string}
func scrResizeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	sess := globalSessions.SessionStart(w, r)

	req := ReqSchema{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil ||
		req.EventType != "screenResize" ||
		req.ResizeFrom == (Dimension{}) ||
		req.ResizeTo == (Dimension{}) {

		errorType := "invalid_request"
		message := "Request sent is invalid. Please try again with a different request."
		res := generateError(errorType, message)
		jsonResponse(w, res, http.StatusBadRequest)
		return
	}

	resizeFrom := make(map[string]string)
	resizeFrom["width"] = req.ResizeFrom.Width
	resizeFrom["height"] = req.ResizeFrom.Height

	resizeTo := make(map[string]string)
	resizeTo["width"] = req.ResizeTo.Width
	resizeTo["height"] = req.ResizeTo.Height

	sess.Set("ResizeFrom", resizeFrom)
	sess.Set("ResizeTo", resizeTo)
	sess.Set("WebsiteUrl", req.WebsiteUrl)
	sess.Set("SessionId", req.SessionId)

	fmt.Printf("%+v\n", sess.Get("items"))

	res := MetricResSchema{true}
	jsonResponse(w, res, http.StatusAccepted)

}

//Handles POST /screen-resize endpoint logic
//returns 201 Accepted with body {Created: bool}
//400 Bad Request will have body {Created: bool, Errortype: string, Message: string}
func cPasteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	sess := globalSessions.SessionStart(w, r)

	req := ReqSchema{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil ||
		req.EventType != "copyAndPaste" ||
		req.CopyPaste == (CopyPaste{}) {

		errorType := "invalid_request"
		message := "Request sent is invalid. Please try again with a different request."
		res := generateError(errorType, message)
		jsonResponse(w, res, http.StatusBadRequest)
		return
	}

	formId := req.CopyPaste.FormId
	pasted := req.CopyPaste.Pasted

	copyPaste := map[string]bool{formId: pasted}

	sess.Set("CopyAndPaste", copyPaste)
	sess.Set("WebsiteUrl", req.WebsiteUrl)
	sess.Set("SessionId", req.SessionId)

	fmt.Printf("%+v\n", sess.Get("items"))

	res := MetricResSchema{true}
	jsonResponse(w, res, http.StatusAccepted)
}
