# maxbook2
A method to download pdf files from max.book118.com

本程序在linux版本Ubuntu18.04下开发，理论上支持大多数linux版本，但暂未测试。
运行该程序需要系统支持如下命令：
convert
unionpdf

功能说明：
版本v1.0
程序默认将同目录下所有html文件检索一遍，支持模板能够转换成PDF的会放在download文件夹下。

使用方法： </br>
    1. 打开文档所在网址，https://max.book118.com/html/2020/0109/8105051126002072.shtm </br>
    2. 点击页面中的[同意并开始全文预览] </br>
    3. 此时在新弹开窗口的上方使用chrome右键开发者模式功能找到这个iframe对应的 </br>
    网址：/index.php?g=Home&m=NewView&a=index&aid=8105051126002072&v=20191216， </br>
    加上域名得到完整网址：https://max.book118.com/index.php?g=Home&m=NewView&a=index&aid=8105051126002072&v=20191216 </br>
    4. 使用新窗口将真实网址打开，然后预览所有页面，然后将页面另存为xxx.html（xxx为自定义名称） </br>
    5. 将xxx.html放在maxbook程序相同目录 </br>
    6. 使用命令行运行“./maxbook”，即可在download文件夹中生成你想要的pdf文件 </br>
