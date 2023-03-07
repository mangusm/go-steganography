package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func binary(s string) string {
	res := ""
	for _, c := range s {
		res = fmt.Sprintf("%s%.8b", res, c)
	}
	return res
}

func binToString(s []byte) string {
	output := make([]byte, len(s)/8)
	for i := 0; i < len(output); i++ {
		val, err := strconv.ParseInt(string(s[i*8:(i+1)*8]), 2, 8)
		if err == nil {
			output[i] = byte(val)
		}
	}
	return string(output)
}

func main() {
	filepath := os.Args[1]
	test := load(filepath)
	save(filepath, test)
}

func load(filePath string) *image.YCbCr {
	s := strings.Split(filePath, ".")
	imgType := s[len(s)-1]
	switch imgType {
	case "jpeg", "jpg":
		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	case "png":
		image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	}

	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file:", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
	}
	return img.(*image.YCbCr)
}

func save(filePath string, img *image.YCbCr) {

	message := "Super secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret messageSuper secret message"
	lenHeader := fmt.Sprintf("%063b", uint64(len(message)))
	messageBin := binary(message)

	imgFile, err := os.Create("OUT_" + filePath)
	pxY := []uint8{}
	pxCb := []uint8{}
	pxCr := []uint8{}

	// TODO:
	// 1. convert each Y, Cb and Cr value to binary
	// 2. For each bit in lenHeader:
	// 3. Set the least-significant bit (LSB) of each color to the
	//    corresponding 0 or 1 in lenHeader
	// 4. Do the same thing for each bit in messageBin
	// Decode:
	// 1. Read the LSB of each color value 8 times to get the lenHeader
	// 2. Read the LSB of the next {lenHeader} color values
	// 3. Convert the result back to a string using binToString

	offset := 0
	for i := 0; i < len(lenHeader); i += 3 {
		digitsToWrite := lenHeader[i:]

		if img.Y[i]%2 != uint8(digitsToWrite[0])%2 {
			pxY = append(pxY, img.Y[i]+1)
		} else {
			pxY = append(pxY, img.Y[i])
		}

		if img.Cb[i]%2 != uint8(digitsToWrite[1])%2 {
			pxCb = append(pxCb, img.Cb[i]+1)
		} else {
			pxCb = append(pxCb, img.Cb[i])
		}

		if img.Cr[i]%2 != uint8(digitsToWrite[2])%2 {
			pxCr = append(pxCr, img.Cr[i]+1)
		} else {
			pxCr = append(pxCr, img.Cr[i])
		}
		offset += 3
	}

	for i := 0; i < len(messageBin); i += 3 {
		digitsToWrite := messageBin[i:Min(i+3, len(messageBin))]

		if len(digitsToWrite) >= 1 {
			if img.Y[i+offset]%2 != uint8(digitsToWrite[0])%2 {
				pxY = append(pxY, img.Y[i]+1)
			} else {
				pxY = append(pxY, img.Y[i])
			}
		}

		if len(digitsToWrite) >= 2 {

			if img.Cb[i+offset]%2 != uint8(digitsToWrite[1])%2 {
				pxCb = append(pxCb, img.Cb[i]+1)
			} else {
				pxCb = append(pxCb, img.Cb[i])
			}
		}

		if len(digitsToWrite) >= 3 {
			if img.Cr[i+offset]%2 != uint8(digitsToWrite[2])%2 {
				pxCr = append(pxCr, img.Cr[i]+1)
			} else {
				pxCr = append(pxCr, img.Cr[i])
			}
		}
		offset += 3
	}

	for i := offset; i < len(img.Y); i++ {
		pxY = append(pxY, img.Y[i])
		pxCb = append(pxCb, img.Cb[i])
		pxCr = append(pxCr, img.Cr[i])
	}
	fmt.Println(len(img.Y) - len(pxY))
	created := image.YCbCr{
		Y:              pxY,
		Cb:             pxCb,
		Cr:             pxCr,
		YStride:        img.YStride,
		CStride:        img.CStride,
		SubsampleRatio: img.SubsampleRatio,
		Rect:           img.Rect,
	}

	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	png.Encode(imgFile, created.SubImage(created.Rect))

}
