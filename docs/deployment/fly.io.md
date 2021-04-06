# Deploy LogPaste to fly.io

fly.io is the best host I've found for LogPaste. It fits in the free tier, and you get a free SSL certificate.

## Pre-requisites

You'll need a fly.io account. You should have `fly` [already installed](https://fly.io/docs/fly/installing/) and authenticated on your machine.

## Clone the LogPaste repo

To begin, clone the LogPaste repository:

```bash
git clone https://github.com/mtlynch/logpaste.git && \
  cd logpaste
```

## Set your environment variables

Create environment variables for your AWS settings:

```bash
AWS_ACCESS_KEY_ID=YOUR-ACCESS-ID
AWS_SECRET_ACCESS_KEY=YOUR-SECRET-ACCESS-KEY
AWS_REGION=YOUR-REGION
DB_REPLICA_URL=s3://your-bucket-name/db
```

## Create your app

Next, create your app on fly.io:

```bash
RANDOM_SUFFIX="$(head /dev/urandom | tr -dc 'a-z0-9' | head -c 6 ; echo '')"
APP_NAME="logpaste-${RANDOM_SUFFIX}"

./dev-scripts/make-fly-config "${APP_NAME}"

fly init "${APP_NAME}" --nowrite
```

## Save AWS credentials

Use the `fly secrets set` command to securely save your AWS credentials to your app:

```bash
fly secrets set \
  "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
  "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}"
```

## Deploy

Finally, it's time to deploy your app. Run it once with `CREATE_NEW_DB='true'` so that LogPaste can bootstrap its database.

```bash
# Change this to the latest Docker image tag
LOGPASTE_IMAGE="mtlynch/logpaste:0.1.1"
fly deploy \
  --env "AWS_REGION=${AWS_REGION}" \
  --env "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  --env "CREATE_NEW_DB='true'" \
  --image "${LOGPASTE_IMAGE}"
```

After that command succeeds, deploy it without the `CREATE_NEW_DB` parameter. On all future deployments, deploy with this command:

```bash
fly deploy \
  --env "AWS_REGION=${AWS_REGION}" \
  --env "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  --image "${LOGPASTE_IMAGE}"

echo "Your LogPaste instance is now ready at: https://${APP_NAME}.fly.dev/"
```

You can test the instance with this command:

```bash
echo "hello, new fly.io instance!" | \
  curl -F '_=<-' "https://${APP_NAME}.fly.dev/"
```
