package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// currentVersion fetches the latest release tag from the GitHub API
func currentVersion() (string, error) {
	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return "", errors.New("GITHUB_REPOSITORY environment variable is not set")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "bump-version")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	tag, ok := data["tag_name"].(string)
	if !ok {
		return "", errors.New("tag_name not found in GitHub API response")
	}
	return tag, nil
}

// bump increases the specified version component and resets lower components to zero
func bump(version string, component string) (string, error) {
	components := []string{"major", "minor", "patch"}
	elements := strings.Split(version, ".")
	index := -1
	for i, c := range components {
		if c == component {
			index = i
			break
		}
	}

	if index < 0 {
		return "", fmt.Errorf("provided version component (%s) is not one of 'major', 'minor', 'patch'", component)
	}
	if index >= len(elements) {
		return "", fmt.Errorf("provided version component (%s) is not part of the provided version (%s)", component, version)
	}

	val, err := strconv.Atoi(elements[index])
	if err != nil {
		return "", fmt.Errorf("invalid number in version: %v", err)
	}
	elements[index] = strconv.Itoa(val + 1)

	for i := index + 1; i < len(elements); i++ {
		elements[i] = "0"
	}

	return strings.Join(elements, "."), nil
}

func main() {
	version, err := currentVersion()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	newVersion, err := bump(version, "minor")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Bumped version:", newVersion)
	}
}