+++
title = "Folder HTTP API "
description = "Grafana Folder HTTP API"
keywords = ["grafana", "http", "documentation", "api", "folder"]
aliases = ["/http_api/folder/"]
type = "docs"
[menu.docs]
name = "Folder"
parent = "http_api"
+++

# Folder API

## Identifier (id) vs unique identifier (uid)

The identifier (id) of a folder is an auto-incrementing numeric value and is only unique per Grafana install.

The unique identifier (uid) of a folder can be used for uniquely identify folders between multiple Grafana installs. It's automatically generated if not provided when creating a folder. The uid allows having consistent URL's for accessing folders and when syncing folders between multiple Grafana installs. This means that changing the title of a folder will not break any bookmarked links to that folder.

The uid can have a maximum length of 40 characters.

## A note about the General folder

The General folder (id=0) is special and is not part of the Folder API which means
that you cannot use this API for retrieving information about the General folder.

## Get all folders

`GET /api/folders`

Returns all folders that the authenticated user has permission to view.

**Example Request**:

```http
GET /api/folders?limit=10 HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk
```

**Example Response**:

```http
HTTP/1.1 200
Content-Type: application/json

[
  {
    "id":1,
    "uid": "nErXDvCkzz",
    "title": "Departmenet ABC",
    "url": "/dashboards/f/nErXDvCkzz/department-abc",
    "hasAcl": false,
    "canSave": true,
    "canEdit": true,
    "canAdmin": true,
    "createdBy": "admin",
    "created": "2018-01-31T17:43:12+01:00",
    "updatedBy": "admin",
    "updated": "2018-01-31T17:43:12+01:00",
    "version": 1
  }
]
```

## Get folder by uid

`GET /api/folders/:uid`

Will return the folder given the folder uid.

**Example Request**:

```http
GET /api/folders/nErXDvCkzzh HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk
```

**Example Response**:

```http
HTTP/1.1 200
Content-Type: application/json

{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet ABC",
  "url": "/dashboards/f/nErXDvCkzz/department-abc",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
```

Status Codes:

- **200** – Found
- **401** – Unauthorized
- **403** – Access Denied
- **404** – Folder not found

## Create folder

`POST /api/folders`

Creates a new folder.

**Example Request**:

```http
POST /api/folders HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk

{
  "uid": "nErXDvCkzz",
  "title": "Department ABC"
}
```

JSON Body schema:

- **uid** – Optional [unique identifier](/http_api/folder/#identifier-id-vs-unique-identifier-uid).
- **title** – The title of the folder.

**Example Response**:

```http
HTTP/1.1 200
Content-Type: application/json

{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet ABC",
  "url": "/dashboards/f/nErXDvCkzz/department-abc",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
```

Status Codes:

- **200** – Created
- **400** – Errors (invalid json, missing or invalid fields, etc)
- **401** – Unauthorized
- **403** – Access Denied

## Update folder

`PUT /api/folders/:uid`

Updates an existing folder identified by uid.

**Example Request**:

```http
PUT /api/folders/nErXDvCkzz HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk

{
  "title":"Department DEF",
  "version": 1
}
```

JSON Body schema:

- **uid** – Provide another [unique identifier](/http_api/folder/#identifier-id-vs-unique-identifier-uid) than stored to change the unique identifier.
- **title** – The title of the folder.
- **version** – Provide the current version to be able to update the folder. Not needed if `overwrite=true`.
- **overwrite** – Set to true if you want to overwrite existing folder with newer version.

**Example Response**:

```http
HTTP/1.1 200
Content-Type: application/json

{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet DEF",
  "url": "/dashboards/f/nErXDvCkzz/department-def",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
```

Status Codes:

- **200** – Updated
- **400** – Errors (invalid json, missing or invalid fields, etc)
- **401** – Unauthorized
- **403** – Access Denied
- **404** – Folder not found
- **412** – Precondition failed

The **412** status code is used for explaing that you cannot update the folder and why.
There can be different reasons for this:

- The folder has been changed by someone else, `status=version-mismatch`

 The response body will have the following properties:

```http
HTTP/1.1 412 Precondition Failed
Content-Type: application/json; charset=UTF-8
Content-Length: 97

{
  "message": "The folder has been changed by someone else",
  "status": "version-mismatch"
}
```

## Delete folder

`DELETE /api/folders/:uid`

Deletes an existing folder identified by uid together with all dashboards stored in the folder, if any. This operation cannot be reverted.

**Example Request**:

```http
DELETE /api/folders/nErXDvCkzz HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk

```

**Example Response**:

```http
HTTP/1.1 200
Content-Type: application/json

{
  "message":"Folder deleted"
}
```

Status Codes:

- **200** – Deleted
- **401** – Unauthorized
- **403** – Access Denied
- **404** – Folder not found

## Get folder by id

`GET /api/folders/id/:id`

Will return the folder identified by id.

**Example Request**:

```http
GET /api/folders/id/1 HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk
```

**Example Response**:

```http
HTTP/1.1 200
Content-Type: application/json

{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet ABC",
  "url": "/dashboards/f/nErXDvCkzz/department-abc",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
```

Status Codes:

- **200** – Found
- **401** – Unauthorized
- **403** – Access Denied
- **404** – Folder not found