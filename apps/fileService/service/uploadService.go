package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hichat.zozoo.net/apps/fileService/models"
	"math/rand"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type (
	UploadService struct {
		model *models.FilesModel
	}
	//上传响应结构体
	UploadResponse struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
	}
)

const (
	//最大图片
	ImgMaxSize = 5 * 1024 * 1024
	//最大音频
	AudioMaxSize = 5 * 1024 * 1024
	//最大视频
	VoiceMaxSize = 20 * 1024 * 1024
	//最大文件
	FileMaxSize = 50 * 1024 * 1024
)

var (
	//图片类型
	Types = map[string]string{
		"jpg":  "img",
		"jpeg": "img",
		"png":  "img",
		"gif":  "img",
		"webp": "img",
		"mp3":  "mp3",
		"mp4":  "mp4",
	}
)

func NewUploadService(m *models.FilesModel) *UploadService {
	return &UploadService{
		m,
	}
}

//文件上传
func (u *UploadService) Upload(f *multipart.FileHeader, c *gin.Context) (rsp *UploadResponse, err error) {
	var (
		suffix   string //文件后缀
		fileType string //文件类型
		path     string
	)

	//判断文件大小是否超过限制
	if fileType, suffix, err = u.limitSize(f.Filename, f.Size); err != nil {
		return nil, err
	}

	//定义保存的文件路径
	path = "static/" + fileType + "/" + time.Now().Format("20060102")
	if err = u.createDir(path); err != nil {
		return nil, err
	}

	//获取新的path
	if path, err = u.getFileName(path, suffix); err != nil {
		return nil, err
	}

	//将文件移入指定文件夹内
	if err = c.SaveUploadedFile(f, path); err != nil {
		return nil, err
	}

	path = "/" + path

	static := &models.Files{
		Uuid: c.GetString("uuid"),
		Name: f.Filename,
		Type: fileType,
		Path: path,
		Size: uint16(f.Size),
	}

	//保存数据库
	if err = u.model.Insert(static); err != nil {
		return nil, err
	}

	//组织返回数据
	rsp = new(UploadResponse)
	rsp.Type = fileType
	rsp.Path = path
	rsp.Name = f.Filename
	return rsp, err
}

//限制文件大小
func (u *UploadService) limitSize(name string, size int64) (fileType string, suffix string, err error) {
	var (
		info  []string
		exist bool
	)

	//将文件名拆分成为数组
	info = strings.Split(name, ".")
	suffix = info[len(info)-1]

	//判断文件类型
	if fileType, exist = Types[suffix]; exist == false {
		fileType = "file"
	}

	//文件类型转小写
	fileType = strings.ToLower(fileType)

	//判断文件是否超出限制
	switch fileType {
	case "img": //限制图片大小
		if size >= ImgMaxSize {
			return "", "", errors.New("图片大小不能超过" + fmt.Sprint(ImgMaxSize) + "限制")
		}
		break
	case "mp3": //限制音频大小
		if size >= AudioMaxSize {
			return "", "", errors.New("音频大小不能超过" + fmt.Sprint(AudioMaxSize) + "限制")
		}
		break
	case "mp4": //限制视频大小
		if size >= VoiceMaxSize {
			return "", "", errors.New("视频大小不能超过" + fmt.Sprint(VoiceMaxSize) + "限制")
		}
		break
	default: //限制文件大小
		if size >= FileMaxSize {
			return "", "", errors.New("文件大小不能超过" + fmt.Sprint(FileMaxSize) + "限制")
		}
	}

	return fileType, suffix, nil
}

//生成文件名
func (u *UploadService) randName(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	//判断文件名是否存在，存在则重新生成
	return string(result)
}

//创建文件夹
func (u *UploadService) createDir(path string) (err error) {
	_, err = os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		//不存在，创建
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			//创建失败
			return err
		}

		return nil
	}
	return err
}

//生成新的文件名
func (u *UploadService) getFileName(path string, suffix string) (name string, err error) {
	for {
		name = u.randName(10)
		name = "/" + name + "." + suffix

		//组织文件路径
		path = path + name

		//判断文件名是否存在，存在则重新获取
		file, err := u.model.FindByName(path)
		if err != nil {
			return "", err
		}

		if file == nil {
			return path, nil
		}
	}
}
