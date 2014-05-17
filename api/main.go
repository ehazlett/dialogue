package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/dialogue/auth"
	"github.com/ehazlett/dialogue/db"
)

var (
	listenAddress    string
	rethinkDbAddress string
	rethinkDbName    string
	enableDebug      bool
	sessionKey       string
	log              = logrus.New()
)

func init() {
	flag.StringVar(&listenAddress, "l", ":3000", "Listen address (i.e. 127.0.0.1:3000)")
	flag.StringVar(&rethinkDbAddress, "rethink-address", "127.0.0.1:28015", "RethinkDB Address")
	flag.StringVar(&rethinkDbName, "rethink-name", "dialogue", "RethinkDB Name")
	flag.BoolVar(&enableDebug, "debug", false, "Enable debug logging")
	flag.StringVar(&sessionKey, "session-key", "dialogue-key", "Secret Session Key")
}

func main() {
	flag.Parse()
	log.Info("Dialogue API")
	var (
		sig = make(chan os.Signal)
	)
	signal.Notify(sig, os.Interrupt)

	// init db
	db, err := db.NewRethinkdbSession(rethinkDbAddress, rethinkDbName)
	if err != nil {
		log.Fatalf("Unable to initialize database: %s", err)
	}

	// init auth
	auth := auth.NewAuthenticator(0) // allow specifying cost?

	// launch api
	api, err := NewApi(listenAddress, db, auth, sessionKey)
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
