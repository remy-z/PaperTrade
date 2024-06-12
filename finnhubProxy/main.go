package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/remy-z/PaperTrade/finnhubProxy/api"
)

var (
	port         string
	finnhubToken string
)

func main() {
	SetENV()

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", finnhubToken)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	server := api.NewAPIServer("0.0.0.0:"+port, finnhubClient)
	server.Run()
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

	// assign the environment variables using the os.Getenv method
	port = os.Getenv("PORT")
	finnhubToken = os.Getenv("FINNHUB_TOKEN")
}
