# sha1cd

A Go implementation of SHA1 with counter-cryptanalysis, which detects
collision attacks. 

The `cgo/lib` code is a carbon copy of the [original code], based on
the award winning [white paper] by Marc Stevens.

The Go implementation is largely based off Go's generic sha1.
At present no SIMD optimisations have been implemented.

## Usage

`sha1cd` can be used as a drop-in replacement for `crypto/sha1`:

```golang
import "github.com/pjbgf/sha1cd"

func test(){
	data := []byte("data to be sha1 hashed")
	h := sha1cd.Sum(data)
	fmt.Printf("hash: %q\n", hex.EncodeToString(h))
}
```

To obtain information as to whether a collision was found, use the
func `CollisionResistantSum`.

```golang
import "github.com/pjbgf/sha1cd"

func test(){
	data := []byte("data to be sha1 hashed")
	h, col  := sha1cd.CollisionResistantSum(data)
	if col {
		fmt.Println("collision found!")
	}
	fmt.Printf("hash: %q", hex.EncodeToString(h))
}
```

Note that the algorithm will automatically avoid collision, by 
extending the SHA1 to 240-steps, instead of 80 when a collision
attempt is detected. Therefore, inputs that contains the unavoidable
bit conditions will yield a different hash from `sha1cd`, when compared
with results using `crypto/sha1`. Valid inputs will have matching the outputs.

## References
- https://shattered.io/
- https://github.com/cr-marcstevens/sha1collisiondetection
- https://csrc.nist.gov/Projects/Cryptographic-Algorithm-Validation-Program/Secure-Hashing#shavs

## Use of the Original Implementation
- https://github.com/git/git/commit/28dc98e343ca4eb370a29ceec4c19beac9b5c01e
- https://github.com/libgit2/libgit2/pull/4136

[original code]: https://github.com/cr-marcstevens/sha1collisiondetection
[white paper]: https://marc-stevens.nl/research/papers/C13-S.pdf
