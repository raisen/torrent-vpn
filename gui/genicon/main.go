package main

import (
    "image"
    "image/png"
    "log"
    "os"

    "github.com/srwiley/oksvg"
    "github.com/srwiley/rasterx"
)

// genicon renders assets/icon.svg to assets/icon.png at 512x512 for app packaging
func main() {
    svgPath := "assets/icon.svg"
    pngPath := "assets/icon.png"

    f, err := os.Open(svgPath)
    if err != nil {
        log.Fatalf("open svg: %v", err)
        return
    }
    defer f.Close()

    icon, err := oksvg.ReadIconStream(f)
    if err != nil {
        log.Fatalf("parse svg: %v", err)
        return
    }

    // Target size
    w, h := 512, 512
    icon.SetTarget(0, 0, float64(w), float64(h))

    rgba := image.NewRGBA(image.Rect(0, 0, w, h))
    scanner := rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())
    raster := rasterx.NewDasher(w, h, scanner)

    // Render at scale 1
    icon.Draw(raster, 1.0)

    out, err := os.Create(pngPath)
    if err != nil {
        log.Fatalf("create png: %v", err)
        return
    }
    defer out.Close()

    if err := png.Encode(out, rgba); err != nil {
        log.Fatalf("encode png: %v", err)
        return
    }

    log.Printf("wrote %s (%dx%d)", pngPath, w, h)
}
