package main

import (
	"NaiveDistributedMemCache/server"
	"flag"
	"log"
)

var db = map[string]string{
	"Tom":  "123",
	"Jack": "456",
	"Sam":  "789",
}

func createGroups(self string, peers map[string]string) map[string]*server.Group {
	group := server.NewGroup("score", 3,
		func(key string) []byte {
			log.Println("Search key in db", key)
			if value, ok := db[key]; ok {
				return []byte(value)
			}
			return nil
		})
	group.RegisterPeers(self, peers)
	groups := map[string]*server.Group{
		"score": group,
	}
	return groups
}


func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Http server listen")
	flag.Parse()
	peers := map[string]string{
		"8081": "http://localhost:8081",
		"8082": "http://localhost:8082",
		"8083": "http://localhost:8083",
	}
	groups := createGroups(port, peers)
	server.StartServer(peers[port][7:], groups)
}
