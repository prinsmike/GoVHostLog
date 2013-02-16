package main

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
)

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
		var l = make(map[string]interface{})
		l["_id"] = id
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
