# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o backend-services

# final stage
FROM alpine
WORKDIR /app
EXPOSE 8000
COPY --from=build-env /src/backend-services /app/
ENTRYPOINT ./backend-services