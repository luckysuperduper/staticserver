package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckysuperduper/staticserver/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

var (
	ssl     bool
	gzip    bool
	cache   bool
	blue    *color.Color
	certPEM = `-----BEGIN CERTIFICATE-----
MIIB2DCCATmgAwIBAgIBATAKBggqhkjOPQQDBDAWMRQwEgYDVQQKEwtEZXZlbG9w
bWVudDAeFw0yMDEwMTMxMzI1MzBaFw00ODAyMjkxMzI1MzBaMBYxFDASBgNVBAoT
C0RldmVsb3BtZW50MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAUaO86ZrWUD2y
nJu6UUOSg/e4V01Disg+NxAqeIQ7w7TFEBQhlRZK7KMt1wOqhKgyrNryUn6jFTXL
WTGwX/9H+roACqoKqn0kcjNUHBs+MHO2zhAzqDsyhNN3XbgdC68OuFKO/2ISBMTd
zxFhG9PC8S9lFCi3ToSwqbBQhe1E6kxXt0CjNTAzMA4GA1UdDwEB/wQEAwIFoDAT
BgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMAoGCCqGSM49BAMEA4GM
ADCBiAJCAPC59IDKyxiKHmWEmKoEvTh4fg2ipJYkZpP/gam/fF1183ZUIo2rNsoi
Wj12ElYYopRoBNi64fFkdZd5mOXf0jOlAkIAwqhs/siSPKIFLcCVIpBgFH2TX71v
Xbf4OiO6TBjOnKgreJK5eB/CsrU7u9tLp37xxRAKg+xDwwrU7VL2m58sabc=
-----END CERTIFICATE-----`
	certKEY = `-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIBH7pQIBNrx/j2DUbH2VZd8ffxu32t0u4YQYf35NrJYtWrNTVeTUoh
keeRkQrQzXMSOyavW6ce4Jls9L1/7/CosrmgBwYFK4EEACOhgYkDgYYABABRo7zp
mtZQPbKcm7pRQ5KD97hXTUOKyD43ECp4hDvDtMUQFCGVFkrsoy3XA6qEqDKs2vJS
fqMVNctZMbBf/0f6ugAKqgqqfSRyM1QcGz4wc7bOEDOoOzKE03dduB0Lrw64Uo7/
YhIExN3PEWEb08LxL2UUKLdOhLCpsFCF7UTqTFe3QA==
-----END EC PRIVATE KEY-----`
)

func init() {
	blue = color.New(color.FgBlue, color.Bold)
}

func main() {
	serverHTTP := http.Server{
		Addr: ":8080",
	}

	serverSSL := http.Server{
		Addr: ":8433",
	}

	// questions
	gzip = doYouWant("gzip")
	ssl = doYouWant("ssl")
	cache = doYouWant("cache")

	// gzip logic
	if ssl {
		serverHTTP.Handler = http.HandlerFunc(redirectTLS)

		if gzip {
			serverSSL.Handler = new(middleware.GzipMiddleware)
		}
	} else {
		if gzip {
			serverHTTP.Handler = new(middleware.GzipMiddleware)
		}
	}

	// ssl logic
	if ssl {
		createCertFiles()
		defer deleteCertFiles()
	}

	// cache logic
	if cache {
		http.Handle("/", middleware.Cache(http.FileServer(http.Dir("."))))
	} else {
		http.Handle("/", http.FileServer(http.Dir(".")))
	}

	fmt.Println(gzip, ssl, cache)

	if ssl {
		go func() {
			log.Fatalln(serverSSL.ListenAndServeTLS("certFile", "keyFile"))
		}()
	}
	log.Fatalln(serverHTTP.ListenAndServe())
}

func createCertFiles() {
	files := map[string]string{"certFile": certPEM, "keyFile": certKEY}
	for file, content := range files {
		f, err := os.Create(file)
		if err != nil {
			log.Fatalln(err)
		}
		f.WriteString(content)
		f.Close()
	}
}

func deleteCertFiles() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c
	os.RemoveAll("certFile")
	os.RemoveAll("keyFile")
}

func doYouWant(option string) (yes bool) {
	answer := ""

	// ask question
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		blue.Printf("Activate %s? (y/N) ", option)
	} else {
		fmt.Printf("Activate %s? (y/N) ", option)
	}

	// get answer
	_, err := fmt.Fscanln(os.Stdin, &answer)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected newline") {
			//	ignore error
		} else {
			log.Fatalln(err)
		}
	}

	if answer != "" {
		if strings.ToLower(answer) == "y" {
			yes = true
		}
	}

	return
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://localhost:8443"+r.RequestURI, http.StatusMovedPermanently)
}
