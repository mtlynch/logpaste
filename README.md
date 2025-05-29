# logpaste

[![CircleCI](https://circleci.com/gh/mtlynch/logpaste.svg?style=svg)](https://circleci.com/gh/mtlynch/logpaste)
[![Docker Pulls](https://img.shields.io/docker/pulls/mtlynch/logpaste.svg?maxAge=604800)](https://hub.docker.com/r/mtlynch/logpaste/)
[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](LICENSE)

A minimalist web service for uploading and sharing log files.

[![LogPaste animated demo](https://raw.githubusercontent.com/mtlynch/logpaste/master/.readme-assets/demo.gif)](https://raw.githubusercontent.com/mtlynch/logpaste/master/.readme-assets/demo.gif)

## Features

- Accept text uploads from command-line, JavaScript, and web UI
- Simple to deploy
  - Runs in a single Docker container
  - Fits in the free tier of Heroku
- Easy database management
  - Syncs persistent data to any S3-compatible cloud storage provider
- Customizable UI without changing source code

## Demo

- <http://logpaste.com>

## Run LogPaste

### From source

```bash
PORT=3001 go run cmd/logpaste/main.go
```

### From Docker

To run LogPaste within a Docker container, mount a volume from your local system to store the LogPaste sqlite database.

```bash
docker run \
  -e "PORT=3001" \
  -p 3001:3001/tcp \
  --volume "${PWD}/data:/app/data" \
  --name logpaste \
  mtlynch/logpaste
```

### From Docker + cloud data replication

If you specify settings for an S3 bucket, LogPaste will use [Litestream](https://litestream.io/) to automatically replicate your data to S3.

You can kill the container and start it later, and it will restore your data from the S3 bucket and continue as if there was no interruption.

```bash
LITESTREAM_ACCESS_KEY_ID=YOUR-ACCESS-ID
LITESTREAM_SECRET_ACCESS_KEY=YOUR-SECRET-ACCESS-KEY
LITESTREAM_REGION=YOUR-REGION
DB_REPLICA_URL=s3://your-bucket-name/db

docker run \
  -e "PORT=3001" \
  -e "LITESTREAM_ACCESS_KEY_ID=${LITESTREAM_ACCESS_KEY_ID}" \
  -e "LITESTREAM_SECRET_ACCESS_KEY=${LITESTREAM_SECRET_ACCESS_KEY}" \
  -e "LITESTREAM_REGION=${LITESTREAM_REGION}" \
  -e "DB_REPLICA_URL=${DB_REPLICA_URL}" \
  -p 3001:3001/tcp \
  --name logpaste \
  mtlynch/logpaste
```

Some notes:

- Only run one Docker container for each S3 location
  - LogPaste can't sync writes across multiple instances.

### With custom site settings

LogPaste offers some options to customize the text for your site. Here's an example that uses a custom title, subtitle, and footer:

```bash
docker run \
  -e "PORT=3001" \
  -p 3001:3001/tcp \
  --name logpaste \
  mtlynch/logpaste \
  -title 'My Cool Log Pasting Service' \
  -subtitle 'Upload all your logs for FooBar here' \
  -footer '<h2>Notice</h2><p>Only cool users can share logs here.</p>' \
  -showdocs=false \
  -perminutelimit 5
```

## Parameters

### Command-line flags

| Flag              | Meaning                                             | Default Value                                          |
| ----------------- | --------------------------------------------------- | ------------------------------------------------------ |
| `-title`          | Title to display on homepage                        | `"LogPaste"`                                           |
| `-subtitle`       | Subtitle to display on homepage                     | `"A minimalist, open-source debug log upload service"` |
| `-footer`         | Footer to display on homepage (may include HTML)    |                                                        |
| `-showdocs`       | Whether to display usage documentation on homepage  | `true`                                                 |
| `-perminutelimit` | Number of pastes to allow per IP per minute         | `0` (no limit)                                         |
| `-maxsize`        | Max file size users can upload                      | `2` (2 MiB)                                            |
| `-prefix`         | Set prefix path for homepage(i.e. `ex.com/prefix/`) | `` (empty string)                                      |

### Docker environment variables

You can adjust behavior of the Docker container by passing these parameters with `docker run -e`:

| Environment Variable           | Meaning                                                                                           |
| ------------------------------ | ------------------------------------------------------------------------------------------------- |
| `PORT`                         | TCP port on which to listen for HTTP connections (defaults to 3001)                               |
| `LP_BEHIND_PROXY`              | Set to `y` if running behind an HTTP proxy to improve logging                                     |
| `DB_REPLICA_URL`               | S3 URL where you want to replicate the LogPaste datastore (e.g., `s3://mybucket.mydomain.com/db`) |
| `LITESTREAM_REGION`            | AWS region where your S3 bucket is located                                                        |
| `LITESTREAM_ACCESS_KEY_ID`     | AWS access key ID for an IAM role with access to the bucket where you want to replicate data.     |
| `LITESTREAM_SECRET_ACCESS_KEY` | AWS secret access key for an IAM role with access to the bucket where you want to replicate data. |
| `LP_PREFIX`                    | Set prefix path for main functionality (i.e. via curl you'll get `https://ex.com/prefix/id`       |

### Docker build args

If you rebuild the Docker image from source, you can adjust the build behavior with `docker build --build-arg`:

| Build Arg            | Meaning                                                                     | Default Value |
| -------------------- | --------------------------------------------------------------------------- | ------------- |
| `litestream_version` | Version of [Litestream](https://litestream.io/) to use for data replication | `v0.3.9`      |

## Deployment

LogPaste is easy to deploy to cloud services. Here are some places it works well:

- [fly.io](docs/deployment/fly.io.md) (recommended)
- [Heroku](docs/deployment/heroku.md)
- [Google Cloud Run](docs/deployment/cloud-run.md)
- [Amazon LightSail](docs/deployment/lightsail.md)

## Further reading

- ["How Litestream Eliminated My Database Server for $0.03/month"](https://mtlynch.io/litestream/): Explains the motivation behind LogPaste
