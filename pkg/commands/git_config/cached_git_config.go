package git_config

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type IGitConfig interface {
	Get(string) string
	GetBool(string) bool
}

type CachedGitConfig struct {
	cache  map[string]string
	getKey func(string) (string, error)
	log    *logrus.Entry
}

func NewStdCachedGitConfig(log *logrus.Entry) *CachedGitConfig {
	return NewCachedGitConfig(getGitConfigValue, log)
}

func NewCachedGitConfig(getKey func(string) (string, error), log *logrus.Entry) *CachedGitConfig {
	return &CachedGitConfig{
		cache:  make(map[string]string),
		getKey: getKey,
		log:    log,
	}
}

func (self *CachedGitConfig) Get(key string) string {
	if value, ok := self.cache[key]; ok {
		self.log.Debugf("using cache for key " + key)
		return value
	}

	value := self.getAux(key)
	self.cache[key] = value
	return value
}

func (self *CachedGitConfig) getAux(key string) string {
	value, err := self.getKey(key)
	if err != nil {
		self.log.Debugf("Error getting git config value for key: " + key + ". Error: " + err.Error())
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
