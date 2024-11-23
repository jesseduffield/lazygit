package keyring

type mockProvider struct {
	mockStore map[string]map[string]string
	mockError error
}

// Set stores user and pass in the keyring under the defined service
// name.
func (m *mockProvider) Set(service, user, pass string) error {
	if m.mockError != nil {
		return m.mockError
	}
	if m.mockStore == nil {
		m.mockStore = make(map[string]map[string]string)
	}
	if m.mockStore[service] == nil {
		m.mockStore[service] = make(map[string]string)
	}
	m.mockStore[service][user] = pass
	return nil
}

// Get gets a secret from the keyring given a service name and a user.
func (m *mockProvider) Get(service, user string) (string, error) {
	if m.mockError != nil {
		return "", m.mockError
	}
	if b, ok := m.mockStore[service]; ok {
		if v, ok := b[user]; ok {
			return v, nil
		}
	}
	return "", ErrNotFound
}

// Delete deletes a secret, identified by service & user, from the keyring.
func (m *mockProvider) Delete(service, user string) error {
	if m.mockError != nil {
		return m.mockError
	}
	if m.mockStore != nil {
		if _, ok := m.mockStore[service]; ok {
			if _, ok := m.mockStore[service][user]; ok {
				delete(m.mockStore[service], user)
				return nil
			}
		}
	}
	return ErrNotFound
}

// DeleteAll deletes all secrets for a given service
func (m *mockProvider) DeleteAll(service string) error {
	if m.mockError != nil {
		return m.mockError
	}
	delete(m.mockStore, service)
	return nil
}

// MockInit sets the provider to a mocked memory store
func MockInit() {
	provider = &mockProvider{}
}

// MockInitWithError sets the provider to a mocked memory store
// that returns the given error on all operations
func MockInitWithError(err error) {
	provider = &mockProvider{mockError: err}
}
