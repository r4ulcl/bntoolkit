package utils

import (
	"encoding/binary"
	"errors"
	"log"
	"math/rand"
	"net"
)

//Node &net.TCPAddr{IP: from.IP, Port: int(port)},
type Node struct {
	Ip   string
	Port uint16
}

//DecodeNodes from DHT query
func DecodeNodes(s string) ([]*Node, error) {
	if len(s)%26 != 0 {
		return nil, errors.New("invalid length")
	}
	var Nodes []*Node
	for i := 0; i < len(s); i += 26 {
		//id := s[i : i+20]
		ip := net.IP([]byte(s[i+20 : i+24])).String()
		port := binary.BigEndian.Uint16([]byte(s[i+24 : i+26]))
		//addr := ip + ":" + strconv.Itoa(int(port))
		//Nodes = append(Nodes, &Node{id: id, addr: addr})
		Nodes = append(Nodes, &Node{Ip: ip, Port: port})
	}

	return Nodes, nil
}

//Decodepeer from DHT query
func Decodepeer(s []string, hash string) ([]*Node, error) {
	var Nodes []*Node

	for _, element := range s {
		ip := net.IP([]byte(element[:4])).String()
		port := binary.BigEndian.Uint16([]byte(element[4:6]))
		Nodes = append(Nodes, &Node{Ip: ip, Port: port})
		//quitado por el import de verbose //addr := ip + ":" + strconv.Itoa(int(port)) //juntamos ip:port
		//if verbose == 1 {
		log.Printf("hash %s con Node %s:%d\n", hash, ip, port)
		//}
	}

	return Nodes, nil
}

//RandomID to DHT
func RandomID() []byte {
	id := make([]byte, 20)
	rand.Read(id[:])
	return id
}

//BytesToString return string
func BytesToString(data []byte) string {
	return string(data[:])
}
