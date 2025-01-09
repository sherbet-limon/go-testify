package main

import (
    "net/http"
    "strconv"
    "strings"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
    "moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
    countStr := req.URL.Query().Get("count")
    if countStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("count missing"))
        return
    }

    count, err := strconv.Atoi(countStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong count value"))
        return
    }

    city := req.URL.Query().Get("city")

    cafe, ok := cafeList[city]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong city value"))
        return
    }

    if count > len(cafe) {
        count = len(cafe)
    }

    answer := strings.Join(cafe[:count], ",")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(answer))
}

func TestMainHandlerRequest(t *testing.T) {
    req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)
    assert.Equal(t, http.StatusOK, responseRecorder.Code, "Полученный статус не соответствует 200")

    body := responseRecorder.Body.String()
    assert.NotEmpty(t,body, "Тело ответа не должно быть пустым")
    
}
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    totalCount := 4
    req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)
    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)
    assert.Equal(t, http.StatusOK, responseRecorder.Code, "Полученный статус не соответствует 200")

    body := responseRecorder.Body.String()
    assert.NotEmpty(t,body, "Тело ответа не должно быть пустым")
    list := strings.Split(body, ",")
    assert.Len(t, list, totalCount)
}
func main() {
}