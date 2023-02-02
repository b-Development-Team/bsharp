use super::*;

pub struct Stream {
    source: String,
    index: usize,

    file: usize,
    line: usize,
    col: usize,
}

impl Stream {
    pub fn new(source: String, file: usize) -> Stream {
        Stream {
            source,
            index: 0,
            line: 0,
            col: 0,
            file,
        }
    }

    pub fn peek(&self) -> Option<char> {
        self.source.chars().nth(self.index)
    }

    pub fn eat(&mut self) -> Option<char> {
        let next = self.peek();
        self.index += 1;
        if let Some(v) = next {
            if v == '\n' {
                self.line += 1;
                self.col = 0;
            }
        }
        next
    }

    pub fn eof(&self) -> bool {
        self.peek().is_none()
    }

    pub fn pos(&self) -> Pos {
        Pos {
            file: self.file,
            line: self.line,
            col: self.col,
        }
    }
}
