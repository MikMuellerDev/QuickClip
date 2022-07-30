# QuickClip
 A self hosted, simple alternative to PasteBin.

## Install
### Docker 🐋
The easiest way of installing QuickClip is most likely using `Docker`.
To get started, using following `docker-compose.yml` file is recommended.

```yml
version: "3.7"
services:
  quickclip:
    image: mikmuellerdev/quickclip
    container_name: quickclip
    ports:
      - 8081:8081
    volumes:
      - /tmp/quickclip-config:/app/QuickClip/config
 ```

Make sure to change common variables, like the port or the volumes.
After configuration, run the following command to start the service

```bash
sudo docker-compose up -d
```
### Conventional
Navigate to the `release` section of this repository, then download the according tarball that fits your system architecture.
To check your system's architecture, run 
```bash
uname -a
```
After downloading is complete, *unpack* this tarball using a command like 
```bash
tar -xvf Quickclip-0.0.0.tar.gz
```
Then navigate in the newly unpacked `QuickClip` directory and it's subfolder `bin`. Here you can just run the executable.

## Issues
On Debian `./QuickClip-0.4.0-Beta: /lib/x86_64-linux-gnu/libc.so.6: version GLIBC_2.32' not found (required by ./QuickClip-0.4.0-Beta)` can happen when running the executable on an old system.
## Fix
It is likely that this has already been fixed, if not you know what do to.
```
sudo apt install golang
```
```
cd QuickClip
```
```
make
```
Your custom-built executable can now be found in `./build/QuicklipX.X.X`.
