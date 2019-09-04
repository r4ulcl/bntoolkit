package dht

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/etix/goscrape"
)

//ScrapeTrackers scrape common trackers for the hashes
func ScrapeTrackers(db *sql.DB, debug bool, verbose bool, projectName string) {
	//https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_udp.txt
	var trackers [34]string
	trackers[0] = "udp://tracker.coppersurfer.tk:6969/announce"
	trackers[1] = "udp://tracker.internetwarriors.net:1337/announce"
	trackers[2] = "udp://tracker.opentrackr.org:1337/announce"
	trackers[3] = "udp://9.rarbg.to:2710/announce"
	trackers[4] = "udp://exodus.desync.com:6969/announce"
	trackers[5] = "udp://tracker1.itzmx.com:8080/announce"
	trackers[6] = "udp://explodie.org:6969/announce"
	trackers[7] = "udp://ipv4.tracker.harry.lu:80/announce"
	trackers[8] = "udp://open.demonii.si:1337/announce"
	trackers[9] = "udp://denis.stalker.upeer.me:6969/announce"
	trackers[10] = "udp://thetracker.org:80/announce"
	trackers[11] = "udp://bt.xxx-tracker.com:2710/announce"
	trackers[12] = "udp://tracker.torrent.eu.org:451/announce"
	trackers[13] = "udp://tracker.tiny-vps.com:6969/announce"
	trackers[14] = "udp://tracker.port443.xyz:6969/announce"
	trackers[15] = "udp://retracker.lanta-net.ru:2710/announce"
	trackers[16] = "udp://tracker.uw0.xyz:6969/announce"
	trackers[17] = "udp://tracker.iamhansen.xyz:2000/announce"
	trackers[18] = "udp://open.stealth.si:80/announce"
	trackers[19] = "udp://zephir.monocul.us:6969/announce"
	trackers[20] = "udp://tracker.vanitycore.co:6969/announce"
	trackers[21] = "udp://tracker.cyberia.is:6969/announce"
	trackers[22] = "udp://tracker4.itzmx.com:2710/announce"
	trackers[23] = "udp://tracker1.wasabii.com.tw:6969/announce"
	trackers[24] = "udp://tracker.swateam.org.uk:2710/announce"
	trackers[25] = "udp://tracker.kamigami.org:2710/announce"
	trackers[26] = "udp://tracker.filepit.to:6969/announce"
	trackers[27] = "udp://tracker.dler.org:6969/announce"
	trackers[28] = "udp://torrentclub.tech:6969/announce"
	trackers[29] = "udp://pubt.in:2710/announce"
	trackers[30] = "udp://bittracker.ru:6969/announce"
	trackers[31] = "udp://amigacity.xyz:6969/announce"
	trackers[32] = "udp://tracker.justseed.it:1337/announce"
	trackers[33] = "udp://packages.crunchbangplusplus.org:6969/announce"
	infohash, err := utils.GetHashes(db, debug, verbose)
	if err != nil {
		log.Fatal(err)
	}
	for _, tracker := range trackers {
		waitGroup.Add(1)
		go scrapeTracker(db, debug, verbose, tracker, infohash)
	}
	waitGroup.Wait()
}

func scrapeTracker(db *sql.DB, debug bool, verbose bool, tracker string, infohash [][]byte) {
	defer func() {
		waitGroup.Done()
	}()

	// A list of infohash to scrape, at most 74 infohash can be scraped at once.
	// Be sure to provide infohash that are 40 hexadecimal characters long only.

	// Create a new instance of the library and specify the torrent tracker to use.
	s, err := goscrape.New(tracker)
	if err != nil {
		if debug {
			log.Printf("Error: %v", err)
		}
	}
	s.SetTimeout(time.Second * 5)
	s.SetRetryLimit(2)

	// Connect to the tracker and scrape the list of infohash in only two UDP round trips.
	res, err := s.Scrape(infohash...)
	if err != nil {
		if debug {
			log.Printf("Error: %v", err)
		}
	}

	// Loop over the results and print them.
	// Result are guaranteed to be in the same order they were requested.
	for _, r := range res {
		if debug {
			log.Println("scrape tracker infohash:\t", string(r.Infohash))
		}
		if r.Seeders > 0 || r.Leechers > 0 || r.Completed > 0 {
			sum := r.Seeders + r.Leechers + r.Completed
			if sum > 0 {
				utils.SetTrueValid(db, debug, verbose, string(r.Infohash))
				utils.SetLen(db, debug, verbose, int(sum), string(r.Infohash))

			}
			if debug {
				fmt.Println("Infohash:\t", string(r.Infohash))
				fmt.Println("Seeders:\t", r.Seeders)
				fmt.Println("Leechers:\t", r.Leechers)
				fmt.Println("Completed:\t", r.Completed)
				fmt.Println("")
			}
		}
	}
}

