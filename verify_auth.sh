#!/bin/bash

# Start server
go run ./cmd/api &
PID=$!
sleep 5

# 1. Anonymous Access (Expect 401)
echo "Testing Anonymous Access..."
CODE=$(curl -s -o /dev/null -w "%{http_code}" localhost:4000/v1/animes)
if [ "$CODE" -eq 401 ]; then
    echo "PASS: Anonymous access denied (401)"
else
    echo "FAIL: Anonymous access got $CODE"
fi

sleep 1

# 2. Register User
echo "Registering User..."
EMAIL="test_auth_$(date +%s)@example.com"
RESPONSE=$(curl -s -X POST localhost:4000/v1/users -d "{\"username\": \"testuser\", \"email\": \"$EMAIL\", \"password\": \"password123\"}")
ID=$(echo $RESPONSE | jq -r '.user.id')
echo "User ID: $ID"

sleep 1

# 3. Manually Activate User
echo "Manually Activating User..."
psql "postgres://animeverse:verse1@localhost/animeverse?sslmode=disable" -c "UPDATE users SET activated = true WHERE id = '$ID';"

sleep 1

# 4. Login (Get Token)
echo "Logging in..."
RESPONSE=$(curl -s -X POST localhost:4000/v1/tokens/authentication -d "{\"email\": \"$EMAIL\", \"password\": \"password123\"}")
AUTH_TOKEN=$(echo $RESPONSE | jq -r '.authentication_token.token')

if [ "$AUTH_TOKEN" == "null" ]; then
    echo "FAIL: Login failed"
    echo $RESPONSE
    kill $PID
    exit 1
fi

sleep 1

# 5. Access Protected Resource (Read)
echo "Testing Read Access..."
CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $AUTH_TOKEN" localhost:4000/v1/animes)
if [ "$CODE" -eq 200 ]; then
    echo "PASS: Read access granted (200)"
else
    echo "FAIL: Read access got $CODE"
fi

sleep 1

# 6. Access Protected Resource (Write)
echo "Testing Write Access..."
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Authorization: Bearer $AUTH_TOKEN" localhost:4000/v1/animes -d '{}')
if [ "$CODE" -eq 403 ]; then
    echo "PASS: Write access denied (403)"
else
    echo "FAIL: Write access got $CODE"
fi

kill $PID
