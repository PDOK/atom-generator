# Workflow Publication New Atom INSPIRE Service

> Note: this workflow only applies to *new* *INSPIRE* Atom Services. For non-INSPIRE Atom services datamanagement is used to manage/host Atom services/ 

Steps below describe the workflow on how to create a new Atom INSPIRE Service.

1. Create `values.json` ATOM service configuration file, for an example see the existing kaderrichtlijnwater [configuration](http://*server hostname*:8080/minio/atom-delivery/rws/inspire/kaderrichtlijnwater/02-07-2019T17:42/) or the `aanlevering_test` folder in this repo.

2. Upload `values.json` config file with associated downloads (supported filetypes zip, gml or geopackage) (references in values.json file to downloads need to be resolable) to the atom-delivery bucket. The atom-delivery bucket is structured more or less like the atom bucket (from which the atom feeds are served) so:

- `atom-delivery/{data_provider}/{inspire?}/{service_id}`

This atom delivery dataset folder can contain multiple deliveries, each delivery lives in a folder with the datetimestamp (local time) of the delivery with the following format: `%d-%m-%YT%H:%M`. This is required for the atom-generator application to determine the latest delivery. 

To copy the `values.json` file with the associated download (the *atom delivery*), run the following (assuming the downloads and `values.json` file are located in the folder `my-delivery`):

```
mc cp -r my-delivery/ *server hostname*/atom-delivery/{$data_provider}/${inspire}/${service_id}$(date +%d-%m-%YT%H:%M)
```

3. Run the atom-generator CLI tool:

    generate-atom atom-delivery/rws/inspire/kaderrichtlijnwater atom http://geodata.nationaalgeoregister.nl

> Note: usage of the parameter `--force` will overwrite an existing Atom service feed if it already exists. Without `--force` the existing Atom service feed will not be overwritten. 

> Note: the `atom` and `atom-delivery` bucket for production lives on `*server hostname*`, therefore the atom-generator application needs to be configured for this minio connection. Use the `set-minio-config` command to configure the application. 

The atom-generator application will get the latest delivery from the `atom-delivery/rws/inspire/kaderrichtlijnwater` folder and generate the Atom service feeds in the `atom` bucket, according to the new PDOK URL strategy.


## Example kaderrichtlijnwater

- kaderrichtlijnwater [delivery](http://*server hostname*:8080/minio/atom-delivery/rws/inspire/kaderrichtlijnwater/02-07-2019T17:42/)
- kaderrichtlijnwater [atom feed](http://*server hostname*:8080/minio/atom/rws/inspire/kaderrichtlijnwater/atom/v1/)
- command to generate atom feed:

```
atom-generator atom-delivery/rws/inspire/kaderrichtlijnwater atom http://geodata.nationaalgeoregister.nl
```
