# graphql-files-s3

[![Build Status](https://travis-ci.org/graphql-services/graphql-files-s3.svg?branch=master)](https://travis-ci.org/graphql-services/graphql-files-s3)

Service for handling file uploads/downloads in cooperation with GraphQL API.

Works seamlessly with [graphql-files service](https://github.com/graphql-services/graphql-files) or with any other GraphQL API containing following query/mutation/input:

```
query getFile($id: ID!) {
    result: file(id: $id) {
        id
        name
        size
        contentType
    }
}
mutation createFile($input: FileCreateInput!) {
    result: createFile(input:$input) {
        id
        name
        size
        contentType
    }
}

input FileCreateInput {
  id: ID
  name: String
  size: Int
  contentType: String
}
```

# How it works

Instead of using GraphQL API server to handle uploads, You can start this service along with the API. All transfers between client and S3 is done using presigned URLs (upload/download).
The main advantage of this approach is that You don't have to expose access tokens in URL to download the file as You can use the presigned URL directly and it will expire automatically.

### Uploading

1. send request to this service `POST /upload`
   - with `application/json` content like `{filename: "myfile.png", size: 42, contentType: "image/png"}`
   - optionally provide `Authorization` header which will be passed to GraphQL API
1. service generates UUID for given file and generates S3 presigned URL for uploading
1. service sends create file mutation to Your GraphQL API
1. as response the data from mutation along with the `uploadURL` is returned back
1. proceed with the upload directly to S3 using presigned URL (returned in `uploadURL`)

### Downloading

The process is similar as for uploading:

1. send request to graphql-files-s3 service `GET /{id}`
1. service send query getFile query with given id (if file not found, 404 is returned)
1. service generates S3 presignedURL for downloading the file
1. as response the `url` is returned
1. you can directly open the URL (return in `url`)

## Using custom fields

You can also provide custom fields for mutation input. The way it works is that You can provide additional values in querystring and they will be passed to createFile `$input`. For example:

```
POST /upload?blah=foo
{
    filename:"profile.png",
    size:3563,
    contentType:"image/png"
}
```

Will send input:

```
{
    filename:"profile.png",
    size:3563,
    contentType:"image/png",
    blah:"foo"
}
```

This can be used for example if You need to pass additional information to Your GraphQL API.
