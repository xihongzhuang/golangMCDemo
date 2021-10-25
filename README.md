# golangMCDemo
This is a coding exercise for Microsoft
This is just a simple RESTful API Service. Authentication and Authorization are out of scope here.
It hosts a RESTful API Server for application metadata in a inmemory database.
Supported endpoints:
http://localhost:8080/api/appmetadata
Support standard CRUD Method: GET, POST, PUT, PATCH, DELETE 

# build a docker image

```
cd cmd/simple-rest-service
make build
```
note: dockerize this API service will allow this service to be deployed by kubernetes to the cloud

# start the docker instance that hosts this service

```
make run
```
# Insert two valid entries in the database in the docker

```
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata1.yaml --url "http://localhost:8080/api/appmetadata"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata2.yaml --url "http://localhost:8080/api/appmetadata"
```
Upon success of each curl POST, an entry with an uniquely assigned id will be created in the API server, and respond in yaml.
 
# Verify the following invalid input should return errors
```
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/invalid_appmetadata1.yaml --url "http://localhost:8080/api/appmetadata"
should response : "invalid maintainer email"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/invalid_appmetadata2.yaml --url "http://localhost:8080/api/appmetadata"
should response : empty version
```
# Get All entries in the database so far

```
curl http://localhost:8080/api/appmetadata
```

# Update full context of an APP metadata entry
```
curl -X PUT --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata2.yaml --url "http://localhost:8080/api/appmetadata/juEYZbNwpF5arCJJMpDQBc"
```

# Patch the APP metadata (partially update)

```
curl -X PATCH --header 'content-type: application/x-yml' --data-binary @./data/valid_metadata1Patch.yaml --url "http://localhost:8080/api/appmetadata/abu24BqqnPYdViRFBNihvA"
```
# Delete an Entry
```
curl -X DELETE "http://localhost:8080/api/appmetadata/$entry_id"
```

# Unit Test outside the docker:

```
go test -v ./internal/api-service
```

# Query by conditions

```
curl "http://localhost:8080/api/appmetadata?title=Valid&version=1.0.1"
```

# field contains substr, syntax: field=in:substr

search for records in which the company field contains substring "microsoft"
```
curl "http://localhost:8080/api/appmetadata?company=in:Microsoft"
```
search for records whose version=1.0.1
```
curl "http://localhost:8080/api/appmetadata?version=1.0.1"
```
