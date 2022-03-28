package main

import (
    "fmt"
    "warehouse-system/pkg/postgres"
)

func main() {
    fmt.Println("Query Worker")
    postgres.NewClient()
}
