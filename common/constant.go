package common

const Relatdb_Version = "Relatdb 1.0"
const Protocol_Version = 10

/*
 * Capabilities
 */
const (
	CLIENT_LONG_PASSWORD     = 0x00000001 //使用长密码认证协议。
	CLIENT_FOUND_ROWS        = 0x00000002 //返回受影响的行数而不是匹配的行数。
	CLIENT_LONG_FLAG         = 0x00000004 // 使用长 flag 字节。
	CLIENT_CONNECT_WITH_DB   = 0x00000008 // 为新的连接指定数据库。
	CLIENT_NO_SCHEMA         = 0x00000010 // 不允许“数据库名.表名.列名”这样的语法。这是对于ODBC的设置。当使用这样的语法时解析器会产生一个错误，这对于一些ODBC的程序限制bug来说是有用的。
	CLIENT_COMPRESS          = 0x00000020 // 使用压缩协议。
	CLIENT_ODBC              = 0x00000040 // ODBC 客户端。
	CLIENT_LOCAL_FILES       = 0x00000080 // 可以从本地读写文件。
	CLIENT_IGNORE_SPACE      = 0x00000100 // 忽略空格。允许在函数名后使用空格。所有函数名可以预留字。
	CLIENT_PROTOCOL_41       = 0x00000200 // 使用新的认证协议。
	CLIENT_INTERACTIVE       = 0x00000400 // 交互式客户端。允许使用关闭连接之前的不活动交互超时的描述，而不是等待超时秒数。客户端的会话等待超时变量变为交互超时变量。
	CLIENT_SSL               = 0x00000800 // 使用SSL。这个设置不应该被应用程序设置，他应该是在客户端库内部是设置的。可以在调用mysql_real_connect()之前调用mysql_ssl_set()来代替设置。
	CLIENT_IGNORE_SIGPIPE    = 0x00001000 // 忽略 SIGPIPE 信号。阻止客户端库安装一个SIGPIPE信号处理器。这个可以用于当应用程序已经安装该处理器的时候避免与其发生冲突。
	CLIENT_TRANSACTIONS      = 0x00002000 // 支持事务和多语句。
	CLIENT_RESERVED          = 0x00004000 // 保留字段。
	CLIENT_SECURE_CONNECTION = 0x00008000 // 支持加密连接。
	CLIENT_MULTI_STATEMENTS  = 0x00010000 // 通知服务器客户端可以发送多条语句（由分号分隔）。如果该标志为没有被设置，多条语句执行。
	CLIENT_MULTI_RESULTS     = 0x00020000 // 通知服务器客户端可以处理由多语句或者存储过程执行生成的多结果集。当打开CLIENT_MULTI_STATEMENTS时，这个标志自动的被打开。
)
