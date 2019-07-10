<?php

require '/vendor/autoload.php';
$config = include("./config.php");

/**
 * @return Memcached
 */
function memcached_client_with(
    array $servers,
    string $id = null
) : Memcached {
    $memcache = $id
        ? new Memcached($id)
        : new Memcached();
    $memcache->setOption(
        Memcached::OPT_DISTRIBUTION,
        Memcached::DISTRIBUTION_CONSISTENT
    );
    $memcache->setOption(
        Memcached::OPT_CONNECT_TIMEOUT,
        20
    );
    if (count($memcache->getServerList()) == 0) {
        if ($id) {
            error_log('adding servers to memache client');
        }
        foreach ($servers as $server) {
            $memcache->addServer($server, 11211);
        }
    }
    return $memcache;
}

/**
 * @return array
 */
function get_keys($n = 100) : array
{
    $keys = [];
    for ($i = 0; $i < $n; $i++) {
        $keys[] = sprintf("%02d", $i);
    }
    return $keys;
}

/**
 * @return array
 */
function get_keys_in(
    Memcached $client,
    array $keys
) : array {
    $items = $client->getMulti($keys);
    $keys = $items
        ? array_keys($items)
        : [];
    sort($keys);
    return $keys;
}

function shuffle_assoc(array $array) : array {
    $keys = array_keys($array);
    shuffle($keys);
    $new = [];
    foreach ($keys as $key) {
        $new[$key] = $array[$key];
    }
    return $new;
}

function get_servers_for(
    Memcached $client,
    array $config,
    array $keys
) : array {
    $name_by_host = array_flip($config);
    $keys_by_server = array_fill_keys(array_keys($config), []);
    foreach ($keys as $key) {
        $server = $client->getServerByKey($key);
        $name = $name_by_host[$server['host']];
        $keys_by_server[$name][] = $key;
    }
    return $keys_by_server;
}

$config = shuffle_assoc($config);
$cluster = memcached_client_with(
    $config,
    'cluster'
);
$cluster->flush();

$keys = get_keys(10);
$cluster->setMulti(array_combine($keys, $keys));

$by_server = [];
foreach ($config as $name => $ip) {
    $by_server[$name] = get_keys_in(
        memcached_client_with([$ip]),
        $keys
    );
}

ksort($by_server);

$debug = get_servers_for(
    $cluster,
    $config,
    $keys
);

$smarty = new Smarty();
$smarty->setTemplateDir('/templates');
$smarty->assign('by_server', $by_server);
// $smarty->assign('debug', json_encode($config));
$smarty->display('index.tpl');