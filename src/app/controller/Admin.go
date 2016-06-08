package controller

import(
		"github.com/martini-contrib/render"
	)


func AdminMain(r render.Render) {
	r.HTML(200, "hello", "amdin")
}