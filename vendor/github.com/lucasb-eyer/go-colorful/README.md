go-colorful
===========

[![go reportcard](https://goreportcard.com/badge/github.com/lucasb-eyer/go-colorful)](https://goreportcard.com/report/github.com/lucasb-eyer/go-colorful)

A library for playing with colors in Go. Supports Go 1.13 onwards.

Why?
====
I love games. I make games. I love detail and I get lost in detail.
One such detail popped up during the development of [Memory Which Does Not Suck](https://github.com/lucasb-eyer/mwdns/),
when we wanted the server to assign the players random colors. Sometimes
two players got very similar colors, which bugged me. The very same evening,
[I want hue](http://tools.medialab.sciences-po.fr/iwanthue/) was the top post
on HackerNews' frontpage and showed me how to Do It Right™. Last but not
least, there was no library for handling color spaces available in go. Colorful
does just that and implements Go's `color.Color` interface.

What?
=====
Go-Colorful stores colors in RGB and provides methods from converting these to various color-spaces. Currently supported colorspaces are:

- **RGB:** All three of Red, Green and Blue in [0..1].
- **HSL:** Hue in [0..360], Saturation and Luminance in [0..1]. For legacy reasons; please forget that it exists.
- **HSV:** Hue in [0..360], Saturation and Value in [0..1]. You're better off using HCL, see below.
- **Hex RGB:** The "internet" color format, as in #FF00FF.
- **Linear RGB:** See [gamma correct rendering](http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/).
- **CIE-XYZ:** CIE's standard color space, almost in [0..1].
- **CIE-xyY:** encodes chromacity in x and y and luminance in Y, all in [0..1]
- **CIE-L\*a\*b\*:** A *perceptually uniform* color space, i.e. distances are meaningful. L\* in [0..1] and a\*, b\* almost in [-1..1].
- **CIE-L\*u\*v\*:** Very similar to CIE-L\*a\*b\*, there is [no consensus](http://en.wikipedia.org/wiki/CIELUV#Historical_background) on which one is "better".
- **CIE-L\*C\*h° (HCL):** This is generally the [most useful](http://vis4.net/blog/posts/avoid-equidistant-hsv-colors/) one; CIE-L\*a\*b\* space in polar coordinates, i.e. a *better* HSV. H° is in [0..360], C\* almost in [-1..1] and L\* as in CIE-L\*a\*b\*.
- **CIE LCh(uv):** Called `LuvLCh` in code, this is a cylindrical transformation of the CIE-L\*u\*v\* color space. Like HCL above: H° is in [0..360], C\* almost in [-1..1] and L\* as in CIE-L\*u\*v\*.
- **HSLuv:** The better alternative to HSL, see [here](https://www.hsluv.org/) and [here](https://www.kuon.ch/post/2020-03-08-hsluv/). Hue in [0..360], Saturation and Luminance in [0..1].
- **HPLuv:** A variant of HSLuv. The color space is smoother, but only pastel colors can be included. Because the valid colors are limited, it's easy to get invalid Saturation values way above 1.0, indicating the color can't be represented in HPLuv beccause it's not pastel.

For the colorspaces where it makes sense (XYZ, Lab, Luv, HCl), the
[D65](http://en.wikipedia.org/wiki/Illuminant_D65) is used as reference white
by default but methods for using your own reference white are provided.

A coordinate being *almost in* a range means that generally it is, but for very
bright colors and depending on the reference white, it might overflow this
range slightly. For example, C\* of #0000ff is 1.338.

Unit-tests are provided.

Nice, but what's it useful for?
-------------------------------

- Converting color spaces. Some people like to do that.
- Blending (interpolating) between colors in a "natural" look by using the right colorspace.
- Generating random colors under some constraints (e.g. colors of the same shade, or shades of one color.)
- Generating gorgeous random palettes with distinct colors of a same temperature.

What not (yet)?
===============
There are a few features which are currently missing and might be useful.
I just haven't implemented them yet because I didn't have the need for it.
Pull requests welcome.

- Sorting colors (potentially using above mentioned distances)

So which colorspace should I use?
=================================
It depends on what you want to do. I think the folks from *I want hue* are
on-spot when they say that RGB fits to how *screens produce* color, CIE L\*a\*b\*
fits how *humans perceive* color and HCL fits how *humans think* colors.

Whenever you'd use HSV, rather go for CIE-L\*C\*h°. for fixed lightness L\* and
chroma C\* values, the hue angle h° rotates through colors of the same
perceived brightness and intensity.

How?
====

### Installing
Installing the library is as easy as

```bash
$ go get github.com/lucasb-eyer/go-colorful
```

The package can then be used through an

```go
import "github.com/lucasb-eyer/go-colorful"
```

### Basic usage

Create a beautiful blue color using different source space:

```go
// Any of the following should be the same
c := colorful.Color{0.313725, 0.478431, 0.721569}
c, err := colorful.Hex("#517AB8")
if err != nil {
    log.Fatal(err)
}
c = colorful.Hsv(216.0, 0.56, 0.722)
c = colorful.Xyz(0.189165, 0.190837, 0.480248)
c = colorful.Xyy(0.219895, 0.221839, 0.190837)
c = colorful.Lab(0.507850, 0.040585,-0.370945)
c = colorful.Luv(0.507849,-0.194172,-0.567924)
c = colorful.Hcl(276.2440, 0.373160, 0.507849)
fmt.Printf("RGB values: %v, %v, %v", c.R, c.G, c.B)
```

And then converting this color back into various color spaces:

```go
hex := c.Hex()
h, s, v := c.Hsv()
x, y, z := c.Xyz()
x, y, Y := c.Xyy()
l, a, b := c.Lab()
l, u, v := c.Luv()
h, c, l := c.Hcl()
```

Note that, because of Go's unfortunate choice of requiring an initial uppercase,
the name of the functions relating to the xyY space are just off. If you have
any good suggestion, please open an issue. (I don't consider XyY good.)

### The `color.Color` interface
Because a `colorful.Color` implements Go's `color.Color` interface (found in the
`image/color` package), it can be used anywhere that expects a `color.Color`.

Furthermore, you can convert anything that implements the `color.Color` interface
into a `colorful.Color` using the `MakeColor` function:

```go
c, ok := colorful.MakeColor(color.Gray16{12345})
```

**Caveat:** Be aware that this latter conversion (using `MakeColor`) hits a
corner-case when alpha is exactly zero. Because `color.Color` uses pre-multiplied
alpha colors, this means the RGB values are lost (set to 0) and it's impossible
to recover them. In such a case `MakeColor` will return `false` as its second value.

### Comparing colors
In the RGB color space, the Euclidian distance between colors *doesn't* correspond
to visual/perceptual distance. This means that two pairs of colors which have the
same distance in RGB space can look much further apart. This is fixed by the
CIE-L\*a\*b\*, CIE-L\*u\*v\* and CIE-L\*C\*h° color spaces.
Thus you should only compare colors in any of these space.
(Note that the distance in CIE-L\*a\*b\* and CIE-L\*C\*h° are the same, since it's the same space but in cylindrical coordinates)

![Color distance comparison](doc/colordist/colordist.png)

The two colors shown on the top look much more different than the two shown on
the bottom. Still, in RGB space, their distance is the same.
Here is a little example program which shows the distances between the top two
and bottom two colors in RGB, CIE-L\*a\*b\* and CIE-L\*u\*v\* space. You can find it in `doc/colordist/colordist.go`.

```go
package main

import "fmt"
import "github.com/lucasb-eyer/go-colorful"

func main() {
	c1a := colorful.Color{150.0 / 255.0, 10.0 / 255.0, 150.0 / 255.0}
	c1b := colorful.Color{53.0 / 255.0, 10.0 / 255.0, 150.0 / 255.0}
	c2a := colorful.Color{10.0 / 255.0, 150.0 / 255.0, 50.0 / 255.0}
	c2b := colorful.Color{99.9 / 255.0, 150.0 / 255.0, 10.0 / 255.0}

	fmt.Printf("DistanceRgb:       c1: %v\tand c2: %v\n", c1a.DistanceRgb(c1b), c2a.DistanceRgb(c2b))
	fmt.Printf("DistanceLab:       c1: %v\tand c2: %v\n", c1a.DistanceLab(c1b), c2a.DistanceLab(c2b))
	fmt.Printf("DistanceLuv:       c1: %v\tand c2: %v\n", c1a.DistanceLuv(c1b), c2a.DistanceLuv(c2b))
	fmt.Printf("DistanceCIE76:     c1: %v\tand c2: %v\n", c1a.DistanceCIE76(c1b), c2a.DistanceCIE76(c2b))
	fmt.Printf("DistanceCIE94:     c1: %v\tand c2: %v\n", c1a.DistanceCIE94(c1b), c2a.DistanceCIE94(c2b))
	fmt.Printf("DistanceCIEDE2000: c1: %v\tand c2: %v\n", c1a.DistanceCIEDE2000(c1b), c2a.DistanceCIEDE2000(c2b))
}
```

Running the above program shows that you should always prefer any of the CIE distances:

```bash
$ go run colordist.go
DistanceRgb:       c1: 0.3803921568627451	and c2: 0.3858713931171159
DistanceLab:       c1: 0.32048458312798056	and c2: 0.24397151758565272
DistanceLuv:       c1: 0.5134369614199698	and c2: 0.2568692839860636
DistanceCIE76:     c1: 0.32048458312798056	and c2: 0.24397151758565272
DistanceCIE94:     c1: 0.19799168128511324	and c2: 0.12207136371167401
DistanceCIEDE2000: c1: 0.17274551120971166	and c2: 0.10665210031428465
```

It also shows that `DistanceLab` is more formally known as `DistanceCIE76` and
has been superseded by the slightly more accurate, but much more expensive
`DistanceCIE94` and `DistanceCIEDE2000`.

Note that `AlmostEqualRgb` is provided mainly for (unit-)testing purposes. Use
it only if you really know what you're doing. It will eat your cat.

### Blending colors
Blending is highly connected to distance, since it basically "walks through" the
colorspace thus, if the colorspace maps distances well, the walk is "smooth".

Colorful comes with blending functions in RGB, HSV and any of the LAB spaces.
Of course, you'd rather want to use the blending functions of the LAB spaces since
these spaces map distances well but, just in case, here is an example showing
you how the blendings (`#fdffcc` to `#242a42`) are done in the various spaces:

![Blending colors in different spaces.](doc/colorblend/colorblend.png)

What you see is that HSV is really bad: it adds some green, which is not present
in the original colors at all! RGB is much better, but it stays light a little
too long. LUV and LAB both hit the right lightness but LAB has a little more
color. HCL works in the same vein as HSV (both cylindrical interpolations) but
it does it right in that there is no green appearing and the lighthness changes
in a linear manner.

While this seems all good, you need to know one thing: When interpolating in any
of the CIE color spaces, you might get invalid RGB colors! This is important if
the starting and ending colors are user-input or random. An example of where this
happens is when blending between `#eeef61` and `#1e3140`:

![Invalid RGB colors may crop up when blending in CIE spaces.](doc/colorblend/invalid.png)

You can test whether a color is a valid RGB color by calling the `IsValid` method
and indeed, calling IsValid will return false for the redish colors on the bottom.
One way to "fix" this is to get a valid color close to the invalid one by calling
`Clamped`, which always returns a nearby valid color. Doing this, we get the
following result, which is satisfactory:

![Fixing invalid RGB colors by clamping them to the valid range.](doc/colorblend/clamped.png)

The following is the code creating the above three images; it can be found in `doc/colorblend/colorblend.go`

```go
package main

import "fmt"
import "github.com/lucasb-eyer/go-colorful"
import "image"
import "image/draw"
import "image/png"
import "os"

func main() {
    blocks := 10
    blockw := 40
    img := image.NewRGBA(image.Rect(0,0,blocks*blockw,200))

    c1, _ := colorful.Hex("#fdffcc")
    c2, _ := colorful.Hex("#242a42")

    // Use these colors to get invalid RGB in the gradient.
    //c1, _ := colorful.Hex("#EEEF61")
    //c2, _ := colorful.Hex("#1E3140")

    for i := 0 ; i < blocks ; i++ {
        draw.Draw(img, image.Rect(i*blockw,  0,(i+1)*blockw, 40), &image.Uniform{c1.BlendHsv(c2, float64(i)/float64(blocks-1))}, image.Point{}, draw.Src)
        draw.Draw(img, image.Rect(i*blockw, 40,(i+1)*blockw, 80), &image.Uniform{c1.BlendLuv(c2, float64(i)/float64(blocks-1))}, image.Point{}, draw.Src)
        draw.Draw(img, image.Rect(i*blockw, 80,(i+1)*blockw,120), &image.Uniform{c1.BlendRgb(c2, float64(i)/float64(blocks-1))}, image.Point{}, draw.Src)
        draw.Draw(img, image.Rect(i*blockw,120,(i+1)*blockw,160), &image.Uniform{c1.BlendLab(c2, float64(i)/float64(blocks-1))}, image.Point{}, draw.Src)
        draw.Draw(img, image.Rect(i*blockw,160,(i+1)*blockw,200), &image.Uniform{c1.BlendHcl(c2, float64(i)/float64(blocks-1))}, image.Point{}, draw.Src)

        // This can be used to "fix" invalid colors in the gradient.
        //draw.Draw(img, image.Rect(i*blockw,160,(i+1)*blockw,200), &image.Uniform{c1.BlendHcl(c2, float64(i)/float64(blocks-1)).Clamped()}, image.Point{}, draw.Src)
    }

    toimg, err := os.Create("colorblend.png")
    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    defer toimg.Close()

    png.Encode(toimg, img)
}
```

#### Generating color gradients
A very common reason to blend colors is creating gradients. There is an example
program in [doc/gradientgen.go](doc/gradientgen/gradientgen.go); it doesn't use any API
which hasn't been used in the previous example code, so I won't bother pasting
the code in here. Just look at that gorgeous gradient it generated in HCL space:

!["Spectral" colorbrewer gradient in HCL space.](doc/gradientgen/gradientgen.png)

### Getting random colors
It is sometimes necessary to generate random colors. You could simply do this
on your own by generating colors with random values. By restricting the random
values to a range smaller than [0..1] and using a space such as CIE-H\*C\*l° or
HSV, you can generate both random shades of a color or random colors of a
lightness:

```go
random_blue := colorful.Hcl(180.0+rand.Float64()*50.0, 0.2+rand.Float64()*0.8, 0.3+rand.Float64()*0.7)
random_dark := colorful.Hcl(rand.Float64()*360.0, rand.Float64(), rand.Float64()*0.4)
random_light := colorful.Hcl(rand.Float64()*360.0, rand.Float64(), 0.6+rand.Float64()*0.4)
```

Since getting random "warm" and "happy" colors is quite a common task, there
are some helper functions:

```go
colorful.WarmColor()
colorful.HappyColor()
colorful.FastWarmColor()
colorful.FastHappyColor()
```

The ones prefixed by `Fast` are faster but less coherent since they use the HSV
space as opposed to the regular ones which use CIE-L\*C\*h° space. The
following picture shows the warm colors in the top two rows and happy colors
in the bottom two rows. Within these, the first is the regular one and the
second is the fast one.

![Warm, fast warm, happy and fast happy random colors, respectively.](doc/colorgens/colorgens.png)

Don't forget to initialize the random seed! You can see the code used for
generating this picture in `doc/colorgens/colorgens.go`.

### Getting random palettes
As soon as you need to generate more than one random color, you probably want
them to be distinguishible. Playing against an opponent which has almost the
same blue as I do is not fun. This is where random palettes can help.

These palettes are generated using an algorithm which ensures that all colors
on the palette are as distinguishible as possible. Again, there is a `Fast`
method which works in HSV and is less perceptually uniform and a non-`Fast`
method which works in CIE spaces. For more theory on `SoftPalette`, check out
[I want hue](http://tools.medialab.sciences-po.fr/iwanthue/theory.php). Yet
again, there is a `Happy` and a `Warm` version, which do what you expect, but
now there is an additional `Soft` version, which is more configurable: you can
give a constraint on the color space in order to get colors within a certain *feel*.

Let's start with the simple methods first, all they take is the amount of
colors to generate, which could, for example, be the player count. They return
an array of `colorful.Color` objects:

```go
pal1, err1 := colorful.WarmPalette(10)
pal2 := colorful.FastWarmPalette(10)
pal3, err3 := colorful.HappyPalette(10)
pal4 := colorful.FastHappyPalette(10)
pal5, err5 := colorful.SoftPalette(10)
```

Note that the non-fast methods *may* fail if you ask for way too many colors.
Let's move on to the advanced one, namely `SoftPaletteEx`. Besides the color
count, this function takes a `SoftPaletteSettings` object as argument. The
interesting part here is its `CheckColor` member, which is a boolean function
taking three floating points as arguments: `l`, `a` and `b`. This function
should return `true` for colors which lie within the region you want and `false`
otherwise. The other members are `Iteration`, which should be within [5..100]
where higher means slower but more exact palette, and `ManySamples` which you
should set to `true` in case your `CheckColor` constraint rejects a large part
of the color space.

For example, to create a palette of 10 brownish colors, you'd call it like this:

```go
func isbrowny(l, a, b float64) bool {
    h, c, L := colorful.LabToHcl(l, a, b)
    return 10.0 < h && h < 50.0 && 0.1 < c && c < 0.5 && L < 0.5
}
// Since the above function is pretty restrictive, we set ManySamples to true.
brownies := colorful.SoftPaletteEx(10, colorful.SoftPaletteSettings{isbrowny, 50, true})
```

The following picture shows the palettes generated by all of these methods
(sourcecode in `doc/palettegens/palettegens.go`), in the order they were presented, i.e.
from top to bottom: `Warm`, `FastWarm`, `Happy`, `FastHappy`, `Soft`,
`SoftEx(isbrowny)`. All of them contain some randomness, so YMMV.

![All example palettes](doc/palettegens/palettegens.png)

Again, the code used for generating the above image is available as [doc/palettegens/palettegens.go](https://github.com/lucasb-eyer/go-colorful/blob/master/doc/palettegens/palettegens.go).

### Sorting colors
TODO: Sort using dist fn.

### Using linear RGB for computations
There are two methods for transforming RGB<->Linear RGB: a fast and almost precise one,
and a slow and precise one.

```go
r, g, b := colorful.Hex("#FF0000").FastLinearRgb()
```

TODO: describe some more.

### Want to use some other reference point?

```go
c := colorful.LabWhiteRef(0.507850, 0.040585,-0.370945, colorful.D50)
l, a, b := c.LabWhiteRef(colorful.D50)
```

### Reading and writing colors from databases

The type `HexColor` makes it easy to store colors as strings in a database. It
implements the [https://godoc.org/database/sql#Scanner](database/sql.Scanner)
and [database/sql/driver.Value](https://godoc.org/database/sql/driver.Value)
interfaces which provide automatic type conversion.

Example:

```go
var hc HexColor
_, err := db.QueryRow("SELECT '#ff0000';").Scan(&hc)
// hc == HexColor{R: 1, G: 0, B: 0}; err == nil
```

FAQ
===

### Q: I get all f!@#ed up values! Your library sucks!
A: You probably provided values in the wrong range. For example, RGB values are
expected to reside between 0 and 1, *not* between 0 and 255. Normalize your colors.

### Q: Lab/Luv/HCl seem broken! Your library sucks!
They look like this:

<img height="150" src="https://user-images.githubusercontent.com/3779568/28646900-6548040c-7264-11e7-8f12-81097a97c260.png">

A: You're likely trying to generate and display colors that can't be represented by RGB,
and thus monitors. When you're trying to convert, say, `HCL(190.0, 1.0, 1.0).RGB255()`,
you're asking for RGB values of `(-2105.254  300.680  286.185)`, which clearly don't exist,
and the `RGB255` function just casts these numbers to `uint8`, creating wrap-around and
what looks like a completely broken gradient. What you want to do, is either use more
reasonable values of colors which actually exist in RGB, or just `Clamp()` the resulting
color to its nearest existing one, living with the consequences:
`HCL(190.0, 1.0, 1.0).Clamp().RGB255()`. It will look something like this:

<img height="150" src="https://user-images.githubusercontent.com/1476029/29596343-9a8c62c6-8771-11e7-9026-b8eb8852cc4a.png">

[Here's an issue going in-depth about this](https://github.com/lucasb-eyer/go-colorful/issues/14),
as well as [my answer](https://github.com/lucasb-eyer/go-colorful/issues/14#issuecomment-324205385),
both with code and pretty pictures. Also note that this was somewhat covered above in the
["Blending colors" section](https://github.com/lucasb-eyer/go-colorful#blending-colors).

### Q: In a tight loop, conversion to Lab/Luv/HCl/... are slooooow!
A: Yes, they are.
This library aims for correctness, readability, and modularity; it wasn't written with speed in mind.
A large part of the slowness comes from these conversions going through `LinearRgb` which uses powers.
I implemented a fast approximation to `LinearRgb` called `FastLinearRgb` by using Taylor approximations.
The approximation is roughly 5x faster and precise up to roughly 0.5%,
the major caveat being that if the input values are outside the range 0-1, accuracy drops dramatically.
You can use these in your conversions as follows:

```go
col := // Get your color somehow
l, a, b := XyzToLab(LinearRgbToXyz(col.LinearRgb()))
```

If you need faster versions of `Distance*` and `Blend*` that make use of this fast approximation,
feel free to implement them and open a pull-request, I'll happily accept.

The derivation of these functions can be followed in [this Jupyter notebook](doc/LinearRGB Approximations.ipynb).
Here's the main figure showing the approximation quality:

![approximation quality](doc/approx-quality.png)

More speed could be gained by using SIMD instructions in many places.
You can also get more speed for specific conversions by approximating the full conversion function,
but that is outside the scope of this library.
Thanks to [@ZirconiumX](https://github.com/ZirconiumX) for starting this investigation,
see [issue #18](https://github.com/lucasb-eyer/go-colorful/issues/18) for details.

### Q: Why would `MakeColor` ever fail!?
A: `MakeColor` fails when the alpha channel is zero. In that case, the
conversion is undefined. See [issue 21](https://github.com/lucasb-eyer/go-colorful/issues/21)
as well as the short caveat note in the ["The `color.Color` interface"](README.md#the-colorcolor-interface)
section above.

Who?
====

This library was developed by Lucas Beyer with contributions from
Bastien Dejean (@baskerville), Phil Kulak (@pkulak) and Christian Muehlhaeuser (@muesli).

It is now maintained by makeworld (@makeworld-the-better-one).


## License

This repo is under the MIT license, see [LICENSE](LICENSE) for details.
