package dht

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	bencode "github.com/jackpal/bencode-go"

	"github.com/r4ulcl/bntoolkit/utils"
)

// struct para el find node
type nodeRQ struct {
	y string
	t string
	q string
	a auxNRQ
}

type auxNRQ struct {
	id     string
	target string
}

type nodeRR struct {
	T string
	Y string
	R auxNRR
}

type auxNRR struct {
	Samples string
	Values  string
	Nodes   string
}

/*
"router.utorrent.com:6881",
		"router.bittorrent.com:6881",
		"dht.transmissionbt.com:6881",
		"dht.aelitis.com:6881",     // Vuze
		"router.silotis.us:6881",   // IPv6
		"dht.libtorrent.org:25401", // @arvidn
*/

// http://www.bittorrent.org/beps/bep_0051.html
func sampleInfohashes(debug bool, verbose bool, n *utils.Node) ([]*utils.Node, []string, error) {
	//random id and target
	var hashList []string
	addr := n.Ip + ":" + strconv.Itoa(int(n.Port)) //juntamos ip:port
	if debug {
		log.Printf("sampleInfohashes %s", addr)
	}
	//id := nodeRQ{"q", "aa", "sample_infohashes", auxNRQ{id: "kjhgrtyhgfr45tyhgft6", target: "agr45yu763efrgthyji8"}}
	id := nodeRQ{"q", "aa", "sample_infohashes", auxNRQ{id: string(utils.RandomID()), target: string(utils.RandomID())}}

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, id)
	if err != nil && debug {
		log.Printf("could not marshal: %v\n", err)
	}

	text := buf.String()

	/////////////////////

	p := make([]byte, 2048)
	timeoutDuration := time.Second // / 2

	//conn, err := net.Dial("udp", "router.bittorrent.com:6881")
	conn, err := net.DialTimeout("udp", addr, timeoutDuration)
	if err != nil {
		log.Println("Some error", err)
		return nil, nil, err
	}
	fmt.Fprintf(conn, text)

	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil && debug {
		log.Printf("Some error %v\n", err)

	}

	if debug {
		log.Println(string(p))
	}

	conn.Close()

	////////////////

	var i nodeRR
	r := strings.NewReader(string(p))

	err = bencode.Unmarshal(r, &i)
	if err != nil && debug {
		log.Printf("Error Sample infohashes, could not unmarshal: %v\n", err)
	}

	encodedStr := hex.EncodeToString([]byte(i.R.Samples))

	aux := "" //use var buffer bytes.Buffer
	contador := 0
	for _, r := range encodedStr {
		c := string(r)
		aux += c
		contador++
		if contador == 40 {
			hashList = append(hashList, aux)
			aux = ""
			contador = 0
		}
	}

	nodes, err := utils.DecodeNodes(i.R.Nodes)
	if err != nil {
		log.Printf("Error decoding nodes, could not unmarshal: %v\n", err)
	}
	if debug {
		for node := range nodes {
			log.Printf("IP: %s", nodes[node].Ip)
		}
	}

	return nodes, hashList, nil
}
