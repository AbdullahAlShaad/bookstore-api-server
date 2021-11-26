# syntax=docker/dockerfile:1

# Base Image
FROM golang:1.16-buster AS build  

# Default destination for all subsequent command
# Directory inside the image
WORKDIR /app

# Copy dependencies into the image
COPY go.mod ./
COPY go.sum ./

#download dependencies in the image
RUN go mod download 

#copy all the .go files in the image
COPY . ./ 

# Building the binary
RUN go build -o /bookstore-api-server-docker .

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bookstore-api-server-docker /bookstore-api-server-docker 

EXPOSE 8081

USER nonroot:nonroot

CMD ["/bookstore-api-server-docker"]


