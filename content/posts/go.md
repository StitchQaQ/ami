+++
date = '2025-02-13T10:20:34+08:00'
draft = false
title = 'Go语言实战笔记'
lastmod = '2025-10-30T15:30:00+08:00'
categories = ["Golang"]
toc = true
+++

> 💡 本文记录 Go 语言学习过程中的核心概念和最佳实践，涵盖并发编程、数据结构、指针操作等关键知识点。

## 引用类型

通道(channel)、 映射(map)、 切片(slice)

> ⚠️ **提示**：引用类型在传递时传递的是引用，修改会影响原始数据。

## sync 包

sync 包提供同步 goroutine 的功能。 在Go语言中， main终止则代表了程序终止， main函数终止前还会关闭所有之前启动且运行的 goroutine。 
并发写程序时，最佳做法是 在main函数返回前，清理并终止所有之前启动的 goroutine。编写启动和终止都清晰的程序，减少bug，防止资源异常。

sync 包的 WaitGroup 会跟踪所有启动的 goroutine。 WaitGroup 是一个计数信号量，用来统计每个 goroutine 是否完成了工作。
具体做法是: 将 WaitGroup 的变量值设置为要启动的 goroutine 数量。 每个 goroutine 完成工作后， 递减 WaitGroup 变量的计数值。
当这个值为0，就知道所有的 goroutine 都完成了工作。

> 💡 **最佳实践**：使用 `defer wg.Done()` 确保 goroutine 完成时一定会调用 Done，避免程序永久阻塞。

### 命名规范

1. 命名接口时，如果接口内就一个方法，需要以 er 结尾。 如 Matcher。 如果接口类型内部声明了多个方法，其名字需要与其行为关联

## 最佳实践

### 接收者类型选择

如果声明函数的时候带有接收者，则意味着声明了一个方法。这个方法会和指定的接收者类型绑定在一起。 
值类型和引用类型作为接收者的方法，在调用时编译器会自动编码/解码对应的类型，然后传递给方法。例如
```golang
// 方法声明为使用 defaultMatcher 类型的值作为接收者
func (m defaultMatch) Search(feed *Feed, stringTerm string)

// 声明一个指向 defaultMatcher类型值的指针
dm := new(defaultMatch)

// 编译器会解开dm指针的引用，使用对应的值调用方法
dm.Search(feed, "test")

// 方法声明为使用指向defaultMatcher类型值的指针作为接收者
func (m *defaultMatcher) Search(feed *Feed, searchTerm string)

// 声明一个defaultMatcher类型的值
var dm defaultMatch

// 编译器会自动生成指针引用dm值，使用指针调用方法
dm.Search(feed, "test")”

```

### 接口类型的方法调用规则

与直接通过值或者指针调用不同， 如果通过接口类型的值调用方法，规则有很大不同。 
1. 使用指针作为接收者声明的方法，只能在接口类型的值是一个指针的时候被调用。
2. 使用值作为接收者声明的方法，在接口类型的值为值或者指针时，都可以被调用。

> ⚠️ **重要**：接口调用与直接调用规则不同！指针接收者的方法只能通过指针调用接口。

```golang
// 方法声明为使用指向defaultMatcher类型值的指针作为接收者
func (m *defaultMatcher) Search(feed *Feed, searchTerm string)

// 通过interface类型的值来调用方法
var dm defaultMatcher
var matcher Matcher = dm     // 将值赋值给接口类型
matcher.Search(feed, "test") // 使用值来调用接口方法

> go build
cannot use dm (type defaultMatcher) as type Matcher in assignment

// 方法声明为使用defaultMatcher类型的值作为接收者
func (m defaultMatcher) Search(feed *Feed, searchTerm string)

// 通过interface类型的值来调用方法
var dm defaultMatcher
var matcher Matcher = &dm    // 将指针赋值给接口类型
matcher.Search(feed, "test") // 使用指针来调用接口方法

> go build
Build Successful

```

## GO 命令

### go vet
可以帮助检测代码的常见错误。 例如： Printf类函数调用时，类型匹配错误参数、错误的结构题标签、没有指定字段名的结构字面量

### go fmt
fmt 工具会将开发人员的代码布局成和go源代码类似的风格。 But,现代IDE在保存代码时会自动格式化。使用 vim 编程建议使用。

### go doc
有两种方式，一直是直接在命令行获取文档. 比如直接查看 `archive/tar`包的相关文档，在命令行中执行: `go doc tar`。 
另一种：在Terminal中输入`godoc -http=:6060` 可以直接生成浏览器文档。查看更多细节。 需要注意的是: godoc 不再包含在 Go 的标准工具链中，需要手动安装

手动安装 godoc（适用于 Go 1.16 及以上版本）
```bash
go install golang.org/x/tools/cmd/godoc@latest
```

> 💡 **提示**：在保存代码前运行 `go vet` 和 `go fmt` 可以避免很多常见错误。大多数 IDE 都支持自动运行这些工具。



## 数组、切片、映射

数组是切片和映射的基础数据结构

在Go语言中，数组是一个长度固定的数据类型，用于存储一段具有相同类型的元素的连续块。 数组存储的类型可以是内置类型，如整数或者字符串，也可以是某种结构类型。数组声明完之后，数据类型和数组长度就不能改变了

> 📝 **注意**：数组是值类型，切片和映射是引用类型。使用切片比数组更灵活，推荐优先使用切片。

## 指针

golang中是值传递。 传址也是传递地址的副本。

- `&` 是地址运算符，它位于值类型之前，返回存储改值的内存位置的地址
- `*` 是间接寻址运算符。位于指针变量之前，返回指向的值

> ⚠️ **记住**：Go 中一切都是值传递。即使传递指针，传递的也是指针的副本，只不过它们指向同一块内存。

## Channel

在Golang中， channel分有缓存通道和无缓冲通道。 

### 无缓冲通道

无缓冲通道：在写入channel后，就会阻塞等待接受者读取。Go中有大量依赖阻塞特性的场景，比如
```golang
// 示例： 创建一个生成器
func generateInteger() func() int {
    ch := make(chan int)
    cnt := 0

    go func () {
        for {
            ch <- cnt
            cnt++
        }
    }{}
}

func main() {

    generate := generateInteger()
    fmt.Println(generate()) // 0
    fmt.Println(generate()) // 1
    fmt.Println(generate()) //
}

```

> 💡 **核心理解**：我一开始有个疑问，在初始化generateInteger()之后，for循环一直在执行中，不会导致超时或者OOM么。但事实上是，循环中的操作会被 channel 的阻塞特性"限流"，不会无限制占用内存。

### 使用 Channel 监控协程完成

并且很多Go程序都会用到这个特性。下面是通过channel监控子协程的完成。不用通过time.Sleep等不可靠的方式。
```golang
func task(done chan<-struct{}) {
    // 这里执行子任务

    done <- struct{}{} //任务完成，发送信号
}


func main() {
    done := make(chan struct{})

    go task(done)

    <-done  // 这里就会阻塞等待协程完成任务

    fmt.Println("任务完成")

}

```

> ✅ **最佳实践**：使用空结构体 `struct{}{}` 作为信号传递，因为它不占用内存空间。使用单向 channel (`chan<-` 或 `<-chan`) 可以提高类型安全性。

---

> 📚 **参考资源**
> - 《Go语言实战》
> - Go 官方文档：https://go.dev/doc/

