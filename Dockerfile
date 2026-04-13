FROM alpine@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

# ca-certificates is needed to download files from S3
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Install the program
ADD credentials-sync /app/credentials-sync
WORKDIR /app

ENTRYPOINT ["/app/credentials-sync"]