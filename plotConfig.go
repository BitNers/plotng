package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	TargetDirectory []string
	TempDirectory   []string
	NumberOfPlots   int
	Fingerprint     string
}

type PlotConfig struct {
	configPath    string
	currentConfig *Config
	lastMod       time.Time
	lock          sync.RWMutex
}

func (pc *PlotConfig) Init() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		pc.ProcessConfig()
		for range ticker.C {
			pc.ProcessConfig()
		}
	}()
}

func (pc *PlotConfig) ProcessConfig() {
	if fs, err := os.Lstat(pc.configPath); err != nil {
		log.Printf("Failed to open config file [%s]: %s\n", pc.configPath, err)
	} else {
		if pc.lastMod != fs.ModTime() {
			if f, err := os.Open(pc.configPath); err != nil {
				log.Printf("Failed to open config file [%s]: %s\n", pc.configPath, err)
			} else {
				decoder := json.NewDecoder(f)
				var newConfig Config
				if err := decoder.Decode(&newConfig); err != nil {
					log.Printf("Failed to process config file [%s]: %s\n", pc.configPath, err)
				} else {
					pc.lock.Lock()
					pc.currentConfig = &newConfig
					pc.lock.Unlock()
					log.Printf("New configuration loaded")
				}
				f.Close()
			}
			pc.lastMod = fs.ModTime()
		}
	}
}