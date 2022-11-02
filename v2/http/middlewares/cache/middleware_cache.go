package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oaago/cloud/redis"
	"github.com/oaago/server/oaa"
)

type CacheMiddleWare oaa.GlobalMiddleware

func NewCacheMiddleWare() CacheMiddleWare {
	return CacheMiddleWare{}
}

type CacheParam struct {
	MapParam map[string]interface{}
}

type CacheCondition struct {
	TimeUnit time.Duration
	Duration uint
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var needCache bool

func (CacheMiddleWare) Cache(c *oaa.Ctx) {
	//约定url请求参数包含 interfaceCache = "need"的为需要缓存的接口
	needCache = true
	context := c.Context.Copy()
	paramMap := map[string]interface{}{}
	var cache CacheParam
	cache.setMapParam(paramMap)
	query := context.Request.URL.Query()
	getResult := query.Get("interfaceCache")
	if getResult != "need" {
		needCache = false
		c.Context.Next()
	}
	if !needCache {
		return
	}
	cacheCondition := buildCacheCondition(query)
	body := bindParam(cache, query, context)
	bindRequestBody(c, context, body, cache)
	md5String := getMd5Key(cache, paramMap, context)
	redisKey := "interfaceCache:" + context.Request.RequestURI + ":" + md5String
	result, err := GetValue(redisKey)
	if err == nil && result != nil {
		c.Return(200, result)
	}
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Context.Writer}
	c.Context.Writer = blw
	c.Context.Next()
	statusCode := c.Context.Writer.Status()
	if statusCode >= 200 && statusCode <= 300 {
		resultStr := blw.body.String()
		r := oaa.Result{}
		err := json.Unmarshal([]byte(resultStr), &r)
		if err != nil {
			return
		}
		marshal, _ := json.Marshal(r.Data)
		SetValue(redisKey, marshal, cacheCondition.TimeUnit*time.Duration(cacheCondition.Duration)) //nolint:errcheck
	}

}

func GetValue(key string) (interface{}, error) {
	var value map[string]interface{}
	cmdResult := redis.Client.Get(key)

	if cmdResult.Err() != nil && cmdResult.Err().Error() != "redis: nil" {
		return GetValue(key)
	} else if cmdResult.Err() != nil {
		return nil, cmdResult.Err()
	}
	val := cmdResult.Val()
	err := json.Unmarshal([]byte(val), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func SetValue(key string, value interface{}, expire time.Duration) error {
	statusCmd := redis.Client.Set(key, value, expire)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func buildCacheCondition(query url.Values) CacheCondition {
	duration := query.Get("duration")
	atoi, err := strconv.Atoi(duration)
	var cacheCondition CacheCondition
	if err != nil {
		cacheCondition.Duration = 5
		cacheCondition.TimeUnit = time.Second
	} else {
		cacheCondition.Duration = uint(atoi)
		timeUnit := query.Get("timeUnit")
		if "second" == timeUnit {
			cacheCondition.TimeUnit = time.Second
		} else if "millSecond" == timeUnit {
			cacheCondition.TimeUnit = time.Millisecond
		} else if "minute" == timeUnit {
			cacheCondition.TimeUnit = time.Minute
		}
	}
	return cacheCondition
}

func bindParam(cache CacheParam, query url.Values, context *gin.Context) io.ReadCloser {
	cache.bindParam(query)
	form := context.Request.Form
	cache.bindParam(form)
	postForm := context.Request.PostForm
	cache.bindParam(postForm)
	body := context.Request.Body
	return body
}

func getMd5Key(cache CacheParam, paramMap map[string]interface{}, context *gin.Context) string {
	mapParam := cache.MapParam
	filedArray := make([]string, 1+len(mapParam))
	if mapParam != nil {
		for key := range mapParam {
			filedArray = append(filedArray, key)
		}
	}
	bubbleSort(&filedArray)
	//遍历所有字段 获取字段值后进行json序列化 并且根据字段名首字母进行排序
	paramMaps := make([]map[string]interface{}, 1+len(filedArray))
	for i := 0; i < len(filedArray); i++ {
		for key, value := range mapParam {
			if filedArray[i] == key {
				param := make(map[string]interface{}, 1)
				//获取属性的字段
				marshal, _ := json.Marshal(value)
				param[key] = string(marshal)
				paramMaps[i] = param
			}
		}
	}
	marshal, _ := json.Marshal(paramMap)
	data := []byte(context.Request.RequestURI + string(marshal))
	m5 := md5.New()
	m5.Write(data)
	md5String := hex.EncodeToString(m5.Sum(nil))
	return md5String
}

//绑定请求体
func bindRequestBody(c *oaa.Ctx, context *gin.Context, body io.ReadCloser, cache CacheParam) {
	if context.Request != nil && body != nil {
		// 把requestBody的内容读取出来
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(context.Request.Body)
			cache.setRequestBody(string(bodyBytes))
		}
	}
}

func (cache *CacheParam) setMapParam(mapParam map[string]interface{}) {
	cache.MapParam = mapParam
}

func (cache *CacheParam) setRequestBody(body string) {
	cache.MapParam["requestBody"] = body
}

//绑定url参数或form
func (cache *CacheParam) bindParam(query url.Values) {
	if query != nil {
		for s := range query {
			getResult := query.Get(s)
			if getResult != "" {
				cache.MapParam[s] = getResult
			}
		}
	}
}

//结构体字段排序
func bubbleSort(arr *[]string) {
	for i := 0; i < len(*arr)-1; i++ {
		for j := 0; j < len(*arr)-1-i; j++ {
			temp := ""
			if (*arr)[j] > (*arr)[j+1] {
				temp = (*arr)[j]
				(*arr)[j] = (*arr)[j+1]
				(*arr)[j+1] = temp
			}
		}
	}
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
