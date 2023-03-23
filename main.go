package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple/auth"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple/recon"
)

const REMOVE_AFTER_VIEWING = false

var AP_CACHE map[string]recon.ReconAP

func signalHandler(signal os.Signal, stop context.CancelFunc) {
	if signal == os.Interrupt || signal == os.Kill {
		fmt.Println("")
		stop()
	}
}

func main() {
	token, err := auth.Login("root", "root")
	if err != nil {
		logrus.Errorf("AUTH   | failed to log into pineapple API: %v", err)
		return
	}
	logrus.Info("AUTH   | successfully authenticated")

	// Clear previous scan data
	backend.ClearScans(token)

	// Attempt to start the scan until it does
	backend.StartupScan(token)

	ctx, stop := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)

	go func() {
		for {
			s := <-sigCh
			signalHandler(s, stop)
		}
	}()

	wg := sync.WaitGroup{}

	apCache := make(map[string]recon.ReconAP)
	apEmit := make(chan recon.ReconAPExt)

	wg.Add(1)
	go backend.PollSSIDs(ctx, &wg, token, apCache, apEmit)
	wg.Add(1)
	go backend.StartAPI(ctx, &wg, apEmit)

	wg.Wait()
	logrus.Info("shutdown successful")
}
