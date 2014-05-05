package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
)

var (
	listenAddress    string
	listenPort       int
	rethinkDbAddress string
	rethinkDbName    string
	enableDebug      bool
	log              = logrus.New()
)

func init() {
	flag.StringVar(&listenAddress, "l", "", "Listen address")
	flag.IntVar(&listenPort, "p", 3000, "Listen port")
	flag.StringVar(&rethinkDbAddress, "rethink-address", "127.0.0.1:28015", "RethinkDB Address")
	flag.StringVar(&rethinkDbName, "rethink-name", "dialogue", "RethinkDB Name")
	flag.BoolVar(&enableDebug, "debug", false, "Enable debug logging")
}

func main() {
	flag.Parse()
	log.Info("Dialogue API")
	var (
		sig = make(chan os.Signal)
	)
	signal.Notify(sig, os.Interrupt)

	// init db
	db, err := NewRethinkdbSession(rethinkDbAddress, rethinkDbName)
	if err != nil {
		log.Fatalf("Unable to initialize database: %s", err)
	}
	// launch api
	api, err := NewApi(listenAddress, listenPort, db)
	if err != nil {
		log.Fatal("Unable to spawn API server")
	}
	go api.Run()

	// watch for shutdown
	for {
		select {
		case <-sig:
			log.Info("Shutting down Dialogue API")
			return
		}
	}
}
