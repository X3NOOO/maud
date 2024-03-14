# maud

```
                      __
  __ _  ___ ___ _____/ /
 /  ' \/ _ `/ // / _  /
/_/_/_/\_,_/\_,_/\_,_/
```

Dead man's switch.

## Features

- Email support by default
- Simple architecture
- Easy development of new delivery methods

## Installation

### Docker (recommended)

1. `git clone https://github.com/X3NOOO/maud && cd ./maud`
2. `bash ./build.sh config`
3. `docker build -t maud . && docker run -d --restart unless-stopped --name maud maud`

### Raw binary

1. `git clone https://github.com/X3NOOO/maud && cd ./maud`
2. `bash ./build.sh config release run`
