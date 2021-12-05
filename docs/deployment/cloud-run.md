# Deploy LogPaste to Google Cloud Run

TODO: Finish this, still a work in progress.

## Requirements

* [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)

## Set project variables

```bash
GCP_PROJECT="your-gcp-project"  # Replace with your GCP project ID
GCS_BUCKET="bucketname"         # Replace with your bucket name (you can create it later)
GCP_REGION="us-east1"           # Replace with your desired GCP region
```

### Authenticate

You need to authenticate gcloud and configure docker to push to Google Container Registry.

```bash
gcloud auth login && \
  gcloud auth configure-docker
```

### Specify GCP Project

Next, configure gcloud to remember your GCP project to save you from typing it in every command:

```bash
gcloud config set project "${GCP_PROJECT}"
```

### Enable services

Your GCP project will need the [Cloud Run API](https://cloud.google.com/run/docs/reference/rest) enabled, so ensure this API is activated in your project:

```bash
gcloud services enable run.googleapis.com
```

### Create a Google Cloud Storage bucket (optional)

If you haven't already created the GCS bucket to provide LogPaste's persistent storage, create it with the command below:

```bash
gsutil mb "gs://${GCS_BUCKET}"
```

## Push to Google Container Registry

Cloud Run can't run Docker images from external image repositories, so you'll need to copy the official LogPaste image to Google Container Registry:

```bash
LOGPASTE_VERSION="mtlynch/logpaste:0.2.5"   # LogPaste version on DockerHub
LOGPASTE_GCR_IMAGE_NAME="logpaste"          # Change to whatever name you prefer
LOGPASTE_GCR_URL="gcr.io/${GCP_PROJECT}/${LOGPASTE_GCR_IMAGE_NAME}"
```

```bash
docker pull "${LOGPASTE_VERSION}" && \
  docker tag "${LOGPASTE_VERSION}" "${LOGPASTE_GCR_URL}" && \
  docker push "${LOGPASTE_GCR_URL}"
```

## Deploy

Finally, it's time to deploy your image to Cloud Run:

```bash
# You can leave this as-is or change to a different name.
GCR_SERVICE_NAME="logpaste"

# Limit to a single instance, as multiple workers will generate sync
# conflicts with the database.
MAX_INSTANCES="1"

gcloud beta run deploy \
  "${GCR_SERVICE_NAME}" \
  --image "${LOGPASTE_GCR_URL}" \
  --set-env-vars "DB_REPLICA_URL=gcs://${GCS_BUCKET}/db" \
  --allow-unauthenticated \
  --region "${GCP_REGION}" \
  --execution-environment gen2 \
  --no-cpu-throttling \
  --max-instances "${MAX_INSTANCES}"
```

If the deploy is successful, you'll see a message like the following:

```text
Service [logpaste] revision [logpaste-00002-cos] has been deployed and is serving 100 percent of traffic.
Service URL: https://logpaste-abc123-ue.a.run.app
```

Your LogPaste instance will serve at the URL listed after to "Service URL."

Cloud Run will shut down your LogPaste instance within a few minutes of each HTTP request. This is normal. LogPaste persists all of its data in Google Cloud Storage before it shuts down, and it will start up again on the next HTTP request it receives with all the same data.

## Run as service account (optional)

The instructions above launch LogPaste a Cloud Run service under the default Compute service credentials. You can improve security by running LogPaste under a service account with a stricter set of permissions:

```bash
SERVICE_ACCOUNT_NAME="logpaste"
SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com"

# Create a new service account.
gcloud iam service-accounts create \
  "${GCR_SERVICE_NAME}" \
  --description="LogPaste microservice, which can write to Cloud Storage"

# Give service account write access to GCS.
gcloud projects add-iam-policy-binding \
  "${GCP_PROJECT}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/storage.admin"

# Update LogPaste Cloud Run service to use your service account.
gcloud run services update "${GCR_SERVICE_NAME}" \
  --service-account "${SERVICE_ACCOUNT_EMAIL}" \
  --region "${GCP_REGION}"
```

You may receive this error after executing `gcloud run services`:

```text
cannot fetch generations: googleapi: Error 403: logpaste@logpaste.iam.gserviceaccount.com does not have storage.objects.list access to the Google Cloud Storage bucket., forbidden
```

On GCP, IAM permissions sometimes take a few minutes to go into effect. Wait 2-3 minutes and try running the `gcloud run services update` command again.

## Set custom domain (optional)

If you have a domain name you'd like to use (e.g., `logpaste.example.com`) instead of Cloud Run's auto-generated URL, you can assign it to your cloud service.

First, you'll need to [verify to GCP that you own the domain](https://cloud.google.com/run/docs/mapping-custom-domains#command-line). Then, create a mapping from the domain to your LogPaste instance:

```bash
CUSTOM_DOMAIN="logpaste.example.com" # Replace with your domain name
```

```bash
gcloud beta run domain-mappings create \
  --service "${GCR_SERVICE_NAME}" \
  --domain "${CUSTOM_DOMAIN}" \
  --region "${GCP_REGION}"
```

The command will show you a DNS record to add to your DNS. When you add the record, wait a few minutes for records to propagate and Google Cloud Platform to generate your certificate. When it's done, you'll be able to access the Cloud Run instance from your custom domain.

## Acknowledgments

Thanks to [Steren Giannini](https://github.com/steren) from the Google Cloud Run team for his help with these instructions.
