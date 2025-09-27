##### Intro
* Host your filesystem on your network, or use it as a static file server for your frontends.
* If running in filesystem mode, it will be accessible at http://network_ip:8000, if running as static file server, it will be accessible at http://network_ip:8001
* uploads will work with PIN in network mode, but will not work with PIN in static file server mode


##### Features
1. Upload Multiple Files
2. Create/Delete Folders
3. Move/Delete Files
4. List files/View it as filesystem
5. Get Signed Url to Download Files
6. Use as static file server for your frontends


##### Usage
1. Archivus uses a UploadDirectory as its root directory
2. Each user has his own master directory inside the UploadDirectory
3. User can create/manage his own filesystem inside the master directory


##### Installation
1. Create a config.yaml file in the root directory, copy the config.yaml.template file from Config folder
2. Run Go build


