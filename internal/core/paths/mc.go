package paths

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/w1lam/Packages/utils"
)

func detectOrPromptMinecraftDir() (string, error) {
	if dir := defaultMinecraftDir(); dir != "" {
		return dir, nil
	}

	fmt.Println("Minecraft directory not found.")
	fmt.Print("Please enter the path to your Minecraft/server directory: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if !utils.CheckFileExists(input) {
		return "", fmt.Errorf("directory does not exist")
	}

	return input, nil
}

func defaultMinecraftDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		dir := filepath.Join(home, "AppData", "Roaming", ".minecraft")
		if utils.CheckFileExists(dir) {
			return dir
		}
	case "linux":
		dir := filepath.Join(home, ".minecraft")
		if utils.CheckFileExists(dir) {
			return dir
		}
	case "darwin":
		dir := filepath.Join("Library", "Application Support", "minecraft")
		if utils.CheckFileExists(dir) {
			return dir
		}
	}

	return ""
}
