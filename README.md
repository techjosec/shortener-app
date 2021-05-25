# Shortener Web APP

### This APP is a shortnener and redirection tool. 

-----------

### CONFIG

- Add the following **.env** file in the root app path with the values:

```
APP_PORT=8000
REDIS_URL=redis://localhost:6379
```

-----------

### REQUESTS

- **POST**

> Creates a new redirection sending a JSON object with an **url** property who's the url to be redirected:
> ```shell
> curl --request POST \
>   --url http://localhost:8000/api/redirect \
>   --header 'Content-Type: application/json' \
>   --data '{
> 	"url": "https://google.com.ar"
> }'
> ```
> 
> Returns a redirect json, with the redirection **code** property:
> ```json
> {
>   "code": "5MVx6x3Mg",
>   "url": "https://google.com.ar",
>   "created_at": 1621894689
> }
> ```

- **GET**

> GET's redirected to an URL from a given string **code**:
> ```bash
> curl --request GET \
>   --url http://localhost:8000/api/redirect/{code}
> ```
