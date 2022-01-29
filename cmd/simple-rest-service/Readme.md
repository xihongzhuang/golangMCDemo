##build a docker image 
```
cd cmd/simple-rest-service
make build
```
note: dockerize this API service will allow this service to be deployed by kubernetes to the cloud

##start the docker instance that hosts this service
```
make run
```
##Insert two valid entries in the database in the docker
```
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata1.yaml --url "http://localhost:8080/api/appmetadata"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata2.yaml --url "http://localhost:8080/api/appmetadata"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/invalid_appmetadata1.yaml --url "http://localhost:8080/api/appmetadata"
should response : "invalid maintainer email"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/invalid_appmetadata2.yaml --url "http://localhost:8080/api/appmetadata"
should response : empty version
```
##Get All entries in the database so far
```
curl http://localhost:8080/api/appmetadata
```
##Update full context of an APP metadata entry
```
curl -X PUT --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata2.yaml --url "http://localhost:8080/api/appmetadata/juEYZbNwpF5arCJJMpDQBc"
```
##Patch the APP metadata (partially update)
```
curl -X PATCH --header 'content-type: application/x-yml' --data-binary @./data/valid_metadata1Patch.yaml --url "http://localhost:8080/api/appmetadata/abu24BqqnPYdViRFBNihvA"
```
##Test outside docker:
```
go test -v ./internal/api-service
```
##Query by conditions
```
curl "http://localhost:8080/api/appmetadata?title=Valid&version=1.0.1"
```
##field contains substr, syntax: field=in:substr
```
search for records in which the company field contains substring "microsoft"
curl "http://localhost:8080/api/appmetadata?company=in:Microsoft"
```
