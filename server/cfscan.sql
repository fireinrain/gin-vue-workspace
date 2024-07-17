CREATE TABLE asn_info
(
    id                 INTEGER PRIMARY KEY, -- id编号 主键
    asn_name           TEXT NOT NULL,       -- ASN名称（不能为空）
    full_name          TEXT NOT NULL,       -- ASN全名（不能为空）
    ipv4_counts        INTEGER,             -- ipv4数量
    ipv6_counts        INTEGER,             -- ipv6数量
    peers_counts       INTEGER,             -- 节点数量
    ipv4_peers         INTEGER,             -- ipv4节点数量
    ipv6_peers         INTEGER,             -- ipv6节点数量
    prefixes_counts    INTEGER,             -- 前缀数量
    ipv4_prefixies     INTEGER,             -- ipv4前缀数量
    ipv6_prefixies     INTEGER,             -- ipv6前缀数量
    regional_registry  TEXT,                -- 地区登记
    traffic_bandwidth  TEXT,                -- 流量估算
    allocation_status  TEXT,                -- 分配状态 未分配 已分配
    traffic_ratio      TEXT,                -- 流量比率
    allocation_date    TEXT,                -- 分配日期
    website            TEXT,                -- 网站地址
    allocation_country TEXT,                -- 分配国家
    ipv4CIDR           TEXT,                --ipv4CIDR 数据
    last_cidr_update   DATETIME,            --cidr最新更新日期
    created_at         DATETIME,            -- 创建日期（默认为创建时的时间戳）
    updated_at         DATETIME,            -- 更新日期（默认为当前时间戳）
    deleted_at         DATETIME,            -- 删除日期
    enable             INTEGER DEFAULT 1    -- 是否开启（默认为1）
);

-- 为 asn_name 创建索引
CREATE INDEX idx_asn_name ON asn_info (asn_name);

-- 创建触发器以自动更新 updated_at 字段
-- CREATE TRIGGER update_asn_info_timestamp
--     AFTER UPDATE ON asn_info
-- BEGIN
--     UPDATE asn_info SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
-- END;


-- 下一张表
-- create table asn_info
-- (
--     id                 INTEGER
--         primary key,
--     asn_name           TEXT not null,
--     full_name          TEXT not null,
--     ipv4_counts        INTEGER,
--     ipv6_counts        INTEGER,
--     peers_counts       INTEGER,
--     ipv4_peers         INTEGER,
--     ipv6_peers         INTEGER,
--     prefixes_counts    INTEGER,
--     ipv4_prefixies     INTEGER,
--     ipv6_prefixies     INTEGER,
--     regional_registry  TEXT,
--     traffic_bandwidth  TEXT,
--     allocation_status  TEXT,
--     traffic_ratio      TEXT,
--     allocation_date    TEXT,
--     website            TEXT,
--     allocation_country TEXT,
--     created_at         DATETIME default CURRENT_TIMESTAMP,
--     updated_at         DATETIME default CURRENT_TIMESTAMP,
--     deleted_at         DATETIME,
--     enable             INTEGER  default 1,
--     ipv4CIDR           TEXT,
--     last_cidr_update   DATETIME
-- );
--
-- create index idx_asn_name
--     on asn_info (asn_name);


create table submit_scan
(
    scan_desc        TEXT,    --扫描描述
    scan_type        INTEGER, --1 表示扫描ASN 2表示扫描ASN列表 3表示单个IP 4表示多个ip
    asn_number       TEXT,    --ASN编号
    ipinfo_type      INTEGER, --ip信息类型ip或是cidr 1是ip 2是cidr
    ipinfo_list      TEXT,    --ip信息文本
    ip_batch_size    INTEGER, --cidr ip 批量大小
    enable_tls       INTEGER, -- 0表示不开启 1表示开启
    scan_ports       TEXT,    --扫描端口集合
    scan_rate        INTEGER, --扫描速率 默认10000
    ipcheck_thread   INTEGER, --ip检测线程数
    enable_speedtest INTEGER, --是否开启测速 0表示关闭 1表示开启
    scan_status      INTEGER, --扫描状态
    scan_result      TEXT     --扫描结果
);


create table schedule_task
(
    task_desc        TEXT, --定时描述
    asn_number  TEXT, --asn编号
    asn_desc    TEXT, -- ASN描述
    crontab_str TEXT, --定时表达式
    task_config TEXT, --任务配置信息 json配置
    enable      TEXT, --是否开启
    task_status TEXT--任务状态 0禁用1空闲2队列中3运行中

);

