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
                    let detail = fn_string(&f, &ir);

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
                for f in builtin_funcs() {
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

// builtin param type
enum Bpt {
    TYPE,
    NUMBER,
    BOX,
    VARIABLE,
    STRUCT,
    BOOL,
    FIELD,
    OPT(Box<Bpt>),
    VARIADIC(Box<Bpt>),
    BLOCK,
    ANY,
    ARRAY,
    ENUM,
    TUPLE,
    STRING,
}

impl Bpt {
    fn name(&self) -> String {
        match self {
            Bpt::TYPE => "TYPE".to_string(),
            Bpt::NUMBER => "NUMBER".to_string(),
            Bpt::BOX => "BOX".to_string(),
            Bpt::VARIABLE => "VARIABLE".to_string(),
            Bpt::STRUCT => "STRUCT".to_string(),
            Bpt::BOOL => "BOOL".to_string(),
            Bpt::FIELD => "FIELD".to_string(),
            Bpt::OPT(v) => v.name() + "?",
            Bpt::VARIADIC(v) => v.name() + "...",
            Bpt::BLOCK => "BLOCK".to_string(),
            Bpt::ANY => "ANY".to_string(),
            Bpt::ARRAY => "ARRAY".to_string(),
            Bpt::ENUM => "ENUM".to_string(),
            Bpt::TUPLE => "TUPLE".to_string(),
            Bpt::STRING => "$STRING".to_string(),
        }
    }
}

struct BuiltinFunc {
    name: &'static str,
    params: Vec<Bpt>,
    ret: Option<Bpt>,
}

fn builtin_funcs() -> Vec<BuiltinFunc> {
    vec![
        BuiltinFunc {
            name: "ARRAY",
            params: vec![Bpt::TYPE, Bpt::OPT(Box::new(Bpt::NUMBER))],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "STRUCT",
            params: vec![Bpt::VARIADIC(Box::new(Bpt::FIELD))],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "TUPLE",
            params: vec![Bpt::VARIADIC(Box::new(Bpt::TYPE))],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "ENUM",
            params: vec![Bpt::VARIADIC(Box::new(Bpt::TYPE))],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "FIELD",
            params: vec![Bpt::VARIABLE, Bpt::TYPE],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "CHAR",
            params: vec![],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "INT",
            params: vec![],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "FLOAT",
            params: vec![],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "BOOL",
            params: vec![],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "BX",
            params: vec![],
            ret: Some(Bpt::TYPE),
        },
        BuiltinFunc {
            name: "PARAM",
            params: vec![Bpt::VARIABLE, Bpt::TYPE],
            ret: None,
        },
        BuiltinFunc {
            name: "RETURN",
            params: vec![Bpt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "ARR",
            params: vec![Bpt::VARIADIC(Box::new(Bpt::ANY))],
            ret: None,
        },
        BuiltinFunc {
            name: "+",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "-",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "/",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "*",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "^",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "XOR",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "SHIFT",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "BOR",
            params: vec![Bpt::NUMBER, Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: ">",
            params: vec![Bpt::ANY, Bpt::ANY],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: ">=",
            params: vec![Bpt::ANY, Bpt::ANY],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "<",
            params: vec![Bpt::ANY, Bpt::ANY],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "<=",
            params: vec![Bpt::ANY, Bpt::ANY],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "=",
            params: vec![Bpt::ANY, Bpt::ANY],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "NEQ",
            params: vec![Bpt::ANY, Bpt::ANY],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "NOT",
            params: vec![Bpt::BOOL],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "&",
            params: vec![Bpt::BOOL, Bpt::BOOL],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "|",
            params: vec![Bpt::BOOL, Bpt::BOOL],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "DEFINE",
            params: vec![Bpt::VARIABLE, Bpt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "WHILE",
            params: vec![Bpt::BOOL, Bpt::BLOCK],
            ret: None,
        },
        BuiltinFunc {
            name: "LEN",
            params: vec![Bpt::ARRAY],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "APPEND",
            params: vec![Bpt::ARRAY, Bpt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "GET",
            params: vec![Bpt::STRUCT, Bpt::FIELD],
            ret: Some(Bpt::ANY),
        },
        BuiltinFunc {
            name: "GET",
            params: vec![Bpt::ENUM, Bpt::NUMBER],
            ret: Some(Bpt::ANY),
        },
        BuiltinFunc {
            name: "GET",
            params: vec![Bpt::ARRAY, Bpt::NUMBER],
            ret: Some(Bpt::ANY),
        },
        BuiltinFunc {
            name: "GET",
            params: vec![Bpt::TUPLE, Bpt::NUMBER],
            ret: Some(Bpt::ANY),
        },
        BuiltinFunc {
            name: "NEW",
            params: vec![Bpt::TYPE, Bpt::ANY],
            ret: Some(Bpt::ANY),
        },
        BuiltinFunc {
            name: "PRINT",
            params: vec![Bpt::STRING],
            ret: None,
        },
        BuiltinFunc {
            name: "IF",
            params: vec![Bpt::BOOL, Bpt::BLOCK, Bpt::BLOCK],
            ret: None,
        },
        BuiltinFunc {
            name: "BOX",
            params: vec![Bpt::ANY],
            ret: Some(Bpt::BOX),
        },
        BuiltinFunc {
            name: "PEEK",
            params: vec![Bpt::BOX, Bpt::TYPE],
            ret: Some(Bpt::BOOL),
        },
        BuiltinFunc {
            name: "UNBOX",
            params: vec![Bpt::BOX],
            ret: Some(Bpt::ANY),
        },
        BuiltinFunc {
            name: ":",
            params: vec![Bpt::FIELD, Bpt::ANY],
            ret: Some(Bpt::FIELD),
        },
        BuiltinFunc {
            name: "SET",
            params: vec![Bpt::STRUCT, Bpt::VARIADIC(Box::new(Bpt::FIELD))],
            ret: None,
        },
        BuiltinFunc {
            name: "SET",
            params: vec![Bpt::ENUM, Bpt::NUMBER, Bpt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "SET",
            params: vec![Bpt::TUPLE, Bpt::NUMBER, Bpt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "SET",
            params: vec![Bpt::ARRAY, Bpt::NUMBER, Bpt::ANY],
            ret: None,
        },
        BuiltinFunc {
            name: "MATCH",
            params: vec![Bpt::ANY, Bpt::VARIADIC(Box::new(Bpt::BLOCK))],
            ret: None,
        },
        BuiltinFunc {
            name: "CASE",
            params: vec![Bpt::ANY, Bpt::BLOCK],
            ret: None,
        },
        BuiltinFunc {
            name: "TOI",
            params: vec![Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "TOF",
            params: vec![Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
        BuiltinFunc {
            name: "TOC",
            params: vec![Bpt::NUMBER],
            ret: Some(Bpt::NUMBER),
        },
    ]
}
