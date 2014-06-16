package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/slogsdon/docker-build-service/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	VERSION         = "0.0.1"
	DEFAULT_COMMAND = "serve"
)

var (
	portFlag string
	hostFlag string
)

func init() {
	const (
		defaultPort = "8000"
		defaultHost = "0.0.0.0"
	)

	flag.StringVar(&portFlag, "port", defaultPort, "")
	flag.StringVar(&portFlag, "p", defaultPort, "")
	flag.StringVar(&hostFlag, "host", defaultHost, "")
	flag.StringVar(&hostFlag, "h", defaultHost, "")
}

func main() {
	flag.Parse()

	command := DEFAULT_COMMAND

	args := flag.Args()
	if len(args) == 1 {
		command = args[0]
	}

	switch command {
	case "help":
		showHelp()
	case "serve":
		serve()
	case "version":
		fmt.Println("docker-build-service version", VERSION)
	default:
		showHelp()
	}
}

func showHelp() {
	fmt.Printf(`NAME:
   docker-build-service - Provides a build service for code using docker containers

USAGE:
   docker-build-service [options] command

VERSION:
   %s

COMMANDS:
  serve
    Start the server

  help
    Show this help message.

  version
    Show version

OPTIONS:
	--port, -p
    Set port number to use. 
    Port can also be set via PORT environment variable.

  --host, -h
    Set host number to use. 
    Host can also be set via HOST environment variable.
`, VERSION)
}

func serve() {
	handleSigInt()

	http.HandleFunc("/compile", func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w)
		handlers.Compile(w, r)
	})

	log.Printf("Listening on %v. Press Ctrl-C to exit", portFlag)
	log.Printf("Go to http://%v/", hostFlag+":"+portFlag)
	log.Fatal(http.ListenAndServe(hostFlag+":"+portFlag, nil))
}

func handlePanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		json, err := json.Marshal(map[string]interface{}{"error": r})
		if err == nil {
			fmt.Fprint(w, string(json))
		}
	}
}

func handleSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Print("Stopping the service...")
		os.Exit(0)
	}()
}
