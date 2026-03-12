package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestParseConfig(t *testing.T) {
	// Reset Viper to ensure clean state across tests
	viper.Reset()

	// Create a temporary directory to store test config files
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir) // Automatically clean up temp directory after tests

	// Define test cases with expected configurations
	tests := []struct {
		name string  // Test case name
		want *Config // Expected configuration after parsing (created by cfg function)
	}{
		{"default", cfg(t, tmpDir)}, // Test default config loaded from temp directory
	}
	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset Viper before each test
			viper.Reset()
			// Add temp directory as config path before parsing
			viper.AddConfigPath(tmpDir)
			viper.SetConfigName("config")
			// Parse the configuration with an empty path (use config in temp dir)
			got := ParseConfig("")
			// Compare the parsed config with the expected config
			if !reflect.DeepEqual(got, tt.want) {
				// Marshal both configs to JSON for easier comparison in error message
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.want)
				t.Errorf("ParseConfig(\"\") = %s, want %s", gotJSON, wantJSON)
			}
		})
	}
}

func cfg(t *testing.T, tmpDir string) *Config {
	// **1. Test Configuration Setup**

	// **a. Set the configuration format to YAML:**
	viper.SetConfigType("yaml")
	// **b. Define the test YAML configuration as a byte slice:**
	// This variable holds the test configuration data in a compact format.
	var yamlCfg = []byte(`
address: 1.2.3.4
port: 42
prefix: /oh-de-lally
tls:
  keyFile: ` + tmpDir + `/robin.pem
  certFile: ` + tmpDir + `/tuck.pem
dir: /sherwood/forest
realm: uk
users:
  lj:
    password: 123
    subdir: /littlejohn
    permissions: "crud"
  srf:
    password: 234
    subdir: /sheriff
    permissions: "crud"
log:
  error: true
`)
	// **2. Temporary File Creation**

	// **a. Create temporary files for TLS key and certificate:**
	// These files will be populated with empty data to represent the TLS configuration.
	err := os.WriteFile(filepath.Join(tmpDir, "config.yaml"), yamlCfg, 0600)
	if err != nil {
		t.Errorf("error writing test config. error = %v", err)
	}

	// **3. Configuration Reading and Parsing**

	// **a. Read the configuration from the temporary file:**
	err = viper.ReadConfig(bytes.NewBuffer(yamlCfg))
	if err != nil {
		t.Errorf("error reading test config. error = %v", err)
	}
	// **b. Allocate and unmarshal the configuration data:**
	// Create a new Config instance and populate it with the parsed data.
	var resultCfg = &Config{
		Users: make(map[string]*UserInfo),
	}
	err = viper.Unmarshal(resultCfg)
	if err != nil {
		log.WithError(err).Error("Error unmarshalling config in test")
		t.Errorf("error unmarshalling config. error = %v", err)
		return nil
	}

	// **4. User Permissions Processing**
	for user := range viper.GetStringMap("Users") {
		permissions := viper.GetString(fmt.Sprintf("Users.%s.permissions", user)) // Access specific user permissions
		if resultCfg.Users[user] == nil {
			resultCfg.Users[user] = &UserInfo{}
		}
		resultCfg.Users[user].Crud = &CrudType{Crud: permissions} // Set user's CRUD permissions object
		err := FormatCrud(context.Background(), user, resultCfg)  // Further process and validate permissions
		if err != nil {
			log.WithError(err).WithField("user", user).Error("Error parsing crud string from config file") // log error with context
		}
		log.WithFields(logrus.Fields{"user": user,
			"crud": resultCfg.Users[user].Crud}).Info("Parsed crud string from config file") // Log parsed permissions
	}

	// **5. Config Path and Dummy Files (Optional)**

	// **a. Add the temporary directory to the config path:**
	// This ensures viper prioritizes the test configuration.
	viper.AddConfigPath(tmpDir)

	// add dummy cert and key file
	_, err = os.OpenFile(filepath.Join(tmpDir, "robin.pem"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("error creating key file. error = %v", err)
		return nil
	}

	_, err = os.OpenFile(filepath.Join(tmpDir, "tuck.pem"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Errorf("error creating cert file. error = %v", err)
		return nil
	}
	viper.AddConfigPath(tmpDir)
	// **6. Return the Config Instance**
	// Return the populated Config instance for further use in the test case.
	return resultCfg
}

func TestGenHashFromPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		cost      int
		wantError bool
	}{
		{
			name:      "valid password with default cost",
			password:  "testpassword",
			cost:      10,
			wantError: false,
		},
		{
			name:      "valid password with high cost",
			password:  "testpassword",
			cost:      12,
			wantError: false,
		},
		{
			name:      "empty password",
			password:  "",
			cost:      10,
			wantError: false,
		},
		{
			name:      "password under bcrypt limit",
			password:  "testpassword12345",
			cost:      10,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GenHashFromPassword(tt.password, tt.cost)

			if (err != nil) != tt.wantError {
				t.Errorf("GenHashFromPassword() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if hash == "" {
					t.Error("GenHashFromPassword() returned empty hash")
				}

				if len(hash) != 60 {
					t.Errorf("GenHashFromPassword() returned hash of length %d, want 60", len(hash))
				}

				if hash[0] != '$' {
					t.Errorf("GenHashFromPassword() returned hash starting with %q, want '$'", hash[0])
				}

				expectedPrefix := fmt.Sprintf("$2a$%02d$", tt.cost)
				if !strings.HasPrefix(hash, expectedPrefix) {
					t.Errorf("GenHashFromPassword() returned hash with prefix %q, want %q", hash[:8], expectedPrefix)
				}
			}
		})
	}
}

func TestHandleConfigUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	cfg := &Config{
		Dir:   tmpDir,
		Users: map[string]*UserInfo{},
	}

	err := os.WriteFile(configFile, []byte("address: \"127.0.0.1\"\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	cfg.handleConfigUpdate(fsnotify.Event{Name: configFile})
}

func TestUpdateConfig(t *testing.T) {
	subdirOld := "/old"
	subdirNew := "/new"
	subdirTest := "/testuser"

	tests := []struct {
		name          string
		existingUsers map[string]*UserInfo
		newUser       *UserInfo
		wantError     bool
	}{
		{
			name:          "add new user",
			existingUsers: map[string]*UserInfo{},
			newUser: &UserInfo{
				Password:    "$2a$10$testhash",
				Subdir:      &subdirTest,
				Permissions: "crud",
			},
			wantError: false,
		},
		{
			name: "update existing user",
			existingUsers: map[string]*UserInfo{
				"testuser": {
					Password:    "$2a$10$oldhash",
					Subdir:      &subdirOld,
					Permissions: "r",
				},
			},
			newUser: &UserInfo{
				Password:    "$2a$10$newhash",
				Subdir:      &subdirNew,
				Permissions: "crud",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Users: tt.existingUsers,
			}

			newCfg := &Config{
				Users: map[string]*UserInfo{"testuser": tt.newUser},
			}

			updateConfig(cfg, newCfg)

			if _, ok := cfg.Users["testuser"]; !ok {
				t.Errorf("updateConfig() failed to add/update user testuser")
			}

			if len(tt.existingUsers) > 0 {
				if cfg.Users["testuser"].Password != tt.newUser.Password {
					t.Errorf("updateConfig() password = %q, want %q", cfg.Users["testuser"].Password, tt.newUser.Password)
				}
			}
		})
	}
}

func TestCreateBaseAndUserDirectoriesIfNeeded(t *testing.T) {
	tmpDir := t.TempDir()

	user1Subdir := "/user1"
	user2Subdir := "/user2"

	cfg := &Config{
		Dir: tmpDir,
		Users: map[string]*UserInfo{
			"user1": {Subdir: &user1Subdir},
			"user2": {Subdir: &user2Subdir},
		},
	}

	os.RemoveAll(cfg.Dir)

	cfg.createBaseAndUserDirectoriesIfNeeded()

	if _, err := os.Stat(cfg.Dir); os.IsNotExist(err) {
		t.Error("createBaseAndUserDirectoriesIfNeeded() failed to create base directory")
	}

	for _, user := range cfg.Users {
		userPath := filepath.Join(cfg.Dir, *user.Subdir)
		if _, err := os.Stat(userPath); os.IsNotExist(err) {
			t.Errorf("createBaseAndUserDirectoriesIfNeeded() failed to create user directory %s", userPath)
		}
	}

	cfg2 := &Config{
		Dir:   filepath.Join(tmpDir, "test2"),
		Users: nil,
	}
	os.RemoveAll(cfg2.Dir)

	cfg2.createBaseAndUserDirectoriesIfNeeded()

	cfg3 := &Config{
		Dir: filepath.Join(tmpDir, "test3"),
		Users: map[string]*UserInfo{
			"emptyuser": {Subdir: new(string)},
		},
	}
	os.RemoveAll(cfg3.Dir)

	cfg3.createBaseAndUserDirectoriesIfNeeded()
}
