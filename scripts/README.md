# 数据库初始化

包含数据库初始化脚本 `init.sql`，用于创建项目所需的表结构。

## 使用方法

```bash
# 1. 创建数据库
mysql -u root -p -e "CREATE DATABASE homebox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. 执行初始化脚本
mysql -u root -p homebox < scripts/init.sql
```

或在 MySQL 客户端中执行：

```sql
CREATE DATABASE homebox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE homebox;
SOURCE scripts/init.sql;
```

## 数据库表结构

| 表名 | 说明 |
|------|------|
| users | 用户账户信息 |
| spaces | 用户空间 |
| tag_groups | 标签分组 |
| tags | 标签 |
| space_members | 空间成员关系 |
| space_invites | 空间邀请 |

### users（用户表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint unsigned | 用户ID（主键）|
| name | varchar(50) | 用户名 |
| email | varchar(255) | 邮箱（唯一）|
| password | varchar(255) | 密码（加密）|
| avatar | varchar(255) | 头像URL |
| status | tinyint | 状态：1-正常，0-禁用 |
| created_at | datetime(3) | 创建时间 |
| updated_at | datetime(3) | 更新时间 |
| deleted_at | datetime(3) | 删除时间 |

### spaces（空间表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint unsigned | 空间ID（主键）|
| user_id | bigint unsigned | 拥有者用户ID |
| name | varchar(100) | 空间名称 |
| icon | varchar(50) | 空间图标（emoji）|
| description | varchar(500) | 空间描述 |
| created_at | datetime(3) | 创建时间 |
| updated_at | datetime(3) | 更新时间 |
| deleted_at | datetime(3) | 删除时间 |

### tag_groups（标签分组表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint unsigned | 分组ID（主键）|
| space_id | bigint unsigned | 空间ID |
| name | varchar(100) | 分组名称 |
| created_at | datetime(3) | 创建时间 |
| updated_at | datetime(3) | 更新时间 |
| deleted_at | datetime(3) | 删除时间 |

### tags（标签表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint unsigned | 标签ID（主键）|
| space_id | bigint unsigned | 空间ID |
| name | varchar(100) | 标签名称 |
| group_id | bigint unsigned | 所属分组ID |
| color | varchar(20) | 标签颜色 |
| description | varchar(500) | 标签描述 |
| order_no | int | 排序序号 |
| created_at | datetime(3) | 创建时间 |
| updated_at | datetime(3) | 更新时间 |
| deleted_at | datetime(3) | 删除时间 |

### space_members（空间成员表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint unsigned | 成员ID（主键）|
| space_id | bigint unsigned | 空间ID |
| user_id | bigint unsigned | 用户ID |
| is_owner | tinyint(1) | 是否为所有者 |
| pin | tinyint(1) | 是否固定 |
| joined_at | datetime(3) | 加入时间 |
| created_at | datetime(3) | 创建时间 |
| updated_at | datetime(3) | 更新时间 |
| deleted_at | datetime(3) | 删除时间 |

### space_invites（空间邀请表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint unsigned | 邀请ID（主键）|
| space_id | bigint unsigned | 空间ID |
| invitee_id | bigint unsigned | 被邀请人用户ID |
| inviter_id | bigint unsigned | 邀请人用户ID |
| status | varchar(20) | 状态：pending/accepted/rejected |
| created_at | datetime(3) | 创建时间 |
| updated_at | datetime(3) | 更新时间 |
| deleted_at | datetime(3) | 删除时间 |

## 注意事项

- 所有表使用 `utf8mb4` 字符集
- 所有表支持软删除（`deleted_at` 字段）
- 时间字段精度为毫秒级（`datetime(3)`）
