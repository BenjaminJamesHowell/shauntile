package manager

import (
	"os"

	"github.com/jezek/xgb/xproto"
)

type key struct {
	modifier int;
	key int;
}

var configCallback func()
var maps map[key]func() = make(map[key]func())

func Config(callback func()) {
	configCallback = callback
}

func loadConfig() {
	configCallback()
}

func Map(modifier int, keyCode int, callback func()) {
	key := key {
		modifier,
		keyCode,
	}

	maps[key] = callback
}

func Logout() {
	if connection == nil {
		return
	}

	connection.Close()
	os.Exit(0)
}

func grabKeys(rootWindow xproto.Window) {
	for key := range maps {
		xproto.GrabKey(
			connection,
			true,
			rootWindow,
			uint16(key.modifier),
			xproto.Keycode(key.key),
			xproto.GrabModeAsync,
			xproto.GrabModeAsync,
		)
	}
}

