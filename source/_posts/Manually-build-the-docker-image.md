---
title: Manually build the docker image
date: 2024-08-17 16:27:49
tags:
- Docker
- Nginx
- Rtmp
categories:
- Doc
keywords:
- Docker
- Rtmp
copyright: Guader
copyright_author_href:
copyright_info:
---

I had built [nginx][nginx] with [nginx-rtmp-module][nginx-rtmp-module] about serval years ago.
Chrome decided to deprecated flash in that year, than the usual video stream and player can't work anymore. (I don't sure rtsp/rtmp or whatever, its a long time.)
And recently I want build a docker image for the live server.  
This is the foreword.


## Build nginx with nginx-rtmp-module in server

This isn't hard, follow the documentation, step and step, and succeed.


## Build image by tiangolo

I search to rtmp in dockerhub, this one's repo had Dockerfile.  
So I followed his file at first.  
And the image size is **851.17MB**.   
But his image is about just more than 300MB, I think that's ridiculous.  
And I found his [Dockerfile][tiangolo-Dockerfile] last updated in *2022/9/25* .   
I have to say, good job.


## Manually build the docker image

Then I check the [nginx-rtmp-module][nginx-rtmp-module] and [nginx][nginx] for some help.
Then I found [arut/wiki/Dockerfile][arut-Dockerfile] and [docker-nginx/modules][docker-nginx-modules] will be useful.  
However, they can't be built in my computer, cause of the GFW.  

Then is how to manually build the image.

1. Choose the base image

    Obviously, I choose alpine.

    `docker pull alpine:3.20.2`

2. Create the container and into shell

    `docker run -it --name tempcontainer alpine:3.20.2 /bin/sh` 

3. Execute the command

    ```bash
    apk add --update build-base git bash gcc make g++ zlib-dev linux-headers pcre-dev openssl-dev
    git clone https://github.com/arut/nginx-rtmp-module.git
    git clone https://github.com/nginx/nginx.git
    cd nginx
    ./auto/configure \
        --sbin-path=/usr/local/sbin/nginx \
        --conf-path=/etc/nginx/nginx.conf \
        --error-log-path=/var/log/nginx/error.log \
        --pid-path=/var/run/nginx/nginx.pid \
        --lock-path=/var/lock/nginx/nginx.lock \
        --http-log-path=/var/log/nginx/access.log \
        --http-client-body-temp-path=/tmp/nginx-client-body \
        --with-http_ssl_module \
        --with-threads \
        --with-ipv6 \ 
        --add-module=../nginx-rtmp-module 
    make && make install
    ```

4. Copy your conf to container

    From another shell, not in container

    `docker cp nginx.conf tempcontainer:/etc/nginx/nginx.conf`

5. Commit the container as image

    `docker commit tempcontainer my-last-image`


This is kind of more complex then `docker build ...` .  
But its completely controllable.  
Like I may can't clone the repo, then I can clone more times, until I get that repo.  
Also the apk pkg, I can wait like 5 minutes, then I will get that at last.

And now my image is **149.28MB** , even within a `ffmpeg` .


If interested, just `docker pull guaderxx/nginx-rtmp` to get that.




[nginx]: https://nginx.org/
[nginx-rtmp-module]: https://github.com/arut/nginx-rtmp-module
[tiangolo-Dockerfile]: https://github.com/tiangolo/nginx-rtmp-docker/blob/master/Dockerfile
[arut-Dockerfile]: https://github.com/arut/nginx-rtmp-module/wiki/Building-a-docker-image-with-nginx--rtmp-module#alpine
[docker-nginx-modules]: https://github.com/nginxinc/docker-nginx/blob/master/modules/README.md
