# Issues in the XZ file format

During the development of the xz package for Go a number of issues with
the xz file format were observed. They are documented here to help a
later development of an improved format.

# xz file format

## Consistency

## General

Packets should either be constant size or should have encoded the size
in the header. File header and footer are constant size and the block
header has the size encoded. The index doesn't fulfill the criteria even
when its size is included in the footer.

## Index

The index doesn't have the size in the header. So in a stream you are
forced to read the whole index to identify its length.

The index should have made optional. This would require to remove the
index size from the footer and include its own footer in the index.

## Padding

The padding should allow direct mapping of the CRC values into memory, but it
wastes bytes bearing no information. This is certainly not optimal for a
compression format. It is argued alignment makes it faster to read and
write the checksum values, but the time spent there is much less than on
encoding and decoding itself.

## Filters for each block

Filters should have been defined in front of blocks. This way they
would not need to be repeated.

# LZMA2 

## Consistent header byte.

LZMA2 consists of a series of chunks with a header byte. The header byte
has a different format depending on whether it is an uncompressed or
compressed chunk. This has the consequence a complete reset of state,
properties and dictionary is not possible with an uncompressed chunk.
The encoder has to keep a state variable tracking a dictionary reset in
an uncompressed chunk to ensure that the flags are added in the first
compressed chunk to follow. This complicates the implementation of the
encoder and decoder.

## Dictionary capacity is not encoded

LZMA2 doesn't encode the dictionary capacity, so LZMA2 doesn't work
standalone.
