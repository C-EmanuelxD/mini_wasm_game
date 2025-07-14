package main

import (
	"cart/w4"
)

type projetil struct {
	x, y int
	vivo bool
}

type entidade struct {
	vida int
	vivo bool
	x, y int
}

type nave_type struct {
	width, height int
	sprite        []byte
}

var vilao = nave_type{
	width:  8,
	height: 8,
	sprite: []byte{0xf4, 0x1f, 0x3c, 0x3c, 0x2d, 0x38, 0x49, 0x61, 0x55, 0x55, 0xdd, 0xdd, 0xff, 0xff, 0x77, 0x77},
}

var nave = nave_type{
	width:  8,
	height: 8,
	sprite: []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xa8, 0x2a, 0xa0, 0x0a, 0x80, 0x02, 0x01, 0x40, 0x01, 0x40, 0x29, 0x68},
}

var naveX, naveY = 79, 135
var tiros = make([]projetil, 10)
var delayTiro int
var jogoIniciado bool

var viloes = [3]entidade{
	{vida: 3, vivo: true, x: 10, y: 10},
	{vida: 3, vivo: true, x: 79, y: 10},
	{vida: 3, vivo: true, x: 150, y: 10},
}
var direcaoVilao = [3]bool{false, true, false}

//go:export update
func update() {
	w4.PALETTE[0] = 0x008
	*w4.DRAW_COLORS = 0
	w4.Rect(0, 0, 160, 160)

	w4.PALETTE[0] = 0x008
	w4.PALETTE[1] = 0xA020F0
	w4.PALETTE[2] = 0xFFFF00
	w4.PALETTE[3] = 0xFFFFFF

	var gamepad = *w4.GAMEPAD1

	if !jogoIniciado {
		*w4.DRAW_COLORS = 4
		w4.Text("Controles:", 10, 40)
		w4.Text("< > para mover", 10, 55)
		w4.Text("^ para atirar", 10, 70)
		*w4.DRAW_COLORS = 2
		w4.Text("Pressione X \npara jogar!", 35, 90)
		if gamepad&w4.BUTTON_1 != 0 {
			jogoIniciado = true
		}
		return
	}

	if gamepad&w4.BUTTON_RIGHT != 0 && naveX <= 150 {
		naveX += 2
	}
	if gamepad&w4.BUTTON_LEFT != 0 && naveX >= 2 {
		naveX -= 2
	}

	if gamepad&w4.BUTTON_UP != 0 && delayTiro == 0 {
		for i := 0; i < len(tiros); i++ {
			if !tiros[i].vivo {
				tiros[i] = projetil{x: naveX + 3, y: naveY, vivo: true}
				delayTiro = 12
				break
			}
		}
	}
	if delayTiro > 0 {
		delayTiro--
	}

	for i := 0; i < len(tiros); i++ {
		if tiros[i].vivo {
			tiros[i].y -= 2
			if tiros[i].y <= 0 {
				tiros[i].vivo = false
			}
		}
	}

	for i := 0; i < 3; i++ {
		if !viloes[i].vivo {
			continue
		}
		if !direcaoVilao[i] {
			viloes[i].x += 1
		} else {
			viloes[i].x -= 1
		}
		if viloes[i].x >= 151 {
			direcaoVilao[i] = true
			viloes[i].y += 8
		}
		if viloes[i].x <= 0 {
			direcaoVilao[i] = false
			viloes[i].y += 8
		}
	}

	for i := 0; i < len(tiros); i++ {
		for j := 0; j < 3; j++ {
			if viloes[j].vivo && tiros[i].vivo {
				if colisao(tiros[i].x, tiros[i].y, 2, 2, viloes[j].x, viloes[j].y, vilao.width, vilao.height) {
					tiros[i].vivo = false
					viloes[j].vida--
					if viloes[j].vida <= 0 {
						viloes[j].vivo = false
					}
				}
			}
		}
	}

	for i := 0; i < 3; i++ {
		if viloes[i].vivo && viloes[i].y >= naveY {
			*w4.DRAW_COLORS = 2
			w4.Text("Que pena!\no espaco e\ncomunista!", 30, 59)
			*w4.DRAW_COLORS = 4
			w4.Text("Pressione R para\ntentar de novo!", 25, 95)
			return
		}
	}

	tela()
}

func tela() {
	*w4.DRAW_COLORS = 0x0002
	w4.Blit(&nave.sprite[0], naveX, naveY, 8, 8, w4.BLIT_2BPP)

	for i := 0; i < 3; i++ {
			if viloes[i].vivo {
				*w4.DRAW_COLORS = 0x234
				w4.Blit(&vilao.sprite[0], viloes[i].x, viloes[i].y, 8, 8, w4.BLIT_2BPP)

				*w4.DRAW_COLORS = 3
				w4.Rect(viloes[i].x, viloes[i].y+9, 9, 2)

				vidaWidth := viloes[i].vida * 3
				*w4.DRAW_COLORS = 2
				w4.Rect(viloes[i].x, viloes[i].y+9, uint(vidaWidth), 2)
			}
		}

	*w4.DRAW_COLORS = 4
	for i := 0; i < len(tiros); i++ {
		if tiros[i].vivo {
			w4.Rect(tiros[i].x, tiros[i].y, 2, 2)
		}
	}

	vitoria := true
	for i := 0; i < 3; i++ {
		if viloes[i].vivo {
			vitoria = false
			break
		}
	}
	if vitoria {
		*w4.DRAW_COLORS = 2
		w4.Text("Parabens voce\n democratizou\n   o espaco!", 30, 59)
		*w4.DRAW_COLORS = 4
		w4.Text("Pressione R para\ntentar de novo!", 25, 95)
	}
}

func colisao(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}
