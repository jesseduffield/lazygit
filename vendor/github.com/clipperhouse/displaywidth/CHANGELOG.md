# Changelog

## [0.11.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.10.0...v0.11.0)

### Added
- New `ControlSequences8Bit` option to treat 8-bit ECMA-48 (C1) escape sequences as zero-width. (#22)

### Changed
- Upgraded uax29 dependency to v2.7.0 for 8-bit escape sequence support in the grapheme iterator.
- Truncation now validates that preserved trailing escape sequences are zero-width, preventing edge cases where non-zero-width sequences could leak into output.

### Note
- `ControlSequences8Bit` is deliberately ignored by `TruncateString` and `TruncateBytes`, because C1 byte values (0x80–0x9F) overlap with UTF-8 multi-byte encoding.

## [0.10.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.9.0...v0.10.0)

### Added
- New `ControlSequences` option to treat ECMA-48/ANSI escape sequences as zero-width. (#20)
- `TruncateString` and `TruncateBytes` now preserve trailing ANSI escape sequences (such as SGR resets) when `ControlSequences` is true, preventing color bleed in terminal output.

### Changed
- Removed `stringish` dependency; generic type constraints are now inline `~string | []byte`.
- Upgraded uax29 dependency to v2.6.0 for ANSI escape sequence support in the grapheme iterator.

## [0.9.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.8.0...v0.9.0)

### Changed
- Unicode 17 support: East Asian Width and emoji data updated to Unicode 17.0.0. (#18)
- Upgraded uax29 dependency to v2.5.0 (Unicode 17 grapheme segmentation).

## [0.8.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.7.0...v0.8.0)

### Changed
- Performance: ASCII fast path that applies to any run of printable
  ASCII. 2x-10x faster for ASCII text vs v0.7.0. (#16)
- Upgraded uax29 dependency to v2.4.0 for Unicode 16 support. Text that includes
  Indic_Conjunct_Break may segment differently (and more correctly). (#15)

## [0.7.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.6.2...v0.7.0)

### Added
- New `TruncateString` and `TruncateBytes` methods to truncate strings to a
  maximum display width, with optional tail (like an ellipsis). (#13)

## [0.6.2]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.6.1...v0.6.2)

### Changed
- Internal: reduced property categories for simpler trie.

## [0.6.1]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.6.0...v0.6.1)

### Changed
- Perf improvements: replaced the ASCII lookup table with a simple
  function. A bit more cache-friendly. More inlining.
- Bug fix: single regional indicators are now treated as width 2, since that
  is what actual terminals do.

## [0.6.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.5.0...v0.6.0)

### Added
- New `StringGraphemes` and `BytesGraphemes` methods, for iterating over the
widths of grapheme clusters.

### Changed
- Fast ASCII lookups

## [0.5.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.4.1...v0.5.0)

### Added
- Unicode 16 support
- Improved emoji presentation handling per Unicode TR51

### Changed
- Corrected VS15 (U+FE0E) handling: now preserves base character width (no-op) per Unicode TR51
- Performance optimizations: reduced property lookups

### Fixed
- VS15 variation selector now correctly preserves base character width instead of forcing width 1

## [0.4.1]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.4.0...v0.4.1)

### Changed
- Updated uax29 dependency
- Improved flag handling

## [0.4.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.3.1...v0.4.0)

### Added
- Support for variation selectors (VS15, VS16) and regional indicator pairs (flags)

## [0.3.1]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.3.0...v0.3.1)

### Added
- Fuzz testing support

### Changed
- Updated stringish dependency

## [0.3.0]

[Compare](https://github.com/clipperhouse/displaywidth/compare/v0.2.0...v0.3.0)

### Changed
- Dropped compatibility with go-runewidth
- Trie implementation cleanup
