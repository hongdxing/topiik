# Topiik
A midleware that for both Key/Value store and Event broker

## How Topiik works
### Controller Plane
Controller Plane maintain health of Topiik cluster, and forward requests to Workers
### Workers
Workers in charge of processing commands, managing memory, and persisting data

## Dev environment

## Minimum PROD node setting

![alt text](src/resource/minimum_prod_architecture.png)