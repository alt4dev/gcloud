package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSetupOptions(t *testing.T) {
	// Check that setting options via a JSON file updates options
	content := struct {
		Token string `json:"token"`
		Mode string `json:"mode"`
		Source string `json:"source"`
	}{"some random token from alt4.dev","testing", "api"}
	contentJson, _ := json.Marshal(content)
	// Write the file
	fileName := fmt.Sprintf("/tmp/alt4_key_file_%d.json", time.Now().UnixNano())
	err := ioutil.WriteFile(fileName, contentJson, 0644)
	if err != nil {
		t.Error(err)
		return
	}
	os.Setenv("ALT4_CONFIG", fileName)
	// Set options and confirm that options were updated accordingly
	setupOptions()

	if options.AuthToken != content.Token {
		t.Error("Tokens don't match after options setup")
		return
	}

	if options.Mode != content.Mode {
		t.Error("Modes don't match after options setup")
		return
	}

	if options.Source != content.Source {
		t.Error("Sinks don't match after options setup")
		return
	}

	// Test that having individual environment variables set overrides the config file
	os.Setenv("ALT4_AUTH_TOKEN", "A second token")
	os.Setenv("ALT4_MODE", "release")
	os.Setenv("ALT4_SOURCE", "javascript")

	setupOptions()
	if options.AuthToken != "A second token" {
		t.Error("Token not overridden by env")
		return
	}

	if options.Mode != "release" {
		t.Error("Mode not overridden by env")
		return
	}

	if options.Source != "javascript" {
		t.Error("Source not overridden by env")
		return
	}

	// Test that setting options from the code overrides options set from env
	SetAuthToken("Some new auth token")
	if options.AuthToken != "Some new auth token" {
		t.Error("Token not overridden")
		return
	}

	SetMode("testing")
	if options.Mode != "testing" {
		t.Error("Mode not overridden")
		return
	}

	SetSource("default")
	if options.Source != "default" {
		t.Error("Sink not overridden")
		return
	}
}


