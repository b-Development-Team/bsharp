/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';

import {
	LanguageClient,
	LanguageClientOptions,
	ServerOptions,
	TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
	// If the extension is launched in debug mode then the debug server options are used
	// Otherwise the run options are used
	const serverOptions: ServerOptions = {
		command: context.asAbsolutePath("../lsp"),
    args: [],
    transport: TransportKind.stdio, // also tried every other option
	};

	// Options to control the language client
	const clientOptions: LanguageClientOptions = {
		// Register the server for plain text documents
		documentSelector: [{ scheme: "file", pattern: "*.bsp"}, {scheme: "file", "language": "plaintext"}]
	};

	// Create the language client and start the client.
	client = new LanguageClient(
		'language ServerExample',
		'Language Server Example',
		serverOptions,
		clientOptions
	);

	client.onReady().then(() => {
    console.log("Example client ready"); // will NOT be logged
	});


	// Start the client. This will also launch the server
	client.start();
}

export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	}
	return client.stop();
}
