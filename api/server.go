package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ankeshnirala/go/authentication/constants"
	"github.com/ankeshnirala/go/authentication/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	logger     *log.Logger
	listenAddr string
	mongoStore storage.MongoStorage
}

func NewServer(logger *log.Logger, listenAddr string, mongoStore storage.MongoStorage) *Server {
	return &Server{
		logger:     logger,
		listenAddr: listenAddr,
		mongoStore: mongoStore,
	}
}

// type MiddlewareFunc func(http.Handler) http.Handler

func (s *Server) Start() error {

	router := mux.NewRouter()

	router.HandleFunc(constants.SIGNUP_PATH, MakeHTTPHandleFunc(s.SignupHandler))
	router.HandleFunc(constants.LOGIN_PATH, MakeHTTPHandleFunc(s.LoginHandler))

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(Authenticate)
	protected.HandleFunc(constants.LOGHISTORY_PATH, MakeHTTPHandleFunc(s.LogHistoryHandler))

	srv := &http.Server{
		Addr:         s.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Println(err)
		}
	}()

	gracefulShutdown(s.logger, srv)

	return nil
}

func gracefulShutdown(l *log.Logger, s *http.Server) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown...", sig)

	tc, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		l.Fatal(err)
	}
	s.Shutdown(tc)
}
