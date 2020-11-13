# message-board-backend

to run tests:
docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name test_db -d postgres
go test
docker stop test_db
docker rm test_db

