use super::*;

#[derive(Debug, Clone)]
pub struct IRNode {
    pub data: IRNodeData,
    pub range: Pos,
    pub pos: Pos,
}

impl IRNode {
    pub fn new(data: IRNodeData, range: Pos, pos: Pos) -> Self {
        Self { data, range, pos }
    }

    pub fn void() -> Self {
        Self {
            data: IRNodeData::Void,
            range: Pos::default(),
            pos: Pos::default(),
        }
    }

    pub fn invalid(pos: Pos) -> Self {
        Self {
            data: IRNodeData::Invalid,
            range: pos,
            pos,
        }
    }

    pub fn typ(&self) -> Type {
        match &self.data {
            IRNodeData::Math(v, _, _) => v.typ(),
            IRNodeData::Comparison(_, _, _)
            | IRNodeData::Boolean(_, _, _)
            | IRNodeData::Peek { .. } => Type::from(TypeData::BOOL),
            IRNodeData::Block { .. }
            | IRNodeData::Define { .. }
            | IRNodeData::Print(_)
            | IRNodeData::Import(_)
            | IRNodeData::While { .. }
            | IRNodeData::TypeCase { .. }
            | IRNodeData::Case { .. }
            | IRNodeData::TypeMatch { .. }
            | IRNodeData::Match { .. }
            | IRNodeData::If { .. }
            | IRNodeData::Void
            | IRNodeData::Append { .. }
            | IRNodeData::StructOp { .. }
            | IRNodeData::SetStruct { .. }
            | IRNodeData::Return(_) => Type::void(),
            IRNodeData::Len(_) => Type::from(TypeData::INT),
            IRNodeData::GetEnum { typ, .. } => typ.clone(),
            IRNodeData::GetStruct { strct, field } => {
                if let TypeData::STRUCT { fields, .. } = &strct.typ().data {
                    for f in fields {
                        if f.name == *field {
                            return f.typ.clone();
                        }
                    }
                }
                unreachable!()
            }
            IRNodeData::Unbox { typ, .. } => typ.clone(),
            IRNodeData::NewArray(typ, _) => typ.clone(),
            IRNodeData::NewEnum(_, enm) => enm.clone(),
            IRNodeData::NewBox(_) => Type::from(TypeData::BOX),
            IRNodeData::NewStruct(typ, _) => typ.clone(),
            IRNodeData::Param { .. } | IRNodeData::Returns(_) => Type::from(TypeData::PARAM),
            IRNodeData::Field { .. } | IRNodeData::Generic { .. } => Type::from(TypeData::FIELD),
            IRNodeData::Type(_) | IRNodeData::TypeInstantiate { .. } => Type::from(TypeData::TYPE),
            IRNodeData::Cast(_, typ) => typ.clone(),
            IRNodeData::Invalid => Type::from(TypeData::INVALID),
            IRNodeData::FnCall { ret_typ, .. } => ret_typ.clone(),
            IRNodeData::Variable(_, typ) => typ.clone(),
            IRNodeData::Int(_) => Type::from(TypeData::INT),
            IRNodeData::Float(_) => Type::from(TypeData::FLOAT),
        }
    }
}

#[derive(Debug, Clone)]
pub enum IRNodeData {
    Block {
        scope: usize,
        body: Vec<IRNode>,
    },
    Print(Vec<IRNode>),
    Math(Box<IRNode>, MathOperator, Box<IRNode>),
    Comparison(Box<IRNode>, ComparisonOperator, Box<IRNode>),
    Boolean(Box<IRNode>, BooleanOperator, Box<IRNode>),
    Define {
        var: usize,
        val: Box<IRNode>,
        edit: bool, // Whether making variable or just changing value
    },
    Variable(usize, Type), // Gets value of variable
    Import(String),
    While {
        cond: Box<IRNode>,
        body: Box<IRNode>,
    },
    TypeCase {
        var: usize,
        typ: Type,
        body: Box<IRNode>,
    },
    Case {
        val: Box<IRNode>,
    },
    TypeMatch {
        val: Box<IRNode>,
        body: Vec<IRNode>, // TypeCase
    },
    Match {
        val: Box<IRNode>,
        body: Vec<IRNode>, // Case
    },
    If {
        cond: Box<IRNode>,
        body: Box<IRNode>, // Doesn't have to be block, can be used as ternary op
        els: Box<IRNode>,
        ret_typ: Type,
    },
    FnCall {
        func: usize,
        args: Vec<IRNode>,
        ret_typ: Type,
    },
    Return(Option<Box<IRNode>>),
    Void,
    Invalid,

    // Composite type ops
    Len(Box<IRNode>),
    Append {
        arr: Box<IRNode>,
        val: Box<IRNode>,
    },
    GetEnum {
        enm: Box<IRNode>,
        typ: Type,
    },
    GetStruct {
        strct: Box<IRNode>,
        field: String,
    },
    StructOp {
        field: String,
        val: Box<IRNode>,
    }, // [:]
    SetStruct {
        strct: Box<IRNode>,
        vals: Box<IRNode>,
    },
    Field {
        name: String,
        typ: Type,
    },
    Peek {
        bx: Box<IRNode>,
        typ: Type,
    },
    Unbox {
        bx: Box<IRNode>,
        typ: Type,
    },

    // Literals
    Int(i64),
    Float(f64),

    // Allocating
    NewArray(Type, Option<Box<IRNode>>), // optional: capacity
    NewEnum(Box<IRNode>, Type),
    NewBox(Box<IRNode>),          // [BOX]
    NewStruct(Type, Vec<IRNode>), // [:] statements

    // Functions
    Param(usize),
    Returns(Type),

    // Types
    Type(Type),     // [INT], [STRUCT], etc.
    Generic(usize), // Type parameters
    TypeInstantiate {
        typ: Type,
        params: Vec<Type>,
    }, // [G <Type> <Params>]
    Cast(Box<IRNode>, Type),
}

#[derive(Debug, Clone)]
pub enum MathOperator {
    ADD,
    SUBTRACT,
    MULTIPLY,
    DIVIDE,
    MODULO,
}

#[derive(Debug, Clone)]
pub enum ComparisonOperator {
    GREATER,
    LESS,
    GREATEREQUAL,
    LESSEQUAL,
    EQUAL,
    NOTEQUAL,
}

#[derive(Debug, Clone)]
pub enum BooleanOperator {
    AND,
    OR,
}
