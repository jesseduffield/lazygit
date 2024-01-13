package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Povray lexer.
var Povray = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "POVRay",
		Aliases:   []string{"pov"},
		Filenames: []string{"*.pov", "*.inc"},
		MimeTypes: []string{"text/x-povray"},
	},
	povrayRules,
))

func povrayRules() Rules {
	return Rules{
		"root": {
			{`/\*[\w\W]*?\*/`, CommentMultiline, nil},
			{`//.*\n`, CommentSingle, nil},
			{`(?s)"(?:\\.|[^"\\])+"`, LiteralStringDouble, nil},
			{Words(`#`, `\b`, `break`, `case`, `debug`, `declare`, `default`, `define`, `else`, `elseif`, `end`, `error`, `fclose`, `fopen`, `for`, `if`, `ifdef`, `ifndef`, `include`, `local`, `macro`, `range`, `read`, `render`, `statistics`, `switch`, `undef`, `version`, `warning`, `while`, `write`), CommentPreproc, nil},
			{Words(`\b`, `\b`, `aa_level`, `aa_threshold`, `abs`, `acos`, `acosh`, `adaptive`, `adc_bailout`, `agate`, `agate_turb`, `all`, `alpha`, `ambient`, `ambient_light`, `angle`, `aperture`, `arc_angle`, `area_light`, `asc`, `asin`, `asinh`, `assumed_gamma`, `atan`, `atan2`, `atanh`, `atmosphere`, `atmospheric_attenuation`, `attenuating`, `average`, `background`, `black_hole`, `blue`, `blur_samples`, `bounded_by`, `box_mapping`, `bozo`, `break`, `brick`, `brick_size`, `brightness`, `brilliance`, `bumps`, `bumpy1`, `bumpy2`, `bumpy3`, `bump_map`, `bump_size`, `case`, `caustics`, `ceil`, `checker`, `chr`, `clipped_by`, `clock`, `color`, `color_map`, `colour`, `colour_map`, `component`, `composite`, `concat`, `confidence`, `conic_sweep`, `constant`, `control0`, `control1`, `cos`, `cosh`, `count`, `crackle`, `crand`, `cube`, `cubic_spline`, `cylindrical_mapping`, `debug`, `declare`, `default`, `degrees`, `dents`, `diffuse`, `direction`, `distance`, `distance_maximum`, `div`, `dust`, `dust_type`, `eccentricity`, `else`, `emitting`, `end`, `error`, `error_bound`, `exp`, `exponent`, `fade_distance`, `fade_power`, `falloff`, `falloff_angle`, `false`, `file_exists`, `filter`, `finish`, `fisheye`, `flatness`, `flip`, `floor`, `focal_point`, `fog`, `fog_alt`, `fog_offset`, `fog_type`, `frequency`, `gif`, `global_settings`, `glowing`, `gradient`, `granite`, `gray_threshold`, `green`, `halo`, `hexagon`, `hf_gray_16`, `hierarchy`, `hollow`, `hypercomplex`, `if`, `ifdef`, `iff`, `image_map`, `incidence`, `include`, `int`, `interpolate`, `inverse`, `ior`, `irid`, `irid_wavelength`, `jitter`, `lambda`, `leopard`, `linear`, `linear_spline`, `linear_sweep`, `location`, `log`, `looks_like`, `look_at`, `low_error_factor`, `mandel`, `map_type`, `marble`, `material_map`, `matrix`, `max`, `max_intersections`, `max_iteration`, `max_trace_level`, `max_value`, `metallic`, `min`, `minimum_reuse`, `mod`, `mortar`, `nearest_count`, `no`, `normal`, `normal_map`, `no_shadow`, `number_of_waves`, `octaves`, `off`, `offset`, `omega`, `omnimax`, `on`, `once`, `onion`, `open`, `orthographic`, `panoramic`, `pattern1`, `pattern2`, `pattern3`, `perspective`, `pgm`, `phase`, `phong`, `phong_size`, `pi`, `pigment`, `pigment_map`, `planar_mapping`, `png`, `point_at`, `pot`, `pow`, `ppm`, `precision`, `pwr`, `quadratic_spline`, `quaternion`, `quick_color`, `quick_colour`, `quilted`, `radial`, `radians`, `radiosity`, `radius`, `rainbow`, `ramp_wave`, `rand`, `range`, `reciprocal`, `recursion_limit`, `red`, `reflection`, `refraction`, `render`, `repeat`, `rgb`, `rgbf`, `rgbft`, `rgbt`, `right`, `ripples`, `rotate`, `roughness`, `samples`, `scale`, `scallop_wave`, `scattering`, `seed`, `shadowless`, `sin`, `sine_wave`, `sinh`, `sky`, `sky_sphere`, `slice`, `slope_map`, `smooth`, `specular`, `spherical_mapping`, `spiral`, `spiral1`, `spiral2`, `spotlight`, `spotted`, `sqr`, `sqrt`, `statistics`, `str`, `strcmp`, `strength`, `strlen`, `strlwr`, `strupr`, `sturm`, `substr`, `switch`, `sys`, `t`, `tan`, `tanh`, `test_camera_1`, `test_camera_2`, `test_camera_3`, `test_camera_4`, `texture`, `texture_map`, `tga`, `thickness`, `threshold`, `tightness`, `tile2`, `tiles`, `track`, `transform`, `translate`, `transmit`, `triangle_wave`, `true`, `ttf`, `turbulence`, `turb_depth`, `type`, `ultra_wide_angle`, `up`, `use_color`, `use_colour`, `use_index`, `u_steps`, `val`, `variance`, `vaxis_rotate`, `vcross`, `vdot`, `version`, `vlength`, `vnormalize`, `volume_object`, `volume_rendered`, `vol_with_light`, `vrotate`, `v_steps`, `warning`, `warp`, `water_level`, `waves`, `while`, `width`, `wood`, `wrinkles`, `yes`), Keyword, nil},
			{Words(``, `\b`, `bicubic_patch`, `blob`, `box`, `camera`, `cone`, `cubic`, `cylinder`, `difference`, `disc`, `height_field`, `intersection`, `julia_fractal`, `lathe`, `light_source`, `merge`, `mesh`, `object`, `plane`, `poly`, `polygon`, `prism`, `quadric`, `quartic`, `smooth_triangle`, `sor`, `sphere`, `superellipsoid`, `text`, `torus`, `triangle`, `union`), NameBuiltin, nil},
			{`[\[\](){}<>;,]`, Punctuation, nil},
			{`[-+*/=]`, Operator, nil},
			{`\b(x|y|z|u|v)\b`, NameBuiltinPseudo, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
			{`[0-9]+\.[0-9]*`, LiteralNumberFloat, nil},
			{`\.[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`\s+`, Text, nil},
		},
	}
}
