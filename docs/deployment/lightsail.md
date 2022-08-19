## Deploy LogPaste to Amazon LightSail

**Warning**: These instructions assume an already initialized datastore on your S3 bucket. Once LogPaste 0.1.2 is released, this assumption won't be necessary.

Amazon LightSail is an attractive option for launching LogPaste to production for two main reasons:

- You can deploy entirely from Amazon's AWS dashboard, so you don't need to install any software.
- You can view LogPaste's server logs directly from your AWS dashboard.
- LightSail offers free TLS certificates for a custom domain.

The downside is that it doesn't fit within AWS's free tier, so hosting will cost $7 per month, as of this writing.

## Pre-requisites

You'll need:

- An Amazon AWS account
- A storage bucket and [IAM credentials](https://aws.amazon.com/iam/) on Amazon S3 or an S3-compatible storage service.

## Create a LogPaste container

Visit [Amazon LightSail's container management dashboard](https://lightsail.aws.amazon.com/ls/webapp/home/containers)

<kbd>

![LightSail management screen](lightsail-images/create-container.png)

</kbd>

Choose a Nano server with 1x scale.

<kbd>

![LightSail server capacity screen](lightsail-images/nano-1x.png)

</kbd>

Click "Set up deployment" and then choose "Specify a custom deployment."

<kbd>

![Screenshot showing where to click specify custom deployment](lightsail-images/set-up-deployment.png)

</kbd>

Enter the following information for your custom deployment:

- Container name: Choose whatever name you want, such as `yourcompany-logpaste`
- Image: `mtlynch/logpaste:0.2.5` (or whatever the [latest LogPaste version](https://github.com/mtlynch/logpaste/releases) is)
- Launch command: _leave blank_
- Environment variables

| Key                            | Value                                               |
| ------------------------------ | --------------------------------------------------- |
| `PORT`                         | 3001                                                |
| `PS_BEHIND_PROXY`              | "y"                                                 |
| `DB_REPLICA_URL`               | The S3 URL of your S3 bucket                        |
| `LITESTREAM_REGION`            | The region of your S3 bucket                        |
| `LITESTREAM_ACCESS_KEY_ID`     | The AWS access key ID from your IAM credentials     |
| `LITESTREAM_SECRET_ACCESS_KEY` | The AWS secret access key from your IAM credentials |

- Open ports: Choose port `3001`, protocol `HTTP`

<kbd>

![Screenshot showing custom values for LightSail container](lightsail-images/container-config.png)

</kbd>

Under "Public Endpoint," select the container you created above:

<kbd>

![Screenshot showing "contoso-logpaste" in dropdown menu](lightsail-images/public-endpoint.png)

</kbd>

Choose a name for your service (it can be the same as your container name):

<kbd>

![Screenshot showing "contoso-logpaste" as name for service](lightsail-images/identify-service.png)

</kbd>

Finally, click "Create container service."

<kbd>

![Screenshot showing "Create container service" button](lightsail-images/create-service.png)

</kbd>

## Completing deployment

As LightSail deploys your image, you'll see a status of "Pending" for the container.

<kbd>

![Screenshot showing LightSail in the process of deploying a container at the "Pending" stage](lightsail-images/container-pending.png)

</kbd>

It will take LightSail about three to five minutes to deploy your instance for the first time. When deployment is complete, it will show a status of "Running."

<kbd>

![Screenshot showing LightSail deployment complete](lightsail-images/container-running.png)

</kbd>

When your container is running, you can access it through the "Public domain" URL.

<kbd>

![Screenshot showing how to access container's public domain URL](lightsail-images/public-domain-url.png)

</kbd>

You can view server logs for your LogPaste instance by clicking "Open log" in the box next to your container.

<kbd>

![Screenshot showing LightSail logs in web dashboard](lightsail-images/view-logs.png)

</kbd>
