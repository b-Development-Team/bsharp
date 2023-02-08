pub enum Type {
    INT,
    FLOAT,
    CHAR,
    BOOL,
    BOX,
    ARRAY(Box<Type>),
    STRUCT(Vec<Field>),
    TUPLE(Vec<Type>),
    ENUM(Vec<Type>),
    INTERFACE(Vec<Type>),

    INVALID,
    TYPE,
}

pub struct Field {
    pub name: String,
    pub typ: Type,
}
