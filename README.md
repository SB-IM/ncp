# Node control protocol

#### Depend
- libgstreamer1.0-dev
- libgstreamer-plugins-base1.0-dev
- gstreamer1.0-plugins-good

For debian
```sh
apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good
```

#### Compile ncpgo
go >= 1.12.4

```sh
git submodule update --init --recursive
make dep
make build
```

