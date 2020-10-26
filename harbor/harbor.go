package harbor

import (
	"bytes"
	"github.com/prodanlabs/kubernetes-client-go-deployments-sample/utils"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func GetImages(username, passwd, registry, repositories string) (string, error) {
	url := "https://" + registry + "/v2/" + repositories + "/tags/list"
	// 初始化客户端请求对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Error.Println(err)
	}
	// 添加自定义请求头
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, passwd)
	// 其它请求头配置
	client := &http.Client{
		// 设置客户端属性
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.Error.Println(err)
	}
	defer resp.Body.Close()
	//io.Copy(os.Stdout, resp.Body)
	buff := bytes.Buffer{}
	io.Copy(&buff, resp.Body)
	value := gjson.Get(buff.String(), "tags")
	strArr := value.Array()
	if len(strArr) == 0 {
		utils.Warn.Println("Failed to get harbor registry image")
		return "0", err
	}
	maxVal := strArr[0].Int()
	maxValIndex := 0
	for i := 0; i < len(strArr); i++ {
		if maxVal < strArr[i].Int() {
			maxVal = strArr[i].Int()
			maxValIndex = i
		}
	}
	images := registry + "/" + repositories + ":" + string(maxVal)
	utils.Info.Printf("The new tag of the %s is: %s, The corresponding index is: %v , Image is: %s\n", repositories, maxVal, maxValIndex, images)
	return images, err
}
