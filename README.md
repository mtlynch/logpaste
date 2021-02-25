# logpaste

[![CircleCI](https://circleci.com/gh/mtlynch/logpaste.svg?style=svg)](https://circleci.com/gh/mtlynch/logpaste)

A minimalist web service for uploading and sharing log files.

## Run locally

```bash
go run main.go
```

## Run in local Docker container

The Docker container automatically replicates to an AWS S3 bucket

```bash
AWS_ACCESS_KEY_ID=YOUR-ACCESS-ID
AWS_SECRET_ACCESS_KEY=YOUR-SECRET-ACCESS-KEY
AWS_REGION=YOUR-REGION
DB_REPLICA_URL=s3://your-bucket-name/db

docker build -t logpaste . && \
docker run \
  -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
  -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" \
  -e "AWS_REGION=${AWS_REGION}" \
  -e "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  -e "CREATE_NEW_DB='true'" \ `# change to false after first run`
  --name logpaste \
  logpaste
```
