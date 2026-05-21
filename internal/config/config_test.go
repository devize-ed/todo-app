package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// saveOriginalState saves the original os.Args and flag.CommandLine to restore after tests
func saveOriginalState() ([]string, *flag.FlagSet) {
	originalArgs := make([]string, len(os.Args))
	copy(originalArgs, os.Args)
	originalFlags := flag.CommandLine
	return originalArgs, originalFlags
}

// restoreOriginalState restores the original os.Args and flag.CommandLine
func restoreOriginalState(originalArgs []string, originalFlags *flag.FlagSet) {
	os.Args = originalArgs
	flag.CommandLine = originalFlags
}

// writeTempConfig writes content to a temp config file and returns its path
func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	assert.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// envKeys are all env vars that influence the Config
var envKeys = []string{
	"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
	"POSTGRES_PORT", "POSTGRES_HOST", "POSTGRES_SSLMODE",
	"SERVER_ADDRESS", "SERVER_TIMEOUT", "LOG_LEVEL",
}

// clearEnv unsets every env var that Config reads
func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range envKeys {
		assert.NoError(t, os.Unsetenv(k))
	}
}

func TestLoad_PriorityEnvOverFile(t *testing.T) {
	originalArgs, originalFlags := saveOriginalState()
	defer restoreOriginalState(originalArgs, originalFlags)

	fileContent := `
postgres:
  user: file-user
  password: file-pass
  db: file-db
  port: 5432
  host: file-host
  ssl_mode: disable
server:
  address: file-host:1234
  timeout: 1s
log_level: info
`

	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "ENV overrides file (top-level LOG_LEVEL only)",
			envVars: map[string]string{
				"POSTGRES_USER": "env-user-ignored",
			},
			expected: Config{
				Postgres: PostgresConfig{
					User:     "file-user",
					Password: "file-pass",
					DB:       "file-db",
					Port:     5432,
					Host:     "file-host",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Addr:    "file-host:1234",
					Timeout: time.Second,
				},
				LogLevel: "error",
			},
		},
		{
			name:    "Only file (no ENV)",
			envVars: map[string]string{},
			expected: Config{
				Postgres: PostgresConfig{
					User:     "file-user",
					Password: "file-pass",
					DB:       "file-db",
					Port:     5432,
					Host:     "file-host",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Addr:    "file-host:1234",
					Timeout: time.Second,
				},
				LogLevel: "info",
			},
		},
		{
			name: "Defaults preserved when file is empty",
			envVars: map[string]string{
				"LOG_LEVEL": "warn",
			},
			expected: Config{
				Postgres: PostgresConfig{
					User:     "postgres",
					Password: "postgres",
					DB:       "postgres",
					Port:     5432,
					Host:     "localhost",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Addr: "localhost:8080",
				},
				LogLevel: "warn",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tc.envVars {
				assert.NoError(t, os.Setenv(k, v))
			}

			var path string
			if tc.name == "Defaults preserved when file is empty" {
				path = writeTempConfig(t, "")
			} else {
				path = writeTempConfig(t, fileContent)
			}

			cfg, err := Load(path)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.Postgres, cfg.Postgres, "Postgres mismatch")
			assert.Equal(t, tc.expected.Server, cfg.Server, "Server mismatch")
			assert.Equal(t, tc.expected.LogLevel, cfg.LogLevel, "Log level mismatch")
		})
	}
}

func TestLoad_EnvironmentVariableParsing(t *testing.T) {
	originalArgs, originalFlags := saveOriginalState()
	defer restoreOriginalState(originalArgs, originalFlags)

	tests := []struct {
		name    string
		envVars map[string]string
		check   func(*testing.T, *Config)
	}{
		{
			name:    "Empty environment variables should use defaults",
			envVars: map[string]string{},
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "localhost:8080", cfg.Server.Addr)
				assert.Equal(t, "postgres", cfg.Postgres.User)
				assert.Equal(t, "postgres", cfg.Postgres.Password)
				assert.Equal(t, "postgres", cfg.Postgres.DB)
				assert.Equal(t, 5432, cfg.Postgres.Port)
				assert.Equal(t, "localhost", cfg.Postgres.Host)
				assert.Equal(t, "disable", cfg.Postgres.SSLMode)
				assert.Equal(t, "debug", cfg.LogLevel)
			},
		},
		{
			name: "LOG_LEVEL env var overrides default",
			envVars: map[string]string{
				"LOG_LEVEL": "info",
			},
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "info", cfg.LogLevel)
				// defaults remain intact for fields with no top-level env tag
				assert.Equal(t, "localhost:8080", cfg.Server.Addr)
				assert.Equal(t, "postgres", cfg.Postgres.User)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tc.envVars {
				assert.NoError(t, os.Setenv(k, v))
			}

			path := writeTempConfig(t, "")
			cfg, err := Load(path)
			assert.NoError(t, err)
			tc.check(t, cfg)
		})
	}
}

