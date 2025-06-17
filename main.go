package main

import "github.com/Harichandra-Prasath/Terminal-OOO/core"

func main() {

	Server := core.NewServer(&core.ServerCfg{Port: 9000})
	if err := Server.Start(); err != nil {
		panic(err)
	}
}
