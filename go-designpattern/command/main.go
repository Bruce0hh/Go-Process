package main

import "go-designpattern/command/command"

func main() {
	tv := &command.TV{}

	onCommand := &command.OnCommand{Device: tv}
	offCommand := &command.OffCommand{Device: tv}

	onButton := &command.Button{Command: onCommand}
	onButton.Press()

	offButton := &command.Button{Command: offCommand}
	offButton.Press()
}
