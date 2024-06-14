package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AdvanceSettings struct {
	ReaderCachedMinutes uint `yaml:"reader_cached_minutes"`
	// JWT configuration
	JWTTimeoutHours    uint   `yaml:"jwt_timeout_hours"`
	JWTMaxRefreshHours uint   `yaml:"jwt_max_refresh_hours"`
	JWTSecretKey       string `yaml:"jwt_secret_key"`
	JWTIdentityKey     string `yaml:"jwt_identity_key`
	JWTHeadName        string `yaml:"jwt_head_name"`
}

type KolibraSettings struct {
	Database          string             `yaml:"database"`
	Library           string             `yaml:"library"`
	Port              uint               `yaml:"port"`
	Host              string             `yaml:"host"`
	BookExtension     []string           `yaml:"book_extension"`
	FileSizeThreshold int64              `yaml:"file_size_threshold"`
	FileNameMethod    FileNameMethodType `yaml:"file_name_method"`
	DefaultTitleRegex string             `yaml:"default_title_regex"`
	Advance           AdvanceSettings    `yaml:"advance"`
}

type FileNameMethodType string

const (
	DIR_AUTHOR  FileNameMethodType = "DIR_AUTHOR"
	FILE_AUTHOR FileNameMethodType = "FILE_AUTHOR"
)

var Settings *KolibraSettings

func load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &Settings)
	log.Printf("Config: %v", Settings)
	return err
}

func LoadConfig() {
	log.Printf("Start to load config from ./config.yaml")
	err := load("config.yaml")
	if err != nil {
		log.Panicf("Failed to load config: %s", err)
	}
}
