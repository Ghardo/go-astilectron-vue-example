package main

import (
	"flag"
	"time"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	Version  = "0.0.1"
	AppName  string
	BuiltAt  string
	debug    = flag.Bool("d", false, "enables the debug mode")
	devtools = flag.Bool("dt", false, "enables the dev tools")
	w        *astilectron.Window
)

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:       *debug,
		MenuOptions: nil,
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]

			go func() {
				time.Sleep(2 * time.Second)

				if *devtools == true {
					astilog.Debug("Open dev tools")
					w.OpenDevTools()
				} else {
					astilog.Debug("Close dev tools")
					w.CloseDevTools()
				}


				if err := bootstrap.SendMessage(w, "sample.message1", "This is the first sample message"); err != nil {
					astilog.Error(errors.Wrap(err, "sending sample.message1 event failed"))
				}

				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "sample.message2", "And this is the second sample message"); err != nil {
					astilog.Error(errors.Wrap(err, "sending sample.message2 event failed"))
				}
				
			}()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(600),
				Width:           astilectron.PtrInt(800),
				Resizable:       astilectron.PtrBool(false),   	   
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
