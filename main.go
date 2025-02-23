package main
import (
    "fmt"
    "time"
)
func main() {
    start := time.Now()
    println("Hello WORLD")
    elapsed := time.Since(start)
    fmt.Printf("Time elapsed:DS", elapsed)
}