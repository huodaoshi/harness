# 何时 Mock

仅在**系统边界** mock：

- 外部 API（支付、邮件等）
- 数据库（有时——更倾向测试库）
- 时间/随机性
- 文件系统（有时）

不要 mock：

- 自己的类/模块
- 内部协作对象
- 任何你掌控的代码

## 为可 Mock 而设计

在系统边界，设计易于 mock 的接口：

**1. 使用依赖注入**

外部依赖由外部传入，而非内部创建：

```typescript
// 易 mock
function processPayment(order, paymentClient) {
  return paymentClient.charge(order.total);
}

// 难 mock
function processPayment(order) {
  const client = new StripeClient(process.env.STRIPE_KEY);
  return client.charge(order.total);
}
```

**2. 优先 SDK 式接口，而非通用 fetcher**

为每个外部操作定义专用函数，而非一个带分支的通用函数：

```typescript
// 好：每个函数可独立 mock
const api = {
  getUser: (id) => fetch(`/users/${id}`),
  getOrders: (userId) => fetch(`/users/${userId}/orders`),
  createOrder: (data) => fetch('/orders', { method: 'POST', body: data }),
};

// 差：mock 里要写条件分支
const api = {
  fetch: (endpoint, options) => fetch(endpoint, options),
};
```

SDK 式做法意味着：

- 每个 mock 返回一种固定形状
- 测试搭建无需条件逻辑
- 更易看出测试覆盖了哪些端点
- 每个端点可有独立类型
