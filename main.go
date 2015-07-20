package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mapper struct {
	Id    bson.ObjectId `bson:"_id"`
	Key   string        `bson:"key"`
	Value int           `bson:"value"`
}

func Favhandler(w http.ResponseWriter, r *http.Request) {
}

func about(w http.ResponseWriter, r *http.Request) {
	log.Println("In about")
	log.Println("PORT: ", os.Getenv("PORT"))
	log.Println("MONGO_URL: ", os.Getenv("MONGO_URL"))
	fmt.Fprintln(w, "About Page: This app is a example go app deployed in docker that connects to mongodb.")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("In handler")
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}
	log.Println("Connecting to mongo: ", uri)
	sess, err := mgo.Dial(uri)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("simple").C("visitors")

	// Query One
	result := Mapper{}
	err = collection.Find(bson.M{"key": "VisitorCount"}).One(&result)
	if err != nil {
		doc := Mapper{Id: bson.NewObjectId(), Key: "VisitorCount", Value: 1}
		err = collection.Insert(doc)
	} else {
		// Update
		colQuerier := bson.M{"key": "VisitorCount"}
		change := bson.M{"$set": bson.M{"value": result.Value + 1}}
		err = collection.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
	}
	fmt.Fprintf(w, "%d visitors have come to this page.", result.Value)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT was not available in the ENV")
		os.Exit(1)
	}
	port = ":" + port
	log.Println("Port: ", port)

	http.HandleFunc("/favicon.ico", Favhandler)
	http.HandleFunc("/hits", handler)
	http.HandleFunc("/about", about)
	http.HandleFunc("/", about)
	http.ListenAndServe(port, nil)
}
