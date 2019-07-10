<?php

/**
 * 
 */
function parse_server_list(string $str): array
{
    $data = json_decode($str, true);
    if (json_last_error() !== JSON_ERROR_NONE || !is_array($data)) {
        throw new Exception("json error");
    }
    return $data;
}

/**
 * 
 */
function create_client(array $servers): Memcached
{
    $memcache = new Memcached('global');
    if (count($memcache->getServerList()) > 0) {
        return $memcache;
    }

    $memcache->setOption(
        Memcached::OPT_DISTRIBUTION,
        Memcached::DISTRIBUTION_CONSISTENT
    );
    $memcache->setOption(
        Memcached::OPT_REMOVE_FAILED_SERVERS,
        true
    );
    $memcache->setOption(
        Memcached::OPT_SERVER_FAILURE_LIMIT,
        2
    );
    $memcache->setOption(
        Memcached::OPT_CONNECT_TIMEOUT,
        10
    );
    foreach ($servers as $server) {
        $memcache->addServer($server, 11211);
    }
    return $memcache;
}

/**
 * 
 */
function get_servers_by_key(
    Memcached $memcache,
    array $keys
): array {
    $by_key = [];
    foreach ($keys as $key) {
        $server = $memcache->getServerByKey($key);
        $by_key[$key] = $server['host'];
    }
    return $by_key;
}

/**
 * 
 */
function get_active_servers(
    Memcached $memcache,
    array $servers
): array {
    $active = [];
    $index = array_flip($servers);
    foreach ($memcache->getServerList() as $server) {
        $host = $server['host'];
        $name = $index[$host];
        $active[$name] = $host;
    }
    return $active;
}

/**
 * 
 */
function array_get(array $array, $key, $default = null)
{
    if (!array_key_exists($key, $array)) {
        return $default;
    }

    return $array[$key];
}