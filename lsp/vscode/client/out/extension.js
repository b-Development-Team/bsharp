"use strict";
/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const cp = require("child_process");
const node_1 = require("vscode-languageclient/node");
let client;
async function activate(context) {
    // If the extension is launched in debug mode then the debug server options are used
    // Otherwise the run options are used
    // Get GOPATH
    let gopath = cp.execSync('go env GOPATH').toString();
    console.log(gopath);
    const serverOptions = {
        command: gopath.trim() + "/bin/bsharp-lsp",
        args: [],
        transport: node_1.TransportKind.stdio, // also tried every other option
    };
    // Options to control the language client
    const clientOptions = {
        // Register the server for plain text documents
        documentSelector: [{ scheme: "file", "language": "bsharp" }]
    };
    // Create the language client and start the client.
    client = new node_1.LanguageClient('bsharp', 'B# Language Server', serverOptions, clientOptions);
    // Start the client. This will also launch the server
    client.start();
}
exports.activate = activate;
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map