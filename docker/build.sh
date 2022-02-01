#!/bin/bash

# mv ../build/*.tar.gz ./
tar -xvf ./app.tar.gz
sudo docker build . -t mikmuellerdev/quickclip