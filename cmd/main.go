package main

import (
	"configurations/test"
	"log"
)

func main() {
	graylog, err := test.TestGraylog()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(graylog)

	cassandra, err := test.TestCassandra()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cassandra)

	//postgresql, err := test.TestPostgreSQL()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(postgresql)

	couchdb, err := test.TestCouchDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(couchdb)

	redis, err := test.TestRedis()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(redis)

}
