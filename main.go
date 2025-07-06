package main

import (
	"cart/w4"
)

type projetil struct{
	x, y int32
	vivo bool
}

type entidade struct{
	vida int
	vivo bool
	x int
	y int

}

type nave_type struct{
	width int
	height int
	sprite []byte
}

var vilao = nave_type{
	width: 8,
	height: 8,
	sprite: []byte{0xf4,0x1f,0x3c,0x3c,0x2d,0x38,0x49,0x61,0x55,0x55,0xdd,0xdd,0xff,0xff,0x77,0x77},
}

var nave = nave_type{
	width: 8,
	height: 8,
	sprite: []byte{0xaa,0xaa,0xaa,0xaa,0xa8,0x2a,0xa0,0x0a,0x80,0x02,0x01,0x40,0x01,0x40,0x29,0x68},
}

var vilao1 = entidade{
		vida: 3,
		vivo: true,
		x: 79,
		y: 10,
	}

var vilaoX, vilaoY uint32 = 79, 10
var naveX, naveY uint32 = 79, 135
var tiro projetil
var onde bool = false
//go:export update
func update() {
	w4.PALETTE[0] = 0x008 //PRETO COR FABRICIO

	*w4.DRAW_COLORS = 0   
	w4.Rect(0, 0, 160, 160)

	w4.PALETTE[0] = 0x008 //PRETO COR FABRICIO
	w4.PALETTE[1] = 0xA020F0 //ROSA
	w4.PALETTE[2] = 0xFFFF00 //AMARELO
	w4.PALETTE[3] = 0xFFFFFF //BRANCO

	

	var gamepad = *w4.GAMEPAD1


	if gamepad & w4.BUTTON_RIGHT != 0{
		if naveX <= 150{ //while navex >= 159
			naveX += 2
		}
	}
	if gamepad & w4.BUTTON_LEFT != 0{
		if naveX >= 2{
			naveX -= 2
		}
		
	}
	if gamepad & w4.BUTTON_UP != 0{
		if tiro.vivo != true{
			tiro.x = int32(naveX)
			tiro.y = int32(naveY)
			tiro.vivo = true
		}
	}
	if tiro.vivo{
		tiro.y -= 2
		if tiro.y <= 0{
			tiro.vivo = false
		}
	}

	if vilao1.vivo{
		if onde == false{
			vilao1.x += 1
		}
		if onde == true{
			vilao1.x -= 1
		}
		if vilao1.x >= 151{
			onde = true
			vilao1.y += 9
		}
		if vilao1.x == 0{
			onde = false
			vilao1.y += 9
		}
	}


	if colisao(int(tiro.x), int(tiro.y), 2, 2, vilao1.x, vilao1.y, vilao.height, vilao.height){
		tiro.vivo = false
		vilao1.vida -= 1
		tiro.x = 0
		tiro.y = 0

		if vilao1.vida == 0{
			vilao1.vivo = false
			vilao1.x = -100
			vilao1.y = -100
		}
	}
	tela()
}


func tela(){
	*w4.DRAW_COLORS = 0x0002
	w4.Blit(&nave.sprite[0], int(naveX), int(naveY), 8, 8, w4.BLIT_2BPP)
	if vilao1.vivo{
		*w4.DRAW_COLORS = 0x234
		w4.Blit(&vilao.sprite[0], int(vilao1.x), int(vilao1.y), 8, 8, w4.BLIT_2BPP)
		if vilao1.y >= int(naveY){
			*w4.DRAW_COLORS = 2
			w4.Text("Que pena!\no espaco e\ncomunista!", 30, 59)
			*w4.DRAW_COLORS = 4
			w4.Text("Pressione R para\ntentar de novo!", 25, 95)
		}
	}else{
		*w4.DRAW_COLORS = 2
		w4.Text("Parabens voce\n democratizou\n   o espaco!", 30, 59)
		*w4.DRAW_COLORS = 4
		w4.Text("Pressione R para\ntentar de novo!", 25, 95)
	}
	if tiro.vivo{
		*w4.DRAW_COLORS = 4
		w4.Rect(int(tiro.x) + 3, int(tiro.y), 2, 2)
	}	
}


func colisao(x1, y1, w1, h1, x2, y2, w2, h2 int) bool{
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}