

# Relatdb

简易数据库实现，用于学习（未完成）

## 主要特性

- **SQL 解析器**：完整的词法分析和语法解析，支持多种 SQL 语句
- **执行引擎**：高效的执行器，支持数据库、表、索引的创建与管理
- **B+ 树索引**：高性能的索引实现，支持唯一索引
- **事务管理**：支持事务的 ACID 特性，包含 undo/redo 日志
- **MySQL 协议兼容**：可以与 MySQL 客户端无缝连接
- **存储引擎**：基于页的存储系统，支持数据的持久化

## 系统架构

### 核心模块

- **parser/** - SQL 解析层
  - 词法分析器 (lexer)
  - 语法分析器 (parser)
  - 抽象语法树 (AST) 定义

- **executor/** - SQL 执行层
  - 执行上下文管理
  - 各类 SQL 语句的执行逻辑
  - 结果集处理

- **meta/** - 元数据管理
  - 数据库定义
  - 表结构定义
  - 字段类型定义
  - 索引描述

- **index/** - 索引实现
  - B+ 树结构
  - 索引条目管理
  - 索引的插入、删除、查找

- **store/** - 存储层
  - 页管理
  - 数据持久化
  - 存储引擎接口

- **server/** - 服务层
  - MySQL 协议实现
  - 连接管理
  - 认证处理

- **transaction/** - 事务处理
  - 事务日志
  - undo/redo 机制
  - 事务状态管理

### 支持的 SQL 语句

#### DDL 语句
- `CREATE DATABASE` / `DROP DATABASE`
- `CREATE TABLE` / `DROP TABLE`
- `CREATE INDEX` / `DROP INDEX`

#### DML 语句
- `INSERT` - 插入数据
- `UPDATE` - 更新数据
- `DELETE` - 删除数据

#### DQL 语句
- `SELECT` - 查询数据
- 支持 `WHERE`、`ORDER BY`、`LIMIT`、`GROUP BY`、`HAVING` 等子句

#### 其他语句
- `USE` - 选择数据库
- `SET` - 设置变量
- `SHOW` - 显示信息（数据库、表、变量等）

## 安装与运行

### 环境要求

- Go 1.16 或更高版本

### 编译运行

```bash
# 克隆项目
git clone https://gitee.com/QQXQQ/Relatdb.git
cd Relatdb

# 编译
go build -o relatdb .

# 运行
./relatdb
```

### 连接测试

使用 MySQL 客户端连接：

```bash
mysql -h 127.0.0.1 -P 3306 -u root -p
```

## 配置说明

可以通过修改源码中的配置常量来调整数据库参数：

- 服务器版本号 (`server/constant.go`)
- 默认用户密码 (`server/constant.go`)
- 页大小 (`store/page.go`)

## 项目结构

```
Relatdb/
├── common/           # 公共工具（缓冲区、常量定义）
├── executor/          # SQL 执行引擎
├── index/            # B+ 树索引实现
├── meta/             # 元数据定义
├── parser/           # SQL 解析器（词法、语法、AST）
├── server/           # MySQL 协议服务器
├── store/            # 存储引擎
├── transaction/      # 事务管理
├── utils/            # 工具函数
├── main.go           # 程序入口
└── test.go           # 测试入口
```

## 使用示例

### 创建数据库和表

```sql
CREATE DATABASE testdb;
USE testdb;

CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100),
    age INT
);
```

### 插入数据

```sql
INSERT INTO users (id, name, age) VALUES (1, '张三', 25);
INSERT INTO users (id, name, age) VALUES (2, '李四', 30);
```

### 查询数据

```sql
SELECT * FROM users;
SELECT name, age FROM users WHERE age > 25 ORDER BY age;
```

### 更新和删除

```sql
UPDATE users SET age = 26 WHERE id = 1;
DELETE FROM users WHERE id = 2;
```

## 技术亮点

1. **B+ 树索引**：实现了完整的 B+ 树结构，支持高效的区间查询和范围扫描
2. **页式存储**：采用固定大小的页进行数据管理，支持数据的顺序读写
3. **事务日志**：使用 Write-Ahead Logging (WAL) 机制，确保数据一致性
4. **MySQL 协议**：完全兼容 MySQL 协议，可使用标准 MySQL 客户端操作

## 贡献指南

欢迎提交 Issue 和 Pull Request。

## 许可证

本项目采用 MIT 许可证，详情请查看 LICENSE 文件。