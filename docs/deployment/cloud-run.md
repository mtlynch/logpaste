# Deploy LogPaste to Google Cloud Run

## Enable Cloud Run API

```bash
gcloud services enable run.googleapis.com
```

## Set environment variables

```bash
GCS_BUCKET="your-logpaste-bucket" # Replace with your bucket name
GCP_REGION="us-east1" # Replace with your desired region
LOGPASTE_IMAGE="gcr.io/logpaste/logpaste:latest"
#LOGPASTE_IMAGE="registry.hub.docker.com/mtlynch/logpaste:latest"
```

## Deploy

```bash
gcloud beta run deploy logpaste \
  --image "${LOGPASTE_IMAGE}" \
  --set-env-vars "DB_REPLICA_URL=gcs://${GCS_BUCKET}/db" \
  --allow-unauthenticated \
  --region "${GCP_REGION}" \
  --execution-environment gen2 \
  --no-cpu-throttling
```
