### Go TCP Server

A simple TCP server written from scratch in Golang. 

For now, it handles the following requests:
```

GET /
HTTP/1.1 200 OK

-----

GET /echo/abc
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 3

abc

-----

GET /user-agent

HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 10

curl/8.4.0

-----

GET /files/<file_name>

HTTP/1.1 200 OK (HTTP/1.1 404 Not Found if file doesn't exist)
Content-Type: application/octet-stream
Content-Length: 244

<!DOCTYPE html>
<html lang="en">
</html>

-----

POST /files/<file_name>
Content-Length: 27

field1=value1&field2=value2

Creates a file with the data equal to request body and name equal to <file_name>

```

### Testing
To test the server locally, you can clone the repository and run `go mod tidy` to download and install the required dependencies.

You can run the server using `go run main.go` and use curl to test the server. Here's an example GET request using cURL:
```
curl -i -X GET http://localhost:4221/echo/abc
```




