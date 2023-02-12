use tokens::Tokenizer;
mod error;
pub use error::FSetError;
pub use parser::*;
pub use tokens::*;

pub struct File {
    pub name: String,
    pub ast: parser::Parser,
}

pub struct FSet {
    pub files: Vec<File>,
}

impl FSet {
    pub fn new() -> FSet {
        FSet { files: Vec::new() }
    }

    pub fn add_file_source(&mut self, name: String, source: String) -> Result<(), FSetError> {
        let ind = self.files.len();
        let mut tok = Tokenizer::new(source, ind);
        tok.tokenize()?;

        let mut p = parser::Parser::new(tok);
        p.parse()?;

        self.files.push(File { name, ast: p });

        Ok(())
    }

    pub fn import(&mut self, dir: &String) -> Result<(), FSetError> {
        panic!("unimplemented")
    }
}
