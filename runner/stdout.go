package runner

import (
	"log"

	"github.com/X3NOOO/maud/types"
)

type Stdout struct{}

func (s Stdout) Fire(sw types.Switch) error {
	log.Printf("%+v\n", sw)
	
	return nil
}
