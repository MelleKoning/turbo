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

## libjpeg-turbo official on Ubuntu

As an example a few steps that can help on Ubuntu (and other debian versions)

1. Step 1: Verifying the installation.

`apt search libjpeg-turbo-official`

should show that the library was installed and can be found by the system. The output should be something like this:

```bash
Sorting... Done
Full Text Search... Done
libjpeg-turbo-official/any,now 3.0.2-20240124 amd64 [installed]
  A SIMD-accelerated JPEG codec that provides both the libjpeg and TurboJPEG APIs

libjpeg-turbo-official32/any 3.0.2-20240124 amd64
  A SIMD-accelerated JPEG codec that provides both the libjpeg and TurboJPEG APIs
```

2. Step 2: finding the `turbolib.h` include file:

Command to search for the `tubojpeg.h` file on the system:

`dpkg -L libjpeg-turbo-official | grep turbojpeg.h`

will output something like this:

`/opt/libjpeg-turbo/include/turbojpeg.h`

This provides us the location of the value for the needed CFLAGS: `#cgo CFLAGS: -I/opt/libjpeg-turbo/include`

3. Step 3: Finding the actual used library and the include file

```bash
pkg-config --cflags --libs libturbojpeg
```

Should give you the include options you need in the code:

`-I/opt/libjpeg-turbo/include -L/opt/libjpeg-turbo/lib64 -lturbojpeg`

The -I option is for the golang CGO `CFLAGS` and the -L option is for the `LDFLAGS`.

The header of `c_interface.go` will then have to look something like this:

```golang
/*
#cgo CFLAGS: -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L/opt/libjpeg-turbo/lib64 -lturbojpeg
#include <turbojpeg.h>
*/
```

When running the code with `go test ./...`, you might still bump into the following error:

```bash
$go test ./...
/tmp/go-build3950251331/b001/turbo.test: error while loading shared libraries: libturbojpeg.so.0: cannot open shared object file: No such file or directory
FAIL    github.com/bmharper/turbo       0.000s
```

This is because golang should know where to find the linked libraries. You should add the earlier found folder to the LD_LIBRARY_PATH as follows:

4. Step 4: Letting golang know the library path

```bash
export LD_LIBRARY_PATH=/opt/libjpeg-turbo/lib64:$LD_LIBRARY_PATH
```

When running, all should be now fine.

```bash
$go test ./... -v
=== RUN   TestCompress
    turbo_test.go:43: Encode return: 48796, <nil>
    turbo_test.go:45: Decode return: 300 x 200, 1200, 240000, <nil>
--- PASS: TestCompress (0.00s)
PASS
```
Have fun using libjpeg-turbo in golang!









