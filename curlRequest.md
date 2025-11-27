///----------------------- Healtcheck ----------------------------------//

curl -i localhost:4000/v1/healthcheck

Output:

{
    "status": "available",
    "system_info": {
        "environment": "development",
        "version": "1.0.0"
    }
}


///----------------------- Error Handling ----------------------------------//

404 Not Found

curl -i localhost:4000/v1/nonexistent

405 Method Not Allowed

curl -i -X POST localhost:4000/v1/healthcheck

//-------------------------------- Testing Animes Ednpoint -------------------------------------------//

// Create a new anime
BODY='{
    "id": 1,
    "title": "Spirited Away",
    "synopsis": "During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits, and where humans are changed into beasts.",
    "cover_image_url": "https://cdn.myanimelist.net/images/anime/6/79597.jpg",
    "total_episodes": 1,
    "status": "Finished Airing",
    "release_date": "2001-07-20",
    "rating": "PG",
    "score": 8.78,
    "genres": ["Adventure", "Supernatural", "Drama"],
    "studios": ["Studio Ghibli"],
    "broadcast_information": "Fridays at 21:00 (JST)"
}'

curl -i -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/animes

BODY='{
    "id": 2,
    "title": "Princess Mononoke",
    "synopsis": "On a journey to find the cure for a Tatarigami''s curse, Ashitaka finds himself in the middle of a war between the forest gods and Tatara, a mining colony. In this quest he also meets San, the Mononoke Hime.",
    "cover_image_url": "https://cdn.myanimelist.net/images/anime/7/75919.jpg",
    "total_episodes": 1,
    "status": "Finished Airing",
    "release_date": "1997-07-12",
    "rating": "PG-13",
    "score": 8.69,
    "genres": ["Action", "Adventure", "Fantasy"],
    "studios": ["Studio Ghibli"],
    "broadcast_information": "Fridays at 21:00 (JST)"
}'

curl -i -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/animes

BODY='{
    "id": 16498,
    "title": "Attack on Titan",
    "synopsis": "Centuries ago, mankind was slaughtered to near extinction by monstrous humanoid creatures called Titans, forcing humans to hide in fear behind enormous concentric walls. What makes these giants truly terrifying is that their taste for human flesh is not born out of hunger but what appears to be out of pleasure.",
    "cover_image_url": "https://cdn.myanimelist.net/images/anime/10/47347.jpg",
    "total_episodes": 87,
    "status": "Finished Airing",
    "release_date": "2013-04-07",
    "rating": "R - 17+ (violence & profanity)",
    "score": 8.53,
    "genres": ["Action", "Military", "Mystery", "Super Power", "Drama", "Fantasy"],
    "studios": ["Wit Studio", "MAPPA"],
    "broadcast_information": "Sundays at 00:10 (JST)"
}'

curl -i -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/animes

BODY='{
    "id": 1535,
    "title": "Death Note",
    "synopsis": "A shinigami, as a god of death, can kill any person—provided they see their victim''s face and write their victim''s name in a notebook called a Death Note. One day, Ryuk, bored by the shinigami lifestyle and interested in seeing how a human would use a Death Note, drops one into the human realm.",
    "cover_image_url": "https://cdn.myanimelist.net/images/anime/9/9453.jpg",
    "total_episodes": 37,
    "status": "Finished Airing",
    "release_date": "2006-10-04",
    "rating": "R - 17+ (violence & profanity)",
    "score": 8.62,
    "genres": ["Mystery", "Police", "Psychological", "Supernatural", "Thriller"],
    "studios": ["Madhouse"],
    "broadcast_information": "Wednesdays at 00:56 (JST)"
}'

curl -i -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/animes

BODY='{
    "id": 5114,
    "title": "Fullmetal Alchemist: Brotherhood",
    "synopsis": "After a horrific alchemy experiment goes wrong in the Elric household, brothers Edward and Alphonse are left in a catastrophic new reality. Ignoring the alchemical principle banning human transmutation, the boys attempted to bring their deceased mother back to life.",
    "cover_image_url": "https://cdn.myanimelist.net/images/anime/1223/96541.jpg",
    "total_episodes": 64,
    "status": "Finished Airing",
    "release_date": "2009-04-05",
    "rating": "R - 17+ (violence & profanity)",
    "score": 9.14,
    "genres": ["Action", "Military", "Adventure", "Comedy", "Drama", "Magic", "Fantasy"],
    "studios": ["Bones"],
    "broadcast_information": "Sundays at 17:00 (JST)"
}'

curl -i -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/animes

// Get an anime by ID
curl -i localhost:4000/v1/animes/2

// Get a non-existent anime
curl -i localhost:4000/v1/animes/999

//delet and anime with specific ID
curl -X DELETE localhost:4000/v1/animes/2

// Get all animes
curl -i localhost:4000/v1/animes

//-------------------------------- Sorting -------------------------------------------//

//-------------------------------- Filtering and Searching Animes -------------------------------------------//

// The `GET /v1/animes` endpoint now supports filtering via URL query string parameters.
// This allows you to search for animes by title and/or genres.

