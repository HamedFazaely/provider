## Providers Service
This service mocks multiple providers. It stores orders status in memory. Current provider names ars : hamed, ali, ahmad, erfan, ilia. It is important to use these names when creating providers in logistics service and set their satus endpoints and pickup endpoints according to these names

## How to run
Build the image:
`docker build -t provider:v1.0.0 .`

Create a docker network: `docker network create -d bridge podro`

Run sms container: `docker run -d --rm --name provider -p 9090:9090 --network=podro provider:v1.0.0`
