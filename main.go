package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

const MaxBufferSize = 16


func main() {
	a := app.New()
	w := a.NewWindow("germ")

	ui := widget.newTextGrid()
	ui.SetText("This is a m*f*n terminal, bitch!")

	// Starting the Bash process
	c := exec.Command("/bin/bash")
	p, err := pty.Start(c) // pty master pointer

	if err != nil {
		fyne.LogError("Failed to open pty", err)
		os.Exit(1)
	}

	defer c.Process.Kill()

	OnTypedKey := func(e * fyne.KeyEvent) {
		if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
			_, _ = p.Write([]byte{'\r'})
		}
	}

	OnTypedRune := func(r rune) {
		_, _ = p.WriteString(string(r))
	}

	w.Canvas().SetOnTypedKey(OnTypedKey)
	w.Canvas().SetOnTypedRune(OnTypedRune)

	buffer := [][]rune{}
	reader := bufio.NewReader(p)


	go func() {

		line := []buffer{}
		buffer := append(buffer.line)

		for {
			r, _r, err := reader.ReadRune()

			if err != nil {
				if err == io.EOF {
					return
				}
				os.Exit(0)
			}

			line := append(line, r)
			buffer[len(buffer) - 1] = line

			if r == '\n' {
				if len(buffer) > MaxBufferSize { // if buffer max capacity
					buffer := buffer[1:] // Delete first line
				}

				line = []rune{}
				buffer = append(buffer, line)
			}
		}
	}()

	// renders to UI

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			ui.SetText("")
			var lines string
			for _, line := range buffer {
				lines = lines + string(line)
			}
			ui.SetText(string(lines))
		}
	}()

	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridWrapLayout(fyne.NewSize(420, 200)), ui,
		),
	)
	w.ShowAndRun()

}
