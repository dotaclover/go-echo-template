# API Documentation

本文档描述当前模板默认提供的 HTTP API，包括健康检查、认证、个人中心和管理员用户管理接口。

## Base URL

开发环境默认地址：

```text
http://localhost:8080
```

## 统一响应结构

### 成功响应

```json
{
  "code": "success",
  "message": "Login successful",
  "data": {},
  "request_id": "f2fbef9c-9d15-4e66-bf4f-6b0540f7f95f"
}
```

### 创建成功响应

```json
{
  "code": "created",
  "message": "User created",
  "data": {},
  "request_id": "f2fbef9c-9d15-4e66-bf4f-6b0540f7f95f"
}
```

### 分页响应

```json
{
  "code": "success",
  "message": "Users fetched",
  "data": {
    "users": []
  },
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  },
  "request_id": "f2fbef9c-9d15-4e66-bf4f-6b0540f7f95f"
}
```

### 错误响应

```json
{
  "code": "validation_error",
  "message": "validation failed",
  "details": {
    "fields": {
      "username": "username is required"
    }
  },
  "request_id": "f2fbef9c-9d15-4e66-bf4f-6b0540f7f95f"
}
```

## 鉴权说明

需要登录的接口统一使用：

```text
Authorization: Bearer <access_token>
```

管理员接口要求：

- 已登录
- 用户角色为 `admin`

## 健康检查

### GET `/health/live`

用于存活探针。

#### Response

```json
{
  "code": "success",
  "message": "Service is live",
  "data": {
    "status": "live"
  }
}
```

### GET `/health/ready`

用于就绪探针，会检查数据库连通性。

#### Response

```json
{
  "code": "success",
  "message": "Service is ready",
  "data": {
    "status": "ready"
  }
}
```

## 认证接口

### POST `/api/v1/auth/register`

注册普通用户。

#### Request Body

```json
{
  "username": "newuser",
  "password": "StrongPass123",
  "email": "newuser@example.com",
  "real_name": "New User"
}
```

#### Response

```json
{
  "code": "created",
  "message": "User registered",
  "data": {
    "user": {
      "id": 1,
      "username": "newuser",
      "email": "newuser@example.com",
      "real_name": "New User",
      "role": "user",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T10:00:00Z"
    }
  }
}
```

### POST `/api/v1/auth/login`

用户登录，返回 access token 与 refresh token。

#### Request Body

```json
{
  "username": "admin",
  "password": "admin123"
}
```

#### Response

```json
{
  "code": "success",
  "message": "Login successful",
  "data": {
    "access_token": "<access-token>",
    "refresh_token": "<refresh-token>",
    "token_type": "Bearer",
    "expires_in": 604800,
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "real_name": "Administrator",
      "role": "admin",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T10:00:00Z"
    }
  }
}
```

### POST `/api/v1/auth/refresh`

使用 refresh token 刷新令牌对。

#### Request Body

```json
{
  "refresh_token": "<refresh-token>"
}
```

#### Response

```json
{
  "code": "success",
  "message": "Token refreshed",
  "data": {
    "access_token": "<new-access-token>",
    "refresh_token": "<new-refresh-token>",
    "token_type": "Bearer",
    "expires_in": 604800
  }
}
```

## 个人中心接口

### GET `/api/v1/auth/profile`

获取当前登录用户资料。

#### Headers

```text
Authorization: Bearer <access_token>
```

#### Response

```json
{
  "code": "success",
  "message": "Profile fetched",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "real_name": "Administrator",
      "role": "admin",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T10:00:00Z"
    }
  }
}
```

### PUT `/api/v1/auth/profile`

更新当前登录用户资料。

#### Request Body

```json
{
  "email": "updated@example.com",
  "real_name": "Updated Name"
}
```

#### Response

```json
{
  "code": "success",
  "message": "Profile updated",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "updated@example.com",
      "real_name": "Updated Name",
      "role": "admin",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T11:00:00Z"
    }
  }
}
```

### POST `/api/v1/auth/password`

修改当前登录用户密码。

#### Request Body

```json
{
  "old_password": "admin123",
  "new_password": "NewStrongPass123"
}
```

#### Response

```json
{
  "code": "success",
  "message": "Password updated",
  "request_id": "f2fbef9c-9d15-4e66-bf4f-6b0540f7f95f"
}
```

## 管理员用户管理

以下接口均要求管理员权限。

### GET `/api/v1/admin/users`

分页查询用户列表。

#### Query Params

- `page`: 页码，默认 `1`
- `page_size`: 每页数量，默认 `20`

#### Response

```json
{
  "code": "success",
  "message": "Users fetched",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "real_name": "Administrator",
        "role": "admin",
        "status": "active",
        "created_at": "2026-04-20T10:00:00Z",
        "updated_at": "2026-04-20T10:00:00Z"
      }
    ]
  },
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

### GET `/api/v1/admin/users/:id`

查询指定用户详情。

#### Response

```json
{
  "code": "success",
  "message": "User fetched",
  "data": {
    "user": {
      "id": 2,
      "username": "user01",
      "email": "user01@example.com",
      "real_name": "User 01",
      "role": "user",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T10:00:00Z"
    }
  }
}
```

### POST `/api/v1/admin/users`

管理员创建用户。

#### Request Body

```json
{
  "username": "ops",
  "password": "StrongPass123",
  "email": "ops@example.com",
  "real_name": "Ops User",
  "role": "admin",
  "status": "active"
}
```

#### Response

```json
{
  "code": "created",
  "message": "User created",
  "data": {
    "user": {
      "id": 3,
      "username": "ops",
      "email": "ops@example.com",
      "real_name": "Ops User",
      "role": "admin",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T10:00:00Z"
    }
  }
}
```

### PUT `/api/v1/admin/users/:id`

管理员更新用户基础资料。

#### Request Body

```json
{
  "email": "ops.updated@example.com",
  "real_name": "Ops Updated",
  "role": "admin",
  "status": "active"
}
```

#### Response

```json
{
  "code": "success",
  "message": "User updated",
  "data": {
    "user": {
      "id": 3,
      "username": "ops",
      "email": "ops.updated@example.com",
      "real_name": "Ops Updated",
      "role": "admin",
      "status": "active",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T11:00:00Z"
    }
  }
}
```

### PATCH `/api/v1/admin/users/:id/status`

管理员更新用户状态。

#### Request Body

```json
{
  "status": "disabled"
}
```

#### Response

```json
{
  "code": "success",
  "message": "User status updated",
  "data": {
    "user": {
      "id": 3,
      "username": "ops",
      "email": "ops.updated@example.com",
      "real_name": "Ops Updated",
      "role": "admin",
      "status": "disabled",
      "created_at": "2026-04-20T10:00:00Z",
      "updated_at": "2026-04-20T11:00:00Z"
    }
  }
}
```

### DELETE `/api/v1/admin/users/:id`

管理员删除用户。

#### Response

```json
{
  "code": "success",
  "message": "User deleted",
  "request_id": "f2fbef9c-9d15-4e66-bf4f-6b0540f7f95f"
}
```

## 常见错误码

- `bad_request`
- `validation_error`
- `unauthorized`
- `forbidden`
- `not_found`
- `conflict`
- `internal_error`
- `too_many_requests`

## 建议调试顺序

- 先调用 `/health/live`
- 再调用 `/health/ready`
- 注册或使用种子账号登录
- 保存 `access_token` 与 `refresh_token`
- 先验证 `/api/v1/auth/profile`
- 再验证管理员接口
