+++
date = '2025-11-04T15:07:11+08:00'
draft = true
title = '一些简单的算法记录'
tags = ["双指针算法", "贪心算法"]
lastmod = '2025-12-03T18:00:00+08:00'
+++

### Leecode 两数相加

#### 题目： 给出两个链表，相同下表的节点数字相加，大于10，往后进1

```golang

    type ListNode struct {
        Val int
        Next *ListNode
    }

    func twoListNodeSum(l1, l2 ListNode) (head *ListNode) {
        // 定义这个是为了让tail始终指向最后一个元素。防止每次在链表尾部追加元素都要遍历一边完成的链表，从 o(n*2) 复杂度 => o(n+1)
        var tail *ListNode

        carry := 0

        for l1 != nil || l2 != nil {
            n1, n2 := 0, 0

            if l1 != nil {
                n1 = l1.Val
                l1 = l1.Next
            }

            if l2 != nil {
                n2 = l2.Val
                l2 = l2.Next
            }

            sum := n1 + n2 + carry
            sum, carry = sum%10, sum/10

            if head == nil {
                head = tail
                tail = &ListNode{Val: sum}
            } else {
                tail.Next = &ListNode{Val: sum}
                tail = tail.Next
            }
        }
        if carry > 0 {
            tail.Next = &ListNode{Val:sum}
        }

        return
    }

```


### 双指针

假设有两个数组，如果要计算两个数组中共同的元素有哪些。最简单的做法是，用python推导式完成即可
```Python

a = [1,2,3,4,5,6]
b = [2,3]

c = [i for i in b if i in a]

```
这种方式虽然简单，但是时间复杂度 O（m*n）。在超大数组的情况下效率低下。 对于大数组，可以使用双指针算法

```python
def find_common_element(arr1, arr2):
    # 先将 数组排序
    sorted_arr1 = sorted(arr1)
    sorted_arr2 = sorted(arr2)

    i = j = 0
    common = []

    while i < len(sorted_arr1) and j < len(sorted_arr2):
        a, b = sorted_arr1[i], sorted_arr2[j]

        if a == b:
            if not common or a != common[-1]:
                common.append(a)
            i += 1
            j += 1
        elif a > b:
            j += 1
        else:
            i += 1
    
    return common

```

#### Golang中的双指针使用示例
双指针有以下几种类型，相撞指针（同一个数组从两侧往中间位移，最终相撞）， 同向指针（快慢指针，删除重复项，滑动窗口）， 分离指针（两个指针分别处理两个不同的数组/链表）

- 两数之和。 
    一个数组中获取相加等于target 的元素位置, 利用的是数组的有序性。前提是需要有序的数组。
```golang

import (
    "fmt",
    "sort"
)


// 先排序，利用slice的有序性 + 双指针
func twoEle(numbers []int, target int) []int {
    left, right := 0, len(numbers) - 1


    for left < right {
        sum := numbers[left] + numbers[right]
        if sum == target {
            return []int{left+1, right + 1}
        } else if sum < target {
            left++
        } else {
            right--
        }
    }
    return []int{}
}


// 不排序，利用 hashMap 
func twoSum(nums []int, target int) []int {
    // 先定义一个map，用来存储nums元素中每个 target - 该元素 = 需要的值，以及对应的下标
    var targetMap map[int]int

    for i, num := range nums {
        targetNum := target - num
        if targetNumIdx, ok := targetMap[targetNum]; ok {
            return [i, targetNumIdx]
        } else {
            targetMap[targetNum] = i
        }
    }

    // 没有符合的
    return []int{}
}



func main() {
    numbers := []int{2,11,7,15}
    target := 9
    sort.Ints(numbers)
    fmt.Println(twoEle(numbers, target))
}

```

- 删除一个数组中的重复项

```golang

func removeDuplicates(numbers []int) int {
    if len(numbers) == 0 {
        return 0
    }

    slow := 0
    for fast := 1; fast < len(numbers); fast++ {
        if numbers[slow] != numbers[fast] {
            slow++
            numbers[slow] = numbers[fast]
        }
    }
    return slow+1
}


func main() {
    numbers := []int{1,1,2,3,4}
    length := removeDuplicates(numbers)
    fmt.Println("有效长度：", length)
    fmt.Println("有效数组", numbers[:length])
}
```

### 贪心算法
``` python

# 给你一个长度为 m 的数组 tasks，其中 tasks[i] 表示第 i 个任务所需的执行时间（单位：秒）。
# 现在有 n 个并发执行的“工作者”（相当于 n 个 CPU 或线程），每个工作者一次只能执行一个任务。
# 所有任务一旦开始不可中断，假设任务可以随意分配给任意工作者。
# 请设计一个算法，计算在最优调度下，所有任务的最短总执行时间。

# tasks = [5, 2, 1, 7, 3, 4]
# n = 3
# 最短完成时间：8
# 机器 1：7 + 1 = 8
# 机器 2：5 + 3 = 8
# 机器 3：4 + 2 = 6
# Max(8, 8, 6) = 8（全部结束时间取最大值）

# 贪心算法
# 1. 将任务按执行时间从小到大排序
# 2. 将任务分配给工作者
# 3. 计算每个工作者的完成时间
# 4. 返回所有工作者的完成时间中的最大值

def min_total_execution_time(tasks, n):
    
    # 1. 计算核心边界值
    total_time = sum(tasks)
    max_task = max(tasks)
    min_possible_T = max(max_task, (total_time + n - 1) // n)  # (a + b -1) // b 等价于向上取整
    
    # 2. 贪心分配任务：用列表记录每个工作者的当前总时间
    workers = [0] * n  # workers[i] 表示第i个工作者已分配的总时间
    
    for task in sorted(tasks, reverse=True):  # 优先分配大任务，更易均衡
        # 找到当前总时间最短的工作者
        min_worker_idx = workers.index(min(workers))
        # 将当前任务分配给该工作者
        workers[min_worker_idx] += task
    
    # 3. 实际分配后的最大时间（应等于理论最小值）
    actual_T = max(workers)
    return actual_T

# 测试示例
tasks = [5, 2, 1, 7, 3, 4]
n = 3
print(min_total_execution_time(tasks, n))  # 输出：8

```

```python

def coin_change(coins, amount):
    #
    coins.sort(reverse=True)
    count = 0
    remaining = amount

    for coin in coins:
        if remaining >= coin:
            num = remaining // coin
            count += num
            remaining -= num * coin
        if remaining == 0:
            break

    return count if remaining == 0 else -1
    
coins = [25, 10, 5, 1]
amount = 37
print(coin_change(coins, amount))

```

```python
# 最多活动数量
def max_activities(activities):
    activities.sort(key=lambda x: x[1])
    count = 1 #至少能参加一个活动
    last_end = activities[0][1]

    for i in range(0, len(activities)):
        start, end = activities[i]
        if start >= last_end:
            count += 1
            last_end = end
    return count

activities = [(1, 4), (3, 5), (0, 6), (5, 7), (3, 9), (5, 9), (6, 10), (8, 11), (8, 12), (2, 14), (12, 16)]
print(max_activities(activities))

```