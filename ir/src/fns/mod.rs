use super::*;

mod composite;
mod control;
mod fns;
mod ops;
mod typs;
mod vars;

impl IR {
    pub fn build_print(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = self.typecheck(pos, &args, &vec![TypeData::DEF(0)])?;
        Ok(IRNode::new(
            IRNodeData::Print(Box::new(args[0].clone())),
            range,
            pos,
        ))
    }

    pub fn build_string(&mut self, pos: Pos, val: &String) -> IRNode {
        let mut nodes = Vec::new();
        for char in val.as_bytes() {
            nodes.push(IRNode::new(IRNodeData::Char(*char), pos, pos));
        }
        IRNode::new(
            IRNodeData::NewArrayLiteral(Type::from(TypeData::DEF(0)), nodes),
            pos,
            pos,
        )
    }
}
