package utility

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BetaGoRobot/BetaGo/betagovar"
	"github.com/golang/freetype/truetype"
	"github.com/lonelyevil/kook"
)

// MicrosoftYaHei  字体类型,未来荧黑
var MicrosoftYaHei *truetype.Font

func init() {
	InitGlowSansSCFontType()
}

// InitGlowSansSCFontType 初始化字体类型
func InitGlowSansSCFontType() {
	fontFile := filepath.Join(betagovar.FontPath, "Microsoft Yahei.ttf")
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return
	}
	MicrosoftYaHei, err = truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
}

// GetOutBoundIP 获取机器人部署的当前ip
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "101.132.154.52:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() (localTime string) {
	timeLocal, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = timeLocal
	localTime = time.Now().Local().Format("2006-01-02 15:04:05")
	return
}

// ForDebug 用于测试
func ForDebug(test ...interface{}) {
	return
}

// IsInSlice 判断机器人是否被at到
//
//	@param target
//	@param slice
//	@return bool
func IsInSlice(target string, slice []string) bool {
	for i := range slice {
		if slice[i] == target {
			return true
		}
	}
	return false
}

// MustAtoI 将字符串转换为int
//
//	@param str
//	@return int
func MustAtoI(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return i
	}
	return i
}

// GetUserInfo 获取用户信息
//
//	@param userID
//	@param guildID
//	@return userInfo
func GetUserInfo(userID, guildID string) (userInfo *kook.User, err error) {
	if guildID != "" {
		userInfo, err = betagovar.GlobalSession.UserView(userID, kook.UserViewWithGuildID(guildID))
	} else {
		userInfo, err = betagovar.GlobalSession.UserView(userID)
	}
	if err != nil {
		return
	}
	return
}

// GetGuildInfo 获取公会信息
//
//	@param guildID
//	@return guildInfo
//	@return err
func GetGuildInfo(guildID string) (guildInfo *kook.Guild, err error) {
	guildInfo, err = betagovar.GlobalSession.GuildView(guildID)
	if err != nil {
		return
	}
	return
}

// GetChannnelInfo  获取频道信息
//
//	@param channelID
//	@return channelInfo
//	@return err
func GetChannnelInfo(channelID string) (channelInfo *kook.Channel, err error) {
	channelInfo, err = betagovar.GlobalSession.ChannelView(channelID)
	if err != nil {
		return
	}
	return
}

// Struct2Map  将结构体转换为map
//
//	@param obj
//	@return map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if timestamp, ok := v.Field(i).Interface().(kook.MilliTimeStamp); ok {
			data[t.Field(i).Name] = time.Unix(int64(timestamp)/1000, 0).Format("2006-01-02 15:04:05")
		} else {
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return data
}

// BuildCardMessageCols 创建卡片消息的列
func BuildCardMessageCols(titleK, titleV string, kvMap map[string]interface{}) (res []interface{}, err error) {
	sectionElements := []interface{}{
		kook.CardMessageSection{
			Text: kook.CardMessageParagraph{
				Cols: 2,
				Fields: []interface{}{
					kook.CardMessageElementKMarkdown{
						Content: titleK,
					},
					kook.CardMessageElementKMarkdown{
						Content: titleV,
					},
				},
			},
		},
	}
	for k, v := range kvMap {
		if fmt.Sprint(v) == "" {
			continue
		}
		if strings.HasPrefix(fmt.Sprint(v), "http") {
			sectionElements = append(sectionElements,
				kook.CardMessageSection{
					Mode: kook.CardMessageSectionModeLeft,
					Text: kook.CardMessageParagraph{
						Cols: 2,
						Fields: []interface{}{
							kook.CardMessageElementKMarkdown{
								Content: "**" + k + "**",
							},
						},
					},
					Accessory: kook.CardMessageElementImage{
						Src:    fmt.Sprint(v),
						Size:   "sm",
						Circle: true,
					},
				},
			)
		} else {
			sectionElements = append(sectionElements,
				kook.CardMessageSection{
					Text: kook.CardMessageParagraph{
						Cols: 2,
						Fields: []interface{}{
							kook.CardMessageElementKMarkdown{
								Content: "**" + k + "**",
							},
							kook.CardMessageElementKMarkdown{
								Content: "`" + fmt.Sprintf("%#v", v) + "`",
							},
						},
					},
				},
			)
		}

	}
	res = sectionElements
	return
}
