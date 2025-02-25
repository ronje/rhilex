package apis

import (
	"fmt"

	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/plugin"
	"github.com/hootrhino/rhilex/typex"

	"github.com/gin-gonic/gin"
)

func InitPluginsRoute() {
	pluginsApi := server.RouteGroup(server.ContextUrl("/plugware"))
	{
		pluginsApi.GET(("/list"), server.AddRoute(Plugins))
		pluginsApi.POST(("/service"), server.AddRoute(PluginService))
		pluginsApi.GET(("/detail"), server.AddRoute(PluginDetail))
	}
}

/*
*
* 插件的服务接口
*
 */

func PluginService(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUID string `json:"uuid" binding:"required"`
		Name string `json:"name" binding:"required"`
		Args any    `json:"args"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	plugin := plugin.Find(form.UUID)
	if plugin != nil {
		result := plugin.Service(typex.ServiceArg{
			Name: form.Name,
			UUID: form.UUID,
			Args: form.Args,
		})
		c.JSON(common.HTTP_OK, common.OkWithData(result.Out))
		return
	}
	c.JSON(common.HTTP_OK, common.Error("plugin not exists:"+form.UUID))
}

/*
*
* 插件详情
*
 */
func PluginDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	plugin := plugin.Find(uuid)
	if plugin != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(plugin.PluginMetaInfo()))
		return
	}
	c.JSON(common.HTTP_OK, common.Error400EmptyObj(fmt.Errorf("no such plugin:%s", uuid)))
}
