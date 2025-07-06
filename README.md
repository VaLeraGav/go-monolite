
# Example Go monolith with embedded microservices

This is a learning project for implementing monolith with embedded microservices

- 400 Bad Request – incorrect data.
- 401 Unauthorized – unauthorized.
- 403 Forbidden – no rights.
- 404 Not Found – the resource was not found.
- 409 Conflict – for example, a unique field is occupied.
- 500 Internal Server Error – if something went wrong on the server.


| Operation | Method | Response Status         | Body                        |
| --------- | ------ | ----------------------- | --------------------------- |
| Create    | POST   | 201 Created             | Created resource            |
| Update    | PUT    | 200 OK / 204 No Content | Updated resource or nothing |
| Delete    | DELETE | 204 No Content / 200 OK | Empty or message            |


<!-- export PATH="$PATH:$(go env GOPATH)/bin" -->