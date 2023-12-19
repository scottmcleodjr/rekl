package main

import (
	"flag"
	"log"

	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/handler"
	"github.com/scottmcleodjr/rekl/tui"
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

	cfg := config.New()
	keyer := cwkeyer.New(cfg, key)
	ui := tui.New()
	ui.SetInputFieldCapture(handler.InputHandler(keyer, ui, cfg))

	go func() {
		for {
			err := keyer.ProcessSendQueue(false)
			if err != nil {
				ui.WriteToEventView(tui.LevelError, err.Error())
			}
		}
	}()

	err = ui.RunApp()
	if err != nil {
		log.Fatal(err)
	}
}
