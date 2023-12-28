package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AsaHero/abclinic/internal/pkg/config"
)

func NewServer(cfg *config.Config, handler http.Handler) (*http.Server, error) {
	readTimeout, err := time.ParseDuration(cfg.Server.ReadTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while parsing server read timeout: %v", err)
	}
	writeTimeout, err := time.ParseDuration(cfg.Server.WriteTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while parsing server write timeout: %v", err)
	}
	idleTimeout, err := time.ParseDuration(cfg.Server.IdleTimeout)
	if err != nil {
		return nil, fmt.Errorf("error while parsing server idle timeout: %v", err)
	}
	return &http.Server{
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Addr:         cfg.Server.Host + cfg.Server.Port,
	}, err
}
