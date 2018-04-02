package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	"strings"

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

func update(client JobUpdatePublisher, state messaging.JobState, jobId string, hostname string, msg string) (*messaging.UpdateMessage, error) {
	updateMessage := &messaging.UpdateMessage{
		Job:     messaging.JobDetails{InvocationID: jobId},
		State:   state,
		Message: msg,
		Sender:  hostname,
	}

	err := client.PublishJobUpdate(updateMessage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("%s (%s) [%s]: %s", jobId, state, hostname, msg)
	return updateMessage, nil
}

type MessagePost struct {
	Hostname string
	Message  string
	State    string
}

func getState(state string) (messaging.JobState, error) {
	switch strings.ToLower(state) {
	case "submitted":
		return messaging.SubmittedState, nil
	case "running":
		return messaging.RunningState, nil
	case "completed":
		return messaging.SucceededState, nil
	case "succeeded":
		return messaging.SucceededState, nil
	case "failed":
		return messaging.FailedState, nil
	default:
		return "", fmt.Errorf("Unknown job state: %s", state)
	}
}

func postUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	out := json.NewEncoder(w)

	var updateMessage MessagePost

	vars := mux.Vars(r)
	jobId := vars["uuid"]

	err := json.NewDecoder(r.Body).Decode(&updateMessage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		out.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	state, err := getState(updateMessage.State)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		out.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	msg, err := update(client, state, jobId, updateMessage.Hostname, updateMessage.Message)
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
	r.Path("/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/status").Methods("POST").HandlerFunc(postUpdate)

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
