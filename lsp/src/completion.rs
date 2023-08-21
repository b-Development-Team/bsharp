use ir::TypeData;

use super::*;

impl Backend {
    pub async fn complete(&self, p: CompletionParams) -> Result<Option<CompletionResponse>> {
        if p.context.is_none() || p.context.as_ref().unwrap().trigger_character.is_none() {
            return Ok(None);
        }
        let trigger = p.context.unwrap().trigger_character.unwrap();
        let ir = &self.state.lock().await.ir;
        match trigger.as_str() {
            "@" => {
                // Functions
                let mut res = Vec::new();
                for f in ir.funcs.iter() {
                    // Detail
                    let mut detail = format!("[{}", f.name);
                    for p in f.params.iter() {
                        //detail.push_str(&format!(" {}", ir.variables[*p].typ.data.fmt(&ir)));
                        detail.push_str(&format!(" {}", ir.variables[*p].name));
                    }
                    detail.push_str("]");
                    if (f.ret_typ.data != TypeData::VOID) && (f.ret_typ.data != TypeData::INVALID) {
                        detail.push_str(" -> ");
                        detail.push_str(&f.ret_typ.data.fmt(&ir));
                    }

                    // Name
                    let mut n = f.name.clone();
                    n.remove(0);

                    // Save
                    res.push(CompletionItem {
                        label: n,
                        kind: Some(CompletionItemKind::METHOD),
                        detail: Some(detail),
                        ..Default::default()
                    });
                }
                Ok(Some(CompletionResponse::Array(res)))
            }
            "!" => {
                // Variables
                let mut res = Vec::new();
                for v in ir.variables.iter() {
                    // Name
                    let mut n = v.name.clone();
                    n.remove(0);

                    // Save
                    res.push(CompletionItem {
                        label: n,
                        kind: Some(CompletionItemKind::VARIABLE),
                        detail: Some(v.typ.data.fmt(&ir)),
                        ..Default::default()
                    });
                }
                Ok(Some(CompletionResponse::Array(res)))
            }
            "$" => {
                // Types
                let mut res = Vec::new();
                for t in ir.types.iter() {
                    let kind = match t.typ.data {
                        TypeData::INT
                        | TypeData::FLOAT
                        | TypeData::BOOL
                        | TypeData::BOX
                        | TypeData::ARRAY(_)
                        | TypeData::TUPLE(_) => CompletionItemKind::VALUE,
                        TypeData::STRUCT(_) => CompletionItemKind::STRUCT,
                        TypeData::ENUM(_) => CompletionItemKind::ENUM,
                        _ => CompletionItemKind::TYPE_PARAMETER,
                    };

                    // Name
                    let mut n = t.name.clone();
                    n.remove(0);

                    // Save
                    res.push(CompletionItem {
                        label: n,
                        kind: Some(kind),
                        detail: Some(t.typ.data.fmt(&ir)),
                        ..Default::default()
                    });
                }
                Ok(Some(CompletionResponse::Array(res)))
            }
            "[" => {
                // Functions
                let mut res = Vec::new();
                for f in builtinFuncs() {
                    // Detail
                    let mut detail = format!("[{}", f.name);
                    for p in f.params.iter() {
                        detail.push_str(&format!(" {}", p.name()));
                    }
                    detail.push_str("]");
                    if f.ret.is_some() {
                        detail.push_str(" -> ");
                        detail.push_str(&f.ret.unwrap().name());
                    }

                    // Save
                    res.push(CompletionItem {
                        label: f.name.to_string(),
                        kind: Some(CompletionItemKind::FUNCTION),
                        detail: Some(detail),
                        ..Default::default()
                    });
                }
                Ok(Some(CompletionResponse::Array(res)))
            }
            _ => Ok(None),
        }
    }
}

// Builtin funcs
enum bPt {
    TYPE,
    NUMBER,
    BOX,
    VARIABLE,
    STRUCT,
    BOOL,
    FIELD,
    OPT(Box<bPt>),
    VARIADIC(Box<bPt>),
    BLOCK,
    ANY,
    ARRAY,
    ENUM,
    TUPLE,
    STRING,
}

// builtin param type
impl bPt {
    fn name(&self) -> String {
        match self {
            bPt::TYPE => "TYPE".to_string(),
            bPt::NUMBER => "NUMBER".to_string(),
            bPt::BOX => "BOX".to_string(),
            bPt::VARIABLE => "VARIABLE".to_string(),
            bPt::STRUCT => "STRUCT".to_string(),
            bPt::BOOL => "BOOL".to_string(),
            bPt::FIELD => "FIELD".to_string(),
            bPt::OPT(v) => v.name() + "?",
            bPt::VARIADIC(v) => v.name() + "...",
            bPt::BLOCK => "BLOCK".to_string(),
            bPt::ANY => "ANY".to_string(),
            bPt::ARRAY => "ARRAY".to_string(),
            bPt::ENUM => "ENUM".to_string(),
            bPt::TUPLE => "TUPLE".to_string(),
            bPt::STRING => "$STRING".to_string(),
        }
    }
}

