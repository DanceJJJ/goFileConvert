goFileConvert
============

目前已经完成
============

    Word、Excel、PPT转码为PDF

    PDF转码为图片

使用Docker快速部署
====

sudo docker pull DanceJJJ/goFileConvert

sudo docker run --name goFileConvert -ti -p 1180:80 DanceJJJ/goFileConvert /root/go/src/github.com/DanceJJJ/goFileConvert

你可以直接访问 http://127.0.0.1:1180/ 使用goFileConvert

部署编译
========

Windows版
----

    准备

        1.安装Libreoffice,下载官方msi包,傻瓜式安装即可 (https://zh-cn.libreoffice.org)

        2.将Libreoffice安装路径下的program文件夹加入PATH中
![](https://github.com/leeli73/goFileView/blob/master/media/win_path.png?raw=true)

        3.安装ImageMagick,官方包,傻瓜式安装即可,安装7.0以上版本 (https://ghostscript.com/download/gsdnld.html)

        4.安装GhostScript,官方包,傻瓜式安装即可 (https://ghostscript.com/download/gsdnld.html)

        5.git clone <https://github.com/DanceJJJ/goFileConvert.git>
    
    编译

        1.cd goFileConvert
        2.go build

    运行

        1. goFileConvert.exe
        2. 访问 http://127.0.0.1:8089/convert?url=被预览文件的url (例如 http://127.0.0.1:8089/perview/onlinePreview?url=http://127.0.0.1:88/test.docx)

    你可以在代码中修改监听的URL、端口等信息。

Linux版
----

    准备

        1.安装Libreoffice:sudo apt install libreoffice

        2.安装ImageMagick:sudo apt install imagemagick

        4.修改ImageMagick的配置,vi etc/ImageMagick-6/policy.xml

            修改
            <policy domain="coder" rights="none" pattern="PDF" />
            为
            <policy domain="coder" rights="read|write" pattern="PDF" />
            下方新增一行
            <policy domain="coder" rights="read|write" pattern="LABEL" />

            wq退出保存

        5.安装字体(如果出现乱码)

            打包一台Windows机器的C:\Windows\Fonts下的所有文件
            发送到Linux机器上
            解压后进入Fonts文件夹，依次执行mkfontscale,mkfontdir,fc-cache

        5.git clone <https://github.com/DanceJJJ/goFileConvert.git>
    
    编译

        1.cd goFileConvert
        2.go build

    运行

        1. ./goFileConvert
        2. 访问 http://127.0.0.1:8089/convert?url=被预览文件的url (例如 http://127.0.0.1:8089/perview/onlinePreview?url=http://127.0.0.1:88/test.docx)

    你可以再代码中修改监听的URL、端口等信息。

在自己的项目中集成
==================

准备
----

    go get github.com/DanceJJJ/goFileConvert

demo
----
```
package main

import(

    "net/http"

    "github.com/DanceJJJ/goFileConvert/convert"

)

func index(w http.ResponseWriter, r \*http.Request) {

    w.Write([]byte("I'm Index"))

}

func main(){

    http.HandleFunc("/convert/",convert.Handle) //绑定到preview的Handle

    http.ListenAndServe(":80", nil)

}
```
