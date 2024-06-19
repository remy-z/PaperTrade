package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/remy-z/PaperTrade/finnhubProxy/api"
)

var (
	port         string
	finnhubToken string
	wilshire     map[string]int
)

func main() {
	setENV()
	loadFiles()
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", finnhubToken)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	server := api.NewAPIServer("0.0.0.0:"+port, finnhubClient, wilshire)
	server.Run()
}

func setENV() {
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

func loadFiles() {
	if _, err := os.Stat("./wilshire.csv"); err == nil {
		log.Printf("File %s exists in the current directory\n", "./wilshire.csv")
		envFile, err := os.Open("./wilshire.csv")
		if err != nil {
			log.Fatalln(err)
		}

		defer envFile.Close()
		wilshire = make(map[string]int)
		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			symbol := strings.Split(scanner.Text(), ",")
			marketCap, _ := strconv.Atoi(symbol[1])
			wilshire[symbol[0]] = marketCap
		}
		fmt.Println(wilshire)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else if os.IsNotExist(err) {
		log.Printf("File %s does not exist in the current directory\n", "./wilshire.csv")
	} else {
		log.Println("Error occurred while checking the file:", "./wilshire.csv")
	}
}
