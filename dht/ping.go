package dht

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"

	bencode "github.com/jackpal/bencode-go"
	"github.com/r4ulcl/bntoolkit/utils"
)

const letterBytes = "abcdef1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func randStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

// struct para el ping request
type pingQ struct {
	t string
	y string
	q string
	a aux
}

// struct para el auxiliar ping request
type aux struct {
	id string
}

// struct para el ping response
type pingR struct {
	T string `json:"t"`
	Y string `json:"y"`
	R auxR
}

// struct aux para el ping response
type auxR struct {
	ID string
}

/*
https://play.golang.org/p/TL7sRRDiUXP
*/
//Funcion para hacer ping a un host de la red
func ping(debug bool, verbose bool, n *utils.Node) {
	//id := identity{42, "Jack", "Why are you ignoring me?", "Daniel"}

	addr := n.Ip + ":" + strconv.Itoa(int(n.Port)) //juntamos ip:port
	log.Printf("Ping %s", addr)

	id := pingQ{"aa", "q", "ping", aux{id: "abcdefghij0123456789"}}
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, id)
	if err != nil {
		log.Printf("could not marshal: %v\n", err)
	}
	log.Printf("Result: %s\n", buf.String())

	text := buf.String()

	/////////////////////

	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "router.bittorrent.com:6881")
	if err != nil {
		log.Print("Some error", err)
		return
	}
	fmt.Fprintf(conn, text)
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		log.Printf("Some error %v\n", err)
	}
	log.Println(string(p))
	conn.Close()

	////////////////

	var i pingR
	r := strings.NewReader(string(p))
	err = bencode.Unmarshal(r, &i)
	if err != nil {
		log.Printf("could not unmarshal: %v\n", err)
	}
	log.Printf("Result: %#+v\n", i)

}
