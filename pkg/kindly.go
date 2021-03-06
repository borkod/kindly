package pkg

import (
	"log"

	"github.com/borkod/kindly/config"
)

// Kindly struct stores kindly config
type Kindly struct {
	cfg    config.Config
	logger *log.Logger
}

// SetConfig sets the kindly struct config
func (k *Kindly) SetConfig(c config.Config) {
	k.cfg = c
}

// SetLogger sets the kindly struct logger
func (k *Kindly) SetLogger(l *log.Logger) {
	k.logger = l
}