struct BuiltinFunc {
    name: &'static str,
    params: Vec<bPt>,
    ret: Option<bPt>,
}

fn builtinFuncs() -> Vec<BuiltinFunc> {
    vec![
        BuiltinFunc {
            name: "ARRAY",
            params: vec![bPt::TYPE, bPt::OPT(Box::new(bPt::NUMBER))],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "STRUCT",
            params: vec![bPt::VARIADIC(Box::new(bPt::FIELD))],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "TUPLE",
            params: vec![bPt::VARIADIC(Box::new(bPt::TYPE))],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "ENUM",
            params: vec![bPt::VARIADIC(Box::new(bPt::TYPE))],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "FIELD",
            params: vec![bPt::VARIABLE, bPt::TYPE],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "CHAR",
            params: vec![],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "INT",
            params: vec![],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "FLOAT",
            params: vec![],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "BOOL",
            params: vec![],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "BX",
            params: vec![],
            ret: Some(bPt::TYPE),
        },
        BuiltinFunc {
            name: "PARAM",
            params: vec![bPt::VARIABLE, bPt::TYPE],
            ret: None,
        },
        BuiltinFunc {
            name: "RETURN",
            params: vec![bPt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "ARR",
            params: vec![bPt::VARIADIC(Box::new(bPt::ANY))],
            ret: None,
        },
        BuiltinFunc {
            name: "+",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "-",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "/",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "*",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "^",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "XOR",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "SHIFT",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "BOR",
            params: vec![bPt::NUMBER, bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: ">",
            params: vec![bPt::ANY, bPt::ANY],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: ">=",
            params: vec![bPt::ANY, bPt::ANY],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "<",
            params: vec![bPt::ANY, bPt::ANY],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "<=",
            params: vec![bPt::ANY, bPt::ANY],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "=",
            params: vec![bPt::ANY, bPt::ANY],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "NEQ",
            params: vec![bPt::ANY, bPt::ANY],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "NOT",
            params: vec![bPt::BOOL],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "&",
            params: vec![bPt::BOOL, bPt::BOOL],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "|",
            params: vec![bPt::BOOL, bPt::BOOL],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "DEFINE",
            params: vec![bPt::VARIABLE, bPt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "WHILE",
            params: vec![bPt::BOOL, bPt::BLOCK],
            ret: None,
        },
        BuiltinFunc {
            name: "LEN",
            params: vec![bPt::ARRAY],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "APPEND",
            params: vec![bPt::ARRAY, bPt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "GET",
            params: vec![bPt::STRUCT, bPt::FIELD],
            ret: Some(bPt::ANY),
        },
        BuiltinFunc {
            name: "GET",
            params: vec![bPt::ENUM, bPt::NUMBER],
            ret: Some(bPt::ANY),
        },
        BuiltinFunc {
            name: "GET",
            params: vec![bPt::ARRAY, bPt::NUMBER],
            ret: Some(bPt::ANY),
        },
        BuiltinFunc {
            name: "GET",
            params: vec![bPt::TUPLE, bPt::NUMBER],
            ret: Some(bPt::ANY),
        },
        BuiltinFunc {
            name: "NEW",
            params: vec![bPt::TYPE, bPt::ANY],
            ret: Some(bPt::ANY),
        },
        BuiltinFunc {
            name: "PRINT",
            params: vec![bPt::STRING],
            ret: None,
        },
        BuiltinFunc {
            name: "IF",
            params: vec![bPt::BOOL, bPt::BLOCK, bPt::BLOCK],
            ret: None,
        },
        BuiltinFunc {
            name: "BOX",
            params: vec![bPt::ANY],
            ret: Some(bPt::BOX),
        },
        BuiltinFunc {
            name: "PEEK",
            params: vec![bPt::BOX, bPt::TYPE],
            ret: Some(bPt::BOOL),
        },
        BuiltinFunc {
            name: "UNBOX",
            params: vec![bPt::BOX],
            ret: Some(bPt::ANY),
        },
        BuiltinFunc {
            name: ":",
            params: vec![bPt::FIELD, bPt::ANY],
            ret: Some(bPt::FIELD),
        },
        BuiltinFunc {
            name: "SET",
            params: vec![bPt::STRUCT, bPt::VARIADIC(Box::new(bPt::FIELD))],
            ret: None,
        },
        BuiltinFunc {
            name: "SET",
            params: vec![bPt::ENUM, bPt::NUMBER, bPt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "SET",
            params: vec![bPt::TUPLE, bPt::NUMBER, bPt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "SET",
            params: vec![bPt::ARRAY, bPt::NUMBER, bPt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "MATCH",
            params: vec![bPt::ANY, bPt::VARIADIC(Box::new(bPt::BLOCK))],
            ret: None,
        },
        BuiltinFunc {
            name: "CASE",
            params: vec![bPt::ANY, bPt::BLOCK],
            ret: None,
        },
        BuiltinFunc {
            name: "TOI",
            params: vec![bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "TOF",
            params: vec![bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
        BuiltinFunc {
            name: "TOC",
            params: vec![bPt::NUMBER],
            ret: Some(bPt::NUMBER),
        },
    ]
}
