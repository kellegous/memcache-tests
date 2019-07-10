<?php

$memcache = new Memcached();
$memcache->setOption(
    Memcached::OPT_DISTRIBUTION,
    Memcached::DISTRIBUTION_CONSISTENT
);
$memcache->setOption(
    Memcached::OPT_CONNECT_TIMEOUT,
    20
);
$memcache->addServer('localhost', 11211);

$memcache->set('00', '00');