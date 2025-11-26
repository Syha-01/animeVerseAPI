#!/bin/bash

# Base URL
BASE_URL="http://localhost:4000/v1"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "--------------------------------------------------"
echo "Starting Smart Caching Verification"
echo "--------------------------------------------------"

# 1. Register User
EMAIL="smartcache_test_$(date +%s)@example.com"
PASSWORD="password123"
echo -e "\n1. Registering user ($EMAIL)..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"SmartCacheUser\", \"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.user.id')
echo "User ID: $USER_ID"

# 2. Get Activation Token (Simulated from DB/Log - for this test we'll assume we can't easily get it without DB access, 
# so we'll rely on the fact that we can't activate easily in this script without parsing logs. 
# HOWEVER, for this specific test, we might need to be authenticated to add to list if we enforced it.
# Let's check if `createUserAnimeListHandler` requires authentication. 
# Looking at routes.go (from memory/context), it might not be strictly enforced yet or we can use the user_id in the body.
# The current handler takes `user_id` in the body.

# 3. Add NEW Anime to List (Smart Cache)
ANIME_ID=$((RANDOM + 10000)) # Random ID to ensure it's new
echo -e "\n2. Adding NEW anime (ID: $ANIME_ID) to list..."

BODY=$(cat <<EOF
{
  "user_id": "$USER_ID",
  "anime_id": $ANIME_ID,
  "status": "Watching",
  "current_episode": 1,
  "score": 8,
  "started_watching_date": "2023-11-25",
  "anime": {
    "title": "Smart Cache Test Anime",
    "synopsis": "This is a test anime created via smart caching.",
    "cover_image_url": "http://example.com/image.jpg",
    "total_episodes": 12,
    "status": "Airing",
    "release_date": "2023-01-01",
    "rating": "PG-13",
    "score": 7.5,
    "genres": ["Action", "Test"],
    "studios": ["Test Studio"],
    "broadcast_information": "Mondays"
  }
}
EOF
)

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/user_anime_list" \
  -H "Content-Type: application/json" \
  -d "$BODY")

HTTP_STATUS=$(echo "$RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')
BODY_ONLY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS:.*//')

if [ "$HTTP_STATUS" == "201" ]; then
  echo -e "${GREEN}Success! List entry created.${NC}"
  echo "$BODY_ONLY" | jq .
else
  echo -e "${RED}Failed to create list entry. Status: $HTTP_STATUS${NC}"
  echo "$BODY_ONLY"
  exit 1
fi

# 4. Verify Anime Exists (by trying to get it)
echo -e "\n3. Verifying anime was created in DB..."
ANIME_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$BASE_URL/animes/$ANIME_ID")
ANIME_STATUS=$(echo "$ANIME_RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

if [ "$ANIME_STATUS" == "200" ]; then
  echo -e "${GREEN}Success! Anime found in DB.${NC}"
else
  echo -e "${RED}Failed! Anime not found. Status: $ANIME_STATUS${NC}"
  exit 1
fi

echo -e "\n--------------------------------------------------"
echo -e "${GREEN}Smart Caching Verification Completed Successfully!${NC}"
echo "--------------------------------------------------"
