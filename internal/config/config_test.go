package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func createTempEnvFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("не удалось создать временный .env: %v", err)
	}
	return path
}

func TestLoadAPIConfig_WithEnvFile(t *testing.T) {
	envPath := createTempEnvFile(t, "API_PORT=9090\n")

	oldWd, _ := os.Getwd()
	defer os.Chdir(filepath.Dir(oldWd))
	os.Chdir(filepath.Dir(envPath))

	cfg, err := LoadAPIConfig()
	if err != nil {
		t.Fatalf("ожидали nil error, получили %v", err)
	}

	if cfg.APIPort != "9090" {
		t.Errorf("ожидали порт 9090, получили %s", cfg.APIPort)
	}
}

func TestLoadAPIConfig_DefaultPort(t *testing.T) {
	t.Setenv("API_PORT", "")

	cfg, err := LoadAPIConfig()
	if err != nil {
		t.Fatalf("ожидали nil error, получили %v", err)
	}

	if cfg.APIPort != "8081" {
		t.Errorf("ожидали дефолтный порт 8081, получили %s", cfg.APIPort)
	}
}

func TestLoadAPIConfig_NoEnvFile(t *testing.T) {
	t.Setenv("API_PORT", "")

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(t.TempDir())

	_, err := LoadAPIConfig()
	if err == nil {
		t.Fatalf("ожидали ошибку, получили nil")
	}

	if !errors.Is(err, ErrEnvNotFound) {
		t.Errorf("ожидали ErrEnvNotFound, получили %v", err)
	}
}

func TestLoadDBConfig_NoParams(t *testing.T) {
	t.Setenv("POSTGRES_PORT", "28371")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_USER", "")

	_, err := LoadDBConfig()
	if err == nil {
		t.Fatalf("ожидали ошибку, получили nil")
	}

	if !errors.Is(err, ErrParamNotFound) {
		t.Errorf("ожиддали ErrParamNotFound, получили %v", err)
	}
}

func TestLoadKafkaConfig_Success(t *testing.T) {
	t.Setenv("KAFKA_PORT", "12345")
	t.Setenv("KAFKA_HOST", "localhost")
	t.Setenv("KAFKA_TOPIC_NAME", "events")
	t.Setenv("KAFKA_GROUP_ID", "group1")
	t.Setenv("KAFKA_DLQ_TOPIC_NAME", "events-dlq")

	cfg, err := LoadKafkaConfig()
	if err != nil {
		t.Fatalf("ожидали nil error, получили %v", err)
	}

	if cfg.KafkaPort != 12345 {
		t.Errorf("ожидали порт 12345, получили %d", cfg.KafkaPort)
	}
	if cfg.KafkaHost != "localhost" {
		t.Errorf("ожидали KafkaHost=localhost, получили %s", cfg.KafkaHost)
	}
	if cfg.KafkaTopicDLQName != "events-dlq" {
		t.Errorf("ожидали DLQ=events-dlq, получили %s", cfg.KafkaTopicDLQName)
	}
}

func TestLoadKafkaConfig_MissingParam(t *testing.T) {
	t.Setenv("KAFKA_PORT", "12345")
	t.Setenv("KAFKA_HOST", "")
	t.Setenv("KAFKA_TOPIC_NAME", "")
	t.Setenv("KAFKA_GROUP_ID", "")
	t.Setenv("KAFKA_DLQ_TOPIC_NAME", "")

	_, err := LoadKafkaConfig()
	if err == nil {
		t.Fatalf("ожидали ошибку, получили nil")
	}

	if !errors.Is(err, ErrKafkaParamNotFound) {
		t.Errorf("ожидали ErrKafkaParamNotFound, получили %v", err)
	}
}
