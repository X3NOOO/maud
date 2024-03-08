package main

import (
	"log"
	"time"

	"github.com/X3NOOO/maud/runner"
)

func init_runners() []runner.Runner {
	return []runner.Runner{runner.Stdout{}, runner.Email{}}
}

func (ctx *maud_context) watchdog() {
	runners := init_runners()
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		switches, err := ctx.db.GetSwitchesToFire()
		if err != nil {
			log.Println(err)
		}

		for _, r := range runners {
			for _, s := range switches {
				go (func() {
					if err := r.Fire(s); err != nil {
						log.Println(err)
					}
				})()
			}
		}

		return
	}
}
