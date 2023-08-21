/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import cp = require('child_process');

import {
	LanguageClient,
	LanguageClientOptions,
	ServerOptions,
	TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export async function activate(context: ExtensionContext) {
	// If the extension is launched in debug mode then the debug server options are used
	// Otherwise the run options are used

	// Get GOPATH
	/*let gopath = cp.execSync('go env GOPATH').toString();
	console.log(gopath);
	const serverOptions: ServerOptions = {
		command:  gopath.trim() + "/bin/bsharp-lsp",
		args: [],
		transport: TransportKind.stdio, // also tried every other option
	};

	// Options to control the language client
	const clientOptions: LanguageClientOptions = {
		// Register the server for plain text documents
		documentSelector: [{scheme: "file", "language": "bsharp"}]
	};

	// Create the language client and start the client.
	client = new LanguageClient(
		'bsharp',
		'B# Language Server',
		serverOptions,
		clientOptions
	);


	// Start the client. This will also launch the server
	client.start();*/
}

export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	}
	return client.stop();
}
