package main

import (
	"log"
	"os"
	"time"

	"github.com/X3NOOO/maud/runners"
)

func (ctx *maud_context) init_runners() []runners.Runner {
	logfile, err := os.OpenFile(ctx.config.Runners.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open logfile")
	}
	return []runners.Runner{
		runners.Logging{
			File: logfile,
		},

		// runners.Email{
		// Host:     ctx.config.Runners.Email.Host,
		// Port:     ctx.config.Runners.Email.Port,
		// Email:    ctx.config.Runners.Email.Email,
		// Password: ctx.config.Runners.Email.Password,
		// },
	}
}

func (ctx *maud_context) watchdog() {
	runners := ctx.init_runners()
	ticker := time.NewTicker(24 * time.Hour)
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
	}
}
