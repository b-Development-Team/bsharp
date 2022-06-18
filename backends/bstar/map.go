package bstar

/*
[FUNC MAP_SET [ARRAY "map" "key" "val"] [BLOCK
  [DEFINE ind [FIND [INDEX [VAR map] 0] [VAR key]]]
  [IF [COMPARE [VAR ind] == -1] [BLOCK
    [DEFINE keys [CONCAT [INDEX [VAR map] 0] [ARRAY [VAR key]]]]
    [DEFINE vals [CONCAT [INDEX [VAR map] 1] [ARRAY [VAR val]]]]
    [RETURN [ARRAY [VAR keys] [VAR vals]]]
  ] ""]
  [DEFINE vals [SETINDEX [INDEX [VAR map] 1] [VAR ind] [VAR val]]]
  [SETINDEX [VAR map] 1 [VAR vals]]
]]
*/
func (b *BStar) mapSetFunc() Node {
	return blockNode(true, constNode("FUNC"), constNode("MAP_SET"), blockNode(true, constNode("ARRAY"), constNode(`"map"`), constNode(`"key"`), constNode(`"val"`)),
		blockNode(false, constNode("BLOCK"),
			// [DEFINE ind [FIND [INDEX [VAR map] 0] [VAR key]]]
			blockNode(false, constNode("DEFINE"), constNode("ind"),
				blockNode(true, constNode("FIND"),
					blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("map")), constNode(0)),
					blockNode(true, constNode("VAR"), constNode("key")),
				),
			),
			// [IF [COMPARE [VAR ind] == -1] [BLOCK
			blockNode(false, constNode("IF"),
				blockNode(true, constNode("COMPARE"),
					blockNode(true, constNode("VAR"), constNode("ind")),
					constNode("=="),
					constNode(-1),
				),
				// [BLOCK
				blockNode(true, constNode("BLOCK"),
					// [DEFINE keys [CONCAT [INDEX [VAR map] 0] [ARRAY [VAR key]]]]
					blockNode(false, constNode("DEFINE"), constNode("keys"),
						blockNode(true, constNode("CONCAT"),
							blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("map")), constNode(0)),
							blockNode(true, constNode("ARRAY"), blockNode(true, constNode("VAR"), constNode("key"))),
						),
					),
					// [DEFINE vals [CONCAT [INDEX [VAR map] 1] [ARRAY [VAR val]]]]
					blockNode(false, constNode("DEFINE"), constNode("vals"),
						blockNode(true, constNode("CONCAT"),
							blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("map")), constNode(1)),
							blockNode(true, constNode("ARRAY"), blockNode(true, constNode("VAR"), constNode("val"))),
						),
					),
					// [RETURN [ARRAY [VAR keys] [VAR vals]]]
					blockNode(false, constNode("RETURN"),
						blockNode(true, constNode("ARRAY"),
							blockNode(true, constNode("VAR"), constNode("keys")),
							blockNode(true, constNode("VAR"), constNode("vals")),
						),
					),
				),
				b.noPrintNode(),
			),
			// [DEFINE vals [SETINDEX [INDEX [VAR map] 1] [VAR ind] [VAR val]]]
			blockNode(false, constNode("DEFINE"), constNode("vals"),
				blockNode(true, constNode("SETINDEX"),
					blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("map")), constNode(1)),
					blockNode(true, constNode("VAR"), constNode("ind")),
					blockNode(true, constNode("VAR"), constNode("val")),
				),
			),
			// [SETINDEX [VAR map] 1 [VAR vals]]
			blockNode(false, constNode("SETINDEX"),
				blockNode(true, constNode("VAR"), constNode("map")),
				constNode(1),
				blockNode(true, constNode("VAR"), constNode("vals")),
			),
		),
		// ] ""]
	)
}

/*
[FUNC MAP_GET [ARRAY "map" "key"] [BLOCK
  [DEFINE ind [FIND [INDEX [VAR map] 0] [VAR key]]]
  [INDEX [INDEX [VAR map] 1] [VAR ind]]
]]
*/
func (b *BStar) mapGetFunc() Node {
	return blockNode(true, constNode("FUNC"), constNode("MAP_GET"), blockNode(true, constNode("ARRAY"), constNode(`"map"`), constNode(`"key"`)),
		blockNode(false, constNode("BLOCK"),
			// [DEFINE ind [FIND [INDEX [VAR map] 0] [VAR key]]]
			blockNode(false, constNode("DEFINE"), constNode("ind"),
				blockNode(true, constNode("FIND"),
					blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("map")), constNode(0)),
					blockNode(true, constNode("VAR"), constNode("key")),
				),
			),
			// [INDEX [INDEX [VAR map] 1] [VAR ind]]
			blockNode(false, constNode("INDEX"),
				blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("map")), constNode(1)),
				blockNode(true, constNode("VAR"), constNode("ind")),
			),
		),
	)
}
