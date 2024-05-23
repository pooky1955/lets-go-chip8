package main

import "os"

func main() {
	filepath := os.Args[1]
	cg := CLIGraphics{}
	cs := CLISpeaker{}
	cg.Init(64, 32)
	cpu := CPU{G: &cg, S: &cs, Speed: 60}
	cpu.Init(filepath)
	// cpu.PrintProgramMemory()
	cpu.Loop()
}
