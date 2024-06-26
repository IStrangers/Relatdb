package server

const SERVER_VERSION = "Relatdb 1.0"
const SERVER_ROOT_USERNAME = "root"
const SERVER_ROOT_PASSWORD = "123456"
const PROTOCOL_VERSION = 10
const MAX_PACKET_SIZE = 16 * 1024 * 1024

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
	CLIENT_PLUGIN_AUTH       = 0x00080000 // 认证插件
)

/*
 * Commond
 */
const (
	/**
	 * none, this is an internal thread state
	 */
	COM_SLEEP = 0

	/**
	 * mysql_close
	 */
	COM_QUIT = 1

	/**
	 * mysql_select_db
	 */
	COM_INIT_DB = 2

	/**
	 * mysql_real_query
	 */
	COM_QUERY = 3

	/**
	 * mysql_list_fields
	 */
	COM_FIELD_LIST = 4

	/**
	 * mysql_create_db (deprecated)
	 */
	COM_CREATE_DB = 5

	/**
	 * mysql_drop_db (deprecated)
	 */
	COM_DROP_DB = 6

	/**
	 * mysql_refresh
	 */
	COM_REFRESH = 7

	/**
	 * mysql_shutdown
	 */
	COM_SHUTDOWN = 8

	/**
	 * mysql_stat
	 */
	COM_STATISTICS = 9

	/**
	 * mysql_list_processes
	 */
	COM_PROCESS_INFO = 10

	/**
	 * none, this is an internal thread state
	 */
	COM_CONNECT = 11

	/**
	 * mysql_kill
	 */
	COM_PROCESS_KILL = 12

	/**
	 * mysql_dump_debug_info
	 */
	COM_DEBUG = 13

	/**
	 * mysql_ping
	 */
	COM_PING = 14

	/**
	 * none, this is an internal thread state
	 */
	COM_TIME = 15

	/**
	 * none, this is an internal thread state
	 */
	COM_DELAYED_INSERT = 16

	/**
	 * mysql_change_user
	 */
	COM_CHANGE_USER = 17

	/**
	 * used by slave server mysqlbinlog
	 */
	COM_BINLOG_DUMP = 18

	/**
	 * used by slave server to getDecoder master table
	 */
	COM_TABLE_DUMP = 19

	/**
	 * used by slave to log connection to master
	 */
	COM_CONNECT_OUT = 20

	/**
	 * used by slave to register to master
	 */
	COM_REGISTER_SLAVE = 21

	/**
	 * mysql_stmt_prepare
	 */
	COM_STMT_PREPARE = 22

	/**
	 * mysql_stmt_execute
	 */
	COM_STMT_EXECUTE = 23

	/**
	 * mysql_stmt_send_long_data
	 */
	COM_STMT_SEND_LONG_DATA = 24

	/**
	 * mysql_stmt_close
	 */
	COM_STMT_CLOSE = 25

	/**
	 * mysql_stmt_reset
	 */
	COM_STMT_RESET = 26

	/**
	 * mysql_set_server_option
	 */
	COM_SET_OPTION = 27

	/**
	 * mysql_stmt_fetch
	 */
	COM_STMT_FETCH = 28

	/**
	 * cobar heartbeat
	 */
	COM_HEARTBEAT = 64

	/**
	 * MORE RESULTS
	 */
	SERVER_MORE_RESULTS_EXISTS = 8
)

/*
 * ErrorCode
 */
