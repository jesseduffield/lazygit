package sha1cd

import "hash"

type CollisionResistantHash interface {
	// CollisionResistantSum extends on Sum by returning an additional boolean
	// which indicates whether a collision was found during the hashing process.
	CollisionResistantSum(b []byte) ([]byte, bool)

	hash.Hash
}
