# Shortener Web APP

### This is a shortnener and redirection tool built in Go Lang.

-----------

### CONFIG

> Add the following **.env** file in the root app path with the values:
> ```
> APP_PORT=8000
> REDIS_URL=redis://localhost:6379
> ```

-----------

### RUN

> Run the app executing the following CLI command over the app rootath:
> ```
> > go run main.go
> ```

-----------

### REQUESTS

- **POST**

> Creates a new redirection posting a JSON object with an **url** property:
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
