# Deploy LogPaste to Google Cloud Run

It's possible to run LogPaste as a microservice on Google Cloud Run, Google's platform for launching Docker containers on-demand.

Google Cloud Run launches your LogPaste instance in response to HTTP requests and then shuts it down during inactivity. This minimizes hosting costs, as you only pay for the time that your instance is running.

Another benefit of running LogPaste on Google Cloud Run is that you don't need to configure S3 credentials. By default, Cloud Run executes LogPaste in a context that has write access to Google Cloud Storage for persistent storage, so LogPaste can read and write data without accruing bandwidth fees or managing credentials for an external service.

## Requirements

You'll need:

* A Google Cloud Platform account
* [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
* Docker

## Set your environment variables

To begin, create environment variables for your GCP  settings:

```bash
GCP_PROJECT="your-gcp-project"  # Replace with your GCP project ID
GCS_BUCKET="bucketname"         # Replace with your bucket name (you can create it later)
GCP_REGION="us-east1"           # Replace with your desired GCP region
```

### Authenticate

To use the Google Cloud SDK, you need to authenticate gcloud and configure docker to push to Google Container Registry:

```bash
gcloud auth login && \
  gcloud auth configure-docker
```

### Specify GCP Project

Next, configure gcloud to remember your GCP project:

```bash
gcloud config set project "${GCP_PROJECT}"
```

### Enable services

Your GCP project will need the [Cloud Run API](https://cloud.google.com/run/docs/reference/rest) enabled, so ensure this API is activated in your project:

```bash
gcloud services enable run.googleapis.com
```

### Create a Google Cloud Storage bucket

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
# conflicts with the data store.
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

If you review your GCP logs, you'll see that the LogPaste server terminates within a few minutes of each HTTP request. This is normal. LogPaste persists all of its data in Google Cloud Storage before it shuts down, and it will start up again on the next HTTP request it receives with all the same data.

## Run as service account (optional)

The instructions above launch LogPaste a Cloud Run service under the Default Compute Service credentials. You can improve security by running LogPaste under a service account with a stricter set of permissions:

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
cannot fetch generations: googleapi: Error 403: logpaste@yourproject.iam.gserviceaccount.com does not have storage.objects.list access to the Google Cloud Storage bucket., forbidden
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
