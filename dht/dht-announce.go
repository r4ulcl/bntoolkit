package dht

import (
	"database/sql"
	"encoding/hex"
	"log"
	"sync"
	"time"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/anacrolix/dht"
)

var waitVar sync.WaitGroup

//DaemonPeers monitor hashes in monitor table
func DaemonPeers(cfgFile string, debug bool, verbose bool, projectName string) {
	db, err := utils.ConnectDb(cfgFile, debug, verbose)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//for {
	infohashes, err := utils.GetMonitor(db, debug, verbose, projectName)
	if err != nil {
		log.Println(err)
		log.Println("Not hashes in the database to monitor, sleep 1 minute")
		time.Sleep(time.Minute)
	}
	for _, hash := range infohashes { //monitor all hashes
		println(hash)
		go GetPeersLib(db, debug, verbose, hash, projectName)
		waitVar.Add(1)
	}

	waitVar.Wait()
	//}

}

//GetPeersLib gets peers from hash and insert into the database, (code based in https://raw.githubusercontent.com/anacrolix/dht/master/cmd/dht-announce/main.go)
func GetPeersLib(db *sql.DB, debug bool, verbose bool, hash string, projectName string) {
	defer waitVar.Done()

	//fmt.Println("starting get peers lib")

	var Infohash [][20]byte

	h, err := hex.DecodeString(hash)

	var aux [20]byte
	copy(aux[:], h)

	Infohash = [][20]byte{aux}

	//fmt.Println(Infohash)

	s, err := dht.NewServer(nil)
	if err != nil {
		log.Fatalf("error creating server: %s", err)
	}
	defer s.Close()
	addrs := make(map[[20]byte]map[string]struct{}, len(Infohash))
	ih := Infohash[0]
	utils.InsertProject(db, debug, verbose, projectName)

	for {
		a, err := s.Announce(ih, 0, true)
		if err != nil {
			log.Printf("error announcing %s: %s", ih, err)
		}
		addrs[ih] = make(map[string]struct{})
		for ps := range a.Peers {
			for _, p := range ps.Peers {
				s := p.String()
				if _, ok := addrs[ih][s]; !ok {
					if debug {
						log.Printf("got peer %s for %x from %s", p, ih, ps.NodeInfo)
					}
					ip := p.IP.String()
					port := p.Port
					utils.InsertIP(db, debug, verbose, ip, projectName)
					utils.InsertDownload(db, debug, verbose, ip, port, hash, projectName)
					addrs[ih][s] = struct{}{}
				}
			}
		}
		time.Sleep(time.Minute) //sleep 1 min before ask again
	}
	/*
		for _, ih := range Infohash {
			a, err := s.Announce(ih, 0, true)
			if err != nil {
				log.Printf("error announcing %s: %s", ih, err)
				continue
			}
			wg.Add(1)
			addrs[ih] = make(map[string]struct{})
			go func() {
				defer wg.Done()
				for ps := range a.Peers {
					for _, p := range ps.Peers {
						s := p.String()
						if _, ok := addrs[ih][s]; !ok {
							if debug {
								log.Printf("got peer %s for %x from %s", p, ih, ps.NodeInfo)
							}
							ip := p.IP.String()
							port := p.Port
							utils.InsertIP(db, debug, verbose, ip, projectName)
							utils.InsertDownload(db, debug, verbose, ip, port, hash)
							addrs[ih][s] = struct{}{}
						}
					}
				}
			}()
		}
	*/

}
