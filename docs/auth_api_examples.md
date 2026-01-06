# Authentication API - Example Requests

---

## 1. POST /api/user/signup - User Registration

**Request:**
```bash
curl -X POST "http://localhost:8080/api/user/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "password": "Password123!",
    "role": "end_user"
  }'
```

**Response (Success):**
```json
{
  "code": 0,
  "message": "User created successfully",
  "data": null
}
```

**Response (Email already exists):**
```json
{
  "code": 1,
  "error_code": "email_already_exists",
  "message": "Email đã tồn tại",
  "internal": "email_already_exists"
}
```

**Response (Invalid input):**
```json
{
  "code": 1,
  "error_code": "invalid_input",
  "message": "Dữ liệu đầu vào không hợp lệ",
  "internal": "invalid_input"
}
```

---

## 2. POST /api/user/login - User Login (Email/Password)

**Request:**
```bash
curl -X POST "http://localhost:8080/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!"
  }'
```

**Response (Success):**
```json
{
  "code": 200,
  "message": "Login successfully",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3MzY0MzI4MDB9..."
}
```

**Response (Invalid credentials):**
```json
{
  "code": 1,
  "error_code": "invalid_email_or_password",
  "message": "Email hoặc mật khẩu không đúng",
  "internal": "invalid_email_or_password"
}
```

**Response (User inactive):**
```json
{
  "code": 1,
  "error_code": "user_inactive",
  "message": "Tài khoản đã bị vô hiệu hóa",
  "internal": "user_inactive"
}
```

---

## 3. POST /api/user/login - Login with Google OAuth

**Request (Default role - end_user):**
```bash
curl -X POST "http://localhost:8080/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFjOTQ1ZGRjMjAxNDYzYjEyZjJlMzFhNTQ0NjdmMDAwOWY0ZjI0M2IiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIxMjM0NTY3ODkwLWFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoiMTIzNDU2Nzg5MC1hYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ei5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjExMDAzMTQ5MDc1MzA5MDcwNzE1IiwiZW1haWwiOiJ1c2VyQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoiWjhURTVNRWJZcFFjdkJROE9XZnNRQSIsImlhdCI6MTY4NDc2NzQwMCwiZXhwIjoxNjg0NzcxMDAwfQ.signature"
  }'
```

**Request (Custom role - restaurant_owner):**
```bash
curl -X POST "http://localhost:8080/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjFjOTQ1ZGRjMjAxNDYzYjEyZjJlMzFhNTQ0NjdmMDAwOWY0ZjI0M2IiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIxMjM0NTY3ODkwLWFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoiMTIzNDU2Nzg5MC1hYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ei5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjExMDAzMTQ5MDc1MzA5MDcwNzE1IiwiZW1haWwiOiJ1c2VyQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoiWjhURTVNRWJZcFFjdkJROE9XZnNRQSIsImlhdCI6MTY4NDc2NzQwMCwiZXhwIjoxNjg0NzcxMDAwfQ.signature",
    "role": "restaurant_owner"
  }'
```

**Response (Success - New user created):**
```json
{
  "code": 200,
  "message": "Login successfully",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZ21haWwuY29tIiwiZXhwIjoxNzM2NDMyODAwfQ..."
}
```

**Response (Invalid Google token):**
```json
{
  "code": 1,
  "error_code": "invalid_google_auth_token",
  "message": "Token Google không hợp lệ",
  "internal": "invalid_google_auth_token"
}
```

**Notes:**
- If `role` is not provided, default role `end_user` will be assigned
- Available roles: `end_user`, `admin`, `restaurant_owner`, `staff`
- Role must exist in database

---

## 4. POST /api/auth/request-reset-password - Request OTP for Password Reset

