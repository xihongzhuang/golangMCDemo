#!/bin/bash -x

cd cmd/simple-rest-service
##build a docker image
#make build

##start the docker instance that hosts this service
#make run

URL="http://localhost:8080/api/appmetadata"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata1.yaml --url "${URL}"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadata2.yaml --url "${URL}"
curl -X POST --header 'content-type: application/x-yml' --data-binary @../../internal/data/valid_metadataMicrosoft.yaml --url "${URL}"
##Get All entries in the database so far
curl ${URL}
#Patch the metadata (partially update)
##curl -X PATCH --header 'content-type: application/x-yml' --data-binary @./data/valid_metadata1Patch.yaml --url "http://localhost:8080/api/appmetadata/abu24BqqnPYdViRFBNihvA"

