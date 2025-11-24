---
title: "广告系统基础"
date: 2025-04-14T09:53:08+08:00
draft: false
description: "广告计费模式与投放策略"
tags: ["广告", "商业模式"]
categories: ["业务知识"]
author: "Heyuuuu"
lastmod: 2025-10-30T18:00:00+08:00
---



#### 广告的计费模式有哪几种
1. **CPM** (Cost Per Mile) 千次展示计费
2. **CPC** (Cost Per Click) 点击收费
3. **CPA** (Cost Per Action) 行动收费 - **广告主**在用户完成某个特定行为才会支付费用
4. **CTP** (Cost Per Time) 按时长收费
5. **CPI** (Cost Per Install) 按安装收费等


#### 现代广告系统架构
- **广告投放平台**（**DSP**, Demand-Side Platform） - **广告主**通过 **DSP** 进行投放，**DSP** 负责竞价、定向、优化等功能
- **供应方平台**（**SSP**, Supply-Side Platform） - **媒体**方通过 **SSP** 管理广告库存，广告库存就是**广告位**。如果自己开发一个 app，想接广告。不需要自己找**广告主**，直接把**广告位**卖给 **SSP** 就行

- **广告交易市场**（**Ad Exchange**） - 一个开放的广告竞价平台，连接 **DSP** 和 **SSP**
- **数据管理平台**（**DMP**, Data Management Platform） - 收集、分析用户数据，提高广告定向能力等

- **广告监测平台**（**Ad Monitor**） - 监测广告投放效果，比如广告投放的曝光量、点击量、转化等，归因分析



### 一个广告的完整请求流程

**Step1: 用户打开应用，就会触发广告请求**

用户行为，打开 **APP**/网站，页面中的**广告位**会自动触发广告展示请求。例如：打开淘宝，触发开屏广告

**Step2: 媒体通过 SSP 提交流量信息到 Ad Exchange**

- **媒体**将广告的基础信息（尺寸、类型、页面内容）和用户匿名标识（**设备ID**、**Cookie**）发送给合作的 **SSP**
- **SSP** 对信息进行初步处理（过滤垃圾流量，设置**广告位**最低竞价等），然后将这些信息打包成"**竞价请求**（**Bid Request**）"，发送给 **Ad Exchange**（可能会同时发送给多个 **Ad Exchange**，以获取更高收益）

**Step3: Ad Exchange 分发竞价请求给 DSP**

- **Ad Exchange** 收到 **SSP** 的**竞价请求**后，会将请求转发给其合作的 **DSP**（**广告主**会通过 **DSP** 参与竞价）
- **竞价请求**包含的关键信息：用户匿名ID、设备类型、地域、**广告位**尺寸、**媒体**类型、**SSP** 设置的底价等

**Step4: DSP 基于用户数据和策略决定是否竞价**

- **DSP** 收到**竞价请求**后，会在 100 毫秒内完成以下操作：
    - 调用 **DMP** 获取用户标签：通过用户匿名ID匹配 **DMP** 中的数据，获取用户兴趣、消费能力、历史行为等标签
    - 匹配**广告主**需求：根据**广告主**在 **DSP** 上设置的定向条件（如"20-25岁男性"、"运动达人"），判断该用户是否符合目标受众
    - 计算出价：如果符合条件，**DSP** 根据**广告主**的预算、出价策略（如"每次点击最高出价2元"）、用户价值（是否高潜用户等）计算一个合理的竞价价格
    - 返回**竞价响应**：**DSP** 将"是否参与竞价"以及"出价金额"通过"**竞价响应**（**Bid Response**）"返回给 **Ad Exchange**

**Step5: Ad Exchange 进行竞价策略，判断获胜者**

- **Ad Exchange** 收集所有的 **DSP** 响应，在 50 毫秒内完成比价：
    - 筛选出价 > **SSP** 底价的竞价
    - 选择出价最高的 **DSP** 作为获胜者（部分 **Ad Exchange** 会结合广告质量评分调整，不完全按照价格）
- **Ad Exchange** 会向获胜 **DSP** 发送"竞价成功通知"，同时向其他 **DSP** 发送"未中标通知"

**Step6: 广告素材投放与信息同步**

- 获胜的 **DSP** 返回广告素材：获胜的 **DSP** 在收到通知后，立即将**广告主**的素材（如图片、视频、跳转链接）通过 **Ad Exchange** 转发给 **SSP**，再由 **SSP** 传递给**媒体**
- **媒体**展示广告：**媒体**将广告素材加载到**广告位**，用户最终看到广告（从用户触发请求到广告展示，全程需 200-300 毫秒内完成，避免影响用户体验）

**Step7: 后续数据追踪和结算**

