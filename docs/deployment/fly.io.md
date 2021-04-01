# Deploy LogPaste to fly.io

fly.io is the best host I've found for LogPaste. It fits in the free tier, and you get a free SSL certificate.

## Pre-requisites

You'll need a fly.io account. You should have `flyctl` [already installed](https://fly.io/docs/flyctl/installing/) and logged in on your machine.

## Set your environment variables

To begin, fill in your AWS settings:

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

# Create fly.toml file
echo "app = \"${APP_NAME}\"" > fly.toml
tail +5 fly.toml.tmpl >> fly.toml

flyctl init "${APP_NAME}" --nowrite
```

## Save AWS credentials

Use the `flyctl secrets set` command to securely save your AWS credentials to your app:

```bash
flyctl secrets set \
  --app "${APP_NAME}" \
  "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
  "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}"
```

## Deploy

Finally, it's time to deploy your app. Run it once with `CREATE_NEW_DB='true'` so that LogPaste can bootstrap its database.

```bash
flyctl deploy \
  --app "${APP_NAME}" \
  --env "AWS_REGION=${AWS_REGION}" \
  --env "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  --env "CREATE_NEW_DB='true'" \
  --image mtlynch/logpaste
```

After that command succeeds, deploy it without `CREATE_NEW_DB`. On all future deployments, you'll want to deploy with this command:

```bash
flyctl deploy \
  --app "${APP_NAME}" \
  --env "AWS_REGION=${AWS_REGION}" \
  --env "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  --image mtlynch/logpaste
```
