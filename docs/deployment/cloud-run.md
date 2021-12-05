# Deploy LogPaste to Google Cloud Run

TODO: Finish this, still a work in progress.

## Requirements

* gcloud SDK

## Set variables

```bash
GCP_PROJECT="your-gcp-project"    # Replace with your GCP project ID
GCS_BUCKET="your-logpaste-bucket" # Replace with your bucket name
GCP_REGION="us-east1"             # Replace with your desired GCP region
```

### Authenticate

You need to authenticate gcloud and configure docker to push to Google Container Registry.

```bash
gcloud auth login && \
  gcloud auth configure-docker
```

### Enable services

```bash
gcloud config set project "${GCP_PROJECT}" && \
  gcloud services enable run.googleapis.com
```

### Create a Google Cloud Storage bucket

If you haven't already created the GCS bucket, create it with the command below:

```bash
gsutil mb "gs://${GCS_BUCKET}"
```

## Push to Google Container Registry

Cloud Run can't run Docker images from external image repositories, so next you'll copy the official LogPaste image to Google Container Registry:

```bash
LOGPASTE_VERSION="mtlynch/logpaste:latest" # LogPaste version on DockerHub
LOGPASTE_GCR_TAG="logpaste"
LOGPASTE_GCR_URL="gcr.io/${GCP_PROJECT}/${LOGPASTE_GCR_TAG}"
```

```bash
docker pull "${LOGPASTE_VERSION}" && \
  docker tag "${LOGPASTE_VERSION}" "${LOGPASTE_GCR_URL}" && \
  docker push "${LOGPASTE_GCR_URL}"
```

## Deploy

Finally, it's time to deploy your image to Cloud Run.

```bash
gcloud beta run deploy logpaste \
  --image "${LOGPASTE_GCR_URL}" \
  --set-env-vars "DB_REPLICA_URL=gcs://${GCS_BUCKET}/db" \
  --allow-unauthenticated \
  --region "${GCP_REGION}" \
  --execution-environment gen2 \
  --no-cpu-throttling
```

If the deploy is successful, you'll see a message like the following:

```text
Service [logpaste] revision [logpaste-00002-cos] has been deployed and is serving 100 percent of traffic.
Service URL: https://logpaste-abc123-ue.a.run.app
```

Your LogPaste instance will serve at the URL listed next to "Service URL." It will shut down a few seconds after each request, but it will start up when it receives its next HTTP request, and it will load all of its data from its persistent GCS bucket.

Thanks to [Steren Giannini](https://github.com/steren) from the Google Cloud Run team for his help with these instructions.
