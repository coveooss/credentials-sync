FROM alpine

# ca-certificates is needed to download files from S3
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Install the program
ADD credentials-sync /app/credentials-sync
WORKDIR /app

ENTRYPOINT ["/app/credentials-sync"]