package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	port         string
	finnhubToken string
)

func main() {
	setupServer()
}

func setupServer() {
	SetENV()
	finhubClient := NewFinnHubClient(finnhubToken)

	manager := NewManager(finhubClient)

	http.Handle("/", http.FileServer(http.Dir("../frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func SetENV() {
	if _, err := os.Stat("./.env"); err == nil {
		log.Printf("File %s exists in the current directory\n", "./.env")
		envFile, err := os.Open("./.env")
		if err != nil {
			log.Fatalln(err)
		}

		defer envFile.Close()

		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			envVar := strings.Split(scanner.Text(), "=")
			os.Setenv(envVar[0], envVar[1])
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else if os.IsNotExist(err) {
		log.Printf("File %s does not exist in the current directory\n", "./.env")
	} else {
		log.Println("Error occurred while checking the file:", "./.env")
	}

	port = os.Getenv("PORT")
	finnhubToken = os.Getenv("FINNHUB_TOKEN")
}
