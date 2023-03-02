package dht

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	bencode "github.com/jackpal/bencode-go"
	"github.com/r4ulcl/bntoolkit/utils"
)

// struct for find node
type fNode struct {
	t string
	y string
	q string
	a auxFN
}

type auxFN struct {
	id     string
	target string
}

type fNodeR struct {
	T string `json:"t"`
	Y string `json:"y"`
	Q string `json:"q"`
	R auxFNR `json:"r"`
}

type auxFNR struct {
	ID     string
	Values string
	Nodes  string
}

/*
"router.utorrent.com:6881",
		"router.bittorrent.com:6881",
		"dht.transmissionbt.com:6881",
		"dht.aelitis.com:6881",     // Vuze
		"router.silotis.us:6881",   // IPv6
		"dht.libtorrent.org:25401", // @arvidn
*/

// First function to search in DHT
func findNodes(debug bool, verbose bool) ([]*utils.Node, error) {
	nodesList := [7]string{"router.bittorrent.com:6881",
		"router.utorrent.com:6881",
		"dht.transmissionbt.com:6881",
		"dht.aelitis.com:6881",
		"router.silotis.us:6881",
		"dht.libtorrent.org:6881",
		"dht.libtorrent.org:25401"}

	var nodes []*utils.Node
	var err error
	for _, node := range nodesList {
		nodesAux, err := findNode(debug, verbose, node)
		if err != nil && debug {
			log.Printf("Error: %v\n", err)
		}
		nodes = append(nodes, nodesAux...)

	}

	return nodes, err
}

func findNode(debug bool, verbose bool, node string) ([]*utils.Node, error) {
	//find_node Query = {"t":"aa", "y":"q", "q":"find_node", "a": {"id":"abcdefghij0123456789", "target":"mnopqrstuvwxyz123456"}}
	//bencoded = d1:ad2:id20:abcdefghij01234567896:target20:mnopqrstuvwxyz123456e1:q9:find_node1:t2:aa1:y1:qe
	//Response = {"t":"aa", "y":"r", "r": {"id":"0123456789abcdefghij", "nodes": "def456..."}}
	//bencoded = d1:rd2:id20:0123456789abcdefghij5:nodes9:def456...e1:t2:aa1:y1:re
	//BytesToString(RandomID())
	/*
		addr := ip + ":" + strconv.Itoa(int(port))
		if debug {
			log.Printf("Find node %s", addr)
		}
	*/
	//id := fNode{"aa", "q", "find_node", auxFN{id: "abcdefghij0123456789", target: "mnopqrstuvwxyz123456"}}
	id := fNode{"aa", "q", "find_node", auxFN{id: string(utils.RandomID()), target: string(utils.RandomID())}}
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, id); err != nil {
		return nil, fmt.Errorf("could not marshal: %v", err)
	}
	timeoutDuration := 5 * time.Second
	conn, err := net.DialTimeout("udp", node, timeoutDuration)
	if err != nil {
		return nil, fmt.Errorf("could not dial %s:  %v", node, err)
	}

	text := buf.String()
	fmt.Fprintf(conn, text)
	p := make([]byte, 2048)

	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var i fNodeR
	r := strings.NewReader(string(p))
	if err := bencode.Unmarshal(r, &i); err != nil {
		return nil, fmt.Errorf("Find nodes error. Could not unmarshal: %v", err)
	}

	nodes, err := utils.DecodeNodes(i.R.Nodes)
	if err != nil {
		return nil, fmt.Errorf("could not decode nodes: %v", err)
	}

	return nodes, nil
}
