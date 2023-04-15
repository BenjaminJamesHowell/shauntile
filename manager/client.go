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
	"time"
	"log"

	"github.com/jezek/xgb/xproto"
)

type position struct {
	X uint32;
	Y uint32;
}

type size struct {
	Width uint32;
	Height uint32;
}

type Client struct {
	Window xproto.Window;
	Position position;
	currentPosition position;
	Size size;
	currentSize size;
	isConfigured bool;
}

var Clients []*Client

func NewClient(
	window xproto.Window,
	xPosition uint32,
	yPosition uint32,
	width uint32,
	height uint32,
) *Client {
	return &Client {
		window,
		position {
			xPosition,
			yPosition,
		},
		position {
			xPosition,
			yPosition,
		},
		size {
			width,
			height,
		},
		size {
			width,
			height,
		},
		false,
	}
}

func (client *Client) Close() {
	if client == nil {
		return
	}

	property, err := xproto.GetProperty(
		connection,
		false,
		client.Window,
		atomWMProtocols,
		xproto.GetPropertyTypeAny,
		0,
		64,
	).Reply()

	if err != nil {
		log.Fatal(err)
	}
	if property == nil {
		err = xproto.SetCloseDownModeChecked(
			connection,
			xproto.CloseDownDestroyAll,
		).Check()

		if err != nil {
			log.Fatal(err)
		}

		err = xproto.DestroyWindowChecked(
			connection,
			client.Window,
		).Check()

		if err != nil {
			log.Fatal(err)
		}
	}

	for propertyValue := property.Value;
	len(propertyValue) >= 4;
	propertyValue = propertyValue[4:] {
		value := xproto.Atom(
			uint32(propertyValue[0]) |
			uint32(propertyValue[1]) << 8 |
			uint32(propertyValue[2]) << 16 |
			uint32(propertyValue[3]) << 24,
		)

		if value == atomWMDeleteWindow {
			currentTime := time.Now().Unix()
			eventData := xproto.ClientMessageDataUnionData32New(
				[]uint32 {
					uint32(atomWMDeleteWindow),
					uint32(currentTime),
					0,
					0,
					0,
				},
			)

			messageEvent := xproto.ClientMessageEvent {
				Format: 32,
				Window: client.Window,
				Type: atomWMProtocols,
				Data: eventData,
			}

			err = xproto.SendEventChecked(
				connection,
				false,
				client.Window,
				xproto.EventMaskNoEvent,
				string(messageEvent.Bytes()),
			).Check()

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (client *Client) Focus() {
	err := xproto.ConfigureWindowChecked(
		connection,
		client.Window,
		xproto.ConfigWindowStackMode,
		[]uint32 {xproto.StackModeAbove },
	).Check()

	if err != nil {
		log.Fatal(err)
	}
}

func (client *Client) GetId() int {
	for index, searchClient := range Clients {
		if searchClient == client {
			return index
		}
	}

	return -1
}

func (client *Client) addToList() {
	Clients = append(Clients, client)
}

func (clientToRemove *Client) removeFromList() {
	var newClients []*Client

	for _, client := range Clients {
		if client == clientToRemove {
			continue
		}

		newClients = append(newClients, client)
	}

	Clients = newClients
}

func (client *Client) reconfigure() {
	var mask uint16
	var values []uint32

	if client.currentPosition.X != client.Position.X ||
		!client.isConfigured {
		mask = mask | xproto.ConfigWindowX
		values = append(values, client.Position.X)
	}

	if client.currentPosition.Y != client.Position.Y ||
		!client.isConfigured {
		mask = mask | xproto.ConfigWindowY
		values = append(values, client.Position.Y)
	}

	if client.currentSize.Width != client.Size.Width ||
		!client.isConfigured {
		mask = mask | xproto.ConfigWindowWidth
		values = append(values, client.Size.Width)
	}

	if client.currentSize.Height != client.Size.Height ||
		!client.isConfigured {
		mask = mask | xproto.ConfigWindowHeight
		values = append(values, client.Size.Height)
	}

	if len(values) > 0 {
		err := xproto.ConfigureWindowChecked(
			connection,
			client.Window,
			mask,
			values,
		).Check()
		if err != nil {
			log.Fatal(err)
		}

		client.currentSize = client.Size
		client.currentPosition = client.Position
		client.isConfigured = true
	}
}

func GetClientByWindow(window xproto.Window) *Client {
	return GetClient(func(client *Client) bool {
		return client.Window == window
	})
}


func GetClient(predicate func(*Client) bool) *Client {
	for _, client := range Clients {
		if predicate(client) {
			return client
		}
	}

	return nil
}

