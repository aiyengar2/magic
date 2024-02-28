ARG BASE_IMAGE

FROM --platform=$BUILDPLATFORM scratch AS build

COPY ./bin ./

# windows

FROM ${BASE_IMAGE} AS windows

ARG OS
ARG ARCH

ENTRYPOINT [ "magic.exe" ]
COPY --from=build magic-$OS-$ARCH.exe ./magic.exe

# linux

FROM ${BASE_IMAGE} AS linux

ARG OS
ARG ARCH

ENTRYPOINT [ "./magic" ]
COPY --from=build magic-$OS-$ARCH ./magic


