package main

import (
	ui "github.com/gizak/termui"
)

// screenFn represents the state of screen
// as a function that returns next screen
type screenFn func(ech *chan ui.Event) screenFn
