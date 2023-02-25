use super::*;

impl IR {
    pub fn build_var(&mut self, pos: Pos, val: String) -> Result<IRNode, IRError> {
        for scope in self.stack.iter().rev() {
            let var = match self.scopes[*scope].vars.get(&val) {
                Some(t) => *t,
                None => {
                    continue;
                }
            };
            return Ok(IRNode::new(
                IRNodeData::Variable(var, self.variables[var].typ.clone()),
                pos,
                pos,
            ));
        }

        Err(IRError::UnknownVariable { pos, name: val })
    }
}
