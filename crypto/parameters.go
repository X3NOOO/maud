package crypto

import "runtime"

const (
	SALT_SIZE uint32 = 16

	ARGON_TIME uint32 = 1
	ARGON_MEMORY uint32 = 64*1024
	ARGON_KEY_LEN uint32 = 32
)

var ARGON_THREADS uint8 = uint8(runtime.NumCPU())