package apiserver

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/UshakovN/speech-to-text/internal/app/recognition"
	"github.com/UshakovN/speech-to-text/internal/app/store"
	"github.com/gorilla/mux"
)

var (
	configRecPath   string = "./configs/recognition.yaml"
	configStorePath string = "./configs/store.yaml"
)

type APIServer struct {
	config *Config
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	log.Println("Server start")

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/", s.handleStart())
	subrouter := s.router.PathPrefix("/recognition").Subrouter()
	subrouter.HandleFunc("/transcription", s.handleRecognition()).Methods("POST")
	subrouter.HandleFunc("/upload", s.handleUpload())
}

func (s *APIServer) configureStore() error {

	if err := s.config.Store.UnmarshalYaml(configStorePath); err != nil {
		return err
	}
	store := store.New(s.config.Store)

	if err := store.Open(); err != nil {
		return err
	}

	s.store = store
	log.Println("Database start")

	return nil
}

func (s *APIServer) handleUpload() http.HandlerFunc {
	htmlFile, err := os.Open("./internal/app/apiserver/template.html")
	if err != nil {
		log.Fatal(err)
	}
	bodyTempl, err := io.ReadAll(htmlFile)
	if err != nil {
		log.Fatal(err)
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write(bodyTempl)
	}
}

func (s *APIServer) handleStart() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Стартовая страница"))
	}
}

func (s *APIServer) handleRecognition() http.HandlerFunc {

	configRec := recognition.NewConfig()
	if err := configRec.UnmarshalYaml(configRecPath); err != nil {
		log.Fatal(err)
	}
	recognitor := recognition.NewYandexClient(configRec)

	return func(rw http.ResponseWriter, r *http.Request) {

		file := r.Body

		/*
			file, _, err := r.FormFile("speech")
			if err != nil {
				log.Fatal(err)
			}
		*/

		out := recognitor.GetTranscription(file)
		rw.Write([]byte(out))
	}
}
