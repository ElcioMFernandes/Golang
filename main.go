package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
    }
    
    apiKey := os.Getenv("API_KEY")
    secretKey := os.Getenv("SECRET_KEY")
    
    if apiKey == "" || secretKey == "" {
        log.Fatal("API_KEY ou SECRET_KEY não definidos no arquivo .env")
    }
    
    baseUrl := "https://api.binance.com"
    endpoint := "/api/v3/account"
    
    timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
    queryString := "timestamp=" + timestamp
    
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString))
    signature := hex.EncodeToString(h.Sum(nil))
    
    fullUrl := baseUrl + endpoint + "?" + queryString + "&signature=" + signature
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", fullUrl, nil)
    if err != nil {
        log.Fatalf("Erro ao criar requisição: %v", err)
    }
    
    req.Header.Add("X-MBX-APIKEY", apiKey)
    
    res, err := client.Do(req)
    if err != nil {
        fmt.Println("Erro na requisição:", err)
        return
    }
    defer res.Body.Close()
    
    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println("Erro ao ler resposta:", err)
        return
    }
    
    fmt.Println(string(body))
}