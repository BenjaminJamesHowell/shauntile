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