**Request:**
```bash
curl -X POST "http://localhost:8080/api/auth/request-reset-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

**Response (Success - OTP sent via email):**
```json
{
  "code": 200,
  "message": "OTP email sent successfully",
  "data": ""
}
```

**Response (Email not found):**
```json
{
  "code": 1,
  "error_code": "email_not_found",
  "message": "Email không tồn tại",
  "internal": "email_not_found"
}
```

**Note:** 
- OTP will be sent to the registered email address
- OTP expires in 5 minutes
- Previous unverified OTP will be invalidated automatically

---

## 5. POST /api/auth/validate-otp - Validate OTP Code

**Request:**
```bash
curl -X POST "http://localhost:8080/api/auth/validate-otp" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'
```

**Response (Success):**
```json
{
  "code": 200,
  "message": "OTP validated successfully",
  "data": ""
}
```

**Response (Invalid OTP):**
```json
{
  "code": 1,
  "error_code": "invalid_otp",
  "message": "OTP không chính xác",
  "internal": "invalid_otp"
}
```

**Response (OTP expired):**
```json
{
  "code": 1,
  "error_code": "otp_expired",
  "message": "OTP đã hết hạn",
  "internal": "otp_expired"
}
```

**Response (OTP not found):**
```json
{
  "code": 1,
  "error_code": "otp_not_found",
  "message": "Không tìm thấy OTP",
  "internal": "otp_not_found"
}
```

---

## 6. POST /api/auth/reset-password - Reset Password

**Request:**
```bash
curl -X POST "http://localhost:8080/api/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "new_password": "NewPassword123!"
  }'
```

**Response (Success):**
```json
{
  "code": 200,
  "message": "Password reset successfully",
  "data": ""
}
```

**Response (OTP not verified):**
```json
{
  "code": 1,
  "error_code": "otp_not_verified",
  "message": "OTP chưa được xác thực",
  "internal": "otp_not_verified"
}
```

**Response (Email not found):**
```json
{
  "code": 1,
  "error_code": "email_not_found",
  "message": "Email không tồn tại",
  "internal": "email_not_found"
}
```

---

## Complete Password Reset Flow Example

### Step 1: Request OTP
```bash
curl -X POST "http://localhost:8080/api/auth/request-reset-password" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

### Step 2: Validate OTP (Check email for OTP code)
```bash
curl -X POST "http://localhost:8080/api/auth/validate-otp" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'
```

### Step 3: Reset Password
```bash
curl -X POST "http://localhost:8080/api/auth/reset-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "new_password": "NewPassword123!"
  }'
```

---

## Error Codes Reference

| Code | HTTP | Message (EN) | Message (VI) |
|------|------|--------------|-------------|
| invalid_input | 400 | Invalid input | Dữ liệu đầu vào không hợp lệ |
| email_already_exists | 400 | Email already exists | Email đã tồn tại |
| invalid_email_or_password | 401 | Invalid email or password | Email hoặc mật khẩu không đúng |
| user_inactive | 403 | User account is inactive | Tài khoản đã bị vô hiệu hóa |
| invalid_google_auth_token | 401 | Invalid Google authentication token | Token Google không hợp lệ |
| email_not_found | 404 | Email not found | Email không tồn tại |
| otp_not_found | 404 | OTP not found | Không tìm thấy OTP |
| otp_expired | 400 | OTP has expired | OTP đã hết hạn |
| invalid_otp | 400 | Invalid OTP | OTP không chính xác |
| otp_not_verified | 400 | OTP not verified | OTP chưa được xác thực |
| otp_already_verified | 400 | OTP already verified | OTP đã được xác thực |

---

## Notes

1. **JWT Token**: The token returned from login should be included in subsequent requests as Bearer token:
   ```bash
   -H "Authorization: Bearer <token>"
   ```

2. **Password Requirements**: 
   - Minimum 6 characters
   - Should contain uppercase, lowercase, numbers, and special characters

3. **OTP Validity**: 
   - OTP is valid for 5 minutes
   - Only one active OTP per email can exist
   - OTP is 6 digits numeric

4. **Email Configuration**:
   - Make sure `mail.api_key` is configured in `config.yaml`
   - SendGrid API integration for email delivery

5. **Role Assignment**:
   - Default role on signup is `end_user`
   - Available roles: `end_user`, `admin`, `restaurant_owner`, `staff`

---

