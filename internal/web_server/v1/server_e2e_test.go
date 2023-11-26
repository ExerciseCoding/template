package v1

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHTTPServer()
	h.addRoute(http.MethodGet, "/user", func(context *Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})

	//handle1 := func(ctx *Context) {
	//	fmt.Println("处理第一件事")
	//}

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})

	h.Get("/order/*", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	h.Post("/form", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	h.Get("/check/:id", func(ctx *Context) {
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不正确"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %d", id)))
	})

	type User struct {
		Name string `json:"name"`
	}
	h.Get("/user/:id", func(ctx *Context) {
		ctx.RespJson(User{
			Name: "Tom",
		})
	})
	h.Get("/user/student/:id", func(ctx *Context) {
		s := SafeContext{
			ctx: ctx,
		}
		s.RespJSONOK(User{
			Name: "Tom",
		})
	})

	h.Start("127.0.0.1:19999")
}

func TestDiff(t *testing.T) {
	// 打开文件
	few := []string{}
	many := []string{}
	file, err := os.Open("few.txt")
	file1, err := os.Open("many.txt")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 创建一个Scanner来读取文件内容
	scanner1 := bufio.NewScanner(file)
	scanner2 := bufio.NewScanner(file1)
	// 逐行读取文件内容
	for scanner1.Scan() {
		line := scanner1.Text()
		few = append(few, line)
	}

	for scanner2.Scan() {
		line := scanner2.Text()
		many = append(many, line)
	}

	// 检查是否有读取错误
	if err := scanner1.Err(); err != nil {
		fmt.Println("读取文件错误:", err)
	}

	if err := scanner2.Err(); err != nil {
		fmt.Println("读取文件错误:", err)
	}

	diffMap := make(map[string]bool)
	for _, num := range many {
		diffMap[num] = true
	}

	for _, num := range few {
		delete(diffMap, num)
	}

	diff := []string{}
	for num := range diffMap {
		diff = append(diff, num)
	}
	fmt.Println("差集:", diff)
}
