# Release Notes v0.5

This release supports multiple xz streams in xz files. The older release
couldn't support those files and there are files (linux kernel
tarballs) that couldn't be decompressed by the code.

The API has changed. Types ReaderConfig, WriterConfig, etc. are
introduced to provide parameters to the readers and writers in the
packages xz and lzma. The old API had multiple inconsistent mechanisms.
Making NewReader or NewWriter a method of the Config types provides more
clarity then the old NewReaderParams and NewWriterParams.

The compression ratio and performance has been improved. An experimental
Binary Tree Matcher has been added, but performance and compression
ratio is poor. It's is not recommended.

