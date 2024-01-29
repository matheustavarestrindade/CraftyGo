package main

import (
	"github.com/matheustavarestrindade/CraftyGo/internal/app/database"
	web_server "github.com/matheustavarestrindade/CraftyGo/internal/app/web"
	"github.com/matheustavarestrindade/CraftyGo/internal/app/workers"
)

func main() {
	db := database.New()
	db.CreateUser("admin", "email", "password", true)

	worker, err := workers.CreateMinecraftServerWorker("1", 1024)
	if err != nil {
		panic(err)
	}

	go worker.Start()
	go func() {
		for {
			select {
			case <-worker.Ctx.Done():
				println("Worker stopped")
				return
			case stdout := <-worker.ProcessStdout:
				println(stdout)
			case stderr := <-worker.ProcessStderr:
				println(stderr)
			}
		}
	}()


	web_server.Start("3000")
}
