# coding-test-book
## Steps to run application
1. Clone the git repo https://github.com/AsnaNesry/coding-test-book.git
2. From the root folder of the project, run `docker build -t coding-test:latest .`
3. Once build is completed, run `docker-compose up`
4. [IMP] Make sure the docker startup log has below
 `go-app  | Database connected and table ensured.`.
5. If not, kill the process and run `docker-compose up` again and ensure the above log is present.

## Steps to test
### Postman
Postman collection is available in the project folder named `CodingTest.postman_collection.json`

### Unit Test
1. Unit test requires actual DB and not mocks. Ensure the db in docker is up while running unit tests.
2. From the project root, run `go test -v -coverprofile cover.out ./...`

### Assumptions
1. For the concurrency task, the PUT end point accepts array of Ids that need to be marked as done,  in the body.
