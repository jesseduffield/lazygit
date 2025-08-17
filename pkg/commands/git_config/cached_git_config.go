package git_config

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type IGitConfig interface {
	// this is for when you want to pass 'mykey' (it calls `git config --get --null mykey` under the hood)
	Get(string) string
	// this is for when you want to pass '--local --get-regexp mykey'
	GetGeneral(string) string
	// this is for when you want to pass 'mykey' and check if the result is truthy
	GetBool(string) bool

	DropCache()
}

type CachedGitConfig struct {
	cache           map[string]string
	runGitConfigCmd func(*exec.Cmd) (string, error)
	log             *logrus.Entry
	mutex           sync.Mutex
}

func NewStdCachedGitConfig(log *logrus.Entry) *CachedGitConfig {
	return NewCachedGitConfig(runGitConfigCmd, log)
}

func NewCachedGitConfig(runGitConfigCmd func(*exec.Cmd) (string, error), log *logrus.Entry) *CachedGitConfig {
	return &CachedGitConfig{
		cache:           make(map[string]string),
		runGitConfigCmd: runGitConfigCmd,
		log:             log,
		mutex:           sync.Mutex{},
	}
}

func (self *CachedGitConfig) Get(key string) string {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if value, ok := self.cache[key]; ok {
		self.log.Debug("using cache for key " + key)
		return value
	}

	value := self.getAux(key)
	self.cache[key] = value
	return value
}

func (self *CachedGitConfig) GetGeneral(args string) string {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if value, ok := self.cache[args]; ok {
		self.log.Debug("using cache for args " + args)
		return value
	}

	value := self.getGeneralAux(args)
	self.cache[args] = value
	return value
}

func (self *CachedGitConfig) getGeneralAux(args string) string {
	cmd := getGitConfigGeneralCmd(args)
	value, err := self.runGitConfigCmd(cmd)
	if err != nil {
		self.log.Debugf("Error getting git config value for args: %s. Error: %v", args, err.Error())
		return ""
	}
	return strings.TrimSpace(value)
}

func (self *CachedGitConfig) getAux(key string) string {
	cmd := getGitConfigCmd(key)
	value, err := self.runGitConfigCmd(cmd)
	if err != nil {
		self.log.Debugf("Error getting git config value for key: %s. Error: %v", key, err.Error())
		return ""
	}
	return strings.TrimSpace(value)
}

func (self *CachedGitConfig) GetBool(key string) bool {
	return isTruthy(self.Get(key))
}

func isTruthy(value string) bool {
	lcValue := strings.ToLower(value)
	return lcValue == "true" || lcValue == "1" || lcValue == "yes" || lcValue == "on"
}

func (self *CachedGitConfig) DropCache() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.cache = make(map[string]string)
}
