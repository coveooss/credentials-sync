FROM alpine@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# ca-certificates is needed to download files from S3
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Install the program
ADD credentials-sync /app/credentials-sync
WORKDIR /app

ENTRYPOINT ["/app/credentials-sync"]