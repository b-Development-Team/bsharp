use std::mem::discriminant;

use super::*;

impl super::Parser {
    pub fn parse(&mut self) -> Result<(), ParseError> {
        while self.tok.has_next() {
            let n = self.parse_node()?;
            self.ast.push(n);
        }
        Ok(())
    }

    fn parse_node(&mut self) -> Result<ASTNode, ParseError> {
        let tok = self.tok.peek().unwrap();
        match tok.data {
            TokenData::COMMENT(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Comment(v),
                    pos: tok.pos,
                })
            }
            TokenData::OPENBRACK => self.parse_stmt(),
            TokenData::STRING(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::String(v),
                    pos: tok.pos,
                })
            }
            TokenData::INTEGER(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Integer(v),
                    pos: tok.pos,
                })
            }
            TokenData::FLOAT(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Float(v),
                    pos: tok.pos,
                })
            }
            TokenData::TYPE(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Type(v),
                    pos: tok.pos,
                })
            }
            TokenData::VARIABLE(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Variable(v),
                    pos: tok.pos,
                })
            }
            TokenData::FUNCTION(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Function(v),
                    pos: tok.pos,
                })
            }
            TokenData::IDENT(v) => match v.as_str() {
                "TRUE" | "FALSE" => {
                    self.tok.eat().unwrap();
                    Ok(ASTNode {
                        data: ASTNodeData::Bool(v == "TRUE"),
                        pos: tok.pos,
                    })
                }
                _ => Err(ParseError::UnexpectedToken(self.tok.eat().unwrap())),
            },
            TokenData::CHAR(v) => {
                self.tok.eat().unwrap();
                Ok(ASTNode {
                    data: ASTNodeData::Char(v),
                    pos: tok.pos,
                })
            }
            _ => Err(ParseError::UnexpectedToken(tok)),
        }
    }

    fn parse_stmt(&mut self) -> Result<ASTNode, ParseError> {
        let mut pos = self.tok.eat().unwrap().pos; // Get [

        let name = match self
            .tok
            .expect(discriminant(&TokenData::IDENT("".to_string())))
        {
            Ok(t) => t,
            Err(_) => {
                // Block
                let mut nodes = Vec::new();
                while self.tok.peek()?.data != TokenData::CLOSEBRACK {
                    let n = self.parse_node()?;
                    nodes.push(n);
                }
                pos.extend(self.tok.eat().unwrap().pos); // Eat ]
                return Ok(ASTNode {
                    data: ASTNodeData::Block(nodes),
                    pos,
                });
            }
        };

        let mut vals = Vec::new();
        while self.tok.peek()?.data != TokenData::CLOSEBRACK {
            let val = self.parse_node()?;
            vals.push(val);
        }

        pos = pos.extend(self.tok.eat()?.pos); // Eat ]
        Ok(ASTNode {
            data: ASTNodeData::Stmt {
                name: match name.data {
                    TokenData::IDENT(v) => v,
                    _ => unreachable!(),
                },
                name_pos: name.pos,
                args: vals,
            },
            pos,
        })
    }
}