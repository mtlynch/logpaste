# Deploy LogPaste to Heroku

Heroku is a fantastic match for LogPaste, as you can run it in a free dyno and pay only pennies per month for data storage on S3.

## Overview

You can deploy an on-demand LogPaste instance to Heroku under their free tier. It means that Heroku shuts your app down after a few hours of inactivity, but LogPaste handles this fine because it syncs all of its data to S3. The downside is that the first request to the LogPaste server after inactivity will . Heroku doesn't offer SSL certificates for dynos in its free tier, so you'll have to upgrade to a paid plan if you want a HTTPS URL for your domain.

## Pre-requisites

You'll need:

* A Heroku account
* The [heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) installed and authenticated on your machine
* Docker installed on your machine


## Set your environment variables

To begin, create environment variables for your AWS settings:

```bash
AWS_ACCESS_KEY_ID=YOUR-ACCESS-ID
AWS_SECRET_ACCESS_KEY=YOUR-SECRET-ACCESS-KEY
AWS_REGION=YOUR-REGION
DB_REPLICA_URL=s3://your-bucket-name/db
```

## Configure your Heroku app

First, log in to the Heroku container registry:

```bash
heroku container:login
```

Next, create a container-based app:

```bash
RANDOM_SUFFIX="$(head /dev/urandom | tr -dc 'a-z0-9' | head -c 6 ; echo '')"
APP_NAME="logpaste-${RANDOM_SUFFIX}"
heroku apps:create "${APP_NAME}" --stack container
```

Assign all the relevant environment variables to your app:

```bash
heroku config:set --app "${APP_NAME}" AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}"
heroku config:set --app "${APP_NAME}" AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}"
heroku config:set --app "${APP_NAME}" AWS_REGION="${AWS_REGION}"
heroku config:set --app "${APP_NAME}" DB_REPLICA_URL="${DB_REPLICA_URL}"
```

## Deploy

Finally, deploy your app:

```bash
# Change this to the latest Docker image tag
LOGPASTE_IMAGE="mtlynch/logpaste:0.1.1"

HEROKU_PROCESS_TYPE="web"
HEROKU_IMAGE_PATH="registry.heroku.com/${APP_NAME}/${HEROKU_PROCESS_TYPE}"
docker pull "${LOGPASTE_IMAGE}" && \
  docker tag "${LOGPASTE_IMAGE}" "${HEROKU_IMAGE_PATH}" && \
  docker push "${HEROKU_IMAGE_PATH}" && \
  heroku container:release "${HEROKU_PROCESS_TYPE}" --app "${APP_NAME}" && \
  echo "Your app is running at http://${APP_NAME}.herokuapp.com"
```
