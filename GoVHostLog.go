package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("apache")
	q := bson.M{"$group": bson.M{"_id": "$serverName", "logCount": bson.M{"$sum": 1}, "bytesReceived": bson.M{"$sum": "$bytesReceived"}, "bytesSent": bson.M{"$sum": "$bytesSent"}}}
	d := bson.D{{"aggregate", "log"}, {"pipeline", []bson.M{q}}}

	var Stats = make(map[string]interface{})
	err = db.Run(d, &Stats)
	if err != nil {
		panic(err)
	}

	for _, v := range Stats["result"].([]interface{}) {
		receivedFix, sentFix, totalFix := "b", "b", "b"
		server := v.(map[string]interface{})["_id"]
		count := v.(map[string]interface{})["logCount"].(int)
		received := float64(v.(map[string]interface{})["bytesReceived"].(int))
		sent := float64(v.(map[string]interface{})["bytesSent"].(int))
		total := received + sent
		if received > 1024 {
			receivedFix = "Kb"
			received = received / 1024
			if received > 1024 {
				receivedFix = "Mb"
				received = received / 1024
				if received > 1024 {
					receivedFix = "Gb"
					received = received / 1024
				}
			}
		}
		if sent > 1024 {
			sentFix = "Kb"
			sent = sent / 1024
			if sent > 1024 {
				sentFix = "Mb"
				sent = sent / 1024
				if sent > 1024 {
					sentFix = "Gb"
					sent = sent / 1024
				}
			}
		}
		if total > 1024 {
			totalFix = "Kb"
			total = total / 1024
			if total > 1024 {
				totalFix = "Mb"
				total = total / 1024
				if total > 1024 {
					totalFix = "Gb"
					total = total / 1024
				}
			}
		}
		fmt.Printf("%s c: %d r: %.2f%s s: %.2f%s t: %.2f%s\n", server, count, received, receivedFix, sent, sentFix, total, totalFix)
	}
}
