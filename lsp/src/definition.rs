use tokens::TokenData;

use super::*;

impl Backend {
    pub async fn definition(
        &self,
        p: GotoDefinitionParams,
    ) -> Result<Option<GotoDefinitionResponse>> {
        // Get token
        let state = self.state.lock().await;
        let path = get_path(&p.text_document_position_params.text_document.uri);
        let f = state.files.get(&path).unwrap();
        let tok = get_tok(p.text_document_position_params.position, f);
        if tok.is_none() {
            return Ok(None);
        }
        let tok = tok.unwrap();

        // Get definition
        // TODO: Variables, fields
        match &tok.data {
            TokenData::FUNCTION(f) => {
                let f = state.ir.funcs.iter().find(|x| x.name == *f);
                if f.is_none() {
                    return Ok(None);
                }
                let f = f.unwrap();
                let res_uri = Url::parse(
                    format!(
                        "file://{}/{}",
                        state.root,
                        state.ir.fset.files[f.definition.file].name.clone()
                    )
                    .as_str(),
                )
                .unwrap();
                let def = GotoDefinitionResponse::Scalar(Location {
                    uri: res_uri,
                    range: pos_range(f.definition),
                });
                Ok(Some(def))
            }
            TokenData::TYPE(t) => {
                let val = state.ir.typemap.get(t);
                if val.is_none() {
                    return Ok(None);
                }
                let valpos = state.ir.types[*val.unwrap()].pos;
                let res_uri = Url::parse(
                    format!(
                        "file://{}/{}",
                        state.root,
                        state.ir.fset.files[valpos.file].name.clone()
                    )
                    .as_str(),
                )
                .unwrap();
                let def = GotoDefinitionResponse::Scalar(Location {
                    uri: res_uri,
                    range: pos_range(valpos),
                });
                Ok(Some(def))
            }
            _ => Ok(None),
        }
    }
}
