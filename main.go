package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/mailgun/log"
)

var (
	config = Config{}
)

func die(err error) {
	log.Errorf("%v", err)
	os.Exit(1)
}

func parseConfig(configtoml *string) error {
	file, err := ioutil.ReadFile(*configtoml)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}
	err = toml.Unmarshal(file, &config)
	if err != nil {
		return fmt.Errorf("Problem parsing config: %v", err)
	}

	return nil
}

func marathonCallback(w http.ResponseWriter, r *http.Request) {
	// log.Infof("callback received from Marathon")
	handleMarathonEvent(r)
	w.WriteHeader(202)
	fmt.Fprintln(w, "Done")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Init([]*log.LogConfig{&log.LogConfig{Name: "console"}})

	configtoml := flag.String("f", "kameni.toml", "Path to config. (default kameni.toml)")
	flag.Parse()

	if err := parseConfig(configtoml); err != nil {
		die(err)
	}

	setupEtcd()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/marathon_callback", marathonCallback)
	router.HandleFunc("/moxy_callback", marathonCallback)

	log.Infof("Running on: %s", config.ListenAddr())

	if err := http.ListenAndServe(config.ListenAddr(), router); err != nil {
		die(fmt.Errorf("Error starting server: %v", err))
	}
}
