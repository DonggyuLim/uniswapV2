/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"sync"

	"github.com/DonggyuLim/uniswap/db"
	"github.com/DonggyuLim/uniswap/rest"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	db.NewDB()
	defer db.Close()
	go rest.Rest(wg)
	wg.Wait()
}
