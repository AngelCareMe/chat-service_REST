@echo off
REM test-api.bat - скрипт для тестирования API на Windows

set BASE_URL=http://localhost:8080/api/v1
set TOKEN=

echo 🚀 Starting API testing...

REM 1. Health check
echo 1. Testing health check...
curl -s -X GET "%BASE_URL%/health"
echo.

REM 2. Register new user
echo 2. Registering new user...
curl -s -X POST "%BASE_URL%/register" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"testuser\", \"email\": \"test@example.com\", \"password\": \"password123\"}"
echo.

REM 3. Login with same user
echo 3. Testing login...
curl -s -X POST "%BASE_URL%/login" ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"test@example.com\", \"password\": \"password123\"}"
echo.

echo.
echo 🎉 Manual testing completed! Check responses above.
echo You can now test other endpoints manually with the tokens received.