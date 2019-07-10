<html>
    <head>
        <title>Keys on Servers</title>
        <link href="/index.css" rel="stylesheet" type="text/css">
        <link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
    </head>
    <body>
        <div class="node-list">
            {foreach $by_server as $name => $keys}
            <div class="node">
                <h2>{$name}</h2>
                <ul class="key-list">
                    {foreach $keys as $key}
                    <li class="key">{$key}</li>
                    {/foreach}
                </ul>
            </div>
            {/foreach}
        </div>
    </body>
</html>