package runner

import (
	"log"

	"github.com/X3NOOO/maud/types"
)

type Email struct{}

func (e Email) Fire(sw types.Switch) error {
	log.Printf("dupa\n")
	
	return nil
}
