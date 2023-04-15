# Shauntile Window Manager
Shauntile is a window manager written in go inspired by the legend himself,
Shaun Donnelly. 
This repository contians the core window manager in the `manager/` directory
and an example config in the `example/` directory.
## Installation
Right now, there is no nice way of installing shauntile because the project is
still incomplete. You will have to create your own config and install shauntile
as a depenceny.
```
$ mkdir shauntile-config
$ cd shauntile-config
$ go mod init myconfig
$ go get github.com/BenjaminJamesHowell/shauntile/manager
```
Your configuration needs to call `manager.Config()` and then `manager.Start()`.
```go
package main

import (
	"github.com/BenjaminJamesHowell/shauntile/manager"
	"github.com/BenjaminJamesHowell/shauntile/manager/keys"
)

func main() {
	manager.Config(func() {
		// Your config here
	})
	manager.Start()
}
```
Then you use `manager.Map()` to create keybindings.
```go
// Within the config function
manager.Map(keys.Super, keys.Num1, func() {
	// Is run on Super+1

	// Get first client
	first := manager.Clients[0]

	// first is nil if there are no windows open
	if first == nil {
		return
	}

	// Focus first client
	first.Focus()
})
```
For more examples of keybindings, see [example/example.go](https://github.com/BenjaminJamesHowell/shauntile/blob/main/example/example.go). The functions used in configs are
mostly inside [manager/client.go](./manager/client.go) and [manager/config.go](./manager/config.go). Because of the current status of the project, there is no real documentation (lol).
## Project Status and Plans
- [-] Basic functionality (creating/closing windows, binding keys, ect.)
- [ ] Multiple workspaces
- [ ] Different window arrangements (split, fullscreen, ect.)

## Credit
Although I did most of the coding for this project, [Jlll1's window manager, btwm](https://github.com/Jlll1/btwm), and [their YouTube series on its
development](https://www.youtube.com/playlist?list=PLjfDSHUGSwofsA789Yc6n_nmknEiA1k-m) were very useful resources.

