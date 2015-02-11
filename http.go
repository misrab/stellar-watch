package main 

import (
	"log"
	"time"
	
	"net/http"
	"encoding/json"

	"github.com/gorilla/schema"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
    // "github.com/coopernurse/gorp"
    "github.com/go-gorp/gorp"
)


var (
	// Recommended as package variable 
	// on http://www.gorillatoolkit.org/pkg/schema
	decoder = schema.NewDecoder()
	// HTTP header defaults
	HEADER_DEFAULTS = map[string]string {
		"Access-Control-Allow-Origin": "*",
		"Content-Type": "application/json",
	}
)


func listenToRequests(dbmap *gorp.DbMap) {
	router := mux.NewRouter()
    
	router.HandleFunc("/amounts", func(res http.ResponseWriter, req *http.Request) {
		// date parameters
		vars := mux.Vars(req)
		from := vars["from"]
		to := vars["to"]
		var amounts []Amounts

		// return all
		if from == "" || to == "" {
    		_, err := dbmap.Select(&amounts, "select * from amounts order by id")
    		Respond(amounts, err, res)
    		return
		}

		// otherwise get range
		_, err := dbmap.Select(&amounts, "select * from amounts order by id where created > $1 and created < $2", from, to)
    	Respond(amounts, err, res)
    	return
    }).Methods("GET")


    http.Handle("/", router)
    log.Println("Listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}





func SetHeaders(res http.ResponseWriter, code int) {
  	for k, v := range HEADER_DEFAULTS {
  		res.Header().Set(k, v)
  	}
  	res.Header().Set("Status", http.StatusText(code))
	res.Header().Set("Date", time.Now().String())	
}


func Respond(i interface{}, err error, res http.ResponseWriter) {
	if err != nil {
		SetHeaders(res, 400)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err2 := json.Marshal(i)
	if err2 != nil {
		SetHeaders(res, 400)
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(res, 200)
	res.Write(js)
}