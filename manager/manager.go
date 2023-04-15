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
	"log"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var connection *xgb.Conn
var screenWidth uint32
var screenHeight uint32
var atomWMProtocols xproto.Atom
var atomWMDeleteWindow xproto.Atom

var Focused *Client

var WM_PROTOCOLS = "WM_PROTOCOLS"
var WM_DELETE_WINDOW = "WM_DELETE_WINDOW"

func Start() {
	var err error
	connection, err = xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}


	connectionInfo := xproto.Setup(connection)
	if connectionInfo == nil {
		log.Fatal("Could not parse connection info")
	}

	screen := connectionInfo.DefaultScreen(connection)
	root := screen.Root

	err = xproto.ChangeWindowAttributesChecked(
		connection,
		root,
		xproto.CwEventMask, []uint32{
			xproto.EventMaskKeyPress |
			xproto.EventMaskStructureNotify |
			xproto.EventMaskSubstructureRedirect,
		}).Check()
	if err != nil {
		log.Fatal(err)
	}

	loadConfig()

	atomWMProtocols = getAtom(WM_PROTOCOLS)
	atomWMDeleteWindow = getAtom(WM_DELETE_WINDOW)

	screenWidth = uint32(screen.WidthInPixels)
	screenHeight = uint32(screen.HeightInPixels)
	grabKeys(root)

	for {
		event, err := connection.WaitForEvent()
		if err != nil {
			log.Print("Non-Fatal Error: ", err)
		}
		if event == nil && err == nil {
			break
		}

		handleEvent(event)
	}
}

func getAtom(name string) xproto.Atom {
	reply, err := xproto.InternAtom(
		connection,
		false,
		uint16(len(name)),
		name,
	).Reply()
	if err != nil {
		log.Fatal(err)
	}

	if reply == nil {
		return 0
	}

	return reply.Atom
}

