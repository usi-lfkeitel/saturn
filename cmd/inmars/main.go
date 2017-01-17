package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	listenAddress string
	modulesPath   string
	authKey       string
)

func init() {
	flag.StringVar(&listenAddress, "listen", "0.0.0.0:80", "Where the server listens for connections. [interface]:port")
	flag.StringVar(&modulesPath, "modules", "./modules", "Location of modules.")
	flag.StringVar(&authKey, "authkey", "", "Client authentication key")
	flag.Parse()
}

func main() {
	if authKey == "" {
		fmt.Println("WARNING: Client running with no authentication key!")
	}

	http.HandleFunc("/module/", request)

	fmt.Println("Starting http server at:", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		fmt.Println("Error starting http server:", err)
		os.Exit(1)
	}
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func checkAuthHeader(r *http.Request) bool {
	if authKey == "" {
		fmt.Println("WARNING: Client running with no authentication key!")
		return true
	}

	// First check for an Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) >= 16 {
		scheme := authHeader[:10]
		key := authHeader[15:]
		if key[0] == '"' {
			key = key[1 : len(key)-1]
		}
		fmt.Printf("%s:%s\n", scheme, key)

		if scheme != "INCA-TOKEN" {
			return false
		}

		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			fmt.Println(err)
			return false
		}

		fmt.Println(string(decoded))
		return string(decoded) == authKey
	}

	// Next check for a URL parameter
	urlParam := r.URL.Query().Get("key")
	return urlParam != "" && urlParam == authKey
}

func request(w http.ResponseWriter, r *http.Request) {
	if !checkAuthHeader(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	url := strings.Split(r.URL.Path, "/")
	if len(url) != 3 {
		http.Error(w, "No module specified.", http.StatusNotAcceptable)
		return
	}

	module := url[2]

	// Execute the command
	cmdPath := fmt.Sprintf("%s/%s.sh", modulesPath, module)
	fmt.Println(cmdPath)
	if !fileExists(cmdPath) {
		http.Error(w, "Requested module doesn't exist.", http.StatusNotAcceptable)
		return
	}

	var output bytes.Buffer
	cmd := exec.Command(cmdPath)
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing '%s': %s\n\tScript output: %s\n", module, err.Error(), output.String())
		http.Error(w, "Unable to execute module.", http.StatusInternalServerError)
		return
	}

	out := struct {
		Error   int         `json:"error"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Data: output.String(),
	}

	jsonOut, err := json.Marshal(out)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonOut)
}
