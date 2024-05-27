# Turbo (TurboJPEG Go wrapper)

This is a very thin wrapper around turbojpeg.

Why not use https://github.com/pixiv/go-libjpeg ?

Because go-libjpeg is built to use the libjpeg-compatible API of either libjpeg or
libjpeg-turbo. That API does not allow one to specify the chroma sub-sampling, so
you're unable to compress to 4:2:0. That was the sole reason for building this
package -- to enable compression to 420 chroma sub-sampling.

This wrapper links explicitly to TurboJPEG. There is no fallback to libjpeg.

### How to use

```go
import "github.com/bmharper/turbo"

func compressImage(width, height int, rgba []byte) {
	raw := turbo.Image{
		Width: width,
		Height: height,
		Stride: width * 4,
		RGBA: rgba,
	}
	params := turbo.MakeCompressParams(turbo.PixelFormatRGBA, turbo.Sampling420, 35, 0)
	jpg, err := turbo.Compress(&raw, params)
}

func decompressImage(jpg []byte) (*Image, error) {
	return turbo.Decompress(jpg)
}
```

# Running with the libjpeg-turbo official (Stable) library

The code was originally created and tested with latest development package. You can also install the official libjpeg-turbo library, which is a stable release of the libjpeg turbo package. Releasing your code base on a stable release can be preferred above relying on the dev version. The following is how to setup and run with the stable libjpeg-turbo release.

Installation instructions for installing the official stable package on a system can be found at [libjpeg-turbo](https://libjpeg-turbo.org/Downloads/YUM). For debian-ubunto this involves trusting and installing the `libjpeg-turbo/gpgkey` on your system first.

After following those instructions the following commands can help finding the location of the library on the system. These locations are important for properly pointing to the C-library in the `c_interface.go`.