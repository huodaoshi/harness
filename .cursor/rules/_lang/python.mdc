---
description: Python 通用编码规范——类型提示、模块化、错误处理、命名
paths:
  - "**/*.py"
---

# Python 编码规范

## 风格

- 遵循 PEP 8（缩进、空格、行长度）
- 行长度上限以项目 lint 配置为准（typical: 88 / 100 / 120）
- 使用 4 空格缩进，不使用 Tab
- 文件编码 UTF-8，使用 LF 行尾

## 类型提示

- 公共函数 / 方法必须有类型注解（参数 + 返回值）
- 复杂类型用 `typing` 模块（或 Python 3.10+ 内置泛型语法）
- 容器类型明确元素类型：`list[str]` / `dict[str, int]`
- 可选类型：`X | None`（3.10+）或 `Optional[X]`
- 配合 `mypy` 或 `pyright` 进行类型检查

## 命名

- 模块 / 包：`lower_snake_case`
- 类：`PascalCase`
- 函数 / 变量：`lower_snake_case`
- 常量：`UPPER_SNAKE_CASE`
- 私有成员：`_leading_underscore`
- 强私有：`__double_leading_underscore`（触发 name mangling）

## 错误处理

- `try/except` 必须指定具体异常类型，不用裸 `except:`
- 不吞错误（至少 log）
- 自定义异常继承 `Exception`，不继承 `BaseException`
- 资源管理用 `with` 语句

## 模块化

- 一个模块单一职责
- `from module import name` 优于 `import module.name`（除非命名冲突）
- 不用通配符导入 `from module import *`
- 循环导入：通过函数内 import 或重构模块边界解决

## 数据类

- 数据容器用 `@dataclass` 或 `pydantic.BaseModel`（按项目约定）
- 不可变配置用 `frozen=True`

## 异步

- 异步函数用 `async def`
- 不在 async 函数中调用阻塞 IO（用 `asyncio.to_thread` 或异步等价 API）
- `asyncio.gather` 管理并发

## 测试

- 测试文件命名 `test_*.py` 或 `*_test.py`（按项目约定）
- 用 `pytest` 框架
- fixtures 用 `@pytest.fixture`

## 工具

- 格式化：`black` 或 `ruff format`
- Lint：`ruff` / `flake8` / `pylint`
- 类型检查：`mypy` / `pyright`
- 依赖管理：`poetry` / `pip-tools` / `uv`（按项目约定）
