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
                let o = match op {
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
                            "INT".to_string(),
                            vec![Node::Tag(
                                "MATH".to_string(),
                                vec![
                                    l,
                                    Node::Ident("*".to_string()),
                                    Node::Tag(
                                        "MATH".to_string(),
                                        vec![Node::Int(2), Node::Ident("^".to_string()), r],
                                    ),
                                ],
                            )],
                        ))
                    }
                };
                if *op == MathOperator::DIVIDE && node.typ(&self.ir).data == TypeData::INT {
                    return Ok(Node::Tag(
                        "INT".to_string(),
                        vec![Node::Tag(
                            "MATH".to_string(),
                            vec![l, Node::Ident(o.to_string()), r],
                        )],
                    ));
                }
                Ok(Node::Tag(
                    "MATH".to_string(),
                    vec![l, Node::Ident(o.to_string()), r],
                ))
            }
            IRNodeData::Boolean(l, op, r) => {
                if op == &BooleanOperator::NOT {
                    let l = self.build_node(l)?;
                    return Ok(Node::Tag(
                        "MATH".to_string(),
                        vec![Node::Int(1), Node::Ident("-".to_string()), l],
                    ));
                }
                let l = self.build_node(l)?;
                let r = self.build_node(r.as_ref().unwrap())?;
                let o = match op {
                    BooleanOperator::AND => "*",
                    BooleanOperator::OR => "+",
                    _ => unreachable!(),
                };
                Ok(Node::Tag(
                    "MATH".to_string(),
                    vec![l, Node::Ident(o.to_string()), r],
                ))
            }
            IRNodeData::Cast(val, _) => self.build_node(val),
            IRNodeData::Variable(var, _) => {
                Ok(Node::Tag("VAR".to_string(), vec![self.fmt_var(*var)?]))
            }
            IRNodeData::Len(val) => {
                let val = self.build_node(val)?;
                Ok(Node::Tag(
                    "LENGTH".to_string(),
                    vec![Node::Tag("HEAPGET".to_string(), vec![val])],
                ))
            }
            IRNodeData::GetArr { arr, ind } => {
                let a = self.build_node(arr)?;
                let ind = self.build_node(ind)?;
                if arr.typ(&self.ir).data == TypeData::DEF(0) {
                    return Ok(Node::Tag(
                        "GETBYTE".to_string(),
                        vec![Node::Tag(
                            "INDEX".to_string(),
                            vec![Node::Tag("HEAPGET".to_string(), vec![a]), ind],
                        )],
                    ));
                }
                Ok(Node::Tag(
                    "INDEX".to_string(),
                    vec![Node::Tag("HEAPGET".to_string(), vec![a]), ind],
                ))
            }
            IRNodeData::Return(val) => {
                if let Some(v) = val {
                    return Ok(Node::Tag("RETURN".to_string(), vec![self.build_node(v)?]));
                } else {
                    return Ok(Node::Tag(
                        "RETURN".to_string(),
                        vec![Node::String("".to_string())],
                    ));
                }
            }
            IRNodeData::NewStruct(t, ops) => {
                let names = self.structfields(t);

                let mut res = vec![Node::Int(0); names.len()];
                for op in ops {
                    if let IRNodeData::StructOp { field, val } = &op.data {
                        let v = self.build_node(&val)?;
                        let ind = names.iter().position(|x| x == field).unwrap();
                        res[ind] = v;
                    }
                }
                Ok(Node::Tag(
                    "HEAPADD".to_string(),
                    vec![Node::ArrayLiteral(res)],
                ))
            }
            IRNodeData::NewArrayLiteral(t, vals) => {
                if t.data == TypeData::DEF(0) {
                    let mut chars = Vec::new();
                    for v in vals {
                        if let IRNodeData::Char(c) = v.data {
                            chars.push(c);
                        } else {
                            return Err(BStarError::UnknownNode(v.clone()));
                        }
                    }
                    Ok(Node::Tag(
                        "HEAPADD".to_string(),
                        vec![Node::String(String::from_utf8(chars).unwrap())],
                    ))
                } else {
                    let mut res = Vec::new();
                    for val in vals {
                        res.push(self.build_node(val)?);
                    }
                    Ok(Node::Tag(
                        "HEAPADD".to_string(),
                        vec![Node::ArrayLiteral(res)],
                    ))
                }
            }
            IRNodeData::NewArray(t, _) => {
                if t.data == TypeData::DEF(0) {
                    return Ok(Node::Tag(
                        "HEAPADD".to_string(),
                        vec![Node::String("".to_string())],
                    ));
                }
                Ok(Node::Tag(
                    "HEAPADD".to_string(),
                    vec![Node::ArrayLiteral(vec![])],
                ))
            }
            IRNodeData::NewBox(val) => {
                let v = self.build_node(val)?;
                Ok(Node::Tag(
                    "ARRAY".to_string(),
                    vec![v, Node::String(self.hashtyp(&val.typ(&self.ir)))],
                ))
            }
            IRNodeData::Unbox { bx, .. } => {
                let v = self.build_node(bx)?;
                Ok(Node::Tag("INDEX".to_string(), vec![v, Node::Int(0)]))
            }
            IRNodeData::FnCall { func, args, .. } => {
                let mut res = Vec::new();
                for arg in args {
                    res.push(self.build_node(arg)?);
                }
                Ok(Node::Tag(self.ir.funcs[*func].name.clone(), res))
            }
            IRNodeData::If {
                cond, body, els, ..
            } => {
                let cond = self.build_node(cond)?;
                let body = self.build_node(body)?;
                let els = self.build_node(els)?;
                Ok(Node::Tag("IF".to_string(), vec![cond, body, els]))
            }
            IRNodeData::GetStruct { strct, field } => {
                let fields = self.structfields(&strct.typ(&self.ir));
                let ind = fields.iter().position(|x| x == field).unwrap() as i64;
                let strct = self.build_node(strct)?;
                Ok(Node::Tag(
                    "INDEX".to_string(),
                    vec![
                        Node::Tag("HEAPGET".to_string(), vec![strct]),
                        Node::Int(ind),
                    ],
                ))
            }
            IRNodeData::Peek { bx, typ } => {
                let bx = self.build_node(bx)?;
                Ok(Node::Tag(
                    "COMPARE".to_string(),
                    vec![
                        Node::Tag("INDEX".to_string(), vec![bx, Node::Int(1)]),
                        Node::Ident("=".to_string()),
                        Node::String(self.hashtyp(typ)),
                    ],
                ))
            }
            IRNodeData::SetArr { arr, ind, val } => {
                let a = self.build_node(arr)?;
                let ind = self.build_node(ind)?;
                let mut v = self.build_node(val)?;
                if val.typ(&self.ir).data == TypeData::CHAR {
                    v = Node::Tag(
                        "INDEX".to_string(),
                        vec![
                            Node::Tag("VAR".to_string(), vec![Node::Ident("c".to_string())]),
                            v,
                        ],
                    );
                }
                let mut f = "SETINDEX".to_string();
                if arr.typ(&self.ir).data == TypeData::DEF(0) {
                    f = "STRSETINDEX".to_string();
                }
                Ok(Node::Tag(
                    "HEAPSET".to_string(),
                    vec![
                        a.clone(),
                        Node::Tag(f, vec![Node::Tag("HEAPGET".to_string(), vec![a]), ind, v]),
                    ],
                ))
            }
            IRNodeData::Append { arr, val } => {
                let arr = self.build_node(arr)?;
                let mut v = self.build_node(val)?;
                if val.typ(&self.ir).data == TypeData::CHAR {
                    v = Node::Tag(
                        "INDEX".to_string(),
                        vec![
                            Node::Tag("VAR".to_string(), vec![Node::Ident("c".to_string())]),
                            v,
                        ],
                    );
                } else {
                    v = Node::ArrayLiteral(vec![v])
                }
                Ok(Node::Tag(
                    "HEAPSET".to_string(),
                    vec![
                        arr.clone(),
                        Node::Tag(
                            "CONCAT".to_string(),
                            vec![Node::Tag("HEAPGET".to_string(), vec![arr]), v],
                        ),
                    ],
                ))
            }
            IRNodeData::SetStruct { strct, vals } => {
                let names = self.structfields(&strct.typ(&self.ir));
                let strct = self.build_node(strct)?;
                let mut res = Vec::new();
                for v in vals {
                    if let IRNodeData::StructOp { field, val } = &v.data {
                        let v = self.build_node(&val)?;
                        let ind = names.iter().position(|x| x == field).unwrap();

                        // Set
                        res.push(Node::Tag(
                            "HEAPSET".to_string(),
                            vec![
                                strct.clone(),
                                Node::Tag(
                                    "SETINDEX".to_string(),
                                    vec![
                                        Node::Tag("HEAPGET".to_string(), vec![strct.clone()]),
                                        Node::Int(ind as i64),
                                        v,
                                    ],
                                ),
                            ],
                        ));
                    }
                }
                Ok(Node::Tag("BLOCK".to_string(), res))
            }
            IRNodeData::NewTuple(_, vals) => {
                let mut res = Vec::new();
                for v in vals {
                    res.push(self.build_node(v)?);
                }
                Ok(Node::Tag(
                    "HEAPADD".to_string(),
                    vec![Node::Tag("ARRAY".to_string(), res)],
                ))
            }
            IRNodeData::GetTuple { tup, ind } => {
                let tup = self.build_node(tup)?;
                Ok(Node::Tag(
                    "INDEX".to_string(),
                    vec![
                        Node::Tag("HEAPGET".to_string(), vec![tup]),
                        Node::Int(*ind as i64),
                    ],
                ))
            }
            IRNodeData::Char(v) => Ok(Node::Int(*v as i64)),
            IRNodeData::Print(v) => {
                let val = self.build_node(v)?;
                Ok(Node::Tag(
                    "CONCAT".to_string(),
                    vec![
                        Node::Tag("HEAPGET".to_string(), vec![val]),
                        Node::String("\n".to_string()),
                    ],
                ))
            }
            IRNodeData::NewEnum(val, _) => {
                let tname = self.hashtyp(&val.typ(&self.ir));
                let val = self.build_node(val)?;
                Ok(Node::Tag(
                    "ARRAY".to_string(),
                    vec![val, Node::String(tname)],
                ))
            }
            IRNodeData::GetEnum { enm, .. } => {
                let v = self.build_node(enm)?;
                Ok(Node::Tag("INDEX".to_string(), vec![v, Node::Int(0)]))
            }
            IRNodeData::TypeMatch { val, body } => {
                let cond = self.build_node(val)?;
                let typval = Node::Tag("INDEX".to_string(), vec![cond.clone(), Node::Int(1)]);
                let varval = Node::Tag("INDEX".to_string(), vec![cond.clone(), Node::Int(0)]);
                let mut res = Node::Tag("BLOCK".to_string(), vec![]);
                for c in body {
                    if let IRNodeData::TypeCase { var, typ, body } = &c.data {
                        let t = self.hashtyp(typ);
                        let body = self.build_node(body)?;
                        let varname = self.fmt_var(*var)?;
                        res = Node::Tag(
                            "IF".to_string(),
                            vec![
                                Node::Tag(
                                    "COMPARE".to_string(),
                                    vec![
                                        Node::String(t),
                                        Node::Ident("=".to_string()),
                                        typval.clone(),
                                    ],
                                ),
                                Node::Tag(
                                    "BLOCK".to_string(),
                                    vec![
                                        Node::Tag(
                                            "DEFINE".to_string(),
                                            vec![varname, varval.clone()],
                                        ),
                                        body,
                                    ],
                                ),
                                res,
                            ],
                        )
                    }
                }
                Ok(res)
            }
            IRNodeData::Match { val, body } => {
                let mut cond = self.build_node(val)?;
                let mut res = Node::Tag("BLOCK".to_string(), vec![]);
                if val.typ(&self.ir).data == TypeData::DEF(0) {
                    // String match
                    cond = Node::Tag("HEAPGET".to_string(), vec![cond]);
                }
                for c in body {
                    if let IRNodeData::Case { val, body } = &c.data {
                        let mut valcond = self.build_node(val)?;
                        if val.typ(&self.ir).data == TypeData::DEF(0) {
                            // TODO: Instead of adding and instantly getting from the heap, just put a string literal
                            // String match
                            valcond = Node::Tag("HEAPGET".to_string(), vec![valcond]);
                        }
                        let body = self.build_node(body)?;
                        res = Node::Tag(
                            "IF".to_string(),
                            vec![
                                Node::Tag(
                                    "COMPARE".to_string(),
                                    vec![valcond, Node::Ident("=".to_string()), cond.clone()],
                                ),
                                body,
                                res,
                            ],
                        )
                    }
                }
                Ok(res)
            }
            _ => Err(BStarError::UnknownNode(node.clone())),
        }
    }
}
