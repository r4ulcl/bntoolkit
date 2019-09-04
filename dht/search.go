package dht

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
)

//SearchDHT searchs in DHT
func SearchDHT(db *sql.DB, debug bool, verbose bool, projectName string) {
	log.Println("Starting the application...")
	//max := 1000
	//data := make(chan *utils.Node, max)

	posibilidades, err := utils.GetPossibles(db, debug, verbose)
	if err != nil {
		log.Println("Error")
	}
	if debug {
		log.Println(posibilidades)
	}
	for i := 0; i < posibilidades; i++ {
		waitGroup.Add(1)
		go workerSearch(db, debug, verbose, i)
	}

	/*
		nodes, err := findNodes()
		if err != nil && verbose  {
			log.Println("Error")
		}
		log.Println(len(nodes))
		for _, v := range nodes {
			data <- v
		}
	*/

	waitGroup.Wait()
}

func workerSearch(db *sql.DB, debug bool, verbose bool, num int) {
	max := 1000
	dataSearch := make(chan *utils.Node, max)

	nodes, err := findNodes(debug, verbose)
	if err != nil {
		log.Printf("Error %v", err)
	}
	if debug {
		log.Println(len(nodes))
	}
	for _, v := range nodes {
		dataSearch <- v
	}

	if debug {
		log.Printf("Goroutine worker %v is now starting...", num)
	}
	defer func() {
		if debug {
			log.Println("Destroying the worker...")
		}
		waitGroup.Done()
	}()

	hash, err := utils.GetHash(db, debug, verbose, num)
	if err != nil {
		log.Printf("Error %v", err)
	} else {
		for start := time.Now(); time.Since(start) < (time.Minute * 5); {
			if len(dataSearch) == 0 {
				break
			}
			value, ok := <-dataSearch
			if !ok {
				log.Println("The channel is closed!")
			}

			nodes, _, err := getPeers(db, debug, verbose, value, strings.ToLower(hash))

			if err != nil {
				if debug {
					log.Printf("Error %v", err)
				}
			}

			if debug {
				log.Printf("len %v", len(dataSearch))
			}
			if len(dataSearch) < 900 {
				for _, v := range nodes {
					dataSearch <- v
				}
			}
		}
	}
}
