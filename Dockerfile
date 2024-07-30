FROM golang:1.22.5

WORKDIR /usr/src/app

COPY . .

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o bin/main cmd/main.go

CMD [ "./bin/main" ]