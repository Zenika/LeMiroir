"use strict";

let express = require('express');
let app = express();
let http = require('http').Server(app);

app.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Headers', 'X-Requested-With')
    next()
})

app.use(express.static('./public/'));
app.use(express.static('./components/'));
app.use(express.static('./node_modules/'));



app.get('/', function (req, res) {
    res.sendfile("public/index.html");
});

/**
 * Server itself
 * @type {http.Server}
 */
let server = app.listen(8080, function () {
    //print few information about the server
    let host = server.address().address;
    let port = server.address().port;
    console.log("Server running and listening @ " + host + ":" + port);
});

/** list of components to be loaded */
let componentsList = [
    {
        "componentName": "github",
        "eltName": "github-item",
        "files": "github.html",
        "propValues": {
            "users": [
                'Gillespie59',
                'GwennaelBuchet',
                'T3kstiil3',
                'RemiEven',
                'looztra',
                'a-cordier',
                'wadendo',
                'NathanDM',
                'Antoinephi',
                'cluster',
                'yyekhlef',
                'gdrouet',
                'Kize',
                'kratisto',
                'Sehsyha',
                'P0ppoff'
            ],
            "client_id": "27b62f039b44ddc08fdf",
            "client_secret": "7b14f465112e87267a72c02d4c3fc58925412dbd"
        }
    }
];

/**
 * Get a list of JSON for all registered components
 * @path /componentsList
 * @HTTPMethod GET
 * @returns {string}
 */
app.get("/componentsList", function (req, res) {
    res.send(componentsList);
});