## Authentication Flow Diagram
┌─────────────────────────────────────────────────────────────┐
│                     FRONTEND (React)                         │
│                                                               │
│  1. User clicks "Login with Google" button                   │
│     ↓                                                         │
│  2. Call Google Sign-In API                                  │
│     ↓                                                         │
│  3. Google shows login dialog                                │
│     ↓                                                         │
│  4. User authenticates with Google                           │
│     ↓                                                         │
│  5. Get ID Token from Google                                 │
│     ↓                                                         │
│  6. Send ID Token to Backend                                 │
└─────────────────────────────────────────────────────────────┘
                            ↓
                            │ POST /api/user/login
                            │ { id_token: "..." }
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  BACKEND (Go/Gin)                            │
│                                                               │
│  1. Receive ID Token from Frontend                           │
│     ↓                                                         │
│  2. Verify ID Token with Google API                          │
│     ↓                                                         │
│  3. Extract user info (email, name, etc.)                    │
│     ↓                                                         │
│  4. Check if user exists in DB                               │
│     ├─ YES → Login user, generate JWT                        │
│     └─ NO → Create new user, generate JWT                    │
│     ↓                                                         │
│  5. Return JWT Token to Frontend                             │
└─────────────────────────────────────────────────────────────┘
                            ↓
                            │ { jwt_token: "..." }
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                     FRONTEND (React)                         │
│                                                               │
│  1. Receive JWT Token from Backend                           │
│     ↓                                                         │
│  2. Save JWT to localStorage/cookies                         │
│     ↓                                                         │
│  3. Use JWT for all future API requests                      │
│     ↓                                                         │
│  4. Redirect to dashboard                                    │
└─────────────────────────────────────────────────────────────┘

### Email/Password Flow
```
┌─────────────────────────────────────────────────────────────────────────┐
│                          USER SIGNUP/LOGIN FLOW                          │
└─────────────────────────────────────────────────────────────────────────┘

SIGNUP FLOW:
┌─────────────────────────────────────────────────────────────────────────┐
│ Frontend                      Backend                      Database      │
│                                                                           │
│ 1. User Input Data            2. Validate Input                          │
│    (email, password, ...)         - Email format                         │
│         │                         - Password strength                    │
│         ├────────POST─────────────>                                      │
│         │  /api/user/signup         │                                    │
│         │                           ├───Hash Password (bcrypt)           │
│         │                           │                                    │
│         │                           ├──Check Email Exists──────────────> │
│         │                           │  SELECT * FROM users               │
│         │                           │  WHERE email = ?                   │
│         │                           │                                    │
│         │                           │ <─────Email Not Found──────────────┤
│         │                           │                                    │
│         │                           ├──Create New User──────────────────>│
│         │                           │  INSERT INTO users (...)           │
│         │                           │                                    │
│         │                           │ <─────User Created─────────────────┤
│         │                           │                                    │
│         │ <──────200 OK──────────────┤                                    │
│         │  {code: 0, message: ...}   │                                    │
│         │                            │                                    │
└─────────────────────────────────────────────────────────────────────────┘

LOGIN FLOW (Email/Password):
┌─────────────────────────────────────────────────────────────────────────┐
│ Frontend                      Backend                      Database      │
│                                                                           │
│ 1. User Input                 2. Validate & Authenticate                │
│    (email, password)              │                                      │
│         │                         ├──Find User──────────────────────────>│
│         ├────────POST─────────────>  SELECT * FROM users                │
│         │  /api/user/login        │  WHERE email = ?                    │
│         │                         │                                      │
│         │                         │ <────User Found────────────────────┤
│         │                         │                                      │
│         │                         ├──Verify Password (bcrypt)            │
│         │                         │  compare(password, hash)             │
│         │                         │                                      │
│         │                         ├──Generate JWT Token (HS256)          │
│         │                         │  {id, email, role, exp: +24h}       │
│         │                         │                                      │
│         │ <──────200 OK──────────────┤                                    │
│         │  {code: 200,              │                                    │
│         │   data: "<JWT_TOKEN>"}    │                                    │
│         │                            │                                    │
│ 3. Save Token to localStorage        │                                    │
│    (for future requests)             │                                    │
│         │                            │                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

### Google OAuth Flow
```
┌──────────────────────────────────────────────────────────────────────────────┐
│                        GOOGLE OAUTH 2.0 LOGIN FLOW                            │
└──────────────────────────────────────────────────────────────────────────────┘

STEP 1: Frontend calls Google Sign-In
┌──────────────────────────────────────────────────────────────────────────────┐
│ User Browser              Frontend (React)           Google Auth Server      │
│      │                           │                          │                │
│      │ Click "Sign in with      │                          │                │
│      │  Google" button          │                          │                │
│      ├──────────────────────────>│                          │                │
│      │                           │                          │                │
│      │                           ├─────OAuth Request───────>│                │
│      │                           │ (Client ID, Scopes, etc) │                │
│      │                           │                          │                │
│      │ <─────Redirect to Google Sign-In────────────────────┤                │
│      │ (Google Login Form)       │                          │                │
│      │                           │                          │                │
└──────────────────────────────────────────────────────────────────────────────┘

