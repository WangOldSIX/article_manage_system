package main

import (
	_ "day_day_fresh/routers"
	_"day_day_fresh/models"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

