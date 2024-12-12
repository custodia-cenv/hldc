package hldcfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func getDeviceType(path string) (string, error) {
	// Ermitteln des Mount-Punktes der Datei
	cmd := exec.Command("df", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute df command: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output from df command")
	}

	// Extrahiere das Device aus der ersten Zeile nach dem Header
	fields := strings.Fields(lines[1])
	if len(fields) < 1 {
		return "", fmt.Errorf("failed to parse df output")
	}
	device := fields[0]

	// Führe diskutil für das ermittelte Gerät aus
	cmd = exec.Command("diskutil", "info", device)
	out.Reset()
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("diskutil error: %v", err)
	}

	// Prüfe die Ausgabe auf bekannte Typen
	output := out.String()
	if strings.Contains(output, "Solid State: Yes") {
		return "SSD", nil
	}
	if strings.Contains(output, "Solid State: No") {
		return "HDD", nil
	}
	if strings.Contains(output, "Fusion Drive") {
		return "FD", nil
	}
	if strings.Contains(output, "Virtual Disk") {
		return "VD", nil
	}

	// Standard-Fallback, falls keine spezifische Zuordnung möglich ist
	return "Unknown", nil
}
