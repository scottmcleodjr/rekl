package main

import (
	"flag"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
)

func main() {

	beep := flag.Bool("beep", false, "If the REKL should use a Beep Key instead of a Serial DTR Key")
	portName := flag.String("port", "/tty/USB0", "Serial port name for Serial DTR Key")
	flag.Parse()

	var key cwkeyer.Key
	var err error
	if *beep {
		key, err = cwkeyer.NewBeepKey(700, 48000, 1200) // Using the suggested values
	} else {
		key, err = cwkeyer.NewSerialDTRKey(*portName, 115200) // Using the suggested value
	}
	if err != nil {
		log.Fatalf("unable to create key: %s", err)
	}

	cfg := newConfig()
	keyer := cwkeyer.New(cfg, key)
	tui := newTUI()

	tui.inputField.SetInputCapture(func(capture *tcell.EventKey) *tcell.EventKey {
		for _, handler := range inputHandlers {
			captureOut, fired := handler(capture, &keyer, tui, cfg)
			if fired {
				return captureOut
			}
		}
		return capture
	})

	go func() {
		for {
			err := keyer.ProcessSendQueue(false)
			if err != nil {
				tui.writeToEventView(levelError, err.Error())
			}
		}
	}()

	err = tui.app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
