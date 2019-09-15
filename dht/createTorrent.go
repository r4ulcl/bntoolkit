package dht

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/anacrolix/tagflag"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

var waitGroup sync.WaitGroup

//FileInfo Path and Lenght of a file
type FileInfo struct {
	Path   []string
	Length int
}

var hashesChan chan string
var inicio = 10

//WorkersTorrents create a hashlist for the file
func WorkersTorrents(cfgFile string, debug bool, verbose bool, file string, projectName string) {
	fmt.Println("file", file)

	max := 1024
	hashesChan = make(chan string, max)

	db, err := utils.ConnectDb(cfgFile, debug, verbose)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	utils.DeletePossibles(db, debug, verbose)
	for i := inicio; i < 60; i++ {
		waitGroup.Add(1)
		go worker(db, debug, verbose, i, file, file, false, projectName)
	}

	waitGroup.Wait()

}

//worker goroutine
func worker(db *sql.DB, debug bool, verbose bool, id int, path string, file string, reverse bool, projectName string) {

	if debug {
		log.Println("Goroutine worker is now starting...")
	}
	defer func() {
		if debug {
			log.Printf("Destroying the worker with id %v...\n", id)
		}
		waitGroup.Done()
	}()

	var (
		builtinAnnounceList = [][]string{
			{"udp://tracker.openbittorrent.com:80"},
			{"udp://tracker.publicbt.com:80"},
			{"udp://tracker.istole.it:6969"},
		}
	)

	log.SetFlags(log.Flags() | log.Lshortfile)
	var args struct {
		AnnounceList []string `name:"a" help:"extra announce-list tier entry"`
		tagflag.StartPos
	}
	//tagflag.Parse(&args, tagflag.Description("Creates a torrent metainfo for the file system rooted at ROOT, and outputs it to stdout."))
	mi := metainfo.MetaInfo{
		AnnounceList: builtinAnnounceList,
	}
	for _, a := range args.AnnounceList {
		mi.AnnounceList = append(mi.AnnounceList, []string{a})
	}
	mi.SetDefaults()
	info := metainfo.Info{
		PieceLength: int64(math.Exp2(float64(id))),
		//PieceLength: int64(16 * 1024 * id),
	}

	if debug {
		log.Println("2^" + strconv.Itoa(id) + ":    " + strconv.FormatInt(info.PieceLength, 10))
	}
	//err := info.BuildFromFilePath(file, reverse)
	err := info.BuildFromFilePath(file)
	if err != nil {
		log.Println(err)
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Println(err)
	}
	hashpossible := mi.HashInfoBytes().HexString()

	if debug {
		log.Println(hashpossible)
	}
	//hashesChan <- hashpossible

	exist, err := utils.CheckExist(db, debug, verbose, hashpossible)
	if err != nil {
		log.Println(err)
	}
	err = utils.InsertPossible(db, debug, verbose, id-inicio, hashpossible, exist, projectName)
	if err != nil {
		log.Println(err)
	}
	output := fmt.Sprintf("/tmp/dat2%d", id)

	f, err := os.Create(output)
	err = mi.Write(f)
	if err != nil {
		log.Println(err)
	}

}

//CrateTorrent create a torrent file
func CrateTorrent(debug bool, verbose bool, file string, outfile string, piecesize int64, tracker string, comment string) {

	if debug {
		log.Println("Goroutine worker is now starting...")
	}
	defer func() {
		if debug {
			log.Printf("Destroying the worker with id..\n")
		}
	}()

	var (
		builtinAnnounceList = [][]string{
			{tracker},
		}
	)

	log.SetFlags(log.Flags() | log.Lshortfile)
	var args struct {
		AnnounceList []string `name:"a" help:"extra announce-list tier entry"`
		tagflag.StartPos
	}

	//tagflag.Parse(&args, tagflag.Description(comment))

	mi := metainfo.MetaInfo{
		AnnounceList: builtinAnnounceList,
	}
	for _, a := range args.AnnounceList {
		mi.AnnounceList = append(mi.AnnounceList, []string{a})
	}

	mi.SetDefaults()
	info := metainfo.Info{
		PieceLength: int64(piecesize),
		//PieceLength: int64(16 * 1024 * id),
	}

	//err := info.BuildFromFilePath(file, false)
	fmt.Println(file)
	err := info.BuildFromFilePath(file)
	if err != nil {
		log.Println(err)
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Println(err)
	}
	hashpossible := mi.HashInfoBytes().HexString()

	if debug {
		log.Println(hashpossible)
	}
	//hashesChan <- hashpossible
	fmt.Println(outfile + ".torrent")
	output := fmt.Sprintf(outfile + ".torrent")

	f, err := os.Create(output)
	err = mi.Write(f)
	if err != nil {
		log.Fatal(err)
	}
}
