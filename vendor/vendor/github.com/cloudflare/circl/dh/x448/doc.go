/*
Package x448 provides Diffie-Hellman functions as specified in RFC-7748.

Validation of public keys.

The Diffie-Hellman function, as described in RFC-7748 [1], works for any
public key. However, if a different protocol requires contributory
behaviour [2,3], then the public keys must be validated against low-order
points [3,4]. To do that, the Shared function performs this validation
internally and returns false when the public key is invalid (i.e., it
is a low-order point).

References:
  - [1] RFC7748 by Langley, Hamburg, Turner (https://rfc-editor.org/rfc/rfc7748.txt)
  - [2] Curve25519 by Bernstein (https://cr.yp.to/ecdh.html)
  - [3] Bernstein (https://cr.yp.to/ecdh.html#validate)
  - [4] Cremers&Jackson (https://eprint.iacr.org/2019/526)
*/
package x448
