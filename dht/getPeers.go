package dht

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	bencode "github.com/jackpal/bencode-go"
)

//struct para el find node
type getPeerQ struct {
	t string
	y string
	q string
	a auxGetPeerQ
}

type auxGetPeerQ struct {
	id        string
	info_hash []byte
}

type getPeerR struct {
	T string
	Y string
	R auxGetPeerR
}

type auxGetPeerR struct {
	ID     string
	Token  string
	Values []string
	Nodes  string
}

func getPeers(db *sql.DB, debug bool, verbose bool, n *utils.Node, hash string) ([]*utils.Node, bool, error) {
	//get_peers Query = {"t":"aa", "y":"q", "q":"get_peers", "a": {"id":"abcdefghij0123456789", "info_hash":"mnopqrstuvwxyz123456"}}
	//bencoded = d1:ad2:id20:abcdefghij01234567899:info_hash20:mnopqrstuvwxyz123456e1:q9:get_peers1:t2:aa1:y1:qe

	//Response with peers = {"t":"aa", "y":"r", "r": {"id":"abcdefghij0123456789", "token":"aoeusnth", "values": ["axje.u", "idhtnm"]}}
	//bencoded = d1:rd2:id20:abcdefghij01234567895:token8:aoeusnth6:valuesl6:axje.u6:idhtnmee1:t2:aa1:y1:re

	//Response with closest nodes = {"t":"aa", "y":"r", "r": {"id":"abcdefghij0123456789", "token":"aoeusnth", "nodes": "def456..."}}
	//bencoded = d1:rd2:id20:abcdefghij01234567895:nodes9:def456...5:token8:aoeusnthe1:t2:aa1:y1:re

	addr := n.Ip + ":" + strconv.Itoa(int(n.Port)) //juntamos ip:port
	if debug {
		log.Printf("Find node %s", addr) //"\xe8B\x13\xa7\x94\xf3\xcc\xd8\x908*T\xa6L\xa6\x8b~\x92T3"
	}
	h, err := hex.DecodeString(hash)
	id := getPeerQ{"aa", "q", "get_peers", auxGetPeerQ{id: "abcdefghij0123456789", info_hash: h}}
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, id); err != nil {
		return nil, false, fmt.Errorf("could not marshal: %v", err)
	}

	conn, err := net.DialTimeout("udp", addr, time.Second*2)
	if err != nil {
		return nil, false, fmt.Errorf("could not dial %s: %v", addr, err)
	}

	text := buf.String()
	fmt.Fprintf(conn, text)
	p := make([]byte, 2048)

	timeoutDuration := time.Second * 2
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, false, fmt.Errorf("could not read: %v", err)
	}
	defer conn.Close()

	var i getPeerR
	r := strings.NewReader(string(p))
	if err := bencode.Unmarshal(r, &i); err != nil {
		if debug {
			log.Printf("could not unmarshal: %v", err)
		}
	}

	nodes, err := utils.DecodeNodes(i.R.Nodes)
	if err != nil {
		return nil, false, fmt.Errorf("could not decode nodes: %v", err)
	}

	if len(i.R.Values) != 0 {
		if debug {
			log.Println(len(i.R.Values))
		}
		nodes, err := utils.Decodepeer(i.R.Values, hash)
		if err != nil {
			return nil, false, fmt.Errorf("could not decode peers: %v", err)
		}
		if len(nodes) > 0 {
			if debug {
				log.Printf("%v     The hash can be good:                        %v\n", addr, hash)
			}
			leng := getPeersLen(debug, verbose, n, hash) //le vuelvo a preguntar por el mismo
			if leng > 0 {
				num := getPeersLen(debug, verbose, n, "bd2a12f3c8bdfab0979506495bfedf26295e777d") //por uno falso
				if num == 0 {
					num2 := getPeersLen(debug, verbose, n, "aaaaaaf3c8bdfab0979506495bfedf22295e277d") //por otro falso
					if num2 == 0 {
						num3 := getPeersLen(debug, verbose, n, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb") //por otro falso mas
						if num3 == 0 {
							utils.SetTrueValid(db, debug, verbose, hash)
							utils.SetLen(db, debug, verbose, len(nodes), hash)
							//if len(nodes) > 10 {
							//setTrueValid(hash)
							//}
						}
					}
				}
			}
		}
		/////////////////////////////////////////////////////////////////////////////////////////////////
		//return nodes, false, nil
	}

	return nodes, false, nil
}

func getPeersLen(debug bool, verbose bool, n *utils.Node, hash string) int {
	addr := n.Ip + ":" + strconv.Itoa(int(n.Port)) //juntamos ip:port

	//Si este no existe nos aseguramos que no devuelve a todo una IP

	h, err := hex.DecodeString(hash)

	id := getPeerQ{"aa", "q", "get_peers", auxGetPeerQ{id: "abcdefghij0123456789", info_hash: h}}
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, id); err != nil {
		return -1
	}

	conn, err := net.DialTimeout("udp", addr, time.Second*2)
	if err != nil {
		return -1
	}

	text := buf.String()
	fmt.Fprintf(conn, text)
	p := make([]byte, 2048)

	timeoutDuration := time.Second * 2
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return -1
	}
	defer conn.Close()

	var i getPeerR
	r := strings.NewReader(string(p))
	if err := bencode.Unmarshal(r, &i); err != nil {
		log.Printf("could not unmarshal: %v", err)
	}

	nodes, err := utils.Decodepeer(i.R.Values, hash)
	if err != nil {
		return -1
	}
	if 1 == 1 {

		if len(nodes) == 0 {
			log.Printf("%v     The hash is good:                        %v\n", addr, hash)
		} else {
			log.Printf("%v     The hash is NOT good:                        %v\n", addr, hash)

		}
	}
	return len(nodes)

}

//SearchHash search a hash and insert peers in DB
func SearchHash(db *sql.DB, debug bool, verbose bool, hash string) {
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
				log.Println("len", len(dataSearch))
			}
			if len(dataSearch) < 900 {
				for _, v := range nodes {
					dataSearch <- v
				}
			}
		}
	}
}
