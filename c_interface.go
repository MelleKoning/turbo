package turbo

/*
#cgo CFLAGS: -I/opt/libjpeg-turbo/include
#cgo LDFLAGS: -L/opt/libjpeg-turbo/lib64 -lturbojpeg -Wl,-rpath=/opt/libjpeg-turbo/lib64

#include <stdlib.h>
#include <stdio.h>
#include <turbojpeg.h>

// Function to compress an image using libjpeg-turbo
unsigned char* compressImage(unsigned char* image, int width, int height, unsigned long* jpegSize) {
    tjhandle handle = tjInitCompress();
    unsigned char* jpegBuf = NULL;

    tjCompress2(handle, image, width, 0, height, TJPF_RGB, &jpegBuf, jpegSize, TJSAMP_444, 80, TJFLAG_FASTDCT);

    tjDestroy(handle);
    return jpegBuf;
}
*/
import "C"
import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"io"
	"os"
	"unsafe"
)

func CompressImage(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Convert image to RGB format
	rgbImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgbImg, rgbImg.Bounds(), img, bounds.Min, draw.Src)

	// Get image data
	imgData := rgbImg.Pix

	// Call C function to compress image
	cImgData := C.CBytes(imgData)
	defer C.free(cImgData)

	var size C.ulong
	jpegData := C.compressImage((*C.uchar)(cImgData), C.int(width), C.int(height), &size)

	defer C.free(unsafe.Pointer(jpegData))

	// Convert C array to Go slice
	jpegBytes := C.GoBytes(unsafe.Pointer(jpegData), C.int(size))

	return jpegBytes, nil
}

func LoadImage(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := DecodeJPEG(file)
	if err != nil {
		panic(err)
	}

	return img
}

func SaveImage(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func DecodeJPEG(r io.Reader) (image.Image, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	cData := C.CBytes(buf.Bytes())
	defer C.free(cData)

	cInfo := C.tjInitDecompress()
	defer C.tjDestroy(cInfo)

	var width C.int
	var height C.int
	var subsamp C.int

	if C.tjDecompressHeader2(cInfo, (*C.uchar)(cData), C.ulong(len(buf.Bytes())), &width, &height, &subsamp) != 0 {
		return nil, errors.New("failed to decompress JPEG header")
	}

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	if C.tjDecompress2(cInfo, (*C.uchar)(cData), C.ulong(len(buf.Bytes())), (*C.uchar)(unsafe.Pointer(&img.Pix[0])), width, 0, height, C.TJPF_RGB, 0) != 0 {
		return nil, errors.New("failed to decompress JPEG image")
	}

	return img, nil
}
