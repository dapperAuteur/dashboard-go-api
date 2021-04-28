# Dashboard App

This is the api to a dashboard app. The dashboard will track expenses, income, assets, tasks, and growth of the user.

It's written in Golang following the Ardan Labs way of building a service.

The goal is to build a tool that adds value to the users life, mine. And to help me get one of those high paying Golang dev jobs!

## First Version

- track and manage expenses

## Run App Locally

#### Start Docker First
docker-compose up -d --remove-orphans
(shut down with this command) docker-compose down


to run mongodb locally (this command isn't needed because Docker starts the db)
`sudo mongod --dbpath /System/Volumes/Data/data/db`

<!-- local ports to use -->
zipkin to trace requests
http://localhost:9411/zipkin/

server/service
http://localhost:8080/v1