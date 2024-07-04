CREATE TABLE asn_info (
                          id INTEGER PRIMARY KEY,                  -- id编号 主键
                          asn_name TEXT NOT NULL,                  -- ASN名称（不能为空）
                          full_name TEXT NOT NULL,                 -- ASN全名（不能为空）
                          ipv4_counts INTEGER,                     -- ipv4数量
                          ipv6_counts INTEGER,                     -- ipv6数量
                          peers_counts INTEGER,                    -- 节点数量
                          ipv4_peers INTEGER,                      -- ipv4节点数量
                          ipv6_peers INTEGER,                      -- ipv6节点数量
                          prefixes_counts INTEGER,                 -- 前缀数量
                          ipv4_prefixies INTEGER,                  -- ipv4前缀数量
                          ipv6_prefixies INTEGER,                  -- ipv6前缀数量
                          regional_registry TEXT,                  -- 地区登记
                          traffic_bandwidth TEXT,                  -- 流量估算
                          allocation_status TEXT,                  -- 分配状态 未分配 已分配
                          traffic_ratio TEXT,                      -- 流量比率
                          allocation_date TEXT,                    -- 分配日期
                          website TEXT,                            -- 网站地址
                          allocation_country TEXT,                 -- 分配国家
                          ipv4CIDR TEXT,                           --ipv4CIDR 数据
                          last_cidr_update DATETIME,               --cidr最新更新日期
                          created_at DATETIME,  -- 创建日期（默认为创建时的时间戳）
                          updated_at DATETIME,  -- 更新日期（默认为当前时间戳）
                          deleted_at DATETIME,                     -- 删除日期
                          enable INTEGER DEFAULT 1                 -- 是否开启（默认为1）
);

-- 为 asn_name 创建索引
CREATE INDEX idx_asn_name ON asn_info(asn_name);

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


create table submit_scan(
    scan_desc TEXT, --扫描描述
    scan_type INTEGER, --1 表示扫描ASN 2表示扫描ASN列表 3表示单个IP 4表示多个ip
    asn_number TEXT, --ASN编号
    ipinfo_type INTEGER, --ip信息类型ip或是cidr 1是ip 2是cidr
    ipinfo_list TEXT, --ip信息文本
    enable_tls INTEGER, -- 0表示不开启 1表示开启
    scan_ports TEXT, --扫描端口集合
    scan_rate INTEGER, --扫描速率 默认10000
    ipcheck_thread INTEGER, --ip检测线程数
    enable_speedtest INTEGER, --是否开启测速 0表示关闭 1表示开启
    scan_status INTEGER, --扫描状态
    scan_result TEXT --扫描结果
);