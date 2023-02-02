pub struct Stream {
    source: String,
    index: usize,
}

impl Stream {
    pub fn new(source: String) -> Stream {
        Stream { source, index: 0 }
    }

    pub fn peek(&self) -> Option<char> {
        self.source.chars().nth(self.index)
    }

    pub fn eat(&mut self) -> Option<char> {
        let next = self.peek();
        self.index += 1;
        next
    }

    pub fn eof(&self) -> bool {
        self.peek().is_none()
    }
}
