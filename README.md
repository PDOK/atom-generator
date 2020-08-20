# Atom-Generator

ATOM Generator is a CLI application to generate ATOM feeds.

Een atom feed is een webstandaard voor het delen van bestanden via http.
De atom in dit geval is een xml bestand met meta data en de werkelijke locatie,
waar een bestand gedownload kan worden.

Een atom feed kan bestaan uit meerdere bestanden.
Deze staan gegroepeerd in de index.xml die weer verwijst naar alle atom xml's,
van specifieke bestanden.

Deze atom generator word gebruikt om de atom xml's te genereren. 
Dit kan alleen voor betsanden die in minio staan.
Minio is de blob store die word gebruikt om de documenten te huisvesten.

De atom generator word op twee manieren gebruikt.
In deze readme refereren wij naar de oude manier en die nieuwe manier.

__Oude situatie:__ De atom generator kopieert de bestanden naar een andere buckit.
En de atom xml bestanden worden ook in deze buckit geplaatst.
De cli signature is niet aangpast voor de nieuwe situatie zodat deze backwards compatible is.

__Nieuwe situatie:__ De atom generator slaat de atom xml bestanden lokaal op.
Deze worden dan direct in de service container beschikbaar.
Hiervoor wordt gebruik gemaakt van de parameter path `=--path=/data/here`.
De bestanden zelf worden niet verplaatst, zij staan al op de juiste locatie.
Hierdoor komen twee parameters van de oude cli signature te vervallen.

## Installeer dependencies for development 

Installeer dependencies met:

```
# install dependencies
$ sudo  apt-get install libxml2-dev libxslt1-dev python-dev
$ PIPENV_VENV_IN_PROJECT=1 pipenv install --python 3.8 --dev
```

Verwijder met:

```
pipenv --rm
```

## Gebruik

De applicatie gebruikt environmental variables voor configuratie. Bij gebruik van
pipenv worden deze variabelen automatisch ingeladen uit het `.env` bestand. Om te testen
zonder environmental variabelen kan dat bestand deze weggehaald worden. 

Start de pipenv shell als volgt: 

```bash
pipenv shell
```

```
$ generate-atom --help
Usage: generate-atom [OPTIONS] [LOCATIONS]... CONFIG_PATH BASE_URL

  Generate Atom service feed.

  [locations]... should be filled with the following parameters (* is
  optional):

  source_bucket       bucket in minio that contains the raw files
  source_path         path to files inside the minio bucket
  *destination_bucket optional new location for the raw files
  *destination_path   path to files iniside the minio bucket

  example: `generate-atom source_bucket /source_path destination_bucket
  /destination_path conf.json http://example.com`

Options:
  --force / --no-force  Overwrite existing atom feed
  --verbose / --silent  Print INFO log messages to stdout
  --path TEXT           Path where atom xml should be stored locally
  --service_path TEXT   prefix path for atom url
  --help                Show this message and exit.
```

### Voorbeeld

_Oude situatie_

`generate-atom deliveries hwh/hydrografie/1/ atom hwh/hydrografie conf.json https://geodata.nationaalgeoregister.nl`

_Nieuw situatie_

`generate-atom atom hwh/hydrografie/1/ conf.json https://geodata.nationaalgeoregister.nl --path=/output --service_path=/hydrografie/v0_1`

Voor meer informatie over specifieke functie:

Een goede voorbeeld kan je vinden op:

[Workflow Publication New Atom INSPIRE Service](manual_create_atom.md)

## Dockerfile

Draai docker container door eerst image te bouwen en dan met `docker run` draaien met dezelfde environmental variabelen:

- `S3_SIGNING_REGION`
- `S3_ENDPOINT_NO_PROTOCOL`
- `S3_ACCESS_KEY`
- `S3_SECRET_KEY`
- `NGR_ENVIRONMENT`

Voer uit op de commandline:
 
```
docker run -e S3_ENDPOINT_NO_PROTOCOL=localhost:8000 -e S3_ACCESS_KEY=my_access_key -e S3_SECRET_KEY=my_secret_key -e S3_SIGNING_REGION=Amsterdam -e NGR_ENVIRONMENT=test generate-atom --help
```

## Config - values.json

Atom implementatie is gebaseerd op de INSPIRE Download Service [Technical Guidance](https://inspire.ec.europa.eu/documents/Network_Services/Technical_Guidance_Download_Services_v3.1.pdf). Zie hieronder voor uitleg documentatie voor de verschillende velden. 


### service_subtitle

Het veld `service_subtitle` mapt naar het veld Resource Abstract in het ISO19119 service metadata record.

### datasets/downloads/download_content

Het veld `datasets/downloads/download_content` is alleen verplicht wanneer een dataset aangeboden wordt in meerdere opgesplitste bestanden, als dit het geval is bevat dit veld een beschrijving van de structuur van de bestanden. In het geval van een enkele download kan het veld achterwege worden gelaten.

### datasets/dataset_bbox

Het veld `datasets/dataset_bbox` wordt gebruikt om voor elke dataset een `georss:polygon` veld aan te maken. CRS van het veld `datasets/dataset_bbox` moet WGS84 (EPSG:4326) zijn. Deze extent moet overeenkomen met de extent van de Geographic Bounding Box van het corresponderende dataset metadata record.
