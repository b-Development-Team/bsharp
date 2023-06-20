use ir::*;
mod node;
mod stmts;
mod util;
use node::*;

pub struct BStar {
    ir: IR,
}

impl BStar {
    pub fn new(ir: IR) -> BStar {
        BStar { ir }
    }

    pub fn fmt_var(&mut self, ind: usize) -> Result<Node, BStarError> {
        let v = &self.ir.variables[ind];
        Ok(Node::Ident(format!(
            "{}{:x}",
            v.name.clone().split_at(1).1,
            ind,
        )))
    }

    pub fn build_fn(&mut self, f: usize) -> Result<Node, BStarError> {
        let body = self.build_node(&self.ir.funcs[f].body.clone())?;
        let mut params = Vec::new();
        for i in 0..self.ir.funcs[f].params.len() {
            params.push(self.fmt_var(self.ir.funcs[f].params[i])?);
        }
        Ok(Node::Tag(
            "FUNC".to_string(),
            vec![
                Node::Ident(self.ir.funcs[f].name.clone()),
                Node::Tag("ARRAY".to_string(), params),
                body,
            ],
        ))
    }

    pub fn build(&mut self) -> Result<Vec<Node>, BStarError> {
        let mut res = Vec::new();
        res.push(Node::Tag(
            "IMPORT".to_string(),
            vec![Node::Ident("bsharplib".to_string())],
        ));
        for i in 0..self.ir.funcs.len() {
            let v = self.build_fn(i)?;
            res.push(v);
        }
        res.push(Node::Tag("@MAIN".to_string(), vec![]));
        Ok(res)
    }
}
