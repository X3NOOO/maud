package runner

import "github.com/X3NOOO/maud/types"

type Runner interface {
	Fire(sw types.Switch) error
}
