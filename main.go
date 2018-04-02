package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"net/http"

	"github.com/cyverse-de/configurate"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/cyverse-de/messaging"
)

var log = logrus.WithFields(logrus.Fields{
	"service": "job-status-listener",
	"art-id":  "job-status-listener",
	"group":   "org.cyverse",
})

var (
	cfgPath = flag.String("config", "", "Path to the configuration file.")
	cfg     *viper.Viper

	client *messaging.Client
)

// JobUpdatePublisher is the interface for types that need to publish a job
// update.
type JobUpdatePublisher interface {
	PublishJobUpdate(m *messaging.UpdateMessage) error
}

func running(client JobUpdatePublisher, job *messaging.JobDetails, hostname string, msg string) (*messaging.UpdateMessage, error) {
	updateMessage := &messaging.UpdateMessage{
		Job:     *job,
		State:   messaging.RunningState,
		Message: msg,
		Sender:  hostname,
	}

	err := client.PublishJobUpdate(updateMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info(msg)
	return updateMessage, nil
}

type MessagePost struct {
	Hostname string
	Message  string
	Job      *messaging.JobDetails
}

func postRunning(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	out := json.NewEncoder(w)

	var updateMessage MessagePost

	err := json.NewDecoder(r.Body).Decode(&updateMessage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		out.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	msg, err := running(client, updateMessage.Job, updateMessage.Hostname, updateMessage.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		out.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	out.Encode(msg)
}

func init() {
	flag.Parse()
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func loadConfig(cfgPath string) {
	var err error
	cfg, err = configurate.Init(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/debug/vars", http.DefaultServeMux)
	r.Path("/running").Methods("POST").HandlerFunc(postRunning)

	return r
}

func main() {
	log.Info("Starting up the job-status-listener service.")
	loadConfig(*cfgPath)

	uri := cfg.GetString("amqp.uri")
	exchange := cfg.GetString("amqp.exchange.name")
	var err error
	client, err = messaging.NewClient(uri, true)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	client.SetupPublishing(exchange)

	r := newRouter()

	listenPortSpec := ":" + "60000"
	log.Infof("Listening on %s", listenPortSpec)
	log.Fatal(http.ListenAndServe(listenPortSpec, r))
}
