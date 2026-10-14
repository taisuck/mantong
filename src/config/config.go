// Package config 는 config.yml 설정 파일을 읽는다.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Log     LogConfig     `yaml:"log"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	Dir   string `yaml:"dir"`
}

type MongoDBConfig struct {
	URI      string        `yaml:"uri"`
	Database string        `yaml:"database"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Load 는 path 의 YAML 파일을 읽어 기본값을 채운 Config 를 반환한다.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	cfg := Config{
		Server:  ServerConfig{Addr: ":8080"},
		Log:     LogConfig{Level: "info", Dir: "logs"},
		MongoDB: MongoDBConfig{Timeout: 10 * time.Second},
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("설정 파일 파싱 실패: %w", err)
	}

	// "8080" 처럼 포트만 지정한 경우 ":8080" 으로 변환한다.
	if addr := strings.TrimSpace(cfg.Server.Addr); addr != "" && !strings.Contains(addr, ":") {
		cfg.Server.Addr = ":" + addr
	}

	if cfg.MongoDB.URI == "" {
		return nil, errors.New("mongodb.uri 가 설정되지 않았습니다")
	}
	if cfg.MongoDB.Database == "" {
		return nil, errors.New("mongodb.database 가 설정되지 않았습니다")
	}
	if (cfg.MongoDB.User == "") != (cfg.MongoDB.Password == "") {
		return nil, errors.New("mongodb.user 와 mongodb.password 는 함께 설정해야 합니다")
	}
	return &cfg, nil
}
