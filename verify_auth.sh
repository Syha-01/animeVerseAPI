#!/bin/bash

# Register a user
echo "Registering user..."
curl -s -X POST -d '{"username": "Test User", "email": "test@example.com", "password": "password123"}' http://localhost:4000/v1/users
echo -e "\n"

# Authenticate user
echo "Authenticating user..."
TOKEN_RESPONSE=$(curl -s -X POST -d '{"email": "test@example.com", "password": "password123"}' http://localhost:4000/v1/tokens/authentication)
echo $TOKEN_RESPONSE
echo -e "\n"

TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.authentication_token.token')

if [ "$TOKEN" == "null" ]; then
    echo "Failed to get token"
    exit 1
fi

echo "Token: $TOKEN"

# Test authenticated request
echo "Testing authenticated request..."
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:4000/v1/healthcheck
echo -e "\n"

# Test invalid token
echo "Testing invalid token..."
curl -s -i -H "Authorization: Bearer invalidtoken" http://localhost:4000/v1/healthcheck
echo -e "\n"

# Test missing token
echo "Testing missing token..."
curl -s -i http://localhost:4000/v1/healthcheck
echo -e "\n"
