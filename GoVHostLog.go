package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"time"
)

type Record struct {
	Server   string
	PostFix  string
	Received float64
	Sent     float64
	Total    float64
}

type Records struct {
	Recs []Record
}

func main() {

	var t = [2]time.Time{}
	tpat := "20060102+1504"
	var err error

	switch {
	case len(os.Args) == 1:
		t[0], err = time.Parse(tpat, tpat)
		if err != nil {
			panic(err)
		}
		t[1] = time.Now()
	case len(os.Args) == 2:
		t[0], err = time.Parse(tpat, os.Args[1])
		if err != nil {
			panic(err)
		}
		t[1] = time.Now()
	case len(os.Args) >= 3:
		t[0], err = time.Parse(tpat, os.Args[1])
		if err != nil {
			panic(err)
		}
		t[1], err = time.Parse(tpat, os.Args[2])
		if err != nil {
			panic(err)
		}
	}

	Stats, err := db(t[0], t[1])
	if err != nil {
		panic(err)
	}

	for _, v := range Stats["result"].([]interface{}) {
		server := v.(map[string]interface{})["_id"]
		count := float64(v.(map[string]interface{})["logCount"].(int))
		received := float64(v.(map[string]interface{})["bytesReceived"].(int))
		fmt.Printf("%T\n", (v.(map[string]interface{})["bytesSent"]))
		sent := float64(v.(map[string]interface{})["bytesSent"].(int))
		total := received + sent

		received, receivedFix := scale(received)
		sent, sentFix := scale(sent)
		total, totalFix := scale(total)

		rpad := padding(received)
		spad := padding(sent)
		tpad := padding(total)
		cpad := padding(count)

		totalStr := fmt.Sprintf("t: %s%.2f%s\t", tpad, total, totalFix)
		receivedStr := fmt.Sprintf("r: %s%.2f%s\t", rpad, received, receivedFix)
		sentStr := fmt.Sprintf("s: %s%.2f%s\t", spad, sent, sentFix)
		countStr := fmt.Sprintf("c: %s%.0f\t", cpad, count)

		fmt.Printf("%s%s%s%s%s\n", totalStr, receivedStr, sentStr, countStr, server)
	}
}

func db(startDate, endDate time.Time) (map[string]interface{}, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("apache")
	var q = []bson.M{
		bson.M{"$match": bson.M{"reqTime": bson.M{"$gt": startDate, "$lt": endDate}}},
		bson.M{"$group": bson.M{"_id": "$serverName", "logCount": bson.M{"$sum": 1}, "bytesReceived": bson.M{"$sum": "$bytesReceived"}, "bytesSent": bson.M{"$sum": "$bytesSent"}}},
		bson.M{"$sort": bson.M{"logCount": 1}},
	}
	d := bson.D{{"aggregate", "log"}, {"pipeline", q}}

	var Stats = make(map[string]interface{})
	err = db.Run(d, &Stats)
	if err != nil {
		return nil, err
	}
	return Stats, nil
}

func scale(val float64) (sval float64, pfix string) {
	pfix = " b"
	sval = val
	switch {
	case sval > 1099511627776:
		pfix = "Tb"
		sval = sval / 1099511627776
		return
	case sval > 1073741824:
		pfix = "Gb"
		sval = sval / 1073741824
		return
	case sval > 1048576:
		pfix = "Mb"
		sval = sval / 1048576
		return
	case sval > 1024:
		pfix = "Kb"
		sval = sval / 1024
		return
	}
	return
}

func padding(val float64) (padding string) {
	switch {
	case val < 10:
		padding = "      "
		return
	case val < 100:
		padding = "     "
		return
	case val < 1000:
		padding = "    "
		return
	case val < 10000:
		padding = "   "
		return
	case val < 100000:
		padding = "  "
		return
	case val < 1000000:
		padding = " "
		return
	}
	return
}
