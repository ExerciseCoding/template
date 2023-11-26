package v2

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type Context struct {
	Req *http.Request

	// Resp 如果用户直接使用这个
	// 那么他就绕开了RespData 和RespStatusCode两个
	// 那么部分 middleware无法运作
	Resp http.ResponseWriter

	// 这个主要是为了middleware 读写用
	RespData       []byte
	RespStatusCode int

	PathParams  map[string]string
	queryValues url.Values
	MatchRoute  string
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJsonAndStatus(http.StatusOK, val)
}

func (c *Context) RespJsonAndStatus(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("web: 未写完数据")
	}
	if err != nil {
		return err
	}
	//c.Resp.WriteHeader(status)
	//c.Resp.Header().Set("Content-Type", "application/json")
	//c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	c.RespStatusCode = status
	c.RespData = data
	return nil
}

func (c *Context) RespJson(val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("web: 未写完数据")
	}
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	return nil
}

func (c *Context) BindJson(val any) error {
	if val == nil {
		return errors.New("web: val不能为Nil")
	}
	if c.Req.Body == nil {
		return errors.New("web: body不能为Nil")
	}
	decoder := json.NewDecoder(c.Req.Body)

	// userNumber => 数字就是用Nbumer 来表示
	// 否则默认是float64
	// decoder.UseNumber()

	// 如果要是有一个未知的字段，就会报错
	// 比如说你的User 只有Name 和 Email两个字段，JSON里面额外多了一个Age字段，那么就会报错
	// decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	if key == "" {
		return "", errors.New("web: key不能为Nil")
	}
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}

	//val, ok := c.Req.Form[key]
	//if !ok {
	//	return "", errors.New("web: key不存在")
	//}
	//return val[0], nil
	return c.Req.FormValue(key), nil
}

func (c *Context) QueryValue(key string) (string, error) {
	// 用户区别不出来是真的有值，但是值恰好是空字符串
	// 还是没有值

	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return "", errors.New("Web: key所对应的值不存在")
	}
	return vals[0], nil
}

func (c *Context) QueryValueV1(key string) StringValue {
	// 用户区别不出来是真的有值，但是值恰好是空字符串
	// 还是没有值

	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("Web: key所对应的值不存在"),
		}
	}
	return StringValue{
		val: vals[0],
		err: nil,
	}
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("Web: key不存在")
	}
	return val, nil
}

func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("Web: key不存在"),
		}
	}
	return StringValue{
		val: val,
		err: nil,
	}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

type SafeContext struct {
	ctx   *Context
	mutex sync.RWMutex
}

func (s *SafeContext) RespJSONOK(val any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.ctx.RespJson(val)
}
