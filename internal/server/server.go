package server

import (
	"fmt"
	"log"
	"net/http"
	"pocketerClient/internal/config"
	"pocketerClient/internal/repository"
	"pocketerClient/pkg/pocket"
	"strconv"
)

type authorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	cfg             *config.Config
	redirectURL     string
}

func NewAuthorizationServer(
	pc *pocket.Client,
	tr repository.TokenRepository,
	cfg *config.Config,
) *authorizationServer {
	return &authorizationServer{
		pocketClient:    pc,
		tokenRepository: tr,
		cfg:             cfg,
		redirectURL:     cfg.TelegramBotURL,
	}
}

func (s *authorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.AuthServerPort),
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *authorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdParam := r.URL.Query().Get("chat_id")
	if chatIdParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatId, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	authResponse, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.tokenRepository.Put(chatId, authResponse.AccessToken, repository.AccessTokens)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf(
		"Succesful authorization with params:\nchat_id=%d\nrequest_token=%s\naccess_token=%s",
		chatId,
		requestToken,
		authResponse.AccessToken,
	)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
