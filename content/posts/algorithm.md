+++
date = '2025-11-04T15:07:11+08:00'
draft = true
title = '一些简单的算法记录'
tags = ["双指针算法", "贪心算法"]
lastmod = '2025-11-04T18:00:00+08:00'
+++


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
    sorted_arr1 = sort(arr1)
    sorted_arr2 = sort(arr2)

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