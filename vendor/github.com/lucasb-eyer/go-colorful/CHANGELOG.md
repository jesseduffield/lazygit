# Changelog
All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The format of this file is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
but only releases after v1.0.3 properly adhere to it.


## [1.2.0] - 2021-01-27
### Added
- HSLuv and HPLuv color spaces (#41, #51)
- CIE LCh(uv) color space, called `LuvLCh` in code (#51)
- JSON and envconfig serialization support for `HexColor` (#42)
- `DistanceLinearRGB` (#53)

### Fixed
- RGB to/from XYZ conversion is more accurate (#51)
- A bug in `XYZToLuvWhiteRef` that only applied to very small values was fixed (#51)
- `BlendHCL` output is clamped so that it's not invalid (#46)
- Properly documented `DistanceCIE76` (#40)
- Some small godoc fixes


## [1.0.3] - 2019-11-11
- Remove SQLMock dependency


## [1.0.2] - 2019-04-07
- Fixes SQLMock dependency


## [1.0.1] - 2019-03-24
- Adds support for Go Modules


## [1.0.0] - 2018-05-26
- API Breaking change in `MakeColor`: instead of `panic`ing when alpha is zero, it now returns a secondary, boolean return value indicating success. See [the color.Color interface](#the-colorcolor-interface) section and [this FAQ entry](#q-why-would-makecolor-ever-fail) for details.


## [0.9.0] - 2018-05-26
- Initial version number after having ignored versioning for a long time :)
