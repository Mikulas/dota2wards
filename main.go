package main

import (
    "fmt"
    "os"
//    "path/filepath"

    "github.com/dotabuff/yasha"
//    "github.com/dotabuff/yasha/dota"

    "github.com/dotabuff/yasha/utils"

    "image"
    "image/draw"
    "image/color"
    "image/png"

    "log"
)

var white color.Color = color.RGBA{255, 255, 255, 255}
var black color.Color = color.RGBA{0,0,0, 255}

const maxCoord = 16384
const maxImg = 1638

type Position struct {
    vec *utils.Vector2
    cellX, cellY int
    cellBits uint
}

func (p Position) PosX() float64 {
    return pos(p.cellX, p.cellBits, p.vec.X)
}

func (p Position) PosY() float64 {
    return pos(p.cellY, p.cellBits, p.vec.Y)
}

func pos(cell int, bits uint, origin float64) float64 {
    return float64(cell * (1 << bits) - maxCoord) + origin / 128
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Expected a .dem file as argument")
    }

    const nameObs = "DT_DOTA_NPC_Observer_Ward"
    const nameSen = "DT_DOTA_NPC_Observer_Ward_TrueSight"

    var wards []Position

    parser := yasha.ParserFromFile(os.Args[1])
    parser.OnEntityCreated = func(entity *yasha.PacketEntity) {
        if entity.Name == nameObs || entity.Name == nameSen {
            var ward Position
            ward.vec = entity.Values["DT_DOTA_BaseNPC.m_vecOrigin"].(*utils.Vector2);
            ward.cellX = entity.Values["DT_DOTA_BaseNPC.m_cellX"].(int)
            ward.cellY = entity.Values["DT_DOTA_BaseNPC.m_cellY"].(int)
            ward.cellBits = uint(entity.Values["DT_BaseEntity.m_cellbits"].(int))

            wards = append(wards, ward)
//            wardType := "obs"
//            if entity.Name == nameSen {
//                wardType = "sen"
//            }
//            fmt.Printf("%d,%d,%v,%f,%f\n", entity.Tick, entity.Values["DT_BaseEntity.m_iTeamNum"], wardType, ward.PosX(), ward.PosY())
        }
    }

    parser.Parse()

    file, err := os.Open("map_1638.png")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    img, _, err := image.Decode(file)
    if err != nil {
        log.Fatal(err)
    }

    const scale = float64(maxImg) / float64(maxCoord)
    const r = 10 // radius
    const b = 5 // border
    for i := range wards {
        x := int((wards[i].PosX() + maxCoord/2) * scale)
        y := int((wards[i].PosY() + maxCoord/2) * scale)
        fmt.Println(x, y)
        draw.Draw(img.(draw.Image), image.Rect(x-r-b, y-r-b, x+r+b, y+r+b), &image.Uniform{black}, image.ZP, draw.Src)
        draw.Draw(img.(draw.Image), image.Rect(x-r, y-r, x+r, y+r), &image.Uniform{white}, image.ZP, draw.Src)
    }

    w, _ := os.Create("output.png")
    defer w.Close()
    png.Encode(w, img)
}
