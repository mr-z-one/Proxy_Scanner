package Server

import (
	"fmt"
	"net/http"
	"os"
	"proxyScanner/Server/routes"
	"proxyScanner/db"
	"proxyScanner/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func StartServer(port int) {
	database := db.GetActiveDatabaseSession()
	if database == nil {
		fmt.Println("[-] failed to connect to database")
		return
	}
	router := mux.NewRouter().StrictSlash(true)
	router.Use(JSONContentTypeMiddleware)
	routes.RegisterRequestRoutes(router)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentPath, _ := os.Getwd()

		fmt.Fprintf(w, utils.ReadeFile(currentPath+"/Server/"+"html/home.html"))

	})

	// Configure server
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("[+] Api Server started on http://localhost:%d\n", port)
	log.Fatal(server.ListenAndServe())

}
func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isJson := strings.Contains(r.URL.Path, "api")
		if isJson {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		} else {
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		}

		next.ServeHTTP(w, r)
	})
}
