package pkg

import (
	"github.com/borkod/kindly/config"
)

// Kindly struct stores kindly config
type Kindly struct {
	cfg config.Config
}

// SetConfig sets the kindly struct config
func (k *Kindly) SetConfig(c config.Config) {
	k.cfg = c
}
