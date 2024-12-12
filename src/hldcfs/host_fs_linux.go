package hldcfs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func getDeviceType(path string) (string, error) {
	// Den tatsächlichen Block-Device-Pfad finden
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Den Mountpunkt der Datei ermitteln
	mountPath := fmt.Sprintf("/proc/mounts")
	mountData, err := ioutil.ReadFile(mountPath)
	if err != nil {
		return "", err
	}

	var device string
	for _, line := range strings.Split(string(mountData), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		if strings.HasPrefix(absPath, parts[1]) {
			device = parts[0]
			break
		}
	}

	if device == "" {
		return "", fmt.Errorf("device not found for path: %s", path)
	}

	// Den Gerätenamen aus /dev extrahieren
	deviceName := filepath.Base(device)
	if strings.HasPrefix(deviceName, "mapper/") {
		deviceName = deviceName[len("mapper/"):]
	}

	rotationalPath := fmt.Sprintf("/sys/block/%s/queue/rotational", deviceName)

	// Rotational-Status lesen
	rotational, err := ioutil.ReadFile(rotationalPath)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(string(rotational)) == "0" {
		return "SSD", nil
	}
	return "HDD", nil
}