func TestLoad_FileParsing(t *testing.T) {
	originalArgs, originalFlags := saveOriginalState()
	defer restoreOriginalState(originalArgs, originalFlags)

	tests := []struct {
		name     string
		content  string
		wantErr  bool
		expected Config
	}{
		{
			name: "Fully populated YAML",
			content: `
postgres:
  user: yaml-user
  password: yaml-pass
  db: yaml-db
  port: 6543
  host: yaml-host
  ssl_mode: require
server:
  address: 0.0.0.0:9090
  timeout: 5s
log_level: warn
`,
			expected: Config{
				Postgres: PostgresConfig{
					User:     "yaml-user",
					Password: "yaml-pass",
					DB:       "yaml-db",
					Port:     6543,
					Host:     "yaml-host",
					SSLMode:  "require",
				},
				Server: ServerConfig{
					Addr:    "0.0.0.0:9090",
					Timeout: 5 * time.Second,
				},
				LogLevel: "warn",
			},
		},
		{
			name: "Partial YAML keeps defaults",
			content: `
postgres:
  user: only-user
log_level: error
`,
			expected: Config{
				Postgres: PostgresConfig{
					User:     "only-user",
					Password: "postgres",
					DB:       "postgres",
					Port:     5432,
					Host:     "localhost",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Addr: "localhost:8080",
				},
				LogLevel: "error",
			},
		},
		{
			name:    "Invalid YAML returns error",
			content: ":::: bad yaml ::::\n  - [unbalanced",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			path := writeTempConfig(t, tc.content)

			cfg, err := Load(path)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.Postgres, cfg.Postgres, "Postgres mismatch")
			assert.Equal(t, tc.expected.Server, cfg.Server, "Server mismatch")
			assert.Equal(t, tc.expected.LogLevel, cfg.LogLevel, "Log level mismatch")
		})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	originalArgs, originalFlags := saveOriginalState()
	defer restoreOriginalState(originalArgs, originalFlags)

	clearEnv(t)
	missing := filepath.Join(t.TempDir(), "does-not-exist.yaml")

	cfg, err := Load(missing)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestGetFilePath_FlagParsing(t *testing.T) {
	originalArgs, originalFlags := saveOriginalState()
	defer restoreOriginalState(originalArgs, originalFlags)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Default config path",
			args:     []string{},
			expected: "/cfg/config.yaml",
		},
		{
			name:     "Custom config path via flag",
			args:     []string{"-config=/custom/path.yaml"},
			expected: "/custom/path.yaml",
		},
		{
			name:     "Custom config path via space-separated flag",
			args:     []string{"-config", "/another/path.yaml"},
			expected: "/another/path.yaml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tc.name, flag.ContinueOnError)
			os.Args = append([]string{"cmd"}, tc.args...)

			got := GetFilePath()
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestPostgresConfig_BuildURI(t *testing.T) {
	tests := []struct {
		name string
		cfg  PostgresConfig
		want string
	}{
		{
			name: "Default-like config",
			cfg: PostgresConfig{
				User:     "postgres",
				Password: "postgres",
				Host:     "localhost",
				Port:     5432,
				DB:       "postgres",
				SSLMode:  "disable",
			},
			want: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		},
		{
			name: "Custom config",
			cfg: PostgresConfig{
				User:     "admin",
				Password: "s3cret",
				Host:     "db.internal",
				Port:     6543,
				DB:       "todoapp",
				SSLMode:  "require",
			},
			want: "postgres://admin:s3cret@db.internal:6543/todoapp?sslmode=require",
		},
		{
			name: "Empty config produces zero-value URI",
			cfg:  PostgresConfig{},
			want: "postgres://:@:0/?sslmode=",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.cfg.BuildURI())
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost:8080", cfg.Server.Addr)
	assert.Equal(t, "postgres", cfg.Postgres.User)
	assert.Equal(t, "postgres", cfg.Postgres.Password)
	assert.Equal(t, "postgres", cfg.Postgres.DB)
	assert.Equal(t, 5432, cfg.Postgres.Port)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "disable", cfg.Postgres.SSLMode)
	assert.Equal(t, "debug", cfg.LogLevel)
}
