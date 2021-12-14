package apiserver

import (
	"log"
	"net/http"

	"github.com/UshakovN/speech-to-text/internal/app/recognition"
	"github.com/gorilla/mux"
)

var (
	configPath string = "./configs/recognition.yaml"
)

type APIServer struct {
	config *Config
	router *mux.Router
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	log.Println("Начало работы")

	s.configureRouter()

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/", s.handleStart())
	subrouter := s.router.PathPrefix("/recognition").Subrouter()
	subrouter.HandleFunc("/transcription", s.handleRecognition()).Methods("POST")
}

func (s *APIServer) handleStart() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Стартовая страница"))
	}
}

func (s *APIServer) handleRecognition() http.HandlerFunc {

	config := recognition.NewConfig()
	if err := config.UnmarshalYaml(configPath); err != nil {
		log.Fatal(err)
	}
	recognitor := recognition.NewYandexClient(config)

	return func(rw http.ResponseWriter, r *http.Request) {
		out := recognitor.GetTranscription(r.Body)
		rw.Write([]byte(out))
	}
}
