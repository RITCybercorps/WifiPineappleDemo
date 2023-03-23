package backend

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple/helpers"
	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple/recon"
)

func ClearScans(token string) {
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
}

func StartupScan(token string) {
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
}

func PrintSSIDs(token string, apCache map[string]recon.ReconAP) {
	apRes, err := recon.ListScanAps(token, 0)
	if err != nil {
		logrus.Errorf("RECON  | failed to list scan aps for scan 0: %v", err)
		return
	}
	for _, ap := range apRes.Aps {
		cachedAp, exists := apCache[ap.Bssid]
		if !exists {
			// AP hasn't been seen before
			apCache[ap.Bssid] = ap
		} else if ap.LastSeen-cachedAp.LastSeen > 30 {
			// If we've seen at least 30 secs after we last saw it
			apCache[ap.Bssid] = ap
		} else {
			// It must have dissapeared, so ignore it
			continue
		}
		logrus.Infof("AP     | \"%s\" - (%d clients)", ap.Ssid, ap.NumClients)
	}
}

func EmitAPs(token string, apCache map[string]recon.ReconAP, apEmit chan recon.ReconAPExt) {
	apRes, err := recon.ListScanAps(token, 0)
	if err != nil {
		logrus.Errorf("RECON  | failed to list scan aps for scan 0: %v", err)
		return
	}
	for _, ap := range apRes.Aps {
		cachedAp, exists := apCache[ap.Bssid]
		if !exists {
			// AP hasn't been seen before
			apCache[ap.Bssid] = ap
		} else if ap.LastSeen-cachedAp.LastSeen > 30 {
			// If we've seen at least 30 secs after we last saw it
			apCache[ap.Bssid] = ap
		} else {
			// It must have dissapeared, so ignore it
			continue
		}
		oui, _ := helpers.LookupOUI(token, strings.Replace(ap.Bssid[0:8], ":", "", 2))

		apEmit <- recon.ReconAPExt{
			ReconAP: ap,
			Vendor:  oui,
		}
		// Wait so animation can play on frontend
		time.Sleep(500 * time.Millisecond)
	}
}

func PollSSIDs(ctx context.Context, wg *sync.WaitGroup, token string, apCache map[string]recon.ReconAP, apEmit chan recon.ReconAPExt) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	// PrintSSIDs(token, apCache)
	EmitAPs(token, apCache, apEmit)
	// logrus.Info("waiting 10 secs...")
	for {
		select {
		case <-ticker.C:
			// PrintSSIDs(token, apCache)
			EmitAPs(token, apCache, apEmit)
			// logrus.Info("waiting 10 secs...")
		case <-ctx.Done():
			ticker.Stop()
			logrus.Warn("stopping pineapple polling...")
			return
		}
	}
}
