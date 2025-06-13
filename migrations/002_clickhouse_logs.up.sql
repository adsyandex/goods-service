CREATE TABLE IF NOT EXISTS goods_logs (
    id UInt64,
    project_id UInt64,
    name String,
    description String,
    priority Int32,
    removed UInt8,
    event_time DateTime,
    action String
) ENGINE = MergeTree()
ORDER BY (event_time, id);