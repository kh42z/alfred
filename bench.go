package main

import (
	"alfred/robot"
	"os"
)


func main() {
	hal :=  robot.NewBot("pong:3000")
	hal.Bench(1000, os.Getenv("ALFRED_CODE"))
}
