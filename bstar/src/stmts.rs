use super::*;

impl BStar {
    pub fn build_node(&mut self, node: &IRNode) -> Result<Node, BStarError> {
        match &node.data {
            IRNodeData::Block { body, .. } => {
                let mut res = Vec::new();
                for i in 0..body.len() {
                    res.push(self.build_node(&body[i])?);
                }
                Ok(Node::Tag("BLOCK".to_string(), res))
            }
            IRNodeData::Define { var, val, .. } => {
                let name = self.fmt_var(*var)?;
                let val = self.build_node(val)?;
                Ok(Node::Tag("DEFINE".to_string(), vec![name, val]))
            }
            IRNodeData::Int(v) => Ok(Node::Int(*v)),
            IRNodeData::Float(v) => Ok(Node::Float(*v)),
            IRNodeData::While { cond, body } => {
                let cond = self.build_node(cond)?;
                let body = self.build_node(body)?;
                Ok(Node::Tag("WHILE".to_string(), vec![cond, body]))
            }
            IRNodeData::Comparison(l, op, r) => {
                let l = self.build_node(l)?;
                let r = self.build_node(r)?;
                let op = match op {
                    ComparisonOperator::EQUAL => "=",
                    ComparisonOperator::NOTEQUAL => "!=",
                    ComparisonOperator::LESS => "<",
                    ComparisonOperator::GREATER => ">",
                    ComparisonOperator::LESSEQUAL => "<=",
                    ComparisonOperator::GREATEREQUAL => ">=",
                };
                Ok(Node::Tag(
                    "COMPARE".to_string(),
                    vec![l, Node::Ident(op.to_string()), r],
                ))
            }
            IRNodeData::Math(l, op, r) => {
                let l = self.build_node(l)?;
                let r = self.build_node(r)?;
                let op = match op {
                    MathOperator::ADD => "+",
                    MathOperator::SUBTRACT => "-",
                    MathOperator::MULTIPLY => "*",
                    MathOperator::DIVIDE => "/",
                    MathOperator::MODULO => "%",
                    MathOperator::POWER => "^",
                    MathOperator::XOR => return Ok(Node::Tag("XOR".to_string(), vec![l, r])),
                    MathOperator::BOR => return Ok(Node::Tag("BOR".to_string(), vec![l, r])),
                    MathOperator::SHIFT => {
                        return Ok(Node::Tag(
                            "MATH".to_string(),
                            vec![
                                l,
                                Node::Ident("*".to_string()),
                                Node::Tag(
                                    "MATH".to_string(),
                                    vec![Node::Int(2), Node::Ident("^".to_string()), r],
                                ),
                            ],
                        ))
                    }
                };
                Ok(Node::Tag(
                    "MATH".to_string(),
                    vec![l, Node::Ident(op.to_string()), r],
                ))
            }
            IRNodeData::Cast(val, _) => self.build_node(val),
            IRNodeData::Variable(var, _) => {
                Ok(Node::Tag("VAR".to_string(), vec![self.fmt_var(*var)?]))
            }
            IRNodeData::Len(val) => {
                let val = self.build_node(val)?;
                Ok(Node::Tag("LEN".to_string(), vec![val]))
            }
            IRNodeData::GetArr { arr, ind } => {
                let arr = self.build_node(arr)?;
                let ind = self.build_node(ind)?;
                Ok(Node::Tag("INDEX".to_string(), vec![arr, ind]))
            }
            IRNodeData::Return(val) => {
                if let Some(v) = val {
                    return Ok(Node::Tag("RETURN".to_string(), vec![self.build_node(v)?]));
                } else {
                    return Ok(Node::Tag("RETURN".to_string(), vec![]));
                }
            }
            _ => Err(BStarError::UnknownNode(node.clone())),
        }
    }
}
