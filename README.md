# logpaste

[![CircleCI](https://circleci.com/gh/mtlynch/logpaste.svg?style=svg)](https://circleci.com/gh/mtlynch/logpaste) [![Docker Pulls](https://img.shields.io/docker/pulls/mtlynch/logpaste.svg?maxAge=604800)](https://hub.docker.com/r/mtlynch/logpaste/) [![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](LICENSE)

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

docker run \
  -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
  -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" \
  -e "AWS_REGION=${AWS_REGION}" \
  -e "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  -e "CREATE_NEW_DB='true'" `# change to false after first run` \
  --name logpaste \
  mtlynch/logpaste
```

## Run with custom site settings

```bash
SITE_TITLE="My Cool Log Pasting Service"
SITE_SUBTITLE="Upload all your logs for FooBar here"
SITE_SHOW_DOCUMENTATION="false" # Hide usage information from homepage

AWS_ACCESS_KEY_ID=YOUR-ACCESS-ID
AWS_SECRET_ACCESS_KEY=YOUR-SECRET-ACCESS-KEY
AWS_REGION=YOUR-REGION
DB_REPLICA_URL=s3://your-bucket-name/db

docker run \
  -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
  -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" \
  -e "AWS_REGION=${AWS_REGION}" \
  -e "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  -e "SITE_TITLE=${SITE_TITLE}" \
  -e "SITE_SUBTITLE=${SITE_SUBTITLE}" \
  -e "SITE_SHOW_DOCUMENTATION=${SITE_SHOW_DOCUMENTATION}" \
  -e "CREATE_NEW_DB='true'" `# change to false after first run` \
  --name logpaste \
  mtlynch/logpaste
```
