package main

import (
	"fmt"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Reading struct {
	K string  "k"
	T int64   "t"
	V float64 "v"
}

var (
	mgoSession   *mgo.Session
	databaseName = "test"
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial("localhost")
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func withCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(databaseName).C(collection)
	return s(c)
}

func Hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello\n")
}

func SearchReading(q interface{}, limit int) (searchResults []Reading, searchErr string) {
	searchErr = ""
	searchResults = []Reading{}
	query := func(c *mgo.Collection) error {
		fn := c.Find(q).Limit(limit).All(&searchResults)
		if limit < 0 {
			fn = c.Find(q).All(&searchResults)
		}
		return fn
	}
	search := func() error {
		return withCollection("reading", query)
	}
	err := search()
	if err != nil {
		searchErr = "Database Error"
	}
	return
}

func GetReadingsForKey(key string, start int64, end int64, limit int) (searchResults []Reading, searchErr string) {
	if end == -1 {
		searchResults, searchErr = SearchReading(bson.M{"k": key, "t": bson.M{"$gte": start}}, limit)
	} else if start == -1 {
		searchResults, searchErr = SearchReading(bson.M{"k": key, "t": bson.M{"$lte": end}}, limit)
	} else {
		searchResults, searchErr = SearchReading(bson.M{"k": key, "t": bson.M{"$gte": start, "$lte": end}}, limit)
	}

	return
}

func ReadingCreate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var key string = req.Form["key"][0]
	mytime := time.Now().Unix()
	val, _ := strconv.ParseFloat(req.Form["val"][0], 64)

	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Panic("Could not connect to DB\n")
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("reading")

	err = c.Insert(&Reading{key, mytime, val})
	if err != nil {
		log.Panic("Could not insert item\n")
	}
	io.WriteString(w, "ok\n")
}

func ReadingRead(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var key = req.Form["key"][0]
	start, err := strconv.ParseInt(req.FormValue("start"), 10, 64)
	if err != nil {
		start = -1
	}
	end, err := strconv.ParseInt(req.FormValue("end"), 10, 64)
	if err != nil {
		end = -1
	}
	fmt.Println("Using start : " + strconv.FormatInt(start, 10) + " and end : " + strconv.FormatInt(end, 10))

	readings, errmsg := GetReadingsForKey(key, start, end, 20)
	if errmsg != "" {
		log.Panic("Could not get readings for key\n")
	}
	io.WriteString(w, "{ \n\"Key\" : \""+key+"\",\n\"Measurements\" : \n[\n")
	for i, r := range readings {
		if i > 0 {
			io.WriteString(w, ",\n")
		}
		io.WriteString(w, "{ \"Time\" : "+strconv.FormatInt(r.T, 10)+", \"Value\" : "+strconv.FormatFloat(r.V, 'f', 1, 64)+" }")
	}
	io.WriteString(w, "\n]\n}")
}

func ReadingHelp(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Service calls available :\n")
	io.WriteString(w, "  create\n")
	io.WriteString(w, "  read\n")
}

func main() {
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/reading", ReadingHelp)
	http.HandleFunc("/reading/create", ReadingCreate)
	http.HandleFunc("/reading/read", ReadingRead)
	err := http.ListenAndServe(":2200", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
