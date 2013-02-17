package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Id            bson.ObjectId "_id,omitempty"
	IP            string        "ip"
	Time          int           "time"
	ReqProtocol   string        "reqProtocol"
	ReqMethod     string        "reqMethod"
	QueryString   string        "queryString"
	LastRequest   int           "lastRequest"
	ReqTime       time.Time     "reqTime"
	Path          string        "path"
	ServerName    string        "serverName"
	BytesReceived int           "bytesReceived"
	BytesSent     int           "bytesSent"
	Referer       string        "Referer"
	UserAgent     string        "userAgent"
}

type ErrorEntry struct {
	Id    bson.ObjectId "_id,omitempty"
	Error string        "error"
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("apache")
	cl := db.C("log")
	ce := db.C("error")

	r := bufio.NewReader(os.Stdin)

	for {
		id := bson.NewObjectId()
		l := new(LogEntry)
		l.Id = id
		line, err := r.ReadBytes('\n')
		if err != nil {
			e := new(ErrorEntry)
			e.Id = id
			e.Error = fmt.Sprintln(err)
			ce.Insert(&e)
			continue
		}
		for i, ch := range line {
			if ch == 92 {
				line[i] = 124
			}
		}
		err = json.Unmarshal(line, &l)
		if err != nil {
			e := new(ErrorEntry)
			e.Id = id
			e.Error = fmt.Sprintln(err)
			ce.Insert(&e)
			continue
		}
		err = cl.Insert(&l)
		if err != nil {
			log.Println(err)
		}
	}
}
