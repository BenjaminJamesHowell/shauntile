// This is an example shauntile config that demonstates various features. Feel
// free to copy this for your own setup.

package main

import (
	"time"
	"math/rand"
	"os/exec"

	"github.com/BenjaminJamesHowell/shauntile/manager"
	"github.com/BenjaminJamesHowell/shauntile/manager/keys"
)

func main() {
	manager.Config(config)
	manager.Start()
}

func config() {
	// Map a key
	manager.Map(
		keys.Super, // Mod key
		keys.P, // Key
		func() {
			exec.Command("dmenu_run").Start() // Callback
		},
	)

	// Run default browser
	manager.Map(
		keys.Super,
		keys.B,
		func() {
			exec.Command("x-www-browser").Start()
		},
	)

	// Run terminal emulator
	manager.Map(
		keys.Super,
		keys.Return,
		func() {
			exec.Command("x-terminal-emulator").Start()
		},
	)

	manager.Map(
		keys.Super | keys.Shift, // Multiple mods
		keys.Q,
		manager.Logout, // Logout
	)

	manager.Map(keys.Super, keys.Q, func() {
		// Get the focused client
		focused := manager.Focused
		if focused == nil {
			return
		}

		// Close a client
		focused.Close()
	})

	manager.Map(keys.Super, keys.R, func() {
		// Get a list of all clients
		clients := manager.Clients
		client := getRandomClient(clients) // FIXME

		// Set focus to a client
		if client != nil {
			client.Focus()
		}
	})

	manager.Map(keys.Super, keys.K, func() {
		nextClient := getNextClient(manager.Focused.GetId(), manager.Clients)
		if nextClient == nil {
			return
		}

		nextClient.Focus()
	})

	manager.Map(keys.Super, keys.J, func() {
		prevClient := getPrevClient(manager.Focused.GetId(), manager.Clients)
		if prevClient == nil {
			return
		}

		prevClient.Focus()
	})
}

func getNextClient(
	current int,
	clients []*manager.Client,
) *manager.Client {
	if len(clients) == 0 {
		return nil
	}

	current++

	if current >= len(clients) {
		current = 0
	}

	return clients[current]
}

func getPrevClient(
	current int,
	clients []*manager.Client,
) *manager.Client {
	if len(clients) == 0 {
		return nil
	}

	current--

	if current <= -1 {
		current = len(clients) - 1
	}

	return clients[current]
}

func getRandomClient(clients []*manager.Client) *manager.Client {
	if len(clients) == 0 { 
		return nil
	}

	rand.Seed(time.Now().Unix())
	return clients[rand.Intn(len(clients))]
}


