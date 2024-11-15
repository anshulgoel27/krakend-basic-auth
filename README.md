# krakend-basic-auth
HTTP Basic authentication middleware for the [KrakenD-CE](https://github.com/krakend/krakend-ce.git)

## Install and test
```bash
git clone https://github.com/krakend/krakend-ce.git
cd krakend-ce

#Modify handler_factory.go
#Add to imports: basicauth "github.com/anshulgoel27/krakend-basic-auth/gin"
#Add to NewHandlerFactory (before "return handlerFactory"): handlerFactory = basicauth.New(handlerFactory, logger)

go get github.com/anshulgoel27/krakend-basic-auth/gin

make build

./krakend run -c ./krakend.json -d

curl --user foo:bar http://localhost:8080/private/test
```

## Example krakend.json
```json
{
    "version": 3,
    "name": "My lovely gateway",
    "port": 8080,
    "cache_ttl": 3600,
    "timeout": "3s",
    "endpoints": [
        {
            "endpoint": "/private/{user}",
            "method": "GET",
            "headers_to_pass": [
                "Authorization",
                "Content-Type"
            ],
            "backend": [
                {
                    "host": [
                        "https://api.github.com"
                    ],
                    "url_pattern": "/",
                    "whitelist": [
                        "authorizations_url",
                        "code_search_url"
                    ]
                }
            ],
            "extra_config": {
                "github.com/anshulgoel27/krakend-basic-auth": {
                    "pass": "bar",
                    "user": "foo"
                }
            }
        },
        {
            "endpoint": "/public/{user}",
            "method": "GET",
            "headers_to_pass": [
                "Authorization",
                "Content-Type"
            ],
            "backend": [
                {
                    "host": [
                        "https://api.github.com"
                    ],
                    "url_pattern": "/",
                    "whitelist": [
                        "authorizations_url",
                        "code_search_url"
                    ]
                }
            ]
        }
    ]
}
```
