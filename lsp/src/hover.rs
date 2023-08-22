use tokens::TokenData;

use super::*;

impl Backend {
    pub async fn handle_hover(&self, p: HoverParams) -> Result<Option<Hover>> {
        // Get file
        let state = self.state.lock().await;
        let path = get_path(&p.text_document_position_params.text_document.uri);
        if !state.files.contains_key(&path) {
            return Ok(None);
        }
        let f = state.files.get(&path).unwrap();
        let pos = p.text_document_position_params.position;

        // Get token
        let tok = get_tok(pos, f);
        if tok.is_none() {
            return Ok(None);
        }
        let tok = tok.unwrap();

        // Get hover
        // TODO: Fields
        match &tok.data {
            // Type
            TokenData::TYPE(t) => {
                let val = state.ir.typemap.get(t);
                if let Some(ind) = val {
                    let typ = &state.ir.types[*ind];
                    let hover = Hover {
                        contents: HoverContents::Markup(MarkupContent {
                            kind: MarkupKind::Markdown,
                            value: format!("```bsharp\n{}\n```", typ.typ.data.fmt(&state.ir)),
                        }),
                        range: Some(pos_range(tok.pos)),
                    };
                    return Ok(Some(hover));
                }
                Ok(None)
            }
            TokenData::FUNCTION(f) => {
                let f = state.ir.funcs.iter().find(|x| x.name == *f);
                if f.is_none() {
                    return Ok(None);
                }
                let f = f.unwrap();
                let hover = Hover {
                    contents: HoverContents::Markup(MarkupContent {
                        kind: MarkupKind::Markdown,
                        value: format!("```bsharp\n{}\n```", fn_string(&f, &state.ir)),
                    }),
                    range: Some(pos_range(tok.pos)),
                };
                Ok(Some(hover))
            }
            TokenData::VARIABLE(v) => {
                // Search for var
                let mut var = state.ir.variables.iter().find(|x| {
                    x.name == *v
                        && state.ir.scopes[x.scope].pos.start_line <= tok.pos.start_line
                        && state.ir.scopes[x.scope].pos.end_line >= tok.pos.end_line
                        && path
                            .ends_with(&state.ir.fset.files[state.ir.scopes[x.scope].pos.file].name)
                });
                if var.is_none() {
                    var = state.ir.variables.iter().find(|x| x.name == *v);
                    if var.is_none() {
                        return Ok(None);
                    }
                }
                let var = var.unwrap();

                let hover = Hover {
                    contents: HoverContents::Markup(MarkupContent {
                        kind: MarkupKind::Markdown,
                        value: format!("```bsharp\n{}\n```", var.typ.data.fmt(&state.ir)),
                    }),
                    range: Some(pos_range(tok.pos)),
                };
                Ok(Some(hover))
            }
            _ => Ok(None),
        }
    }
}
