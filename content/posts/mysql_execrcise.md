+++
date = '2025-11-18T10:31:44+08:00'
draft = false
title = '用户留存问题联系'
tags = ["留存", "MySQL"]
toc = true
+++


### 题目
有一批用户的注册日期和活跃日期数据，计算：
次日留存率（Day 1 Retention）
每日留存率
每日流失率（1 - 当日留存率）
[
  { id: 1, signup: '2024-06-01', activeDays: ['2024-06-01', '2024-06-02', '2024-06-04'] },
  { id: 2, signup: '2024-06-01', activeDays: ['2024-06-01'] },
  { id: 3, signup: '2024-06-01', activeDays: ['2024-06-01', '2024-06-02'] },
  { id: 4, signup: '2024-06-02', activeDays: ['2024-06-02', '2024-06-03', '2024-06-04'] }
]
设计数据表结构，以及给出窗口函数，写出计算次流的Sql

### 参考答案
```mysql

create table if not exists user_signup (
    id int auto_increment primary key ,
    user_id int comment '用户ID',
    signup_date Date comment '注册日期'
);


create table if not exists user_active (
    id int auto_increment primary key ,
    user_id int not null,
    active_date date comment  '活跃日期'
);

create unique index idx_active_date on user_active(user_id, active_date);

-- 插入注册数据
INSERT INTO user_signup (user_id, signup_date)
VALUES (1, '2024-06-01'),
       (2, '2024-06-01'),
       (3, '2024-06-01'),
       (4, '2024-06-02');

-- 插入活跃数据
INSERT INTO user_active (user_id, active_date)
VALUES (1, '2024-06-01'),
       (1, '2024-06-02'),
       (1, '2024-06-04'),
       (2, '2024-06-01'),
       (3, '2024-06-01'),
       (3, '2024-06-02'),
       (4, '2024-06-02'),
       (4, '2024-06-03'),
       (4, '2024-06-04');


-- 计算次日留存率
with signup_users as (
    -- 按注册日期统计每日注册用户数
    select signup_date, count(user_id) as total_signup
    from user_signup
    group by signup_date
),
retained_user as (
    -- 统计注册次日仍活跃的用户数
    select s.signup_date, count(distinct a.user_id) as retained
    from user_signup s
        left join user_active a
        on s.user_id = a.user_id and a.active_date = s.signup_date + interval 1 day
    group by s.signup_date
)

select s.signup_date,
       s.total_signup,
       r.retained,
       round(IFNULL(r.retained / s.total_signup, 0) * 100, 2) as day1_retention_rate
from signup_users s
         left join retained_user r on s.signup_date = r.signup_date
order by s.signup_date;




-- 计算每日留存率
with signup_users as (
    select signup_date
       , count(user_id) as total_signup
    from user_signup
        group by signup_date
),
user_retention as (
    -- 计算每个用户注册后的活跃天数（N）
    select s.user_id
         , s.signup_date
         , DATEDIFF(a.active_date, s.signup_date) as retention_day
    from user_signup s
         left join user_active a
              on s.user_id = a.user_id and datediff(a.active_date, s.signup_date) >= 1
)
-- 按注册和留存天数统计留存率/流失率
select signup_date
     , retention_day
     , total_signup
     , retained
     , round(retained / total_signup * 100, 2)       as retention_rate
     , round((1 - retained / total_signup) * 100, 2) as churn_rate
from (select ur.signup_date
           , ur.retention_day
           , su.total_signup
           , count(distinct ur.user_id) as retained
      from user_retention ur
               left join signup_users su on su.signup_date = ur.signup_date
      where ur.retention_day is not null
      group by ur.signup_date, ur.retention_day, su.total_signup) t
order by signup_date, retention_day;

```


