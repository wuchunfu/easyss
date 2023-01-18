//go:build !with_notray

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/nange/easyss"
	log "github.com/sirupsen/logrus"
)

func Start(ss *easyss.Easyss) {
	st := NewSysTray(ss)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM,
			syscall.SIGQUIT)

		select {
		case sig := <-c:
			log.Infof("[EASYSS-MAIN] got signal to exit: %v", sig)
			systray.Quit()
		case <-st.closing:
			log.Infof("[EASYSS-MAIN] easyss exiting...")
		}
	}()

	systray.Run(st.TrayReady, st.Exit) // system tray management
}
