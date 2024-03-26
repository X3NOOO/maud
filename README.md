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

## Usage

### API

All requests should be JSON-formatted and contain the Content-Type header.

Endpoints:

#### `POST /register`
- Request fields:
  `"nick"`
  `"password"`
- Response fields:
  `"nick"`
  `"authorization_token"`
- Cookies set:
- `Authorization`

#### `POST /login`
- Request fields:
  `"nick"`
  `"password"`
- Response fields:
  `"authorization_token"`
- Cookies set:
- `Authorization`

#### `GET /status`
- Required headers/cookies: `Authorization`<sub>Your authorization token</sub>
- Response fields:
  `"nick"` - username
  `"date"` - last activity date

#### `POST /alive`
- Required headers/cookies: `Authorization`<sub>Your authorization token</sub>
- Response fields:
  `"date"` - last activity date

#### `POST /switches`
- Required headers/cookies: `Authorization`<sub>Your authorization token</sub>
- Request fields:
  `"run_after"` - after how many days of inactivity fire the switch
  `"recipients"` - array of emails that will recieve the email
  `"content"` - content of the email
  `"subject"` - subject of the email
- Response fields:
  `"id"` - id of the switch
  `"subject"`
  `"content"`
  `"run_after"`
  `"recipients"`

#### `GET /switches/{id...}`
- Required headers/cookies: `Authorization`<sub>Your authorization token</sub>
- URL variables:
  `{id...}`<sub>(optional)</sub> - id of the switch to get, if not provided all of the existing switches will be returned
- Response fields:
  - Array of:
    `"id"`
    `"subject"`
    `"content"`
    `"run_after"`
    `"recipients"`

#### `DELETE /switches/{id...}`
- Required headers/cookies: `Authorization`<sub>Your authorization token</sub>
- URL variables:
  `{id...}`<sub>(required)</sub> - id of the switch to delete
- Response fields:
  `"id"`
  `"subject"`
  `"content"`
  `"run_after"`
  `"recipients"`

#### `PATCH /switches/{id...}`
- Required headers/cookies: `Authorization`<sub>Your authorization token</sub>
- URL variables:
  `{id...}`<sub>(required)</sub> - id of the switch to modify
- Request fields:
  `"subject"`<sub>(optional, only modified if provided)</sub>
  `"content"`<sub>(optional, only modified if provided)</sub>
  `"run_after"`<sub>(optional, only modified if provided)</sub>
  `"recipients"`<sub>(optional, only modified if provided)</sub>
- Response fields:
  `"id"`
  `"subject"`
  `"content"`
  `"run_after"`
  `"recipients"`

## Installation

### Docker (recommended)

1. `git clone https://github.com/X3NOOO/maud && cd ./maud`
2. `bash ./build.sh config`
3. `docker compose up --detach`

### Raw binary

1. `git clone https://github.com/X3NOOO/maud && cd ./maud`
2. `bash ./build.sh config release run`
