package main

import (
    "bytes"
    "fmt"
    "image"
    "image/color"
    "image/gif"
    "image/jpeg"
    "image/png"
    "io/ioutil"
    "os"
	"strings"
)
// 图像变黑白
func Preprocess(source string, target string) {
	ff, _ := ioutil.ReadFile(source)
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb)
	newGray := hdImage(m)
	target += strings.Split(source, "/")[len(strings.Split(source, "/")) - 1]
	f , _ := os.Create(target)
	defer f.Close()
	encode(target,f,newGray)
}

func hdImage(m image.Image) *image.RGBA {
    bounds := m.Bounds()
    dx := bounds.Dx()
    dy := bounds.Dy()
    newRgba := image.NewRGBA(bounds)
    for i := 0; i < dx; i++ {
        for j := 0; j < dy; j++ {
            colorRgb := m.At(i, j)
            _, g, _, a := colorRgb.RGBA()
            g_uint8 := uint8(g >> 8)
            a_uint8 := uint8(a >> 8)
            newRgba.SetRGBA(i, j, color.RGBA{g_uint8, g_uint8, g_uint8, a_uint8})
        }
    }
    return newRgba
}

func encode(inputName string, file *os.File, rgba *image.RGBA) {
    if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
        jpeg.Encode(file, rgba, nil)
    } else if strings.HasSuffix(inputName, "png") {
        png.Encode(file, rgba)
    } else if strings.HasSuffix(inputName, "gif") {
        gif.Encode(file, rgba, nil)
    } else {
        fmt.Errorf("不支持的图片格式")
    }
}