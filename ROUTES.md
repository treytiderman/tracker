# Routes


## Public

```
/public/*
```


## Pages

> Auth token cookie required

```
/login

/trackers
/tracker-create

/tracker-info
/tracker-log
/tracker-records
/tracker-history

/entry-view

/settings
/test
```


## API - htmx

> Auth token cookie required

```
GET      /htmx/token   test auth token
POST     /htmx/token   returns auth token

POST     /htmx/tracker       create new tracker
PUT      /htmx/tracker/:id   update tracker by id
DELETE   /htmx/tracker/:id   delete tracker by id

POST     /htmx/entry/:tracker_id   create new entry
PUT      /htmx/entry/:tracker_id   update entry by tracker_id
DELETE   /htmx/entry/:tracker_id   delete entry by tracker_id

POST     /content-upload        upload file
GET      /content/{file_name}   upload file
DELETE   /content/{file_name}   delete file
```


## API - json

> Auth token header required

```
GET      /json/token   test auth token
POST     /json/token   returns auth token

GET      /json/tracker       get list of all trackers
POST     /json/tracker       create new tracker
GET      /json/tracker/:id   get tracker by id
PUT      /json/tracker/:id   update tracker by id
DELETE   /json/tracker/:id   delete tracker by id

GET      /json/entry               get list of all entries
GET      /json/entry/:tracker_id   get entry by tracker_id
POST     /json/entry/:tracker_id   create new entry
PUT      /json/entry/:tracker_id   update entry by tracker_id
DELETE   /json/entry/:tracker_id   delete entry by tracker_id

GET      /json/file         get list of all file names
GET      /json/file/:name   get file
POST     /json/file/:name   upload file
DELETE   /json/file/:name   delete file
```
