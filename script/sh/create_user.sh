curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' --request POST '127.0.0.1:8090/poetry/CreateUser' -d'{
 "username":"test",
 "password":"test123"
}'

curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' --request POST '127.0.0.1:8090/poetry/Login' -d'{
 "username":"test",
 "password":"test123"
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJpc3MiOiJwb2V0cnkiLCJleHAiOjE3NDg1MDUzOTIsIm5iZiI6MTc0ODQxODk5MiwiaWF0IjoxNzQ4NDE4OTkyfQ.CN3Vd5LH0FhEloBDa-rcdLbN3ouMvA5jIAEzaZyL-iw' \
--request POST '127.0.0.1:8090/poetry/CreateGameRoom' -d'{
 "maxPlayers":2,
 "password":"test123"
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc0ODY3NTQxOCwibmJmIjoxNzQ4NTg5MDE4LCJpYXQiOjE3NDg1ODkwMTh9.i1PYHKfR5mapsnn8s1SZJYwBNzxnopyV9n036lpIR8g' \
--request POST '127.0.0.1:8090/poetry/DescribeGameRoom' -d'{
 "roomId":"2a9f0a60-ae7"
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribeTagInfo' -d'{
    "filter":[{
        "name":"category",
        "value":["main"]
    }],
    "limit":200
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribeTagInfo' -d'{
    "filter":[{
        "name":"parent-tag-id",
        "value":["309","345","341"]
    }],
    "limit":200
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribePoetryInfo' -d'{
    "filter":[{
        "name":"tag-id",
        "value":["329"]
    }],
    "limit":200
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribePoetryInfo' -d'{
    "filter":[{
        "name":"title",
        "value":["陈情表"]
    }],
    "limit":200
}'

curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribePoetryInfo' -d'{
    "filter":[{
        "name":"author",
        "value":["李白"]
    }],
    "limit":200
}'

curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribePoetryInfo' -d'{
    "filter":[{
        "name":"dynasty",
        "value":["唐代"]
    }],
    "limit":200
}'


curl -H 'requestId: adfasdfasdfasdfasdf' -H'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1jNTMxZjUwOS01NzciLCJ1c2VyX25hbWUiOiJ1c2VyMSIsImlzcyI6InBvZXRyeSIsImV4cCI6MTc1NDM3ODkzNiwibmJmIjoxNzU0MjkyNTM2LCJpYXQiOjE3NTQyOTI1MzZ9.FBJyn0MUN2U4GK2VSr5DMI2wZFmWFCdMAMYB-4qcEM0' \
--request POST 'http://106.55.49.83:8090/poetry/DescribePoetryInfo' -d'{
    "filter":[{
        "name":"paragraphs",
        "value":["明月几时有"]
    }],
    "limit":200
}'
