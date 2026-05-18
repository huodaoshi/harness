# 好测试与差测试

## 好测试

**集成风格**：经真实接口测试，不 mock 内部部件。

```typescript
// 好：测可观察行为
test("user can checkout with valid cart", async () => {
  const cart = createCart();
  cart.add(product);
  const result = await checkout(cart, paymentMethod);
  expect(result.status).toBe("confirmed");
});
```

特征：

- 测用户/调用方关心的行为
- 仅用公共 API
- 内部重构后仍成立
- 描述 WHAT，而非 HOW
- 每个测试一个逻辑断言

## 差测试

**实现细节测试**：与内部结构耦合。

```typescript
// 差：测实现细节
test("checkout calls paymentService.process", async () => {
  const mockPayment = jest.mock(paymentService);
  await checkout(cart, payment);
  expect(mockPayment.process).toHaveBeenCalledWith(cart.total);
});
```

红旗：

- Mock 内部协作对象
- 测私有方法
- 断言调用次数/顺序
- 行为未变但重构导致测试失败
- 测试名描述 HOW 而非 WHAT
- 绕开接口验证

```typescript
// 差：绕过接口验证
test("createUser saves to database", async () => {
  await createUser({ name: "Alice" });
  const row = await db.query("SELECT * FROM users WHERE name = ?", ["Alice"]);
  expect(row).toBeDefined();
});

// 好：经接口验证
test("createUser makes user retrievable", async () => {
  const user = await createUser({ name: "Alice" });
  const retrieved = await getUser(user.id);
  expect(retrieved.name).toBe("Alice");
});
```
