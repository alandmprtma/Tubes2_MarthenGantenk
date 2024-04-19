// variasi gabungin semua
package main

import (
	"fmt"
	"sync"
)

func sum(numbers []int, result chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    result <- sum // Mengirim nilai sum ke channel result
}

func main() {
    var wg sync.WaitGroup
    numbers := []int{1, 2, 3, 4, 5}
    result := make(chan int)

    // Menambahkan 1 ke WaitGroup
    wg.Add(1)

    // Memulai goroutine untuk menjumlahkan bilangan
    go sum(numbers, result, &wg)

    // Menunggu hingga goroutine selesai
    go func() {
        wg.Wait()
        close(result) // Menutup channel setelah selesai
    }()

    // Menerima hasil dari channel result
    total := <-result
    fmt.Println("Total:", total)
}

// variasi pake wait group
// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func printNumbers(wg *sync.WaitGroup) {
//     defer wg.Done() // Mengurangi nilai WaitGroup setelah goroutine selesai
//     for i := 1; i <= 5; i++ {
//         fmt.Printf("%d ", i)
//     }
// }

// func main() {
//     var wg sync.WaitGroup

//     // Menambahkan 3 ke WaitGroup, karena kita akan menjalankan 3 goroutine
//     wg.Add(3)

//     // Memulai 3 goroutine
//     go printNumbers(&wg)
//     go printNumbers(&wg)
//     go printNumbers(&wg)
    
//     // Menunggu hingga semua goroutine selesai
//     wg.Wait()

//     fmt.Println("Selesai")
// }

// variasi pake channel
// package main

// import (
// 	"fmt"
// 	"time"
// )

// func someTask(id int, data chan int) {
//    for taskId := range data {
//       time.Sleep(2 * time.Second)
//       fmt.Printf("Worker: %d executed Task %d\n", id, taskId)
//    }
// }

// func main() {
//    // Creating a channel
//    channel := make(chan int)

//    // Creating 10.000 workers to execute the task
//    for i := 0; i < 10000; i++ {
//       go someTask(i, channel)
//    }

//    // Filling channel with 100.000 numbers to be executed
//    for i := 0; i < 100000; i++ {
//       channel <- i
//    }

// }

// variasi biasa
// package main

// import (
// 	"fmt"
// 	"time"
// )

// func printNumbers() {
//     for i := 1; i <= 5; i++ {
//         time.Sleep(1 * time.Second)
//         fmt.Printf("%d ", i)
//     }
// }

// func main() {
//     // Memulai goroutine
//     printNumbers()

//     // Membiarkan program tetap berjalan selama goroutine masih berjalan
//     // Jika tidak, program akan selesai sebelum goroutine selesai
//     time.Sleep(6 * time.Second)

//     fmt.Println("Selesai")
// }