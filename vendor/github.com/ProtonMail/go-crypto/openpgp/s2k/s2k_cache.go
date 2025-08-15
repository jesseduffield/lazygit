package s2k

// Cache stores keys derived with s2k functions from one passphrase
// to avoid recomputation if multiple items are encrypted with
// the same parameters.
type Cache map[Params][]byte

// GetOrComputeDerivedKey tries to retrieve the key
// for the given s2k parameters from the cache.
// If there is no hit, it derives the key with the s2k function from the passphrase,
// updates the cache, and returns the key.
func (c *Cache) GetOrComputeDerivedKey(passphrase []byte, params *Params, expectedKeySize int) ([]byte, error) {
	key, found := (*c)[*params]
	if !found || len(key) != expectedKeySize {
		var err error
		derivedKey := make([]byte, expectedKeySize)
		s2k, err := params.Function()
		if err != nil {
			return nil, err
		}
		s2k(derivedKey, passphrase)
		(*c)[*params] = key
		return derivedKey, nil
	}
	return key, nil
}