var data chan *utils.Node
var hashes chan string
var max = 100000

//Crawler DHT
func Crawler(db *sql.DB, debug bool, verbose bool, threads int) {
	if debug {
		log.Println("Starting the application...")
	}

	max = threads * 1000
	data = make(chan *utils.Node, max)
	hashes = make(chan string, max)

	for len(data) == 0 {
		nodes, err := findNodes(debug, verbose)
		if err != nil {
			if debug {
				log.Println(err)
			}
		}
		if debug {
			log.Println("NODES")
			log.Println(len(nodes))
			log.Println(len(hashes))
		}
		for _, v := range nodes {
			data <- v
		}
		fmt.Println("Can't find any node")
	}

	if verbose {
		fmt.Println("Start scrapers ")
	}
	for i := 0; i < threads; i++ {
		waitGroup.Add(1)
		go workerCrawler(debug, verbose)
	}

	//start inserts
	if verbose {
		fmt.Println("Start inserts ")
	}
	for i := 0; i < threads/250; i++ { //pq: deadlock detected
		waitGroup.Add(1)
		go workerInsert(db, debug, verbose)
	}

	waitGroup.Wait()
	close(data)
	close(hashes)

}

func workerCrawler(debug bool, verbose bool) {
	if debug {
		log.Println("Goroutine worker is now starting...")
	}
	defer func() {
		if debug {
			log.Println("Destroying the worker...")
		}
		waitGroup.Done()
	}()
	for {

		value, ok := <-data
		if !ok {
			if verbose {
				log.Println("ERROR: The channel is closed!")
			}
			break
		}

		if debug {
			log.Printf("LEN data, %d LEN hashes %d", len(data), len(hashes))
		}
		nodes, hashList, err := sampleInfohashes(debug, verbose, value)

		if err != nil {
			log.Println(err)
		}
		for _, v := range nodes {
			if len(data) < max/10*9 {
				data <- v
			}
		}

		for _, h := range hashList {
			if len(data) < max/10*9 {
				hashes <- h
			}
		}

		if len(data) == 0 {
			//if the number of threads is too hight can be empty, so it waits
			time.Sleep(30 * time.Second)
			if len(data) == 0 {
				break
			}
		}
	}
}

func workerInsert(db *sql.DB, debug bool, verbose bool) {
	if verbose {
		log.Println("Goroutine workerInsert is now starting...")
	}
	defer func() {
		if verbose {
			log.Println("Destroying the worker...")
		}
		waitGroup.Done()
	}()

	for {
		sql := `INSERT INTO hash(hash ,source, first_seen)
VALUES `
		max := len(hashes)

		for i := 0; i < max; i++ {
			value, ok := <-hashes
			if !ok {
				if debug {
					log.Println("The channel is closed!")
				}
				break
			}
			if !strings.Contains(sql, value) {
				if i == 0 {
					sql += "\n ('" + value + "'," + " 'dht' , " + "current_timestamp)"
				} else {
					sql += ",\n('" + value + "'," + " 'dht' , " + "current_timestamp)"
				}
			}
			if debug {
				log.Println(value)
			}
		}
		sql += " ON CONFLICT DO NOTHING" //UPDATE set last_seen = '" + timeAux + "'::timestamp; "
		if max > 0 {
			//fmt.Println(sql)
			//fmt.Println(max)
			utils.ExecuteDb(db, debug, verbose, sql)
		} else {
			time.Sleep(time.Second)
		}
	}
}
