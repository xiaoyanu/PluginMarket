package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Uploads  UploadsConfig  `yaml:"uploads"`
}

type ServerConfig struct {
	Port   int    `yaml:"port"`
	Mode   string `yaml:"mode"`
	WebURL string `yaml:"webUrl"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"` // 小时
}

type UploadsConfig struct {
	Path    string  `yaml:"path"`
	MaxSize float64 `yaml:"maxSize"` // MB
	Avatars string  `yaml:"avatars"`
	Images  string  `yaml:"images"`
	Frames  string  `yaml:"frames"`
	Titles  string  `yaml:"titles"`
}

var C Config

func Load(filename string) error {
	// 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// 获取可执行文件所在的目录
	exeDir := filepath.Dir(exePath)

	// 拼接配置文件的绝对路径
	configPath := filepath.Join(exeDir, filename)

	// 优先读取可执行文件同级目录下的配置
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 如果在可执行文件同级目录没找到（可能是开发环境下用 go run 运行），
		// 则尝试直接读取文件名（相对于当前工作目录）
		data, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
	}
	return yaml.Unmarshal(data, &C)
}
