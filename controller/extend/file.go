package extend

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"zlsapp/internal/error_code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
)

type File struct {
	service.App
}

func (f *File) Init(g *znet.Engine) {

}

func (f *File) PostUploadImage(c *znet.Context) (interface{}, error) {
	files, err := c.FormFiles("file")
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "上传失败")
	}

	imageDir := "resource/upload/image/"
	dir := c.DefaultFormOrQuery("dir", "")
	dir = zfile.RealPathMkdir(imageDir+dir, true)
	if !zfile.IsSubPath(dir, imageDir) {
		return nil, error_code.InvalidInput.Text("非法存储目录")
	}

	uploads := make(map[string]*multipart.FileHeader, len(files))
	buf := bytes.NewBuffer(nil)

	for _, v := range files {
		f, err := v.Open()
		if err != nil {
			return nil, zerror.InvalidInput.Wrap(err, "文件读取失败")
		}

		if _, err := io.Copy(buf, f); err != nil {
			if err != nil {
				return nil, zerror.InvalidInput.Wrap(err, "文件读取失败")
			}
		}

		f.Close()
		b := buf.Bytes()

		mimeType := strings.Split(zfile.GetMimeType(v.Filename, b), "/")
		if len(mimeType) < 2 {
			return nil, zerror.InvalidInput.Wrap(err, "文件类型错误")
		}

		if mimeType[0] != "image" {
			return nil, error_code.InvalidInput.Text("只支持图片文件")
		}

		ext := "." + mimeType[1]
		id := zstring.Md5Byte(b) + ext
		uploads[id] = v

		buf.Reset()
	}

	for n, f := range uploads {
		err = c.SaveUploadedFile(f, dir+n)
		if err != nil {
			return nil, error_code.ServerError.Text("文件保存失败", err)
		}
	}

	return zarray.Map(zarray.Keys(uploads), func(_ int, p string) string {
		return "/" + zfile.SafePath(dir+p, "resource")
	}), nil
}
