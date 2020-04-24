# backend

## RPC

See:

* [users/rpc/userspb/service.proto](users/rpc/userspb/service.proto)

## Starting service

```
$ make docker-run
```

## Interacting with the service

You can use [evans](https://github.com/ktr0731/evans):

```
$ evans repl -r --host localhost --port 65010

  ______
 |  ____|
 | |__    __   __   __ _   _ __    ___
 |  __|   \ \ / /  / _. | | '_ \  / __|
 | |____   \ V /  | (_| | | | | | \__ \
 |______|   \_/    \__,_| |_| |_| |___/

 more expressive universal gRPC client


userspb.UsersRPC@localhost:65010> call AddUser
user::name (TYPE_STRING) => David BP
user::address::address (TYPE_STRING) => somewhere
user::address::city (TYPE_STRING) => Stockholm
user::address::postcode (TYPE_STRING) => 12345
user::address::country (TYPE_STRING) => Sweden
 PHONE
<repeated> user::contact_details::identifier (TYPE_STRING) => 070123456
 EMAIL
<repeated> user::contact_details::identifier (TYPE_STRING) => asd@somewhere.com

<repeated> user::contact_details::skills (TYPE_STRING) => plumbing
<repeated> user::contact_details::skills (TYPE_STRING) => fixing_printers
<repeated> user::contact_details::skills (TYPE_STRING) =>
{
  "userId": "8cabc918-6e01-4ed3-bec4-4bb0a8f19fa8"

Use the arrow keys to navigate: ‚Üì ‚Üë ‚Üí ‚Üê
? platform:
  ‚ñ∏ PHONE
    EMAIL
    WHATSAPP
    FACEBOOK
‚Üì   TELEGRAM
<repeated> user::contact_details::skills (TYPE_STRING) => eating
<repeated> user::contact_details::skills (TYPE_STRING) => sleeping
<repeated> user::contact_details::skills (TYPE_STRING) =>
{
  "userId": "c9f3af70-362d-4b8f-8c71-764f14b32a8e"
}

userspb.UsersRPC@localhost:65010> call GetUsers
{
  "userId": "c9f3af70-362d-4b8f-8c71-764f14b32a8e",
  "user": {
    "name": "David BP",
    "address": {
      "address": "somewhere",
      "city": "Stockholm",
      "postcode": "12345",
      "country": "Sweden"
    },
    "contactDetails": [
      {
        "identifier": "070123456"
      },
      {
        "platform": "EMAIL",
        "identifier": "asd@somewhere.com"
      }
    ],
    "skills": [
      "plumbing",
      "fixing_printers"
    ]
  }
}
{
  "userId": "c9f3af70-362d-4b8f-8c71-764f14b32a8e",
  "user": {
    "name": "Pi the Dog",
    "address": {
      "address": "somewhere",
      "city": "Stockholm",
      "postcode": "12345",
      "country": "Sweden"
    },
    "contactDetails": [
      {
        "platform": "FACEBOOK",
        "identifier": "@pi"
      }
    ],
    "skills": [
      "eating",
      "sleeping"
    ]
  }
}
```

You can also use it this way:

```
$ evans cli -r --host localhost --port 65010 --call userspb.UsersRPC.AddUser -f data/requests/add_user_david.json
{
  "userId": "8ce07df5-2d06-4b44-b2b0-a1e7df7c77e4"
}

$ evans cli -r --host localhost --port 65010 --call userspb.UsersRPC.AddUser -f data/requests/add_user_pi.json
{
  "userId": "07f3999d-6037-4397-9a97-934af11f208a"
}

$ evans cli -r --host localhost --port 65010 --call userspb.UsersRPC.GetUsers -f data/requests/get_users.json
{
  "userId": "8ce07df5-2d06-4b44-b2b0-a1e7df7c77e4",
  "user": {
    "name": "David BP",
    "address": {
      "address": "somewhere",
      "city": "Stockholm",
      "postcode": "12345",
      "country": "Sweden"
    },
    "contactDetails": [
      {
        "identifier": "070123456"
      },
      {
        "platform": "EMAIL",
        "identifier": "asd@somewhere.com"
      }
    ],
    "skills": [
      "plumbing",
      "fixing_printers"
    ]
  }
}
{
  "userId": "07f3999d-6037-4397-9a97-934af11f208a",
  "user": {
    "name": "Pi the Dog",
    "address": {
      "address": "somewhere",
      "city": "Stockholm",
      "postcode": "12345",
      "country": "Sweden"
    },
    "contactDetails": [
      {
        "platform": "FACEBOOK",
        "identifier": "@pi"
      }
    ],
    "skills": [
      "eating",
      "sleeping"
    ]
  }
}
```
