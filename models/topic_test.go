package models

import (
	//"fmt"
	// "github.com/astaxie/beego/orm"
	"testing"
)

// func init() {
// 	orm.RegisterDataBase("default", "mysql", "root@/gotalk?charset=utf8", 30)
// 	orm.RunSyncdb("default", true, false)
// }

func TestTopicInsertRead(t *testing.T) {
	user := &User{}
	user.Id = 1
	user.Read()
	category := &Category{Id: 1}
	category.Read()
	Content := &Content{}
	topic := Topic{}
	topic.User = user
	topic.Category = category
	topic.Title = "为什么要使用 Go 语言，Go 语言的优势在哪里？"
	topic.Content = Content
	topic.Content.Message = `我尝试来回答你几个问题：
1、Go有什么优势

    可直接编译成机器码，不依赖其他库，glibc的版本有一定要求，部署就是扔一个文件上去就完成了。
    静态类型语言，但是有动态语言的感觉，静态类型的语言就是可以在编译的时候检查出来隐藏的大多数问题，动态语言的感觉就是有很多的包可以使用，写起来的效率很高。
    语言层面支持并发，这个就是Go最大的特色，天生的支持并发，我曾经说过一句话，天生的基因和整容是有区别的，大家一样美丽，但是你喜欢整容的还是天生基因的美丽呢？Go就是基因里面支持的并发，可以充分的利用多核，很容易的使用并发。
    内置runtime，支持垃圾回收，这属于动态语言的特性之一吧，虽然目前来说GC不算完美，但是足以应付我们所能遇到的大多数情况，特别是Go1.1之后的GC。
    简单易学，Go语言的作者都有C的基因，那么Go自然而然就有了C的基因，那么Go关键字是25个，但是表达能力很强大，几乎支持大多数你在其他语言见过的特性：继承、重载、对象等。
    丰富的标准库，Go目前已经内置了大量的库，特别是网络库非常强大，我最爱的也是这部分。
    内置强大的工具，Go语言里面内置了很多工具链，最好的应该是gofmt工具，自动化格式化代码，能够让团队review变得如此的简单，代码格式一模一样，想不一样都很困难。
    跨平台编译，如果你写的Go代码不包含cgo，那么就可以做到window系统编译linux的应用，如何做到的呢？Go引用了plan9的代码，这就是不依赖系统的信息。
    内嵌C支持，前面说了作者是C的作者，所以Go里面也可以直接包含c代码，利用现有的丰富的C库。

2、Go适合用来做什么

    服务器编程，以前你如果使用C或者C++做的那些事情，用Go来做很合适，例如处理日志、数据打包、虚拟机处理、文件系统等。
    分布式系统，数据库代理器等
    网络编程，这一块目前应用最广，包括Web应用、API应用、下载应用、
    内存数据库，前一段时间google开发的groupcache，couchbase的部分组建
    云平台，目前国外很多云平台在采用Go开发，CloudFoundy的部分组建，前VMare的技术总监自己出来搞的apcera云平台。

3、Go成功的项目
nsq：bitly开源的消息队列系统，性能非常高，目前他们每天处理数十亿条的消息
docker:基于lxc的一个虚拟打包工具，能够实现PAAS平台的组建。
packer:用来生成不同平台的镜像文件，例如VM、vbox、AWS等，作者是vagrant的作者
skynet：分布式调度框架
Doozer：分布式同步工具，类似ZooKeeper
Heka：mazila开源的日志处理系统
cbfs：couchbase开源的分布式文件系统
tsuru：开源的PAAS平台，和SAE实现的功能一模一样
groupcache：memcahe作者写的用于Google下载系统的缓存系统
god：类似redis的缓存系统，但是支持分布式和扩展性
gor：网络流量抓包和重放工具
以下是一些公司，只是一小部分：

    http://Apcera.com
    http://Stathat.com
    Juju at Canonical/Ubuntu, presentation
    http://Beachfront.iO at Beachfront Media
    CloudFlare
    Soundcloud
    Mozilla
    Disqus
    http://Bit.ly
    Heroku
    google
    youtube

下面列出来了一些使用的用户
GoUsers - go-wiki - A list of organizations that use Go.
4、Go还存在的缺点
以下缺点是我自己在项目开发中遇到的一些问题：

    Go的import包不支持版本，有时候升级容易导致项目不可运行，所以需要自己控制相应的版本信息
    Go的goroutine一旦启动之后，不同的goroutine之间切换不是受程序控制，runtime调度的时候，需要严谨的逻辑，不然goroutine休眠，过一段时间逻辑结束了，突然冒出来又执行了，会导致逻辑出错等情况。
    GC延迟有点大，我开发的日志系统伤过一次，同时并发很大的情况下，处理很大的日志，GC没有那么快，内存回收不给力，后来经过profile程序改进之后得到了改善。
    pkg下面的图片处理库很多bug，还是使用成熟产品好，调用这些成熟库imagemagick的接口比较靠谱



最后还是建议大家学习Go，这门语言真的值得大家好好学习，因为它可以做从底层到前端的任何工作。

学习Go的话欢迎大家通过我写的书来学习，我已经开源在github：
astaxie/build-web-application-with-golang · GitHub

还有如果你用来做API开发或者网络开发，那么我做的开源框架beego也许适合你，可以适当的来学习一下：
astaxie/beego · GitHub`
	err := topic.Insert()
	if err != nil {
		t.Error(err)
	}
	topic2 := topic
	topic2.Read()
	if topic2.Content.Message != topic.Content.Message {
		t.Error("Content different")
	}
	t.Error(&topic.Content, &topic2.Content)
}

func TestTopicUpdate(t *testing.T) {
	topic := Topic{Id: 1}
	topic.Read()
	newTitle := "New Title"
	newMessage := "New message"
	topic.Title = newTitle
	topic.Content.Message = newMessage
	topic.Update()
	topic.Read()
	if topic.Title != newTitle || topic.Content.Message != newMessage {
		t.Error("Update failed! Content mismatch!")
	}
}

func TestTopicDelete(t *testing.T) {
	topic := Topic{Id: 1}
	topic.Read()
	topic.Delete()
	err := topic.Read()
	if err == nil {
		t.Error(err)
	}
}
