
"file_data_dir": "F:\\tmp",     //指定http文件系统的根目录路径

"listen_host": "192.168.1.119",     //程序启动监听的ip

"listen_port": "7879",      //程序启动监听的端口

"start_white_list": "true",     //是否开启程序白名单,true/false选择。

"white_list": "127.0.0.1,192.168.1.1,192.168.1.119",        //基于start_white_list选项是否开启，如果开启这里必须写上可访问的ip地址,多个ip以逗号分隔

"passwd_admin_white_list": "127.0.0.1,192.168.1.119",       //允许哪些ip访问密码加密保存页面选项,多个ip以逗号分隔

"debug_mode": "debug",      //程序启动选择的模式 debug/release

"access_dir": "/root",      //选择密码加密目录

"hs_dir": "F:\\Go1.20.6\\go-project\\src\\ls\\guide_HS",        //文件删除临时回收站

"interface_name": "eth0",   // 监听的网卡接口(已经弃用)

"start_mon": "false"    // 是否开启监控(已经弃用)
