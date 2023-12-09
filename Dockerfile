FROM alpine:3.8

# This Dockerfile is optimized for go binaries, change it as much as necessary
# for your language of choice.

#RUN apk --no-cache add ca-certificates=20190108-r0 libc6-compat=1.1.19-r10
RUN apk --no-cache add ca-certificates libc6-compat

EXPOSE 9091

COPY ./bin/car-pooling-challenge /
 
ENTRYPOINT [ "/car-pooling-challenge" ]