STEP 2: User authenticates with Google
┌──────────────────────────────────────────────────────────────────────────────┐
│ User Browser              Frontend (React)           Google Auth Server      │
│      │                           │                          │                │
│      │ Enter Email/Password      │                          │                │
│      │ & Grant Permission        │                          │                │
│      ├──────────────────────────>│                          │                │
│      │                           │                          │                │
│      │                           │ <─User Authenticated──────┤               │
│      │                           │                          │                │
│      │ <───Receive ID Token──────┤                          │                │
│      │ (JWT from Google)         │                          │                │
│      │                           │                          │                │
└──────────────────────────────────────────────────────────────────────────────┘

STEP 3: Frontend sends ID Token to Backend
┌──────────────────────────────────────────────────────────────────────────────┐
│ Frontend (React)           Backend Server              Database              │
│      │                           │                         │                 │
│      │ Extract ID Token from     │                         │                 │
│      │ credentialResponse        │                         │                 │
│      │                           │                         │                 │
│      ├──POST /api/user/login────>│                         │                 │
│      │  {id_token: "JWT..."}     │                         │                 │
│      │                           │                         │                 │
│      │                           ├──Verify ID Token──────> │                 │
│      │                           │ Call Google API:        │                 │
│      │                           │ https://oauth2.googleapis│                │
│      │                           │ /tokeninfo?id_token=... │                 │
│      │                           │                         │                 │
│      │                           │ <──Verified─────────────│                 │
│      │                           │ {email, name, sub, ...} │                 │
│      │                           │                         │                 │
│      │                           ├──Check User Exists──────>│                 │
│      │                           │ SELECT * FROM users     │                 │
│      │                           │ WHERE email = ?         │                 │
│      │                           │                         │                 │
│      │                           │ <─User Not Found────────│                 │
│      │                           │ (For first-time login)  │                 │
│      │                           │                         │                 │
│      │                           ├──Create New User───────>│                 │
│      │                           │ INSERT INTO users       │                 │
│      │                           │ email, name, provider   │                 │
│      │                           │ = 'google'              │                 │
│      │                           │                         │                 │
│      │                           │ <─User Created──────────│                 │
│      │                           │                         │                 │
│      │                           ├──Generate JWT Token     │                 │
│      │                           │ (Backend JWT, HS256)    │                 │
│      │                           │ {id, email, role, exp}  │                 │
│      │                           │                         │                 │
│      │ <─────200 OK──────────────┤                         │                 │
│      │ {code: 200,               │                         │                 │
│      │  data: "<BACKEND_JWT>"}   │                         │                 │
│      │                           │                         │                 │
│ Save Backend JWT to localStorage │                         │                 │
│ for authenticated requests       │                         │                 │
│      │                           │                         │                 │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Password Reset Flow
```
┌──────────────────────────────────────────────────────────────────────────────┐
│                           PASSWORD RESET FLOW                                 │
└──────────────────────────────────────────────────────────────────────────────┘

STEP 1: Request OTP
┌──────────────────────────────────────────────────────────────────────────────┐
│ Frontend                Backend                    Database        Email      │
│      │                      │                         │              │        │
│      │ POST /api/auth/      │                         │              │        │
│      │ request-reset-       ├──Find User─────────────>│              │        │
│      │ password             │ WHERE email = ?         │              │        │
│      │ {email: "..."}       │                         │              │        │
│      │                      │ <─User Found────────────│              │        │
│      ├─────────────────────>│                         │              │        │
│      │                      ├─Generate OTP────────────>│              │        │
│      │                      │ 6-digit code, 5min exp  │              │        │
│      │                      │                         │              │        │
│      │                      │ <─OTP Created───────────│              │        │
│      │                      │                         │              │        │
│      │                      ├─Send Email─────────────────────────────>        │
│      │                      │ (via SendGrid)          │              │        │
│      │                      │ Subject: "Reset Code"   │              │        │
│      │                      │ Body: "Your OTP: xxxxxx"│              │        │
│      │                      │                         │              │        │
│      │ <──200 OK────────────┤                         │              │        │
│      │                      │                         │              │        │
└──────────────────────────────────────────────────────────────────────────────┘

STEP 2: Validate OTP
┌──────────────────────────────────────────────────────────────────────────────┐
│ Frontend                Backend                    Database                   │
│      │                      │                         │                       │
│      │ User receives OTP    │                         │                       │
│      │ from email           │                         │                       │
│      │                      │                         │                       │
│      │ POST /api/auth/      │                         │                       │
│      │ validate-otp         ├──Find OTP──────────────>│                       │
│      │ {email, otp: "xxx"}  │ WHERE email = ? AND     │                       │
│      │                      │       type = 'reset'    │                       │
│      ├─────────────────────>│                         │                       │
│      │                      │ <─OTP Found─────────────│                       │
│      │                      │                         │                       │
│      │                      ├──Validate OTP           │                       │
│      │                      │ - Check expiry          │                       │
│      │                      │ - Verify code match     │                       │
│      │                      │ - Check verified flag   │                       │
│      │                      │                         │                       │
│      │                      ├──Mark OTP Verified─────>│                       │
│      │                      │ UPDATE otps SET         │                       │
│      │                      │ verified = true         │                       │
│      │                      │                         │                       │
│      │                      │ <─Updated──────────────│                       │
│      │                      │                         │                       │
│      │ <──200 OK────────────┤                         │                       │
│      │                      │                         │                       │
└──────────────────────────────────────────────────────────────────────────────┘

STEP 3: Reset Password
┌──────────────────────────────────────────────────────────────────────────────┐
│ Frontend                Backend                    Database                   │
│      │                      │                         │                       │
│      │ User enters          │                         │                       │
│      │ new password         │                         │                       │
│      │                      │                         │                       │
│      │ POST /api/auth/      │                         │                       │
│      │ reset-password       ├──Find OTP──────────────>│                       │
│      │ {email, new_pass}    │ WHERE email = ? AND     │                       │
│      │                      │       verified = true   │                       │
│      ├─────────────────────>│                         │                       │
│      │                      │ <─OTP Verified──────────│                       │
│      │                      │                         │                       │
│      │                      ├──Hash New Password      │                       │
│      │                      │ bcrypt(password)        │                       │
│      │                      │                         │                       │
│      │                      ├──Update User────────────>│                       │
│      │                      │ UPDATE users SET        │                       │
│      │                      │ password = ? WHERE      │                       │
│      │                      │ email = ?               │                       │
│      │                      │                         │                       │
│      │                      │ <─Password Updated──────│                       │
│      │                      │                         │                       │
│      │                      ├──Invalidate OTP─────────>│                       │
│      │                      │ DELETE FROM otps        │                       │
│      │                      │ WHERE email = ?         │                       │
│      │                      │                         │                       │
│      │ <──200 OK────────────┤                         │                       │
│      │                      │                         │                       │
│ User can now login with     │                         │                       │
│ new password                │                         │                       │
│      │                      │                         │                       │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## Using JWT Token in Requests

After successful login/authentication, use the returned JWT token in subsequent requests:

```bash
# Include token in Authorization header
curl -X GET "http://localhost:8080/api/user/profile" \
  -H "Authorization: Bearer <JWT_TOKEN_HERE>"
```

**Token Structure (HS256):**
```json
{
  "id": "user-uuid-here",
  "email": "user@example.com",
  "role": "end_user",
  "iat": 1704067200,
  "exp": 1704153600,
  "iss": "smart-restaurant-api"
}
```

---

## Key Integration Points

### Frontend Libraries
- **Google Sign-In**: Use `@react-oauth/google` package
- **HTTP Client**: Use `fetch` or `axios`
- **Token Storage**: Use `localStorage` or `sessionStorage`

### Backend Integration
- **Google Token Verification**: Call `https://oauth2.googleapis.com/tokeninfo?id_token={token}`
- **Password Hashing**: bcrypt with default cost
- **JWT Generation**: HS256 with 24-hour expiry
- **Email Service**: SendGrid API for OTP delivery

### Security Considerations
1. **Never expose JWT secret** - Store in environment variables
2. **Always use HTTPS** - For OAuth redirects and token transmission
3. **Validate on server side** - Never trust client-side validation alone
4. **Set secure cookie flags** - If using cookies for token storage
5. **Implement rate limiting** - On login/OTP endpoints to prevent brute force
6. **Use strong passwords** - Enforce password complexity requirements
