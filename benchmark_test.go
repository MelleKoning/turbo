package turbo

import (
	"image/jpeg"
	"os"
	"testing"
	"time"
)

func TestCompress(t *testing.T) {
	// Load your image here
	img := LoadImage("./image.jpg")

	// Compress the image
	compressedData, err := CompressImage(img)
	if err != nil {
		panic(err)
	}

	// Save the compressed image
	err = SaveImage("compressed.jpg", compressedData)
	if err != nil {
		panic(err)
	}
}

func BenchmarkDefaultGolang(b *testing.B) {
	f, _ := os.OpenFile("./camera1.jpg", os.O_RDONLY, 0)
	// imgBytes, _ := io.ReadAll(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jpegImage, _ := jpeg.Decode(f)
		_, _ = f.Seek(0, 0)
		openFile, _ := os.Create("compressed.jpg")
		_ = jpeg.Encode(openFile, jpegImage, nil)
	}

	ms := float64(b.Elapsed()/time.Duration(b.N)) / 1e6
	b.ReportMetric(ms, "ms/op")
}

func BenchmarkTurboJpeg(b *testing.B) {
	file, err := os.Open("./camera1.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		img, _ := DecodeJPEG(file)
		_, _ = file.Seek(0, 0)
		// Compress the image
		compressedData, err := CompressImage(img)
		if err != nil {
			panic(err)
		}

		// Save the compressed image
		err = SaveImage("compressed.jpg", compressedData)
		if err != nil {
			panic(err)
		}
	}

	ms := float64(b.Elapsed()/time.Duration(b.N)) / 1e6
	b.ReportMetric(ms, "ms/op")
}
