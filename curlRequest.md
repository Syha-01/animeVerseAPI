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


//-------------------------------- Validation -------------------------------------------//

test for empty title

curl -X POST -H "Content-Type: application/json" -d '{"title": ""}' localhost:4000/v1/animes


test for invalid episode count

curl -X POST -H "Content-Type: application/json" -d '{"title": "Test Anime", "total_episodes": -1}' localhost:4000/v1/animes

test for invalid score
curl -X POST -H "Content-Type: application/json" -d '{"title": "Test Anime", "score": 11}' localhost:4000/v1/animes
