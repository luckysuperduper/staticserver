package main

import (
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"github.com/luckysuperduper/staticserver/middleware"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	portHTTP = ":8080"
	portSSL  = ":8043"
	ssl      bool
	gzip     bool
	cache    bool
	blue     *color.Color
	certFile = `-----BEGIN CERTIFICATE-----
MIIDHzCCAgcCFAWtII3r9G7eTl9aH+4jZjPqSlbeMA0GCSqGSIb3DQEBCwUAMEwx
CzAJBgNVBAYTAlJPMRAwDgYDVQQIDAdSb21hbmlhMRUwEwYDVQQHDAxQaWF0cmEg
TmVhbXQxFDASBgNVBAoMC0RldmVsb3BtZW50MB4XDTIwMTAxMzE1MDkwNVoXDTQ4
MDIyOTE1MDkwNVowTDELMAkGA1UEBhMCUk8xEDAOBgNVBAgMB1JvbWFuaWExFTAT
BgNVBAcMDFBpYXRyYSBOZWFtdDEUMBIGA1UECgwLRGV2ZWxvcG1lbnQwggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDhrJN3ZufMLz9uVHzwYXsrZcEFtwGw
g2Uj3sbfzu6CWzv9tm8cyD4I1B6c9H5QHSFYkGyxPnFEokBIndRvhwPlQZy3zh0C
rvgQp/tQmBNpdx5+V1lxx3SqILkbIQ/Fp8qnn6RYyct/KAGGNYF00rVJuYS0HEnl
KcfhEM6vHoqurCFAMSRnkQaZ+l+FGqlMsegdkFbL2rLSFf5vAOZtddZrq2FgyiJq
vzQWwyl079r3m/KrQodxnlxOnyIBekkcUx4AcVLf2tr64oy4qVlWocztPh/Lcbpa
mpb4FlK9opPZ29LHREnG5B9ocGME0ogmJ6A+O0IHiToqvq13qNt85KJdAgMBAAEw
DQYJKoZIhvcNAQELBQADggEBAJFn8xWO+OrWsD6GGI+rBJw/xNArmR29Z4W9MTxZ
1PbVVbGJqEQKaR7e/cgNA8r6m1x5aA7JgNhYfD4pw3/XDRk9oTox01SDeV1S6HRk
TamTbWCIEAxkxO88wifWFWh7IsGf3tpAyzFCE/o82q9MB00IOc41up3NXURu7rlI
atBF3NZwKKXiKFNfTyJsN6twaL0pWnoCARlf03Yrv7jqQKM4P7B384Grh0MZs42R
tHykpqgoYl7SIDxuh1yUDBRkuW2BdNTLaVu8Qp6D1y+LrYGxiWQoom+zDfrekHd8
FaYXxwBGz9vsNy2uSyzOQhZueBg5Qs2RjEI8MG6fH3H9GyA=
-----END CERTIFICATE-----`
	keyFile = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDhrJN3ZufMLz9u
VHzwYXsrZcEFtwGwg2Uj3sbfzu6CWzv9tm8cyD4I1B6c9H5QHSFYkGyxPnFEokBI
ndRvhwPlQZy3zh0CrvgQp/tQmBNpdx5+V1lxx3SqILkbIQ/Fp8qnn6RYyct/KAGG
NYF00rVJuYS0HEnlKcfhEM6vHoqurCFAMSRnkQaZ+l+FGqlMsegdkFbL2rLSFf5v
AOZtddZrq2FgyiJqvzQWwyl079r3m/KrQodxnlxOnyIBekkcUx4AcVLf2tr64oy4
qVlWocztPh/Lcbpampb4FlK9opPZ29LHREnG5B9ocGME0ogmJ6A+O0IHiToqvq13
qNt85KJdAgMBAAECggEBALAZPql5v39xjwnFHAlnx/lBWbHf8I2QuqeW+5FBpJRM
JTAB4AqRpva0r37Cup5BXPgDGw3kL/bitU70+gRdUwjefjBfwfuKFUDKFC37vYoa
zczA1KcYgU0QY+Frlychm93ZkSFHtmfvC+FydyZ2FckF3yu8t1z/kV1rBB1as9VA
PMsVhwl7NTPE9XSqNziINwLU4EwioR1JATz8rSo5sEiyxzj5TvXCaYdGWDztqnJY
9HRQDm7G1NEjh2MealQca/eUGVstlt0MJNzuIgSVLcSXRT1EX/d0km7urvzCD2DA
CfseSHl7LC4vsTDDru6DLOvr4KmPRwPdmPHNpifmZwECgYEA/K3lKDyNKcTFTe9m
hwibUUFbvMRfiLY/3vs/msY4pxo0wKX55L5zka4YbAzBBjO+5SsVkBQ4rnmWZubm
StJMHe3J6fQ/SFCds0Q9j5ggpTq7rbBq5p2/faOpNDPGxmiHf45yBDZE4JiyU4fe
xAuGV/KIvZegVCIEDf4wZmqceNECgYEA5KPTZKKqAQAWNNaFcJBhZjU6+Tc48M1Q
VV2KGZ7BOh5yw4MBYPqj8QM1bMGO4wWQP6xTKESaD/zCSSMUPIJrcAYoslmaFLlu
gjooZPPUvcOqGmFNCxA7Ghz5xkcimtRO7NOAVndGLDPww1YpkhMM84+XyHg9ZQwX
AdUvXzKyc80CgYEAhwFClyUDJ3YDFYj79toaYmfRZCJoCNuXdMQ5T7DpRB80YFpO
EnHPvd6PHewSlgW/0SIb+0dSoaZFPeXQ1dlW4gbTAzWFOlYYbFfhrH9TsfSXok3I
UD+ouLBhD4s6gXgILZcmRCna00XCwe6uj4C43vSvKt2AxHMIR5Gwuofr4oECgYBo
tA1OfJ9VrfB9ae/ZyISSBbZoAj31KFCthxSC/wyFzQPJPOkYvC7vZATHNSx2Ekoo
noXGXwQeZiWi0Imn3CHPP0LLyfShoPlWccOl13OJI112jzB07I3kO3i2sETMmoU6
NvECp8Re4bpT+dU3q7m2n/9mMooLCCpREIuNEO5f0QKBgQDxE41okEfPWzahAF4s
CpbijIv9QNhlMlTQLGHBb2TZIoEYt6cZPazMgS3jHx1l5kDVmEmyaxai4sKBKxJg
OZGxx3reYHzKlWNY65AKSLg39ibYxPmh8w7rMWHNQ43HbnxI0drsAuAGwZ2RMQTK
JUfLSXETUsn62nemwe5G3jMEGQ==
-----END PRIVATE KEY-----`
)

func init() {
	blue = color.New(color.FgBlue, color.Bold)
}

func main() {
	serverHTTP := http.Server{
		Addr: portHTTP,
	}

	serverSSL := http.Server{
		Addr: portSSL,
	}

	// questions
	ssl = doYouWant("ssl")
	gzip = doYouWant("gzip")
	cache = doYouWant("cache")

	// ssl logic
	if ssl {
		cert, err := tls.X509KeyPair([]byte(certFile), []byte(keyFile))
		if err != nil {
			log.Fatalln(err)
		}

		serverSSL.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

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

	// cache logic
	if cache {
		http.Handle("/", middleware.Cache(http.FileServer(http.Dir("."))))
	} else {
		http.Handle("/", http.FileServer(http.Dir(".")))
	}

	if ssl {
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			fmt.Println("Navigate to: " + blue.Sprintf("https://localhost%s", portSSL))
		} else {
			fmt.Println("Navigate to: https://localhost" + portSSL)
		}
		go log.Fatalln(serverSSL.ListenAndServeTLS("", ""))
	} else {
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			fmt.Println("Navigate to: " + blue.Sprintf("http://localhost%s", portHTTP))
		} else {
			fmt.Println("Navigate to: http://localhost" + portHTTP)
		}
	}

	log.Fatalln(serverHTTP.ListenAndServe())
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
	http.Redirect(w, r, "https://localhost"+portSSL+r.RequestURI, http.StatusMovedPermanently)
}
