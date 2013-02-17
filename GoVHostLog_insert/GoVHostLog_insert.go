package main

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Id            bson.ObjectId "_id,omitempty"
	IP            string        "ip"
	LIP           string        "lip"
	RespSize      int           "respSize"
	Time          int           "time"
	FileName      string        "filename"
	ReqProtocol   string        "reqProtocol"
	KeepAlive     string        "keepalive"
	ReqMethod     string        "reqMethod"
	Port          int           "port"
	ProcessID     int           "processId"
	QueryString   string        "queryString"
	OrigRequest   int           "origRequest"
	LastRequest   int           "lastRequest"
	ReqTime       time.Time     "reqTime"
	Path          string        "path"
	ServerName    string        "serverName"
	ConnStatus    string        "connStatus"
	BytesReceived int           "bytesReceived"
	BytesSent     int           "bytesSent"
	Referer       string        "Referer"
	UserAgent     string        "userAgent"
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("apache").C("log")

	dec := json.NewDecoder(os.Stdin)

	for {
		id := bson.NewObjectId()
		l := new(LogEntry)
		l.Id = id
		if err := dec.Decode(&l); err != nil {
			log.Println(err)
			return
		}
		err = c.Insert(&l)
		if err != nil {
			log.Println(err)
		}
	}
}
