package main

import (
    //"fmt"
	"log"
    "image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const height = 480
const width = 640
const sq_side = 50
const acc_g = 2

var posx = 100
var posy = 100
var ds = 2
var vely = 0
var jumping = false

type Root struct {
    Objects []ScreenObject
    You *YouSquare
}

func NewRoot() *Root {
    plat1 := NewPlatform(300, 400, 75, 20)
    plat2 := NewPlatform(40, 200, 75, 20)
    //plat2 := NewPlatform(200, 200, 75, 20)
    //plat3 := NewPlatform(400, 300, 75, 20)
    //plat4 := NewPlatform(425, 100, 75, 20)
    you := NewYouSquare(20)
    objects := []ScreenObject{plat1, plat2}
    return &Root{objects, you}
}

func (self *Root) Update(screen *ebiten.Image) error {
    for _,obj := range self.Objects {
        obj.Update(screen)
    }

    min_y, max_y := self.DetectCollisions()

    self.You.Update(screen, min_y, max_y)
    return nil
}

func (self *Root) DetectCollisions() (float64, float64){
    x := self.You.GetPosX()
    y := self.You.GetPosY()
    h := self.You.GetHeight()
    w := self.You.GetWidth()

    top := y
    bot := y + h
    lt := x
    rt := x + w

    max_y := float64(10000)
    //max_x := float64(0)
    min_y := float64(0)
    //min_x := 10000

    for _,obj := range self.Objects {
        o_x := obj.GetPosX()
        o_y := obj.GetPosY()
        o_h := obj.GetHeight()
        o_w := obj.GetWidth()
        o_top := o_y
        o_bot := o_y + o_h
        o_lt := o_x
        o_rt := o_x + o_w

        vert := false
        horz := false
        over := top <= o_top
        under := bot >= o_bot
        //left := lt < o_lt
        //right := rt > o_rt



        if (top >= o_top && top <= o_bot) || (bot >= o_top && bot <= o_bot) {
            vert = true
        }
        if (rt >= o_lt && rt <= o_rt) || (lt >= o_lt && lt <= o_rt) {
            horz = true
        }

        if over && horz {
            if o_top <= max_y {
                max_y = o_top
            }
        }

        if under && horz {
            if o_bot >= min_y {
                min_y = o_bot
            }
        }


        if vert && horz {
            if over {
                if o_top < max_y {
                    //max_y = o_top
                }
            } else if under {
                if o_bot > min_y {
                    //min_y = o_bot
                }
            }
        }
    }
    return min_y, max_y
}

type ScreenObject interface {
    GetPosX() float64
    GetPosY() float64
    GetWidth() float64
    GetHeight() float64
    Update(screen *ebiten.Image)
}

type Platform struct {
    PosX float64
    PosY float64
    Width float64
    Height float64
}

func NewPlatform(x,y, width, height float64) *Platform {
    return &Platform{x, y, width, height}
}

func (self *Platform) Update(screen *ebiten.Image) {
    gray := color.RGBA{0xaa,0xaa,0xaa, 0xaa}
    ebitenutil.DrawRect(
        screen, float64(self.PosX), float64(self.PosY),
        self.Width, self.Height, gray)
}

func (self *Platform) GetPosX() float64 {
    return self.PosX
}

func (self *Platform) GetPosY() float64 {
    return self.PosY
}

func (self *Platform) GetWidth() float64 {
    return self.Width
}

func (self *Platform) GetHeight() float64 {
    return self.Height
}

type YouSquare struct {
    Side float64
    PosX float64
    PosY float64
    VelY float64
    Speed float64
    Charge float64
    Jumping bool
    Charging bool
}

func NewYouSquare(side float64) *YouSquare {
    return &YouSquare{side, 0, 0, 0, 2, 0, false, false}
}

func (self *YouSquare) GetPosX() float64 {
    return self.PosX
}

func (self *YouSquare) GetPosY() float64 {
    return self.PosY
}

func (self *YouSquare) GetWidth() float64 {
    return self.Side
}

func (self *YouSquare) GetHeight() float64 {
    return self.Side
}


func (self *YouSquare) Update(screen *ebiten.Image, min_y, max_y float64) {
    actual_max_y := float64(height)
    if max_y < actual_max_y {
        actual_max_y = max_y
    }

    actual_min_y := float64(0)
    if min_y > actual_min_y {
        actual_min_y = min_y
    }

    if ebiten.IsKeyPressed(ebiten.KeySpace) && !self.Jumping {
        self.Jumping = true
        self.VelY = -20
    }

    if ebiten.IsKeyPressed(ebiten.KeyF) {
        self.Charging = true
        self.Charge += 1
        if self.Charge >= 10 {
            self.Charge = 10
        }
    } else if !self.Jumping && self.Charging{
        self.Charging = false
        self.VelY = -(20 + self.Charge)
        self.Jumping = true
        self.Charge = 0
    }

    green := 20*uint8(self.Charge)
    if green > 255 {
        green = uint8(255)
    }

    rect_color := color.RGBA{0,green,0xff,0xff}

    if ebiten.IsKeyPressed(ebiten.KeyS) {
        rect_color = color.RGBA{0xff,green,0,0xff}
        self.Speed = 8
    } else {
        self.Speed = 2
    }

    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        self.PosX -= self.Speed
        if self.PosX < 0 {
            self.PosX = 0
        }
    } else if ebiten.IsKeyPressed(ebiten.KeyRight) {
        self.PosX += self.Speed
        if self.PosX + sq_side > width {
            self.PosX = width - self.Side
        }
    }
    if self.PosY + self.Side < actual_max_y {
        self.VelY += acc_g
    }

    self.PosY += self.VelY
    if self.PosY + self.Side >= actual_max_y {
        self.PosY = actual_max_y - self.Side
        self.VelY = 0
        self.Jumping = false
    }

    if self.PosY <= actual_min_y {
        self.PosY = actual_min_y
        self.VelY = 0
        //self.Jumping = false
    }

    ebitenutil.DrawRect(
        screen, float64(self.PosX), float64(self.PosY),
        self.Side, self.Side, rect_color)
}




func main() {
    root := NewRoot()
	if err := ebiten.Run(root.Update, width, height, 2, "Title"); err != nil {
		log.Fatal(err)
	}
}
