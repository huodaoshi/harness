# 面向可测性的接口设计

好接口让测试自然而然：

1. **接受依赖，不要内部创建**

   ```typescript
   // 可测
   function processOrder(order, paymentGateway) {}

   // 难测
   function processOrder(order) {
     const gateway = new StripeGateway();
   }
   ```

2. **返回结果，不要只产生副作用**

   ```typescript
   // 可测
   function calculateDiscount(cart): Discount {}

   // 难测
   function applyDiscount(cart): void {
     cart.total -= discount;
   }
   ```

3. **表面积小**
   - 方法越少 → 所需测试越少
   - 参数越少 → 测试搭建越简单