- 行为数据回传：如果用户点击/转化（下单/下载）广告，**媒体**/**SSP** 会将行为通过 **Ad Exchange** 回传给获胜的 **DSP**，**DSP** 再同步**广告主**，用于效果监控
- 结算：**Ad Exchange** 根据竞价结果，在**广告主**（**DSP**）和**媒体**（**SSP**）之间进行费用结算（通常按 **CPM**/**CPC**/**CPI** 等方式），并抽取一定比例的交易佣金

---

看完上面广告的完整生命周期，再介绍几个专有名词：

- **媒体**（**Publisher**）：一般是指拥有流量的平台（网站、App），提供**广告位**（banner、信息流）
- **SSP**（Supply-Side Platform，供应方平台）：帮助**媒体**管理**广告位**，设置底价，对接 **Ad Exchange**，最大化流量收益
- **Ad Exchange**（广告交易平台）：连接供需双方的"交易市场"，接收 **SSP** 的流量请求，向 **DSP** 发起竞价，最终促成交易
- **DSP**（Demand-Side Platform，需求方平台）：代表**广告主**参与竞价，根据用户标签和投放策略决定是否出价以及出价金额
- **广告主**（**Advertiser**）：需要投放广告的企业主，通过 **DSP** 平台设置预算、定向条件（如用户年龄、性别、兴趣）等
- **DMP**（Data Management Platform，数据管理平台）：提供用户数据标签（行为、偏好），辅助 **DSP** 和 **SSP** 做决策


#### Bid Request Payload

##### Object: Request

<table>
  <tr>
    <td><strong>Attribute&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</strong></td>
    <td><strong>Type&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</strong></td>
    <td><strong>Definition</strong></td>
  </tr>
  <tr>
    <td><code>id</code></td>
    <td>string;&nbsp;required</td>
    <td>Unique ID of the bid request; provided by the exchange.</td>
  </tr>
  <tr>
    <td><code>test</code></td>
    <td>integer;<br/>default&nbsp;0</td>
    <td>Indicator of test mode in which auctions are not billable, where 0 = live mode, 1 = test mode.</td>
  </tr>
  <tr>
    <td><code>tmax</code></td>
    <td>integer</td>
    <td>Maximum time in milliseconds the exchange allows for bids to be received including Internet latency to avoid timeout. This value supersedes any general guidance from the exchange.  If an exchange acts as an intermediary, it should decrease the outbound <code>tmax</code> value from what it received to account for its latency and the additional internet hop.</td>
  </tr>
  <tr>
    <td><code>at</code></td>
    <td>integer;<br/>default&nbsp;2</td>
    <td>Auction type, where 1 = First Price, 2 = Second Price Plus.  Values greater than 500 can be used for exchange-specific auction types.</td>
  </tr>
  <tr>
    <td><code>cur</code></td>
    <td>string&nbsp;array;<br/>default&nbsp;[“USD”]</td>
    <td>Array of accepted currencies for bids on this bid request using ISO-4217 alpha codes. Recommended if the exchange accepts multiple currencies. If omitted, the single currency of “USD” is assumed.</td>
  </tr>
  <tr>
    <td><code>seat</code></td>
    <td>string&nbsp;array</td>
    <td>Restriction list of buyer seats for bidding on this item.  Knowledge of buyer’s customers and their seat IDs must be coordinated between parties beforehand. Omission implies no restrictions.</td>
  </tr>
  <tr>
    <td><code>wseat</code></td>
    <td>integer;<br/>default&nbsp;1</td>
    <td>Flag that determines the restriction interpretation of the <code>seat</code> array, where 0 = block list, 1 = allow list.</td>
  </tr>
  <tr>
    <td><code>cdata</code></td>
    <td>string</td>
    <td>Allows bidder to retrieve data set on its behalf in the exchange’s cookie (refer to <code>cdata</code> in <a href="#object_response">Object: Response</a>) if supported by the exchange. The string must be in base85 cookie-safe characters.</td>
  </tr>
  <tr>
    <td><code>source</code></td>
    <td>object</td>
    <td>A <code>Source</code> object that provides data about the inventory source and which entity makes the final decision. Refer to <a href="#object_source">Object: Source</a>.</td>
  </tr>
  <tr>
    <td><code>item</code></td>
    <td>object&nbsp;array; required</td>
    <td>Array of <code>Item</code> objects (at least one) that constitute the set of goods being offered for sale. Refer to <a href="#object_item">Object: Item</a>.</td>
  </tr>
  <tr>
    <td><code>package</code></td>
    <td>integer</td>
    <td>Flag to indicate if the Exchange can verify that the items offered represent all of the items available in context (e.g., all impressions on a web page, all video spots such as pre/mid/post roll) to support road-blocking, where 0 = no, 1 = yes.</td>
  </tr>
  <tr>
    <td><code>context</code></td>
    <td>object; recommended</td>
    <td>Layer-4 domain object structure that provides context for the items being offered conforming to the specification and version referenced in <code>openrtb.domainspec</code> and <code>openrtb.domainver</code>. <br />
For AdCOM v1.x, the objects allowed here all of which are optional are one of the <code>DistributionChannel</code> subtypes (i.e., <code>Site</code>, <code>App</code>, or <code>Dooh</code>), <code>User</code>, <code>Device</code>, <code>Regs</code>, <code>Restrictions</code>, and any objects subordinate to these as specified by AdCOM.</td>
  </tr>
  <tr>
    <td><code>ext</code></td>
    <td>object</td>
    <td>Optional exchange-specific extensions.</td>
  </tr>
</table>



#### Bid Request Paylaod

##### Object: Request
