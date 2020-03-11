FROM golang as buile
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o surge .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=buile /app/surge .
COPY data data
COPY migration migration
CMD ["./surge"]