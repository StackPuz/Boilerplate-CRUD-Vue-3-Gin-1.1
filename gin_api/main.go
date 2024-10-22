package main

import (
    "app/config"
    "app/router"
)

func main() {
    config.Setup()
    router.Setup().Run(":8080")
}