-- 调度历史表
create table schedule_task_hist(
    schedule_task_id INTEGER, --定时任务ID
    start_time TEXT, --任务开始时间
    end_time TEXT, --任务结束时间
    cost_time INTEGER, --耗时s
    hist_status TEXT, --调度历史状态 1进行中 2已完成 3已超时
    task_result TEXT --任务结果
);

-- {"ip":"206.237.112.12",
-- "port":443,
-- "data_center":"LAX",
-- "region":"North America",
-- "city":"Los Angeles",
-- "latency":"20 ms",
-- "tcp_duration":20000000,
-- "enable_tls":true,
-- "download_speed":"23434 kB/s"}
--反代ip库表
-- 入库需要去重
create table proxy_ips(
  asn_number TEXT, --asn 名称

  ip TEXT, --ip
  port INTEGER, --端口
  data_center TEXT, -- 数据中心
  region TEXT, --区域
  city TEXT, --城市
  latency TEXT, --延迟
  tcp_duration INTEGER, --TCP延迟
  enable_tls BOOLEAN, --是否开启tls
  download_speed TEXT --下载速度

);

-- 添加联合唯一索引
CREATE UNIQUE INDEX idx_ip_port ON proxy_ips (ip, port);
Create INDEX idx_ip ON  proxy_ips(ip,port);




-- 经过检测的IP代理池
-- 从proxy_ips 查询出来检测更新

create table alive_proxy_ips(
     asn_number TEXT, --asn 名称
     asn_desc TEXT, --asn描述
     ip TEXT, --ip
     port INTEGER, --端口
     enable_tls BOOLEAN, --是否开启tls
     geo_distance INTEGER, --物理距离
     data_center TEXT, -- 数据中心
     region TEXT, --区域
     city TEXT, --城市
     latency TEXT, --延迟
     tcp_duration INTEGER, --TCP延迟
     download_speed TEXT, --下载速度
     ttl INTEGER, --存活时间
     desc_str TEXT --ip描述
);

-- 添加联合唯一索引
CREATE UNIQUE INDEX idx_ip_port ON alive_proxy_ips (ip, port);
Create INDEX idx_ip ON  alive_proxy_ips(ip,port);


-- 查询得到统计数据 然后插入统计表
SELECT
    CURRENT_DATE AS date,
    COUNT(DISTINCT asn_number) AS asn_number_count,
    COUNT(DISTINCT ip) AS daily_ip_count,
    COUNT(DISTINCT port) AS distinct_port_count,
    COUNT(DISTINCT data_center) AS distinct_data_center_count,
    COUNT(DISTINCT city) AS distinct_city_count,
    SUM(CASE WHEN enable_tls = TRUE THEN 1 ELSE 0 END) AS tls_enabled_count,
    SUM(CASE WHEN enable_tls = FALSE THEN 1 ELSE 0 END) AS tls_disabled_count
FROM
    proxy_ips;

SELECT
            CURRENT_DATE AS date,
            COUNT(DISTINCT asn_number) AS asn_number_count,
            COUNT(DISTINCT ip) AS daily_ip_count,
            COUNT(DISTINCT port) AS distinct_port_count,
            COUNT(DISTINCT data_center) AS distinct_data_center_count,
            COUNT(DISTINCT city) AS distinct_city_count,
            SUM(CASE WHEN enable_tls = TRUE THEN 1 ELSE 0 END) AS tls_enabled_count,
            SUM(CASE WHEN enable_tls = FALSE THEN 1 ELSE 0 END) AS tls_disabled_count
FROM
    alive_proxy_ips;

-- 统计表 用来统计proxy_ips 和alive_proxy_ips
CREATE TABLE daily_proxy_stats (
   date DATE,                   -- 统计日期，作为主键
   asn_number_count INTEGER ,       -- 不同 ASN 数量
   daily_ip_count INTEGER ,         -- 每日不同 IP 数量
   distinct_port_count INTEGER ,    -- 不同端口数量
   distinct_data_center_count INTEGER , -- 不同数据中心数量
   distinct_city_count INTEGER ,    -- 不同城市数量
   tls_enabled_count INTEGER,      -- 启用 TLS 的 IP 数量
   tls_disabled_count INTEGER,      -- 未启用 TLS 的 IP 数量
   stats_from_alive BOOLEAN       --是否是alive_proxy_ip的数据
);


