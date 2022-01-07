package git_config

type FakeGitConfig struct {
	mockResponses map[string]string
}

func NewFakeGitConfig(mockResponses map[string]string) *FakeGitConfig {
	return &FakeGitConfig{
		mockResponses: mockResponses,
	}
}

func (self *FakeGitConfig) Get(key string) string {
	if self.mockResponses == nil {
		return ""
	}
	return self.mockResponses[key]
}

func (self *FakeGitConfig) GetGeneral(args string) string {
	if self.mockResponses == nil {
		return ""
	}
	return self.mockResponses[args]
}

func (self *FakeGitConfig) GetBool(key string) bool {
	return isTruthy(self.Get(key))
}