// -- Available Parameters --
//
// title (string):
//   Performs a case-insensitive, partial-word search on the anime's title.
//   Example: `?title=titan` will match "Attack on Titan".
//
// genres (comma-separated string):
//   Filters for animes that contain ALL of the specified genres. The search is case-sensitive.
//   Example: `?genres=Action,Drama` will match animes that have BOTH "Action" and "Drama" in their genres list.

// -- Examples --

// Get all animes (no filters)
// This is the base request. It returns all animes in the database.
curl -i "localhost:4000/v1/animes"

// Search for animes with 'titan' in the title
// This uses the `title` parameter to perform a full-text search.
curl -i "localhost:4000/v1/animes?title=titan"

// Get all animes that have both the 'Action' AND 'Drama' genres
// This uses the `genres` parameter with a comma-separated list.
// The database will look for records where the genres array contains both values.
curl -i "localhost:4000/v1/animes?genres=Action,Drama"

// Get animes with 'note' in the title that also have the 'Thriller' genre
// This shows how to combine multiple parameters. The filters are joined with an AND condition.
// Note the use of the '&' to separate different query parameters.
curl -i "localhost:4000/v1/animes?title=note&genres=Thriller"

//-------------------------------- Pagination -------------------------------------------//

curl -i "localhost:4000/v1/animes?page=1&page_size=2"

curl -i "localhost:4000/v1/animes?page=2&page_size=2"



//-------------------------------- Validation -------------------------------------------//

test for empty title

curl -X POST -H "Content-Type: application/json" -d '{"title": ""}' localhost:4000/v1/animes


test for invalid episode count

curl -X POST -H "Content-Type: application/json" -d '{"title": "Test Anime", "total_episodes": -1}' localhost:4000/v1/animes

test for invalid score
curl -X POST -H "Content-Type: application/json" -d '{"title": "Test Anime", "score": 11}' localhost:4000/v1/animes

//-------------------------------- Rate limiting -------------------------------------------//

for i in {1..8}; do curl -i localhost:4000/v1/healthcheck; echo ""; done

run with disabled rate limiting
go run ./cmd/api -port=4000 -env=development -limiter-enabled=false -db-dsn="postgres://animeverse:verse1@localhost/animeverse?sslmode=disable"


//-------------------------------- User Anime List -------------------------------------------//

// 1. Create Entry (Standard)
// Uses the user_id "f86f401b-57ec-4712-adbb-aefe02a717cf" as an example
BODY='{
    "user_id": "f86f401b-57ec-4712-adbb-aefe02a717cf",
    "anime_id": 1,
    "status": "Watching",
    "current_episode": 1,
    "score": 8,
    "started_watching_date": "2023-01-01"
}'

curl -i -X POST -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/user_anime_list


// 1b. Create Entry with Smart Caching (Add New Anime)
// If the anime does not exist in the database, you can provide its details to create it on the fly.
BODY_SMART='{
    "user_id": "f86f401b-57ec-4712-adbb-aefe02a717cf",
    "anime_id": 99999,
    "status": "Watching",
    "current_episode": 1,
    "score": 8,
    "anime": {
        "title": "New Anime from Jikan",
        "synopsis": "This anime was fetched from Jikan and cached automatically.",
        "cover_image_url": "https://example.com/image.jpg",
        "total_episodes": 12,
        "status": "Airing",
        "release_date": "2023-10-01",
        "rating": "PG-13",
        "score": 7.5,
        "genres": ["Action", "Sci-Fi"],
        "studios": ["New Studio"],
        "broadcast_information": "Fridays"
    }
}'

curl -i -X POST -H "Content-Type: application/json" -d "$BODY_SMART" localhost:4000/v1/user_anime_list

// 2. Get Entry
// Replace ":id" with the UUID of the created entry, the UUID will be given by the response when creating an anime
curl -i localhost:4000/v1/user_anime_list/:id

// 3. Update Entry
// Replace ":id" with the UUID of the entry you want to update
BODY='{
    "status": "Completed",
    "score": 9,
    "current_episode": 12,
    "finished_watching_date": "2023-01-05"
}'

curl -i -X PATCH -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/user_anime_list/ff401b-57ec-4712-adbb-aefe02a717cf

// 4. List Entries for user
// Filter by user_id and/or status. Pagination and sorting are also supported.
curl -i "localhost:4000/v1/user_anime_list?user_id=f86f401b-57ec-4712-adbb-aefe02a717cf&status=Completed"

// 5. Delete Entry
// Replace ":id" with the UUID of the entry you want to delete
curl -i -X DELETE localhost:4000/v1/user_anime_list/:id


// 6. Sorting
// Sort by score (descending)
curl -i "localhost:4000/v1/user_anime_list?user_id=f86f401b-57ec-4712-adbb-aefe02a717cf&sort=-score"

// Sort by most recently updated
curl -i "localhost:4000/v1/user_anime_list?user_id=f86f401b-57ec-4712-adbb-aefe02a717cf&sort=-created_at"

