# Relatdb

Relatdb is a lightweight relational database management system implemented in Go. It implements the MySQL protocol and supports standard SQL statements, including DDL (Data Definition Language), DML (Data Manipulation Language), and DQL (Data Query Language).

## Key Features

- **SQL Parser**: Complete lexical and syntactic analysis supporting multiple SQL statements
- **Execution Engine**: Efficient executor supporting creation and management of databases, tables, and indexes
- **B+ Tree Index**: High-performance index implementation supporting unique indexes
- **Transaction Management**: Supports ACID properties with undo/redo logging
- **MySQL Protocol Compatibility**: Seamlessly connects with MySQL clients
- **Storage Engine**: Page-based storage system enabling data persistence

## System Architecture

### Core Modules

- **parser/** - SQL Parsing Layer
  - Lexical analyzer (lexer)
  - Syntax analyzer (parser)
  - Abstract Syntax Tree (AST) definitions

- **executor/** - SQL Execution Layer
  - Execution context management
  - Execution logic for various SQL statements
  - Result set handling

- **meta/** - Metadata Management
  - Database definitions
  - Table schema definitions
  - Field type definitions
  - Index descriptions

- **index/** - Index Implementation
  - B+ tree structure
  - Index entry management
  - Insertion, deletion, and lookup operations

- **store/** - Storage Layer
  - Page management
  - Data persistence
  - Storage engine interface

- **server/** - Server Layer
  - MySQL protocol implementation
  - Connection management
  - Authentication handling

- **transaction/** - Transaction Management
  - Transaction logs
  - Undo/redo mechanisms
  - Transaction state management

### Supported SQL Statements

#### DDL Statements
- `CREATE DATABASE` / `DROP DATABASE`
- `CREATE TABLE` / `DROP TABLE`
- `CREATE INDEX` / `DROP INDEX`

#### DML Statements
- `INSERT` - Insert data
- `UPDATE` - Update data
- `DELETE` - Delete data

#### DQL Statements
- `SELECT` - Query data
- Supports clauses such as `WHERE`, `ORDER BY`, `LIMIT`, `GROUP BY`, and `HAVING`

#### Other Statements
- `USE` - Select database
- `SET` - Set variables
- `SHOW` - Display information (databases, tables, variables, etc.)

## Installation and Execution

### Prerequisites

- Go 1.16 or higher

### Build and Run

```bash
# Clone the project
git clone https://gitee.com/QQXQQ/Relatdb.git
cd Relatdb

# Build
go build -o relatdb .

# Run
./relatdb
```

### Connect and Test

Connect using a MySQL client:

```bash
mysql -h 127.0.0.1 -P 3306 -u root -p
```

## Configuration

Adjust database parameters by modifying configuration constants in the source code:

- Server version number (`server/constant.go`)
- Default username and password (`server/constant.go`)
- Page size (`store/page.go`)

## Project Structure

```
Relatdb/
├── common/           # Common utilities (buffers, constant definitions)
├── executor/          # SQL execution engine
├── index/            # B+ tree index implementation
├── meta/             # Metadata definitions
├── parser/           # SQL parser (lexical, syntactic, AST)
├── server/           # MySQL protocol server
├── store/            # Storage engine
├── transaction/      # Transaction management
├── utils/            # Utility functions
├── main.go           # Program entry point
└── test.go           # Test entry point
```

## Usage Examples

### Create Database and Table

```sql
CREATE DATABASE testdb;
USE testdb;

CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100),
    age INT
);
```

### Insert Data

```sql
INSERT INTO users (id, name, age) VALUES (1, 'Zhang San', 25);
INSERT INTO users (id, name, age) VALUES (2, 'Li Si', 30);
```

### Query Data

```sql
SELECT * FROM users;
SELECT name, age FROM users WHERE age > 25 ORDER BY age;
```

### Update and Delete

```sql
UPDATE users SET age = 26 WHERE id = 1;
DELETE FROM users WHERE id = 2;
```

## Technical Highlights

1. **B+ Tree Index**: Fully implemented B+ tree structure supporting efficient range queries and scans
2. **Page-Based Storage**: Uses fixed-size pages for data management, enabling sequential read/write operations
3. **Transaction Logging**: Employs Write-Ahead Logging (WAL) to ensure data consistency
4. **MySQL Protocol**: Fully compatible with MySQL protocol, allowing standard MySQL clients to interact with the system

## Contribution Guidelines

Issues and pull requests are welcome.

## License

This project is licensed under the MIT License. See the LICENSE file for details.