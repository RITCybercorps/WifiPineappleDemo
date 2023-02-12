package main

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/pineapple/auth"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/pineapple/pineap"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/pineapple/recon"
)

const REMOVE_AFTER_VIEWING = false

func PrintSSIDs(token string) {
	ssids, err := pineap.ListSSIDs(token)
	if err != nil {
		logrus.Errorf("PINEAP | failed to list PineAP ssids: %v", err)
		return
	}
	ssidList := strings.Split(ssids, "\n")
	logrus.Info("PINEAP | SSID List:")
	for _, ssid := range ssidList {
		if ssid != "" {
			logrus.Infof("PINEAP |  | \t%s", ssid)
			if REMOVE_AFTER_VIEWING {
				success, err := pineap.RemoveSSID(token, ssid)
				if err != nil {
					logrus.Errorf("PINEAP | failed to remove PineAP ssid \"%s\": %v", ssid, err)
					return
				}
				if success {
					logrus.Infof("PINEAP |  | \t | removed SSID \"%s\"", ssid)
				} else {
					logrus.Warnf("PINEAP |  | \t | couldn't remove SSID \"%s\"", ssid)
				}
			}
		}
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
	for {
		scanList, err := recon.ListScans(token)
		if err != nil {
			logrus.Errorf("RECON  | failed to list scans")
			time.Sleep(1 * time.Second)
			continue
		}
		// If no scans, keep going
		if len(scanList) == 0 {
			logrus.Info("RECON  | successfully cleared previous scans")
			break
		}
		logrus.Info("RECON  | clearing previous scans...")
		for _, scan := range scanList {
			success, err := recon.DeleteScan(token, scan.ScanID)
			if err != nil || !success {
				logrus.Warnf("RECON  | failed to delete scan %d", scan.ScanID)
			} else {
				logrus.Infof("RECON  | deleted scan %d", scan.ScanID)
			}
		}
	}

	// Attempt to start the scan until it does
	for {
		scanStatus, err := recon.Status(token)
		if err != nil {
			logrus.Errorf("RECON  | failed to get recon status: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if !scanStatus.ScanRunning {
			success, err := recon.StartScan(token, false, recon.ReconCONTINUOUS, recon.Recon2_4GHZ_5GHZ)
			if err != nil || !success {
				// Try it again if it doesn't work
				logrus.Warn("RECON  | failed to start recon scan, trying again...")
				time.Sleep(1 * time.Second)
				continue
			}
			if success {
				logrus.Info("RECON  | started recon scan")
				continue
			}
		}
		// If the scan is running, but not in continuous mode
		if !scanStatus.Continuous {
			logrus.Warn("RECON  | recon scan running, but is misconfigured. stopping...")
			success, err := recon.StopScan(token)
			if err != nil || !success {
				logrus.Warn("RECON  | failed to stop recon scan, trying again...")
				time.Sleep(1 * time.Second)
				continue
			}
			if success {
				logrus.Info("RECON  | recon scan stopped, restarting...")
				time.Sleep(1 * time.Second)
				continue
			}
		}
		// If the scan is running and in continuous mode, we exit the loop
		logrus.Info("RECON  | scan properly configured")
		break
	}

	ticker := time.NewTicker(10 * time.Second)
	PrintSSIDs(token)
	logrus.Info("waiting 10 secs...")
	for range ticker.C {
		PrintSSIDs(token)
		logrus.Info("waiting 10 secs...")
	}
}