const (
	ERR_OPEN_SOCKET               = 3001
	ERR_CONNECT_SOCKET            = 3002
	ERR_FINISH_CONNECT            = 3003
	ERR_REGISTER                  = 3004
	ERR_READ                      = 3005
	ERR_PUT_WRITE_QUEUE           = 3006
	ERR_WRITE_BY_EVENT            = 3007
	ERR_WRITE_BY_QUEUE            = 3008
	ERR_HANDLE_DATA               = 3009
	ERR_EXCEPTION_CAUGHT          = 3010
	ERR_SINGLE_EXECUTE_NODES      = 3011
	NO_TRANSACTION_IN_MULTI_NODES = 3012

	// mysql errorMessage code
	ER_HASHCHK                                 = 1000
	ER_NISAMCHK                                = 1001
	ER_NO                                      = 1002
	ER_YES                                     = 1003
	ER_CANT_CREATE_FILE                        = 1004
	ER_CANT_CREATE_TABLE                       = 1005
	ER_CANT_CREATE_DB                          = 1006
	ER_DB_CREATE_EXISTS                        = 1007
	ER_DB_DROP_EXISTS                          = 1008
	ER_DB_DROP_DELETE                          = 1009
	ER_DB_DROP_RMDIR                           = 1010
	ER_CANT_DELETE_FILE                        = 1011
	ER_CANT_FIND_SYSTEM_REC                    = 1012
	ER_CANT_GET_STAT                           = 1013
	ER_CANT_GET_WD                             = 1014
	ER_CANT_LOCK                               = 1015
	ER_CANT_OPEN_FILE                          = 1016
	ER_FILE_NOT_FOUND                          = 1017
	ER_CANT_READ_DIR                           = 1018
	ER_CANT_SET_WD                             = 1019
	ER_CHECKREAD                               = 1020
	ER_DISK_FULL                               = 1021
	ER_DUP_KEY                                 = 1022
	ER_ERROR_ON_CLOSE                          = 1023
	ER_ERROR_ON_READ                           = 1024
	ER_ERROR_ON_RENAME                         = 1025
	ER_ERROR_ON_WRITE                          = 1026
	ER_FILE_USED                               = 1027
	ER_FILSORT_ABORT                           = 1028
	ER_FORM_NOT_FOUND                          = 1029
	ER_GET_ERRNO                               = 1030
	ER_ILLEGAL_HA                              = 1031
	ER_KEY_NOT_FOUND                           = 1032
	ER_NOT_FORM_FILE                           = 1033
	ER_NOT_KEYFILE                             = 1034
	ER_OLD_KEYFILE                             = 1035
	ER_OPEN_AS_READONLY                        = 1036
	ER_OUTOFMEMORY                             = 1037
	ER_OUT_OF_SORTMEMORY                       = 1038
	ER_UNEXPECTED_EOF                          = 1039
	ER_CON_COUNT_ERROR                         = 1040
	ER_OUT_OF_RESOURCES                        = 1041
	ER_BAD_HOST_ERROR                          = 1042
	ER_HANDSHAKE_ERROR                         = 1043
	ER_DBACCESS_DENIED_ERROR                   = 1044
	ER_ACCESS_DENIED_ERROR                     = 1045
	ER_NO_DB_ERROR                             = 1046
	ER_UNKNOWN_COM_ERROR                       = 1047
	ER_BAD_NULL_ERROR                          = 1048
	ER_BAD_DB_ERROR                            = 1049
	ER_TABLE_EXISTS_ERROR                      = 1050
	ER_BAD_TABLE_ERROR                         = 1051
	ER_NON_UNIQ_ERROR                          = 1052
	ER_SERVER_SHUTDOWN                         = 1053
	ER_BAD_FIELD_ERROR                         = 1054
	ER_WRONG_FIELD_WITH_GROUP                  = 1055
	ER_WRONG_GROUP_FIELD                       = 1056
	ER_WRONG_SUM_SELECT                        = 1057
	ER_WRONG_VALUE_COUNT                       = 1058
	ER_TOO_LONG_IDENT                          = 1059
	ER_DUP_FIELDNAME                           = 1060
	ER_DUP_KEYNAME                             = 1061
	ER_DUP_ENTRY                               = 1062
	ER_WRONG_FIELD_SPEC                        = 1063
	ER_PARSE_ERROR                             = 1064
	ER_EMPTY_QUERY                             = 1065
	ER_NONUNIQ_TABLE                           = 1066
	ER_INVALID_DEFAULT                         = 1067
	ER_MULTIPLE_PRI_KEY                        = 1068
	ER_TOO_MANY_KEYS                           = 1069
	ER_TOO_MANY_KEY_PARTS                      = 1070
	ER_TOO_LONG_KEY                            = 1071
	ER_KEY_COLUMN_DOES_NOT_EXITS               = 1072
	ER_BLOB_USED_AS_KEY                        = 1073
	ER_TOO_BIG_FIELDLENGTH                     = 1074
	ER_WRONG_AUTO_KEY                          = 1075
	ER_READY                                   = 1076
	ER_NORMAL_SHUTDOWN                         = 1077
	ER_GOT_SIGNAL                              = 1078
	ER_SHUTDOWN_COMPLETE                       = 1079
	ER_FORCING_CLOSE                           = 1080
	ER_IPSOCK_ERROR                            = 1081
	ER_NO_SUCH_INDEX                           = 1082
	ER_WRONG_FIELD_TERMINATORS                 = 1083
	ER_BLOBS_AND_NO_TERMINATED                 = 1084
	ER_TEXTFILE_NOT_READABLE                   = 1085
	ER_FILE_EXISTS_ERROR                       = 1086
	ER_LOAD_INFO                               = 1087
	ER_ALTER_INFO                              = 1088
	ER_WRONG_SUB_KEY                           = 1089
	ER_CANT_REMOVE_ALL_FIELDS                  = 1090
	ER_CANT_DROP_FIELD_OR_KEY                  = 1091
	ER_INSERT_INFO                             = 1092
	ER_UPDATE_TABLE_USED                       = 1093
	ER_NO_SUCH_THREAD                          = 1094
	ER_KILL_DENIED_ERROR                       = 1095
	ER_NO_TABLES_USED                          = 1096
	ER_TOO_BIG_SET                             = 1097
	ER_NO_UNIQUE_LOGFILE                       = 1098
	ER_TABLE_NOT_LOCKED_FOR_WRITE              = 1099
	ER_TABLE_NOT_LOCKED                        = 1100
	ER_BLOB_CANT_HAVE_DEFAULT                  = 1101
	ER_WRONG_DB_NAME                           = 1102
	ER_WRONG_TABLE_NAME                        = 1103
	ER_TOO_BIG_SELECT                          = 1104
	ER_UNKNOWN_ERROR                           = 1105
	ER_UNKNOWN_PROCEDURE                       = 1106
	ER_WRONG_PARAMCOUNT_TO_PROCEDURE           = 1107
	ER_WRONG_PARAMETERS_TO_PROCEDURE           = 1108
	ER_UNKNOWN_TABLE                           = 1109
	ER_FIELD_SPECIFIED_TWICE                   = 1110
	ER_INVALID_GROUP_FUNC_USE                  = 1111
	ER_UNSUPPORTED_EXTENSION                   = 1112
	ER_TABLE_MUST_HAVE_COLUMNS                 = 1113
	ER_RECORD_FILE_FULL                        = 1114
	ER_UNKNOWN_CHARACTER_SET                   = 1115
	ER_TOO_MANY_TABLES                         = 1116
	ER_TOO_MANY_FIELDS                         = 1117
	ER_TOO_BIG_ROWSIZE                         = 1118
	ER_STACK_OVERRUN                           = 1119
	ER_WRONG_OUTER_JOIN                        = 1120
	ER_NULL_COLUMN_IN_INDEX                    = 1121
	ER_CANT_FIND_UDF                           = 1122
	ER_CANT_INITIALIZE_UDF                     = 1123
	ER_UDF_NO_PATHS                            = 1124
	ER_UDF_EXISTS                              = 1125
	ER_CANT_OPEN_LIBRARY                       = 1126
	ER_CANT_FIND_DL_ENTRY                      = 1127
	ER_FUNCTION_NOT_DEFINED                    = 1128
	ER_HOST_IS_BLOCKED                         = 1129
	ER_HOST_NOT_PRIVILEGED                     = 1130
	ER_PASSWORD_ANONYMOUS_USER                 = 1131
	ER_PASSWORD_NOT_ALLOWED                    = 1132
	ER_PASSWORD_NO_MATCH                       = 1133
	ER_UPDATE_INFO                             = 1134
	ER_CANT_CREATE_THREAD                      = 1135
	ER_WRONG_VALUE_COUNT_ON_ROW                = 1136
	ER_CANT_REOPEN_TABLE                       = 1137
	ER_INVALID_USE_OF_NULL                     = 1138
	ER_REGEXP_ERROR                            = 1139
	ER_MIX_OF_GROUP_FUNC_AND_FIELDS            = 1140
	ER_NONEXISTING_GRANT                       = 1141
	ER_TABLEACCESS_DENIED_ERROR                = 1142
	ER_COLUMNACCESS_DENIED_ERROR               = 1143
	ER_ILLEGAL_GRANT_FOR_TABLE                 = 1144
	ER_GRANT_WRONG_HOST_OR_USER                = 1145
	ER_NO_SUCH_TABLE                           = 1146
	ER_NONEXISTING_TABLE_GRANT                 = 1147
	ER_NOT_ALLOWED_COMMAND                     = 1148
	ER_SYNTAX_ERROR                            = 1149
	ER_DELAYED_CANT_CHANGE_LOCK                = 1150
	ER_TOO_MANY_DELAYED_THREADS                = 1151
	ER_ABORTING_CONNECTION                     = 1152
	ER_NET_PACKET_TOO_LARGE                    = 1153
	ER_NET_READ_ERROR_FROM_PIPE                = 1154
	ER_NET_FCNTL_ERROR                         = 1155
	ER_NET_PACKETS_OUT_OF_ORDER                = 1156
	ER_NET_UNCOMPRESS_ERROR                    = 1157
	ER_NET_READ_ERROR                          = 1158
	ER_NET_READ_ERRUPTED                       = 1159
	ER_NET_ERROR_ON_WRITE                      = 1160
	ER_NET_WRITE_ERRUPTED                      = 1161
	ER_TOO_LONG_STRING                         = 1162
	ER_TABLE_CANT_HANDLE_BLOB                  = 1163
	ER_TABLE_CANT_HANDLE_AUTO_INCREMENT        = 1164
	ER_DELAYED_INSERT_TABLE_LOCKED             = 1165
	ER_WRONG_COLUMN_NAME                       = 1166
	ER_WRONG_KEY_COLUMN                        = 1167
	ER_WRONG_MRG_TABLE                         = 1168
	ER_DUP_UNIQUE                              = 1169
	ER_BLOB_KEY_WITHOUT_LENGTH                 = 1170
	ER_PRIMARY_CANT_HAVE_NULL                  = 1171
	ER_TOO_MANY_ROWS                           = 1172
	ER_REQUIRES_PRIMARY_KEY                    = 1173
	ER_NO_RAID_COMPILED                        = 1174
	ER_UPDATE_WITHOUT_KEY_IN_SAFE_MODE         = 1175
	ER_KEY_DOES_NOT_EXITS                      = 1176
	ER_CHECK_NO_SUCH_TABLE                     = 1177
	ER_CHECK_NOT_IMPLEMENTED                   = 1178
	ER_CANT_DO_THIS_DURING_AN_TRANSACTION      = 1179
	ER_ERROR_DURING_COMMIT                     = 1180
	ER_ERROR_DURING_ROLLBACK                   = 1181
	ER_ERROR_DURING_FLUSH_LOGS                 = 1182
	ER_ERROR_DURING_CHECKPO                    = 1183
	ER_NEW_ABORTING_CONNECTION                 = 1184
	ER_DUMP_NOT_IMPLEMENTED                    = 1185
	ER_FLUSH_MASTER_BINLOG_CLOSED              = 1186
	ER_INDEX_REBUILD                           = 1187
	ER_MASTER                                  = 1188
	ER_MASTER_NET_READ                         = 1189
	ER_MASTER_NET_WRITE                        = 1190
	ER_FT_MATCHING_KEY_NOT_FOUND               = 1191
	ER_LOCK_OR_ACTIVE_TRANSACTION              = 1192
	ER_UNKNOWN_SYSTEM_VARIABLE                 = 1193
	ER_CRASHED_ON_USAGE                        = 1194
	ER_CRASHED_ON_REPAIR                       = 1195
	ER_WARNING_NOT_COMPLETE_ROLLBACK           = 1196
	ER_TRANS_CACHE_FULL                        = 1197
	ER_SLAVE_MUST_STOP                         = 1198
	ER_SLAVE_NOT_RUNNING                       = 1199
	ER_BAD_SLAVE                               = 1200
	ER_MASTER_INFO                             = 1201
	ER_SLAVE_THREAD                            = 1202
	ER_TOO_MANY_USER_CONNECTIONS               = 1203
	ER_SET_CONSTANTS_ONLY                      = 1204
	ER_LOCK_WAIT_TIMEOUT                       = 1205
	ER_LOCK_TABLE_FULL                         = 1206
	ER_READ_ONLY_TRANSACTION                   = 1207
	ER_DROP_DB_WITH_READ_LOCK                  = 1208
	ER_CREATE_DB_WITH_READ_LOCK                = 1209
	ER_WRONG_ARGUMENTS                         = 1210
	ER_NO_PERMISSION_TO_CREATE_USER            = 1211
	ER_UNION_TABLES_IN_DIFFERENT_DIR           = 1212
	ER_LOCK_DEADLOCK                           = 1213
	ER_TABLE_CANT_HANDLE_FT                    = 1214
	ER_CANNOT_ADD_FOREIGN                      = 1215
	ER_NO_REFERENCED_ROW                       = 1216
	ER_ROW_IS_REFERENCED                       = 1217
	ER_CONNECT_TO_MASTER                       = 1218
	ER_QUERY_ON_MASTER                         = 1219
	ER_ERROR_WHEN_EXECUTING_COMMAND            = 1220
	ER_WRONG_USAGE                             = 1221
	ER_WRONG_NUMBER_OF_COLUMNS_IN_SELECT       = 1222
	ER_CANT_UPDATE_WITH_READLOCK               = 1223
	ER_MIXING_NOT_ALLOWED                      = 1224
	ER_DUP_ARGUMENT                            = 1225
	ER_USER_LIMIT_REACHED                      = 1226
	ER_SPECIFIC_ACCESS_DENIED_ERROR            = 1227
	ER_LOCAL_VARIABLE                          = 1228
	ER_GLOBAL_VARIABLE                         = 1229
	ER_NO_DEFAULT                              = 1230
	ER_WRONG_VALUE_FOR_VAR                     = 1231
	ER_WRONG_TYPE_FOR_VAR                      = 1232
	ER_VAR_CANT_BE_READ                        = 1233
	ER_CANT_USE_OPTION_HERE                    = 1234
	ER_NOT_SUPPORTED_YET                       = 1235
	ER_MASTER_FATAL_ERROR_READING_BINLOG       = 1236
	ER_SLAVE_IGNORED_TABLE                     = 1237
	ER_INCORRECT_GLOBAL_LOCAL_VAR              = 1238
	ER_WRONG_FK_DEF                            = 1239
	ER_KEY_REF_DO_NOT_MATCH_TABLE_REF          = 1240
	ER_OPERAND_COLUMNS                         = 1241
	ER_SUBQUERY_NO_1_ROW                       = 1242
	ER_UNKNOWN_STMT_HANDLER                    = 1243
	ER_CORRUPT_HELP_DB                         = 1244
	ER_CYCLIC_REFERENCE                        = 1245
	ER_AUTO_CONVERT                            = 1246
	ER_ILLEGAL_REFERENCE                       = 1247
	ER_DERIVED_MUST_HAVE_ALIAS                 = 1248
	ER_SELECT_REDUCED                          = 1249
	ER_TABLENAME_NOT_ALLOWED_HERE              = 1250
	ER_NOT_SUPPORTED_AUTH_MODE                 = 1251
	ER_SPATIAL_CANT_HAVE_NULL                  = 1252
	ER_COLLATION_CHARSET_MISMATCH              = 1253
	ER_SLAVE_WAS_RUNNING                       = 1254
	ER_SLAVE_WAS_NOT_RUNNING                   = 1255
	ER_TOO_BIG_FOR_UNCOMPRESS                  = 1256
	ER_ZLIB_Z_MEM_ERROR                        = 1257
	ER_ZLIB_Z_BUF_ERROR                        = 1258
	ER_ZLIB_Z_DATA_ERROR                       = 1259
	ER_CUT_VALUE_GROUP_CONCAT                  = 1260
	ER_WARN_TOO_FEW_RECORDS                    = 1261
	ER_WARN_TOO_MANY_RECORDS                   = 1262
	ER_WARN_NULL_TO_NOTNULL                    = 1263
	ER_WARN_DATA_OUT_OF_RANGE                  = 1264
	WARN_DATA_TRUNCATED                        = 1265
	ER_WARN_USING_OTHER_HANDLER                = 1266
	ER_CANT_AGGREGATE_2COLLATIONS              = 1267
	ER_DROP_USER                               = 1268
	ER_REVOKE_GRANTS                           = 1269
	ER_CANT_AGGREGATE_3COLLATIONS              = 1270
	ER_CANT_AGGREGATE_NCOLLATIONS              = 1271
	ER_VARIABLE_IS_NOT_STRUCT                  = 1272
	ER_UNKNOWN_COLLATION                       = 1273
	ER_SLAVE_IGNORED_SSL_PARAMS                = 1274
	ER_SERVER_IS_IN_SECURE_AUTH_MODE           = 1275
	ER_WARN_FIELD_RESOLVED                     = 1276
	ER_BAD_SLAVE_UNTIL_COND                    = 1277
	ER_MISSING_SKIP_SLAVE                      = 1278
	ER_UNTIL_COND_IGNORED                      = 1279
	ER_WRONG_NAME_FOR_INDEX                    = 1280
	ER_WRONG_NAME_FOR_CATALOG                  = 1281
	ER_WARN_QC_RESIZE                          = 1282
	ER_BAD_FT_COLUMN                           = 1283
	ER_UNKNOWN_KEY_CACHE                       = 1284
	ER_WARN_HOSTNAME_WONT_WORK                 = 1285
	ER_UNKNOWN_STORAGE_ENGINE                  = 1286
	ER_WARN_DEPRECATED_SYNTAX                  = 1287
	ER_NON_UPDATABLE_TABLE                     = 1288
	ER_FEATURE_DISABLED                        = 1289
	ER_OPTION_PREVENTS_STATEMENT               = 1290
	ER_DUPLICATED_VALUE_IN_TYPE                = 1291
	ER_TRUNCATED_WRONG_VALUE                   = 1292
	ER_TOO_MUCH_AUTO_TIMESTAMP_COLS            = 1293
	ER_INVALID_ON_UPDATE                       = 1294
	ER_UNSUPPORTED_PS                          = 1295
	ER_GET_ERRMSG                              = 1296
	ER_GET_TEMPORARY_ERRMSG                    = 1297
	ER_UNKNOWN_TIME_ZONE                       = 1298
	ER_WARN_INVALID_TIMESTAMP                  = 1299
	ER_INVALID_CHARACTER_STRING                = 1300
	ER_WARN_ALLOWED_PACKET_OVERFLOWED          = 1301
	ER_CONFLICTING_DECLARATIONS                = 1302
	ER_SP_NO_RECURSIVE_CREATE                  = 1303
	ER_SP_ALREADY_EXISTS                       = 1304
	ER_SP_DOES_NOT_EXIST                       = 1305
	ER_SP_DROP_FAILED                          = 1306
	ER_SP_STORE_FAILED                         = 1307
	ER_SP_LILABEL_MISMATCH                     = 1308
	ER_SP_LABEL_REDEFINE                       = 1309
	ER_SP_LABEL_MISMATCH                       = 1310
	ER_SP_UNINIT_VAR                           = 1311
	ER_SP_BADSELECT                            = 1312
	ER_SP_BADRETURN                            = 1313
	ER_SP_BADSTATEMENT                         = 1314
	ER_UPDATE_LOG_DEPRECATED_IGNORED           = 1315
	ER_UPDATE_LOG_DEPRECATED_TRANSLATED        = 1316
	ER_QUERY_ERRUPTED                          = 1317
	ER_SP_WRONG_NO_OF_ARGS                     = 1318
	ER_SP_COND_MISMATCH                        = 1319
	ER_SP_NORETURN                             = 1320
	ER_SP_NORETURNEND                          = 1321
	ER_SP_BAD_CURSOR_QUERY                     = 1322
	ER_SP_BAD_CURSOR_SELECT                    = 1323
	ER_SP_CURSOR_MISMATCH                      = 1324
	ER_SP_CURSOR_ALREADY_OPEN                  = 1325
	ER_SP_CURSOR_NOT_OPEN                      = 1326
	ER_SP_UNDECLARED_VAR                       = 1327
	ER_SP_WRONG_NO_OF_FETCH_ARGS               = 1328
	ER_SP_FETCH_NO_DATA                        = 1329
	ER_SP_DUP_PARAM                            = 1330
	ER_SP_DUP_VAR                              = 1331
	ER_SP_DUP_COND                             = 1332
	ER_SP_DUP_CURS                             = 1333
	ER_SP_CANT_ALTER                           = 1334
	ER_SP_SUBSELECT_NYI                        = 1335
	ER_STMT_NOT_ALLOWED_IN_SF_OR_TRG           = 1336
	ER_SP_VARCOND_AFTER_CURSHNDLR              = 1337
	ER_SP_CURSOR_AFTER_HANDLER                 = 1338
	ER_SP_CASE_NOT_FOUND                       = 1339
	ER_FPARSER_TOO_BIG_FILE                    = 1340
	ER_FPARSER_BAD_HEADER                      = 1341
	ER_FPARSER_EOF_IN_COMMENT                  = 1342
	ER_FPARSER_ERROR_IN_PARAMETER              = 1343
	ER_FPARSER_EOF_IN_UNKNOWN_PARAMETER        = 1344
	ER_VIEW_NO_EXPLAIN                         = 1345
	ER_FRM_UNKNOWN_TYPE                        = 1346
	ER_WRONG_OBJECT                            = 1347
	ER_NONUPDATEABLE_COLUMN                    = 1348
	ER_VIEW_SELECT_DERIVED                     = 1349
	ER_VIEW_SELECT_CLAUSE                      = 1350
	ER_VIEW_SELECT_VARIABLE                    = 1351
	ER_VIEW_SELECT_TMPTABLE                    = 1352
	ER_VIEW_WRONG_LIST                         = 1353
	ER_WARN_VIEW_MERGE                         = 1354
	ER_WARN_VIEW_WITHOUT_KEY                   = 1355
	ER_VIEW_INVALID                            = 1356
	ER_SP_NO_DROP_SP                           = 1357
	ER_SP_GOTO_IN_HNDLR                        = 1358
	ER_TRG_ALREADY_EXISTS                      = 1359
	ER_TRG_DOES_NOT_EXIST                      = 1360
	ER_TRG_ON_VIEW_OR_TEMP_TABLE               = 1361
	ER_TRG_CANT_CHANGE_ROW                     = 1362
	ER_TRG_NO_SUCH_ROW_IN_TRG                  = 1363
	ER_NO_DEFAULT_FOR_FIELD                    = 1364
	ER_DIVISION_BY_ZERO                        = 1365
	ER_TRUNCATED_WRONG_VALUE_FOR_FIELD         = 1366
	ER_ILLEGAL_VALUE_FOR_TYPE                  = 1367
	ER_VIEW_NONUPD_CHECK                       = 1368
	ER_VIEW_CHECK_FAILED                       = 1369
	ER_PROCACCESS_DENIED_ERROR                 = 1370
	ER_RELAY_LOG_FAIL                          = 1371
	ER_PASSWD_LENGTH                           = 1372
	ER_UNKNOWN_TARGET_BINLOG                   = 1373
	ER_IO_ERR_LOG_INDEX_READ                   = 1374
	ER_BINLOG_PURGE_PROHIBITED                 = 1375
	ER_FSEEK_FAIL                              = 1376
	ER_BINLOG_PURGE_FATAL_ERR                  = 1377
	ER_LOG_IN_USE                              = 1378
	ER_LOG_PURGE_UNKNOWN_ERR                   = 1379
	ER_RELAY_LOG_INIT                          = 1380
	ER_NO_BINARY_LOGGING                       = 1381
	ER_RESERVED_SYNTAX                         = 1382
	ER_WSAS_FAILED                             = 1383
	ER_DIFF_GROUPS_PROC                        = 1384
	ER_NO_GROUP_FOR_PROC                       = 1385
	ER_ORDER_WITH_PROC                         = 1386
	ER_LOGGING_PROHIBIT_CHANGING_OF            = 1387
	ER_NO_FILE_MAPPING                         = 1388
	ER_WRONG_MAGIC                             = 1389
	ER_PS_MANY_PARAM                           = 1390
	ER_KEY_PART_0                              = 1391
	ER_VIEW_CHECKSUM                           = 1392
	ER_VIEW_MULTIUPDATE                        = 1393
	ER_VIEW_NO_INSERT_FIELD_LIST               = 1394
	ER_VIEW_DELETE_MERGE_VIEW                  = 1395
	ER_CANNOT_USER                             = 1396
	ER_XAER_NOTA                               = 1397
	ER_XAER_INVAL                              = 1398
	ER_XAER_RMFAIL                             = 1399
	ER_XAER_OUTSIDE                            = 1400
	ER_XAER_RMERR                              = 1401
	ER_XA_RBROLLBACK                           = 1402
	ER_NONEXISTING_PROC_GRANT                  = 1403
	ER_PROC_AUTO_GRANT_FAIL                    = 1404
	ER_PROC_AUTO_REVOKE_FAIL                   = 1405
	ER_DATA_TOO_LONG                           = 1406
	ER_SP_BAD_SQLSTATE                         = 1407
	ER_STARTUP                                 = 1408
	ER_LOAD_FROM_FIXED_SIZE_ROWS_TO_VAR        = 1409
	ER_CANT_CREATE_USER_WITH_GRANT             = 1410
	ER_WRONG_VALUE_FOR_TYPE                    = 1411
	ER_TABLE_DEF_CHANGED                       = 1412
	ER_SP_DUP_HANDLER                          = 1413
	ER_SP_NOT_VAR_ARG                          = 1414
	ER_SP_NO_RETSET                            = 1415
	ER_CANT_CREATE_GEOMETRY_OBJECT             = 1416
	ER_FAILED_ROUTINE_BREAK_BINLOG             = 1417
	ER_BINLOG_UNSAFE_ROUTINE                   = 1418
	ER_BINLOG_CREATE_ROUTINE_NEED_SUPER        = 1419
	ER_EXEC_STMT_WITH_OPEN_CURSOR              = 1420
	ER_STMT_HAS_NO_OPEN_CURSOR                 = 1421
	ER_COMMIT_NOT_ALLOWED_IN_SF_OR_TRG         = 1422
	ER_NO_DEFAULT_FOR_VIEW_FIELD               = 1423
	ER_SP_NO_RECURSION                         = 1424
	ER_TOO_BIG_SCALE                           = 1425
	ER_TOO_BIG_PRECISION                       = 1426
	ER_M_BIGGER_THAN_D                         = 1427
	ER_WRONG_LOCK_OF_SYSTEM_TABLE              = 1428
	ER_CONNECT_TO_FOREIGN_DATA_SOURCE          = 1429
	ER_QUERY_ON_FOREIGN_DATA_SOURCE            = 1430
	ER_FOREIGN_DATA_SOURCE_DOESNT_EXIST        = 1431
	ER_FOREIGN_DATA_STRING_INVALID_CANT_CREATE = 1432
	ER_FOREIGN_DATA_STRING_INVALID             = 1433
	ER_CANT_CREATE_FEDERATED_TABLE             = 1434
	ER_TRG_IN_WRONG_SCHEMA                     = 1435
	ER_STACK_OVERRUN_NEED_MORE                 = 1436
	ER_TOO_LONG_BODY                           = 1437
	ER_WARN_CANT_DROP_DEFAULT_KEYCACHE         = 1438
	ER_TOO_BIG_DISPLAYWIDTH                    = 1439
	ER_XAER_DUPID                              = 1440
	ER_DATETIME_FUNCTION_OVERFLOW              = 1441
	ER_CANT_UPDATE_USED_TABLE_IN_SF_OR_TRG     = 1442
	ER_VIEW_PREVENT_UPDATE                     = 1443
	ER_PS_NO_RECURSION                         = 1444
	ER_SP_CANT_SET_AUTOCOMMIT                  = 1445
	ER_NO_VIEW_USER                            = 1446
	ER_VIEW_FRM_NO_USER                        = 1447
	ER_VIEW_OTHER_USER                         = 1448
	ER_NO_SUCH_USER                            = 1449
	ER_FORBID_SCHEMA_CHANGE                    = 1450
	ER_ROW_IS_REFERENCED_2                     = 1451
	ER_NO_REFERENCED_ROW_2                     = 1452
	ER_SP_BAD_VAR_SHADOW                       = 1453
	ER_PARTITION_REQUIRES_VALUES_ERROR         = 1454
	ER_PARTITION_WRONG_VALUES_ERROR            = 1455
	ER_PARTITION_MAXVALUE_ERROR                = 1456
	ER_PARTITION_SUBPARTITION_ERROR            = 1457
	ER_PARTITION_WRONG_NO_PART_ERROR           = 1458
	ER_PARTITION_WRONG_NO_SUBPART_ERROR        = 1459
	ER_CONST_EXPR_IN_PARTITION_FUNC_ERROR      = 1460
	ER_NO_CONST_EXPR_IN_RANGE_OR_LIST_ERROR    = 1461
	ER_FIELD_NOT_FOUND_PART_ERROR              = 1462
	ER_LIST_OF_FIELDS_ONLY_IN_HASH_ERROR       = 1463
	ER_INCONSISTENT_PARTITION_INFO_ERROR       = 1464
	ER_PARTITION_FUNC_NOT_ALLOWED_ERROR        = 1465
	ER_PARTITIONS_MUST_BE_DEFINED_ERROR        = 1466
	ER_RANGE_NOT_INCREASING_ERROR              = 1467
	ER_INCONSISTENT_TYPE_OF_FUNCTIONS_ERROR    = 1468
	ER_MULTIPLE_DEF_CONST_IN_LIST_PART_ERROR   = 1469
	ER_PARTITION_ENTRY_ERROR                   = 1470
	ER_MIX_HANDLER_ERROR                       = 1471
	ER_PARTITION_NOT_DEFINED_ERROR             = 1472
	ER_TOO_MANY_PARTITIONS_ERROR               = 1473
	ER_SUBPARTITION_ERROR                      = 1474
	ER_CANT_CREATE_HANDLER_FILE                = 1475
	ER_BLOB_FIELD_IN_PART_FUNC_ERROR           = 1476
	ER_CHAR_SET_IN_PART_FIELD_ERROR            = 1477
	ER_UNIQUE_KEY_NEED_ALL_FIELDS_IN_PF        = 1478
	ER_NO_PARTS_ERROR                          = 1479
	ER_PARTITION_MGMT_ON_NONPARTITIONED        = 1480
	ER_DROP_PARTITION_NON_EXISTENT             = 1481
	ER_DROP_LAST_PARTITION                     = 1482
	ER_COALESCE_ONLY_ON_HASH_PARTITION         = 1483
	ER_ONLY_ON_RANGE_LIST_PARTITION            = 1484
	ER_ADD_PARTITION_SUBPART_ERROR             = 1485
	ER_ADD_PARTITION_NO_NEW_PARTITION          = 1486
	ER_COALESCE_PARTITION_NO_PARTITION         = 1487
	ER_REORG_PARTITION_NOT_EXIST               = 1488
	ER_SAME_NAME_PARTITION                     = 1489
	ER_CONSECUTIVE_REORG_PARTITIONS            = 1490
	ER_REORG_OUTSIDE_RANGE                     = 1491
	ER_DROP_PARTITION_FAILURE                  = 1492
	ER_DROP_PARTITION_WHEN_FK_DEFINED          = 1493
	ER_PLUGIN_IS_NOT_LOADED                    = 1494
	ER_MULTI_QUERY_TIMEOUT                     = 1495
	ER_MULTI_EXEC_ERROR                        = 1496
)
