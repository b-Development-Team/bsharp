use std::collections::HashMap;
mod compile;
mod completion;
mod hover;

use fset::{FSet, Token};
use ir::{Function, TypeData, IR};
use tokens::Pos;
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
                completion_provider: Some(CompletionOptions {
                    resolve_provider: Some(false),
                    trigger_characters: Some(vec![
                        "!".to_string(),
                        "@".to_string(),
                        "$".to_string(),
                        "[".to_string(),
                        ".".to_string(), // TODO: struct field completion
                    ]),
                    work_done_progress_options: WorkDoneProgressOptions::default(),
                    all_commit_characters: None,
                    completion_item: None,
                }),
                hover_provider: Some(HoverProviderCapability::Simple(true)),
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

    async fn completion(&self, p: CompletionParams) -> Result<Option<CompletionResponse>> {
        self.complete(p).await
    }

    async fn hover(&self, p: HoverParams) -> Result<Option<Hover>> {
        self.handle_hover(p).await
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

// Util funcs
fn get_tok(pos: Position, f: &Vec<Token>) -> Option<&Token> {
    for tok in f {
        if tok.pos.start_line <= pos.line as usize
            && tok.pos.end_line >= pos.line as usize
            && tok.pos.start_col <= pos.character as usize
            && tok.pos.end_col >= pos.character as usize
        {
            return Some(tok);
        }
    }
    None
}

fn pos_range(pos: Pos) -> Range {
    Range {
        start: Position {
            line: pos.start_line as u32,
            character: pos.start_col as u32,
        },
        end: Position {
            line: pos.end_line as u32,
            character: pos.end_col as u32,
        },
    }
}

pub fn fn_string(f: &Function, ir: &IR) -> String {
    let mut desc = format!("[{}", f.name);
    for p in f.params.iter() {
        //detail.push_str(&format!(" {}", ir.variables[*p].typ.data.fmt(&ir)));
        desc.push_str(&format!(
            " [PARAM {} {}]",
            ir.variables[*p].name,
            ir.variables[*p].typ.data.fmt(&ir)
        ));
    }
    if (f.ret_typ.data != TypeData::VOID) && (f.ret_typ.data != TypeData::INVALID) {
        desc.push_str(" [RETURNS ");
        desc.push_str(&f.ret_typ.data.fmt(&ir));
        desc.push_str("]");
    }
    desc.push_str("]");
    desc
}
