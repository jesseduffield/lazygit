# Release Notes v0.4

This release support the compression and decompression to xz files. Note
that only the LZMA filter is supported for the xz format, but this seems
to be the standard setup anyway.

The performance and compression ration is not good compared to the xz
tool written in C. But optimization has not been the target of this
release.

A gxz binary is included that supports the compression and decompression of
xz files.
