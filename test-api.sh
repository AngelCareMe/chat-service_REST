#!/bin/bash

# test-api.sh - скрипт для тестирования API
set -e

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "🚀 Starting API testing..."

# 1. Health check
echo "1. Testing health check..."
curl -s -X GET "http://localhost:8080/health"
echo ""

# 2. Register new user
echo "2. Registering new user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

echo "Register response:"
echo "$REGISTER_RESPONSE"
echo ""

# Extract token from register response
TOKEN=$(echo "$REGISTER_RESPONSE" -r '.data.session.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Failed to get token from registration"
    exit 1
fi

echo "✅ Got auth token: ${TOKEN:0:20}..."
echo ""

# 3. Login with same user
echo "3. Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

echo "Login response:"
echo "$LOGIN_RESPONSE"
echo ""

# Extract token from login response
LOGIN_TOKEN=$(echo "$LOGIN_RESPONSE"-r '.data.session.token')

if [ "$LOGIN_TOKEN" = "null" ] || [ -z "$LOGIN_TOKEN" ]; then
    echo "❌ Failed to get token from login"
    exit 1
fi

echo "✅ Got login token: ${LOGIN_TOKEN:0:20}..."
echo ""

# 4. Get user profile
echo "4. Getting user profile..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "Profile response:"
echo "$PROFILE_RESPONSE" 
echo ""

# 5. Create a message
echo "5. Creating a message..."
MESSAGE_RESPONSE=$(curl -s -X POST "$BASE_URL/messages" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LOGIN_TOKEN" \
  -d '{
    "content": "Hello, this is my first message!"
  }')

echo "Message response:"
echo "$MESSAGE_RESPONSE" 
echo ""

# Extract message ID
MESSAGE_ID=$(echo "$MESSAGE_RESPONSE" -r '.data.id')

if [ "$MESSAGE_ID" = "null" ] || [ -z "$MESSAGE_ID" ]; then
    echo "❌ Failed to get message ID"
    exit 1
fi

echo "✅ Created message with ID: $MESSAGE_ID"
echo ""

# 6. Get all messages
echo "6. Getting all messages..."
ALL_MESSAGES=$(curl -s -X GET "$BASE_URL/messages" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "All messages:"
echo "$ALL_MESSAGES"
echo ""

# 7. Get user's messages
echo "7. Getting user messages..."
USER_MESSAGES=$(curl -s -X GET "$BASE_URL/messages/my" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "User messages:"
echo "$USER_MESSAGES"
echo ""

# 8. Update profile
echo "8. Updating profile..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LOGIN_TOKEN" \
  -d '{
    "username": "testuser_updated"
  }')

echo "Update response:"
echo "$UPDATE_RESPONSE"
echo ""

# 9. Logout
echo "9. Logging out..."
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/logout" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "Logout response:"
echo "$LOGOUT_RESPONSE"
echo ""

echo "🎉 All tests completed successfully!"