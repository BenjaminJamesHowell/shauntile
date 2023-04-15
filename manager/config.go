/* 
	Shauntile: An epic window manager inspired by the legend Shaun Donnelly.
	Copyright (C) 2023 Benjamin Howell

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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

