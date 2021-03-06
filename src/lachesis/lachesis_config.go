package lachesis

import (
	"crypto/ecdsa"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/andrecronje/lachesis/src/node"
	"github.com/andrecronje/lachesis/src/proxy"
	aproxy "github.com/andrecronje/lachesis/src/proxy/app"
	"github.com/sirupsen/logrus"
)

type LachesisConfig struct {
	DataDir     string `mapstructure:"datadir"`
	BindAddr    string `mapstructure:"listen"`
	ServiceAddr string `mapstructure:"service-listen"`
	MaxPool     int    `mapstructure:"max-pool"`
	Store       bool   `mapstructure:"store"`
	LogLevel    string `mapstructure:"log"`

	NodeConfig node.Config `mapstructure:",squash"`

	LoadPeers bool
	Proxy     proxy.AppProxy
	Key       *ecdsa.PrivateKey
	Logger    *logrus.Logger

	Test  bool
	TestN uint64
}

func NewDefaultConfig() *LachesisConfig {
	config := &LachesisConfig{
		DataDir:     DefaultDataDir(),
		BindAddr:    ":1337",
		ServiceAddr: ":8000",
		MaxPool:     2,
		NodeConfig:  *node.DefaultConfig(),
		Store:       false,
		LogLevel:    "info",
		Proxy:       nil,
		Logger:      logrus.New(),
		LoadPeers:   true,
		Key:         nil,
		Test:        false,
		TestN:       ^uint64(0),
	}

	config.Logger.Level = LogLevel(config.LogLevel)
	config.Proxy = aproxy.NewInmemAppProxy(config.Logger)
	config.NodeConfig.Logger = config.Logger

	return config
}

func DefaultBadgerDir() string {
	dataDir := DefaultDataDir()
	if dataDir != "" {
		return filepath.Join(dataDir, "badger_db")
	}
	return ""
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".lachesis")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "BABBLE")
		} else {
			return filepath.Join(home, ".lachesis")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
