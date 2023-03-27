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

    pub fn typ(&self, ir: &IR) -> Type {
        match &self.data {
            IRNodeData::Math(v, _, _) => v.typ(ir),
            IRNodeData::Comparison(_, _, _)
            | IRNodeData::Boolean(_, _, _)
            | IRNodeData::Peek { .. } => Type::from(TypeData::BOOL),
            IRNodeData::Block { .. }
            | IRNodeData::Define { .. }
            | IRNodeData::Print(_)
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
            | IRNodeData::SetArr { .. }
            | IRNodeData::Return(_) => Type::void(),
            IRNodeData::Len(_) => Type::from(TypeData::INT),
            IRNodeData::GetEnum { typ, .. } => typ.clone(),
            IRNodeData::GetStruct { strct, field } => {
                if let TypeData::STRUCT(fields) = &strct.typ(ir).data.concrete(ir) {
                    for f in fields {
                        if f.name == *field {
                            return f.typ.clone();
                        }
                    }
                }
                unreachable!()
            }
            IRNodeData::GetTuple { tup, ind } => match tup.typ(ir).data.concrete(ir) {
                TypeData::TUPLE(fields) => fields[*ind].clone(),
                _ => unreachable!(),
            },
            IRNodeData::Unbox { typ, .. } => typ.clone(),
            IRNodeData::NewArray(typ, _) | IRNodeData::NewArrayLiteral(typ, _) => typ.clone(),
            IRNodeData::NewEnum(_, enm) => enm.clone(),
            IRNodeData::NewBox(_) => Type::from(TypeData::BOX),
            IRNodeData::NewStruct(typ, _) => typ.clone(),
            IRNodeData::NewTuple(typ, _) => typ.clone(),
            IRNodeData::Param { .. } | IRNodeData::Returns(_) => Type::from(TypeData::PARAM),
            IRNodeData::Field { .. } => Type::from(TypeData::FIELD),
            IRNodeData::Type(_) => Type::from(TypeData::TYPE),
            IRNodeData::Cast(_, typ) => typ.clone(),
            IRNodeData::Invalid => Type::from(TypeData::INVALID),
            IRNodeData::FnCall { ret_typ, .. } => ret_typ.clone(),
            IRNodeData::Variable(_, typ) => typ.clone(),
            IRNodeData::Int(_) => Type::from(TypeData::INT),
            IRNodeData::Float(_) => Type::from(TypeData::FLOAT),
            IRNodeData::Char(_) => Type::from(TypeData::CHAR),
            IRNodeData::GetArr { arr, .. } => match arr.typ(ir).data.concrete(ir) {
                TypeData::ARRAY(body) => *body,
                _ => unreachable!(),
            },
        }
    }
}

#[derive(Debug, Clone)]
pub enum IRNodeData {
    Block {
        scope: usize,
        body: Vec<IRNode>,
    },
    Print(Box<IRNode>),
    Math(Box<IRNode>, MathOperator, Box<IRNode>),
    Comparison(Box<IRNode>, ComparisonOperator, Box<IRNode>),
    Boolean(Box<IRNode>, BooleanOperator, Option<Box<IRNode>>),
    Define {
        var: usize,
        val: Box<IRNode>,
        edit: bool, // Whether making variable or just changing value
    },
    Variable(usize, Type), // Gets value of variable
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
        body: Box<IRNode>,
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
        ret_typ: Option<Type>,
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
    GetArr {
        arr: Box<IRNode>,
        ind: Box<IRNode>,
    },
    SetArr {
        arr: Box<IRNode>,
        ind: Box<IRNode>,
        val: Box<IRNode>,
    },
    GetTuple {
        tup: Box<IRNode>,
        ind: usize,
    },
    StructOp {
        field: String,
        val: Box<IRNode>,
    }, // [:]
    SetStruct {
        strct: Box<IRNode>,
        vals: Vec<IRNode>,
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
    Char(u8),

    // Allocating
    NewArray(Type, Option<Box<IRNode>>), // optional: capacity
    NewArrayLiteral(Type, Vec<IRNode>),
    NewEnum(Box<IRNode>, Type),
    NewBox(Box<IRNode>),          // [BOX]
    NewStruct(Type, Vec<IRNode>), // [:] statements
    NewTuple(Type, Vec<IRNode>),

    // Functions
    Param(usize),
    Returns(Type),

    // Types
    Type(Type), // [INT], [STRUCT], etc.
    Cast(Box<IRNode>, Type),
}

#[derive(Debug, Clone, PartialEq, Copy)]
pub enum MathOperator {
    ADD,
    SUBTRACT,
    MULTIPLY,
    DIVIDE,
    MODULO,

    // Bitwise ops
    XOR,
    SHIFT,
    BOR,
}

#[derive(Debug, Clone, PartialEq, Copy)]
pub enum ComparisonOperator {
    GREATER,
    LESS,
    GREATEREQUAL,
    LESSEQUAL,
    EQUAL,
    NOTEQUAL,
}

#[derive(Debug, Clone, PartialEq, Copy)]
pub enum BooleanOperator {
    AND,
    OR,
    NOT,
}
