FROM akshmakov/linuxgsm:base 
MAINTAINER Mathieu Duperre <mduperre@cisco.com>

# create working directory
RUN mkdir -p quake3

# set the working directory
WORKDIR /quake3

# add binary
COPY bin/main /quake3/main

# set the entrypoint
ENTRYPOINT ["/bin/main"]
