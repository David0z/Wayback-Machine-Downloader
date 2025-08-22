package util

import (
	"os/exec"
	"strings"
)

func GetClipboardContent() string {
	cmd := exec.Command("powershell.exe", "-command", "Get-Clipboard")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
