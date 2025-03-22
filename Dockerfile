FROM golang:1.23.7-bookworm AS build

WORKDIR /app

COPY go.mod go.sum main.go ./

RUN go mod download && \
    go build -o test-runner main.go

FROM debian:bookworm

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        ca-certificates \
        git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

#RUN groupadd --gid 1000 noroot \
#    && useradd \
#    --create-home \
#    --home-dir /home/noroot \
#    --uid 1000 \
#    --gid 1000 \
#    --no-log-init noroot

WORKDiR /app

COPY --from=build /app/test-runner ./
COPY find.sh generate_report.sh test_from_list.sh ./

#USER noroot

CMD ["./test-runner"]

