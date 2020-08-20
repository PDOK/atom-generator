# atom-generator

ATOM Generator is a CLI application to generate ATOM feeds.

An atom feed is a [web standard](https://tools.ietf.org/html/rfc5023) for sharing files via http. The atom in this case is an xml file with metadata and the actual location, where a file can be downloaded. An atom feed can consist of several files. These are grouped in the index.xml which in turn refers to all atom xml, of specific files. This atom generator is used to generate the atom xml. This is only possible for files that are in [Minio](https://min.io/). Minio is a blob store that is used to house the documents. The atom generator is used in two ways. In this readme we refer to the old way and the new way.

__Old situation__: The atom generator copies the files to another bucket. And the atom xml files are also placed in this bucket. The cli signature has not been adapted for the new situation so that it is backwards compatible.

__New situation__: The atom generator stores the atom xml files locally. These are then immediately available in the service container. The parameter path `= - path = / data / here` is used for this. The files themselves are not moved, they are already in the correct location. This removes two parameters from the old CLI signature.

## Install dependencies for development

Install dependencies with:

```bash
# install dependencies
$ sudo apt-get install libxml2-dev libxslt1-dev python-dev
$ PIPENV_VENV_IN_PROJECT=1 pipenv install --python 3.8 --dev
```

Remove with:

```pipenv
pipenv --rm
```

## Usage

The application uses environmental variables for configuration. With use of pipenv these variables are automatically loaded from the `.env` file. To test without the use of environmental variables the `.env` file can be deleted.

Start the pipenv shell:

```bash
pipenv shell
```

```bash
$ generate-atom --help
Usage: generate-atom [OPTIONS] [LOCATIONS]... CONFIG_PATH BASE_URL

  Generate Atom service feed.

  [locations]... should be filled with the following parameters (* is
  optional):

  source_bucket       bucket in minio that contains the raw files
  source_path         path to files inside the minio bucket
  *destination_bucket optional new location for the raw files
  *destination_path   path to files inside the minio bucket

  example: `generate-atom source_bucket /source_path destination_bucket /destination_path conf.json http://example.com`

Options:
  --force / --no-force  Overwrite existing atom feed
  --verbose / --silent  Print INFO log messages to stdout
  --path TEXT           Path where atom xml should be stored locally
  --service_path TEXT   prefix path for atom url
  --help                Show this message and exit.
```

### Example

#### _Old situation_

```bash
generate-atom deliveries hwh/hydrografie/1/ atom hwh/hydrografie conf.json https://geodata.nationaalgeoregister.nl
```

#### _New situation_

```bash
generate-atom atom hwh/hydrografie/1/ conf.json https://geodata.nationaalgeoregister.nl --path=/output --service_path=/hydrografie/v0_1
```

For more information about the specific functions:

A good example can be found at:

[Workflow Publication New Atom INSPIRE Service](manual_create_atom.md)

## Docker

```docker
docker build -t pdok/atom-generator .
```

To run the docker container you first need to build it and run it through `docker run` with the same environment variables:

- `S3_SIGNING_REGION`
- `S3_ENDPOINT_NO_PROTOCOL`
- `S3_ACCESS_KEY`
- `S3_SECRET_KEY`
- `NGR_ENVIRONMENT`

Run through the command-line interface:

```docker
docker run -e S3_ENDPOINT_NO_PROTOCOL=localhost:8000 -e S3_ACCESS_KEY=my_access_key -e S3_SECRET_KEY=my_secret_key -e S3_SIGNING_REGION=Amsterdam -e NGR_ENVIRONMENT=test generate-atom --help
```

## Config - values.json

Atom implementation is based on the INSPIRE Download Service [Technical Guidance](https://inspire.ec.europa.eu/documents/Network_Services/Technical_Guidance_Download_Services_v3.1.pdf). See here below the explanation for the different fields.

### service_subtitle

The field `service_subtitle` maps towards the field Resource Abstract in the ISO19119 service metadata record.

### datasets/downloads/download_content

The field `datasets/downloads/download_content` is only mandatory when a dataset is submitted in multiple files. When this is the case this field describes the structure of the files. In case of a single download this field can be ignored.

### datasets/dataset_bbox

The field `datasets/dataset_bbox` will be used to generated a `georss:polygon` for each dataset. The CRS of the field `datasets/dataset_bbox` needs to be WGS84 (EPSG:4326). The extent needs to match the extent of the Geographic Bounding Box of the corresponding dataset metadata record.
