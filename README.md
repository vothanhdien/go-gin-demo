# Gin context must bind with
## Mở đầu
Hôm rồi mình có thấy một đoạn code xử lý khá là lạ.
Kiểu như dùng gin.Context.MustBindWith nhưng sau đó lại xử lý error trả về.
Đại loại như đoạn code sau:
```golang
if err := c.MustBindWith(r,binding.Default(c.Request.Method, c.ContentType())); err != nil {
    _ = c.AbortWithError(http.StatusOK, err)
    return
}
```
Vậy thì response sẽ trả về như thế nào?

## Em yêu khoa học

Đầu tiên mình tạo một service đơn giản và có chứa đoạn code trên. Chi tiết xem main.go


Test success case:
```bash
curl -i --request POST 'localhost:8080/ping' \
--header 'Content-Type: application/json' \
--data-raw '{"name":"dienvt"}'

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 05 Dec 2020 10:11:01 GMT
Content-Length: 26

{"message":"hello dienvt"}
```

Ngon lành! Giờ thử fail case

```bash
curl -i --request POST 'localhost:8080/ping' \
--header 'Content-Type: application/json' \
--data-raw 'im not json' 

HTTP/1.1 400 Bad Request
Date: Sat, 05 Dec 2020 10:18:03 GMT
Content-Length: 0
```

Vậy là response status là 400 (Bad request). Việc mình call c.AbortWithError chả có tách dụng gì.

**Nhưng mà khoan !!!** Nếu như để ý thì log service sau khi thực hiện request trên sẽ như thế này:

```log
[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 400 with 200
[GIN] 2020/12/05 - 17:24:39 | 200 |      98.309µs |       127.0.0.1 | POST     "/ping"
Error #01: invalid character 'i' looking for beginning of value
```

Câu hỏi bây giờ là tại sao response về cho client status 400 nhưng log lại ghi nhận status 200?

Sau một hồi tìm hiểu thì mình nhận ra lý do nằm ở func AbortWithStatus của gin.Context

```golang
// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Writer.WriteHeaderNow()
	c.Abort()
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}
```

Nếu như đào sâu vào thì sẽ thấy hàm `Status` sẽ write xuống respone header và lần thứ 2 sẽ override lại status của lần đầu. Đó là lý do dòng warning hiện lên ở đoạn log phía trên. Vậy thì câu hỏi là tại sao đã override status về 200 nhưng khi response lại vẫn còn status 400.

```golang
func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}
```

Vấn đề nằm ở chỗ  `	c.Writer.WriteHeaderNow()`

```golang
func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}
```

Trước khi write thì responseWriter check thử response đã được set hay chưa. Nhưng trái với func WriteHeader ở trên, responseWriter không cho phép override lại status.

## Kết

Vậy là để trả lời cho câu hỏi đầu tiên, thì câu trả lời sẽ là response status sẽ phụ thuộc vào hàm `AbortWithStatus` đầu tiên được gọi