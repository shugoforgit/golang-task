package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// * 指针
	fmt.Println("task 指针 开始 ------------------------------")
	var a int = 10
	task2_pointer(&a)
	fmt.Println(a)
	fmt.Println("task 指针 结束 ------------------------------")

	// * Goroutine 1
	fmt.Println("task Goroutine 1 开始 ------------------------------")
	task2_goroutine()
	fmt.Println("task Goroutine 1 结束 ------------------------------")
	// * Goroutine 2
	fmt.Println("task Goroutine 2 开始 ------------------------------")
	task2_goroutine_2()
	fmt.Println("task Goroutine 2 结束 ------------------------------")

	// * 面向对象 1
	fmt.Println("task 面向对象 1 开始 ------------------------------")
	rectangle := Rectangle{width: 10, height: 20}
	circle := Circle{radius: 5}
	shapes := []Shape{rectangle, circle}
	for _, shape := range shapes {
		shape.Area()
		shape.Perimeter()
	}
	fmt.Println("task 面向对象 1 结束 ------------------------------")
	// * 面向对象 2
	fmt.Println("task 面向对象 2 开始 ------------------------------")
	employee := Employee{
		Person: Person{
			name: "John",
			age:  30,
		},
		EmployeeID: "123456",
	}
	employee.PrintInfo()
	fmt.Println("task 面向对象 2 结束 ------------------------------")

	// * Channel 1
	fmt.Println("task Channel 1 开始 ------------------------------")
	task4_channel()
	fmt.Println("task Channel 1 结束 ------------------------------")
	// * Channel 2
	fmt.Println("task Channel 2 开始 ------------------------------")
	task4_channel_2()
	fmt.Println("task Channel 2 结束 ------------------------------")

	// * 锁机制 1
	fmt.Println("task 锁机制 1 开始 ------------------------------")
	task5_lock()
	fmt.Println("task 锁机制 1 结束 ------------------------------")
	// * 锁机制 2
	fmt.Println("task 锁机制 2 开始 ------------------------------")
	task5_lock_2()
	fmt.Println("task 锁机制 2 结束 ------------------------------")
}

// * 指针---------------------------------------
func task2_pointer(ptrInt *int) {
	*ptrInt += 10
}

// * Goroutine 1 ----------------------------------
func task2_goroutine() {
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Println(i)
			}
		}
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Println(i)
			}
		}
	}()

	time.Sleep(time.Second * 1)
}

// * Goroutine 2 ----------------------------------
func task2_goroutine_2() {
	type Task struct {
		taskId int
		Exec   func()
	}

	var taskChan = make(chan Task)
	// * 创建调度器
	go func() {
		for task := range taskChan {
			go func(t Task) {
				start := time.Now()
				t.Exec()
				duration := time.Since(start)
				fmt.Printf("task id: %d 执行时间: %v\n", t.taskId, duration)
			}(task)
		}
	}()

	// * 往chan中添加任务
	go func() {
		for i := 1; i <= 5; i++ {
			taskId := i
			taskChan <- Task{
				taskId: i,
				Exec: func() {
					time.Sleep(time.Millisecond * time.Duration(taskId*100))
					fmt.Printf("task id: %d 执行完成\n", taskId)
				},
			}
		}

		close(taskChan)
	}()

	time.Sleep(time.Second * 2)
}

// * 面向对象 1 ----------------------------------
type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() {
	fmt.Printf("矩形的面积为 %.2f\n", r.width*r.height)
}

func (r Rectangle) Perimeter() {
	fmt.Printf("矩形的周长为 %.2f\n", 2*(r.width+r.height))
}

type Circle struct {
	radius float64
}

func (c Circle) Area() {
	fmt.Printf("圆的面积为 %.2f\n", math.Pi*c.radius*c.radius)
}

func (c Circle) Perimeter() {
	fmt.Printf("圆的周长为 %.2f\n", 2*math.Pi*c.radius)
}

// * 面向对象 2 ----------------------------------
type Person struct {
	name string
	age  int
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("Employee ID: %s, Name: %s, Age: %d\n", e.EmployeeID, e.name, e.age)
}

// * Channel 1 ----------------------------------
func task4_channel() {
	var channel = make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			channel <- i
		}
		close(channel)
	}()

	go func() {
		for i := range channel {
			fmt.Println(i)
		}
	}()

	time.Sleep(time.Second * 1)
}

// * Channel 2 ----------------------------------
func task4_channel_2() {
	var channel = make(chan int, 5)

	go func() {
		for i := 1; i <= 100; i++ {
			channel <- i
		}
		close(channel)
	}()

	go func() {
		for i := range channel {
			fmt.Println(i)
		}
	}()

	time.Sleep(time.Second * 1)
}

// * 锁机制 1 ----------------------------------
func task5_lock() {
	var mutex = sync.Mutex{}
	var count = 0

	for i := 0; i < 10; i++ {
		go func() {
			for i := 1; i <= 1000; i++ {
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}()
	}

	time.Sleep(time.Second * 1)
	fmt.Println("task5_lock count: ", count)
}

// * 锁机制 2 ----------------------------------
func task5_lock_2() {
	var count int64

	for i := 0; i < 10; i++ {
		go func() {
			for i := 1; i <= 1000; i++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}

	time.Sleep(time.Second * 1)
	fmt.Println("task5_lock_2 count: ", count)
}
