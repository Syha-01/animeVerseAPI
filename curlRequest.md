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

BODY='{
    "title": "Spirited Away",
    "synopsis": "A classic adventure.",
    "total_episodes": 1,
    "status": "Finished Airing",
    "score": 8.8,
    "genres": ["Adventure", "Supernatural"],
    "studios": ["Studio Ghibli"]
}'

curl -i -H "Content-Type: application/json" -d "$BODY" localhost:4000/v1/animes
