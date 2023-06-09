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
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func handleEvent(event xgb.Event) {
	switch ev := event.(type) {
	case xproto.ConfigureRequestEvent:
		handleConfigureRequestEvent(ev)

	case xproto.DestroyNotifyEvent:
		handleDestroyNotifyEvent(ev)

	case xproto.EnterNotifyEvent:
		handleEnterNotifyEvent(ev)

	case xproto.KeyPressEvent:
		handleKeyPressEvent(ev)

	case xproto.MapRequestEvent:
		handleMapRequestEvent(ev)

	case xproto.UnmapNotifyEvent:
		handleUnmapNotifyEvent(ev)
	}
}

func handleConfigureRequestEvent(
	event xproto.ConfigureRequestEvent,
) {
	configureEvent := xproto.ConfigureNotifyEvent{
		Event: event.Window,
		Window: event.Window,
		AboveSibling: 0,
		X: 0,
		Y: 0,
		Width: event.Width,
		Height: event.Height,
		BorderWidth: event.BorderWidth,
		OverrideRedirect: false,
	}

	xproto.SendEventChecked(
		connection,
		false, 
		event.Window,
		xproto.EventMaskStructureNotify,
		string(configureEvent.Bytes()),
	)
}

func handleDestroyNotifyEvent(
	event xproto.DestroyNotifyEvent,
) {
	GetClientByWindow(event.Window).removeFromList()
}

func handleEnterNotifyEvent(
	event xproto.EnterNotifyEvent,
) {
	client := GetClientByWindow(event.Event)
	if client == nil {
		return
	}

	xproto.SetInputFocusChecked(
		connection,
		xproto.InputFocusPointerRoot,
		client.Window,
		xproto.TimeCurrentTime,
	)

	Focused = client
}

func handleKeyPressEvent(event xproto.KeyPressEvent) {
	mask := event.State & (xproto.ModMaskShift | xproto.ModMask4)
	key := key {
		int(mask),
		int(event.Detail),
	}
	command := maps[key]

	if command != nil {
		command()
	}
}

func handleMapRequestEvent(
	event xproto.MapRequestEvent,
) error {
	if GetClientByWindow(event.Window) != nil {
		return fmt.Errorf("Attempted to remap a window that has already been mapped.")
	}
	err := xproto.ChangeWindowAttributesChecked(
		connection,
		event.Window,
		xproto.CwEventMask,
		[]uint32 { xproto.EventMaskEnterWindow },
	).Check()
	if err != nil {
		return err
	}

	err = xproto.MapWindowChecked(connection, event.Window).Check()
	if err != nil {
		return err
	}

	client := NewClient(
		event.Window,
		0,
		0,
		screenWidth,
		screenHeight,
	)
	client.reconfigure()
	client.addToList()

	return nil
}

func handleUnmapNotifyEvent(
	event xproto.UnmapNotifyEvent,
) {
	GetClientByWindow(event.Window).removeFromList()
}

