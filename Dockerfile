FROM alpine@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

# ca-certificates is needed to download files from S3
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Install the program
ADD credentials-sync /app/credentials-sync
WORKDIR /app

ENTRYPOINT ["/app/credentials-sync"]