use std::collections::HashMap;
mod tokenize;

use fset::{FSet, Token};
use ir::IR;
use tokio::sync::Mutex;
use tower_lsp::jsonrpc::Result;
use tower_lsp::lsp_types::*;
use tower_lsp::{Client, LanguageServer, LspService, Server};

struct State {
    root: String,
    files: HashMap<String, Vec<Token>>,
    ir: IR,
}

struct Backend {
    client: Client,
    state: Mutex<State>,
}

fn get_path(url: &Url) -> String {
    return url.to_string().replacen("file://", "", 1);
}

#[tower_lsp::async_trait]
impl LanguageServer for Backend {
    async fn initialize(&self, p: InitializeParams) -> Result<InitializeResult> {
        let mut state = self.state.lock().await;
        state.root = get_path(&p.root_uri.unwrap());
        Ok(InitializeResult {
            server_info: None,
            capabilities: ServerCapabilities {
                text_document_sync: Some(TextDocumentSyncCapability::Kind(
                    TextDocumentSyncKind::FULL,
                )),
                workspace: Some(WorkspaceServerCapabilities {
                    workspace_folders: Some(WorkspaceFoldersServerCapabilities {
                        supported: Some(true),
                        change_notifications: Some(OneOf::Left(true)),
                    }),
                    file_operations: None,
                }),
                ..ServerCapabilities::default()
            },
        })
    }

    async fn initialized(&self, _: InitializedParams) {
        let mut state = self.state.lock().await;
        state.compile(&self.client).await;
        self.client
            .log_message(MessageType::INFO, "B# Language Server is ready")
            .await;
    }

    async fn did_open(&self, p: DidOpenTextDocumentParams) {
        let mut state = self.state.lock().await;
        state.update_file(&p.text_document.uri, p.text_document.text);
    }

    async fn did_change(&self, p: DidChangeTextDocumentParams) {
        let mut state = self.state.lock().await;
        state.update_file(&p.text_document.uri, p.content_changes[0].text.clone());
    }

    async fn did_save(&self, _: DidSaveTextDocumentParams) {
        let mut state = self.state.lock().await;
        state.compile(&self.client).await;
    }

    async fn shutdown(&self) -> Result<()> {
        Ok(())
    }
}

#[tokio::main]
async fn main() {
    let stdin = tokio::io::stdin();
    let stdout = tokio::io::stdout();

    let (service, socket) = LspService::new(|client| Backend {
        client,
        state: Mutex::new(State {
            root: "".to_string(),
            files: HashMap::new(),
            ir: IR::new(FSet::new()),
        }),
    });

    Server::new(stdin, stdout, socket).serve(service).await;
}
