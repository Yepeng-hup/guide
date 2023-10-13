# guide 公司内部系统集合工具
guide 支持系统url管理，文件管理(支持多文件上传，下载文件(目录上传正在开发中)，在线文件查看，在线删除(可批量)文件及目录，在线创建目录，在线压缩，在线解压(目前只实现了4种解压格式[.gz, .tar.gz, .tar, .zip]))，其他正在开发中^_^

linux run:

```shell
export GUIDE_HOST=0.0.0.0 GUIDE_PORT=7878 GUIDE_FILEDATA_DIR=/opt GUIDE_INTERFACE_NAME=eth0 && ./guide

```

win run:

```shell
set GUIDE_HOST=0.0.0.0 GUIDE_PORT=7878 GUIDE_FILEDATA_DIR=D:\tmp GUIDE_INTERFACE_NAME=以太网 1
./guide.exe

```