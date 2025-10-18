FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM gcr.io/distroless/base
COPY --from=builder /app/server /server
CMD ["/server"]
