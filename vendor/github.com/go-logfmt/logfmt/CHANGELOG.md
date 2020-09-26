# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0] - 2020-01-03

### Changed
- Remove the dependency on github.com/kr/logfmt by [@ChrisHines]
- Move fuzz code to github.com/go-logfmt/fuzzlogfmt by [@ChrisHines]

## [0.4.0] - 2018-11-21

### Added
- Go module support by [@ChrisHines]
- CHANGELOG by [@ChrisHines]

### Changed
- Drop invalid runes from keys instead of returning ErrInvalidKey by [@ChrisHines]
- On panic while printing, attempt to print panic value by [@bboreham]

## [0.3.0] - 2016-11-15
### Added
- Pool buffers for quoted strings and byte slices by [@nussjustin]
### Fixed
- Fuzz fix, quote invalid UTF-8 values by [@judwhite]

## [0.2.0] - 2016-05-08
### Added
- Encoder.EncodeKeyvals by [@ChrisHines]

## [0.1.0] - 2016-03-28
### Added
- Encoder by [@ChrisHines]
- Decoder by [@ChrisHines]
- MarshalKeyvals by [@ChrisHines]

[0.5.0]: https://github.com/go-logfmt/logfmt/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/go-logfmt/logfmt/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/go-logfmt/logfmt/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/go-logfmt/logfmt/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/go-logfmt/logfmt/commits/v0.1.0

[@ChrisHines]: https://github.com/ChrisHines
[@bboreham]: https://github.com/bboreham
[@judwhite]: https://github.com/judwhite
[@nussjustin]: https://github.com/nussjustin
