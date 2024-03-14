package runners

import (
	"fmt"
	"os"

	"github.com/X3NOOO/maud/types"
)

type Logging struct{
	File *os.File
}

func (l Logging) Fire(sw types.Switch) error {
	_, err := fmt.Fprintf(l.File, "%+v\n", sw)
	
	return err
}
