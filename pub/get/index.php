<?php

require_once '../lib.php';

$servers = parse_server_list(
    getenv('MEMCACHE_SERVERS')
);

$memcache = create_client(
    array_values($servers)
);

$keys = explode(',', array_get($_GET, 'keys', ''));

$servers_by_key = get_servers_by_key(
    $memcache,
    $keys
);

$active_servers = get_active_servers(
    $memcache,
    $servers
);

$res = $memcache->getMulti(
    $keys
);

header('Content-Type: application/json;charset=utf8');
printf("%s\n", json_encode([
    'keys' => $keys,
    'active_servers' => $active_servers,
    'servers_by_key' => $servers_by_key,
    'result' => $res,
    'result_code' => $memcache->getResultCode(),
    'result_message' => $memcache->getResultMessage(),
]));