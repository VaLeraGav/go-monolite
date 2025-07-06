


- 400 Bad Request – некорректные данные.
- 401 Unauthorized – неавторизован.
- 403 Forbidden – нет прав.
- 404 Not Found – ресурс не найден.
- 409 Conflict – например, уникальное поле занято.
- 500 Internal Server Error – если что-то пошло не так на сервере.


| Операция | Метод  | Статус                  | Тело ответа                   |
| -------- | ------ | ----------------------- | ----------------------------- |
| Create   | POST   | 201 Created             | Созданный ресурс              |
| Update   | PUT    | 200 OK / 204 No Content | Обновлённый ресурс или ничего |
| Delete   | DELETE | 204 No Content / 200 OK | Пусто или сообщение           |


export PATH="$PATH:$(go env GOPATH)/bin"