package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Record struct {
	ProjectId    int    `json:"ProjectId"`
	IssueId      string `json:"IssueId"`
	PongsCounter string `json:"PongsCounter"`
}

func getID(responseWriter http.ResponseWriter, params httprouter.Params) (int, bool) {
	id, err := strconv.Atoi(params.ByName("Id"))
	if err != nil {
		responseWriter.WriteHeader(400)
		return 0, false
	}
	return id, true
}

func getRecords(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var str string
	if len(request.URL.RawQuery) > 0 {
		str = request.URL.Query().Get("IssueId")
		if str == "" {
			responseWriter.WriteHeader(400)
			return
		}
	}
	dbRowsArray, err := readByName(str)
	if err != nil {
		responseWriter.WriteHeader(500)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(responseWriter).Encode(dbRowsArray); err != nil {
		responseWriter.WriteHeader(500)
	}
}

func getRecord(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, noError := getID(responseWriter, params)
	if !noError {
		return
	}
	rec, err := readIssuesById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			responseWriter.WriteHeader(404)
			return
		}
		responseWriter.WriteHeader(500)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(responseWriter).Encode(rec); err != nil {
		responseWriter.WriteHeader(500)
	}
}

func addRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var record Record
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		w.WriteHeader(400)
		return
	} else {
		if record.IssueId == "" {
			w.WriteHeader(401)
			w.WriteHeader(407)
			return
		} else {
			if record.PongsCounter == "" {
				w.WriteHeader(402)
				return
			}
		}
	}

	if _, err := insert(record.ProjectId, record.IssueId, record.PongsCounter); err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
}

func updateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, ok := getID(w, ps)
	if !ok {
		return
	}
	var rec Record
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil || rec.IssueId == "" || rec.PongsCounter == "" {
		w.WriteHeader(400)
		return
	}
	res, err := update(id, rec.IssueId, rec.PongsCounter)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

func deleteRecord(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, noError := getID(responseWriter, params)
	if !noError {
		return
	}
	if _, err := remove(id); err != nil {
		responseWriter.WriteHeader(500)
	}
	responseWriter.WriteHeader(204)
}