// 7. Pagination
// Get the second page of results, with 5 items per page
curl -i "localhost:4000/v1/user_anime_list?user_id=f86f401b-57ec-4712-adbb-aefe02a717cf&page=2&page_size=5"

// 8. Error Handling Examples

// Invalid Status
BODY='{
    "user_id": "f86f401b-57ec-4712-adbb-aefe02a717cf",
    "anime_id": 1,
    "status": "InvalidStatus",
    "current_episode": 1
}'
curl -i -X POST -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/user_anime_list

// Score out of range (must be 1-10)
BODY='{
    "user_id": "f86f401b-57ec-4712-adbb-aefe02a717cf",
    "anime_id": 1,
    "status": "Watching",
    "score": 11
}'
curl -i -X POST -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/user_anime_list

//-------------------------------- Rate Limiting -------------------------------------------//


for i in {1..8}; do curl -i localhost:4000/v1/healthcheck; echo ""; done


//-------------------------------- Creating Users -------------------------------------------//

// 1. Successful Registration
BODY='{
    "username": "Test User",
    "email": "test@example.com",
    "password": "password123"
}'
curl -i -X POST -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/users

// 2. Duplicate Email Test
// Run the same command again to verify duplicate email handling
curl -i -X POST -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/users

// 3. Validation Error Test (Invalid Data)
BODY_INVALID='{
    "username": "",
    "email": "invalid-email",
    "password": "short"
}'
curl -i -X POST -H "Content-Type: application/json" -d "$BODY_INVALID" localhost:4000/v1/users


// 4. Activate User
// Replace "35GDX654ZCHMPJQR5PNDV3UJIE" with the token received in the email
BODY_ACTIVATE='{
    "token": "35GDX654ZCHMPJQR5PNDV3UJIE"
}'
curl -i -X PUT -H "Content-Type: application/json" -d "$BODY_ACTIVATE" localhost:4000/v1/users/activated


//-------------------------------- Token Types Explanation -------------------------------------------//
// There are two different types of tokens used in this API:
//
// 1. Activation Token (Email Token)
//    - Purpose: Verify your email address.
//    - Source: Sent to your email inbox after registration.
//    - Usage: Sent ONCE to `PUT /v1/users/activated`.
//    - Result: Your account status changes to "activated". It does NOT log you in.
//
// 2. Authentication Token (Bearer Token)
//    - Purpose: Prove your identity for API requests (Login).
//    - Source: Returned by `POST /v1/tokens/authentication` in exchange for email/password.
//    - Usage: Sent in the `Authorization` header for EVERY protected request.
//    - Result: The API knows who you are (User vs Anonymous).

//-------------------------------- Authentication & Anonymous Users -------------------------------------------//

// 1. Authenticate (Login)
// Exchange your credentials for a Bearer token.
// The token is valid for 24 hours.
BODY_LOGIN='{
    "email": "test@example.com",
    "password": "password123"
}'
curl -i -X POST -H "Content-Type: application/json" -d "$BODY_LOGIN" localhost:4000/v1/tokens/authentication

// Response will contain the token:
// {
//     "authentication_token": {
//         "token": "OP35GDX654ZCHMPJQR5PNDV3UJIE",
//         "expiry": "2025-11-27T10:00:00Z"
//     }
// }


// 2. Authenticated Request
// Use the token from the previous step in the Authorization header.
// Replace "OP35GDX654ZCHMPJQR5PNDV3UJIE" with the actual token.
curl -i -H "Authorization: Bearer OP35GDX654ZCHMPJQR5PNDV3UJIE" localhost:4000/v1/healthcheck

// Get Current User Profile
curl -i -H "Authorization: Bearer OP35GDX654ZCHMPJQR5PNDV3UJIE" localhost:4000/v1/users/me


// 3. Anonymous Request
// If you do not provide an Authorization header, the API treats you as an Anonymous User.
// Currently, this allows access to public endpoints like healthcheck.
curl -i localhost:4000/v1/healthcheck


// 4. Invalid Token
// If you provide an invalid or expired token, you will receive a 401 Unauthorized error.
curl -i -H "Authorization: Bearer invalid_token" localhost:4000/v1/healthcheck





//--force db migration--//

migrate -path ./migrations -database "postgres://animeverse:verse1@localhost/animeverse?sslmode=disable" force 7


//-------------------------------- Verified Test Run (Syha) -------------------------------------------//

// 1. Activate Account (Token: LY77HSSIL2IE7JRSSKUHWYXY7U)
curl -i -X PUT -H "Content-Type: application/json" -d '{"token": "LY77HSSIL2IE7JRSSKUHWYXY7U"}' localhost:4000/v1/users/activated

// 2. Authenticate/Login (john@gmail.com)
curl -i -X POST -H "Content-Type: application/json" -d '{"email": "john@gmail.com", "password": "mypassword1234"}' localhost:4000/v1/tokens/authentication

curl -X PUT http://localhost:4000/v1/users/activated \
  -H "Content-Type: application/json" \
  -d '{"token": "LMEBSMPKGIKLOUIKV7MEYGFO6Q"}'