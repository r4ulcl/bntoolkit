package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/anacrolix/torrent"
)

func torrentLib(db *sql.DB, timeout int, debug bool, verbose bool, hash string) {

	c, _ := torrent.NewClient(nil)
	defer func() {
		c.Close()
		waitGroup.Done()
	}()
	t, _ := c.AddMagnet("magnet:?xt=urn:btih:" + hash)

	metadata := t.Info()
	for start := time.Now(); time.Since(start) < (time.Minute * time.Duration(timeout)); {
		metadata = t.Info()
		if metadata != nil {
			//mi := t.Metainfo()
			fmt.Println("Downloading", t.Name())
			t.Drop()
			/*f, err := os.Create(t.Info().Name + ".torrent")
			if err != nil && verbose  {
				log.Fatalf("error creating torrent metainfo file: %s", err)
			}
			defer f.Close()
			err = bencode.NewEncoder(f).Encode(mi)
			if err != nil && verbose  {
				log.Fatalf("error writing torrent metainfo file: %s", err)
			}*/
			SetTrueDownload(db, debug, verbose, hash)
			SetNamePossibles(db, debug, verbose, t.Name(), hash)

			log.Println(hash)
			break
		} else {
			if debug {
				fmt.Println("Esperamos")
			}
			time.Sleep(time.Second * 30)
		}
	}

	//fmt.Printf("====================================================== El hash %v con name %v\n", hash, metadata.Name)
	//swarm := t.KnownSwarm()
	//for _, s := range swarm {
	//	fmt.Println(s)
	//}
}
