package v4

import (
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	lru "github.com/hashicorp/golang-lru"
)

type FileUploader struct {
	FileField string
	// 为什么要用户传
	// 要考虑文件名冲突的问题
	// 所以很多时候，目标文件名字，都是随机的
	DstPathFunc func(header *multipart.FileHeader) string
}

func (u FileUploader) Handle() HandlerFunc {
	if u.FileField == "" {
		u.FileField = "file"
	}
	if u.DstPathFunc == nil {
		// 设置默认值
	}
	return func(ctx *Context) {
		// 第一步：读到文件内容
		// 第二步：计算出目标路径
		// 第三步：保存文件
		// 第四步：返回响应
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}

		defer file.Close()
		// 将目标路径计算的逻辑交给用户
		dst := u.DstPathFunc(fileHeader)
		// dst 上不存在的路径全部建立起来
		if err := os.MkdirAll(dst, 0o777); err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		// O_WRONLY 写入数据
		// O_TRUNC 如果文件本身存在，清空数据
		// O_CREATE 创建一个新的
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dstFile.Close()
		// buf会影响性能
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK

	}
}

type FileUploaderOption func(uploader *FileUploader)

func NewFileUploader(opts ...FileUploaderOption) *FileUploader {
	res := &FileUploader{
		FileField: "file",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", uuid.New().String())
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (u FileUploader) HandleFunc(ctx *Context) {
	//  文件上传逻辑
}

type FileDownloader struct {
	Dir string
}

func (d FileDownloader) Handle() HandlerFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		dst := filepath.Join(d.Dir, req)

		req = filepath.Clean(req)
		// 做一个校验，防止相对路径引起攻击者下载了你的系统文件
		//dst, err = filepath.Abs(dst)
		//if strings.Contains(dst, d.Dir) {
		//
		//}
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Disposition", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)


// 两个层面上
// 1.大文件不缓存
// 2.控制住了缓存的文件的数量
// 所以，做多消耗多少内存? size(cache) * maxSize
type StaticResourceHandler struct {
	dir string
	extensionContentTypeMap map[string]string
	cache *lru.Cache
	maxSize int

}

func NewStaticResourceHandler(dir string, opts...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	// 总共缓存 key-value
	c, err := lru.New(1000)
	if err != nil {
		return nil, err
	}
	res := &StaticResourceHandler{
		dir: dir,
		cache: c,
		// 10 MB, 文件大小超过这个值，就不会缓存
		maxSize: 1024 * 1024 * 10,
		extensionContentTypeMap: map[string]string{
			"jpeg":"image/jpeg",
			"jpe":"image/jpeg",
			"jpg":"image/jpeg",
			"png":"image/png",
			"pdf":"image/pdf",
		},
	}

	for _, opt := range opts  {
		opt(res)
	}
	return res, nil
}



func StaticWithMaxFileSize(maxSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}


func StaticWithCache(c *lru.Cache) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = c
	}
}

func StaticWithExtension(extMp map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMp {
			h.extensionContentTypeMap[ext] = contentType
		}
	}
}
func (s *StaticResourceHandler) Handle(ctx *Context) {
	// 无缓存
	// 1.拿到目标文件
	// 2.定位到目标文件，并且读出来
	// 3.返回给前端

	// 有缓存

	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}

	dst := filepath.Join(s.dir, file)
	ext := filepath.Ext(dst)[1:]
	header := ctx.Resp.Header()
	if data, ok := s.cache.Get(file); ok {
		contentType := s.extensionContentTypeMap[ext]
		// 可能的有文本文件，图片，多媒体(视频,音频)
		header.Set("Content-Type", contentType)
		header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))
		ctx.RespData = data.([]byte)
		ctx.RespStatusCode = 200
		return
	}


	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
		return
	}
	// 大文件不缓存
	if len(data) <=  s.maxSize {
		s.cache.Add(file, data)
	}

	// 可能的有文本文件，图片，多媒体(视频,音频)
	header.Set("Content-Type", s.extensionContentTypeMap[ext])
	header.Set("Content-Length", strconv.Itoa(len(data)))
	ctx.RespData = data
	ctx.RespStatusCode = 200
}

