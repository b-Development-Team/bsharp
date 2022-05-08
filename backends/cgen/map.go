package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type mapFns struct {
	structName string
	hashFn     string
	compareFn  string
	freeFn     string
}

func (c *CGen) addMapFns(typ *types.MapType) mapFns {
	name := typName(typ)
	exists1 := isDynamic(typ.KeyType)
	exists2 := isDynamic(typ.ValType)
	_, exists := c.addedFns[name]
	if exists {
		freeFn := name + "_free"
		if !exists1 && !exists2 {
			freeFn = "NULL"
		}
		return mapFns{
			structName: name,
			hashFn:     name + "_hash",
			compareFn:  name + "_compare",
			freeFn:     freeFn,
		}
	}

	// Add struct
	fmt.Fprintf(c.globalfns, "struct %s {\n%s%s key;\n%s%s val;\n};\n\n", name, c.Config.Tab, c.CType(typ.KeyType), c.Config.Tab, c.CType(typ.ValType))

	// Add compare fn
	fmt.Fprintf(c.globalfns, "int %s_compare(const void *va, const void *vb, void *udata) {\n%sconst struct %s *a = va;\n%sconst struct %s *b = vb;\n%s", name, c.Config.Tab, name, c.Config.Tab, name, c.Config.Tab)
	switch typ.KeyType.BasicType() {
	case types.INT, types.FLOAT, types.BYTE:
		c.globalfns.WriteString("return (a->key == b->key) ? 0 : 1;")

	case types.STRING:
		c.globalfns.WriteString("return (string_cmp(a->key, b->key) == 0) ? 0 : 1;")
	}
	c.globalfns.WriteString("\n}\n\n")

	// Add hash fn
	fmt.Fprintf(c.globalfns, "static inline uint64_t %s_hash(const void *item, uint64_t seed0, uint64_t seed1) {\n%sconst struct %s *v = item;\n%s", name, c.Config.Tab, name, c.Config.Tab)
	switch typ.KeyType.BasicType() {
	case types.INT, types.FLOAT, types.BYTE:
		c.globalfns.WriteString("return *(uint64_t*)(v->key);")

	case types.STRING:
		c.globalfns.WriteString("return hashmap_sip(v->key->data, v->key->len, seed0, seed1);")
	}
	c.globalfns.WriteString("\n}\n\n")

	// Add free fn
	freeFn := "NULL"
	if exists1 || exists2 {
		freeFn = name + "_free"
		fmt.Fprintf(c.globalfns, "void %s_free(void *item) {\n%sconst struct %s *v = item;\n%s", name, c.Config.Tab, name, c.Config.Tab)
		if exists1 {
			c.globalfns.WriteString(c.FreeCode("v->key", typ.KeyType))
			if exists2 {
				c.globalfns.WriteString("\n" + c.Config.Tab)
			}
		}
		if exists2 {
			c.globalfns.WriteString(c.FreeCode("v->val", typ.ValType))
		}
		c.globalfns.WriteString("\n}\n\n")
	}

	// Return
	c.addedFns[name] = struct{}{}
	return mapFns{
		structName: name,
		hashFn:     name + "_hash",
		compareFn:  name + "_compare",
		freeFn:     freeFn,
	}
}

func (c *CGen) addMake(n *ir.MakeNode) (*Code, error) {
	if types.ARRAY.Equal(n.Type()) {
		name := c.GetTmp("arr")
		pre := fmt.Sprintf("array* %s = array_new(sizeof(%s), 1);", name, c.CType(n.Type().(*types.ArrayType).ElemType))
		c.stack.Add(c.FreeCode(name, n.Type()))
		return &Code{
			Pre:   pre,
			Value: name,
		}, nil
	}

	if types.STRUCT.Equal(n.Type()) {
		name := c.GetTmp("struct")
		t := typName(n.Type())
		pre := fmt.Sprintf("struct %s* %s = malloc(sizeof(struct %s));", t, name, t)
		pre = JoinCode(pre, fmt.Sprintf("%s->refs = 1;", name))
		// Initialize all fields
		typ := n.Type().(*types.StructType)
		for i, f := range typ.Fields {
			pre = JoinCode(pre, fmt.Sprintf("%s->f%d = %s;", name, i, c.ZeroValue(f.Type)))
		}

		c.stack.Add(c.FreeCode(name, n.Type()))
		return &Code{
			Pre:   pre,
			Value: name,
		}, nil
	}

	// Map
	m := n.Type().(*types.MapType)
	fns := c.addMapFns(m)
	name := c.GetTmp("map")
	pre := fmt.Sprintf("map* %s = map_new(sizeof(struct %s), %s, %s, %s);", name, fns.structName, fns.compareFn, fns.hashFn, fns.freeFn)
	c.stack.Add(c.FreeCode(name, n.Type()))
	return &Code{
		Pre:   pre,
		Value: name,
	}, nil
}

func (c *CGen) addSet(n *ir.SetNode) (*Code, error) {
	m, err := c.AddNode(n.Map)
	if err != nil {
		return nil, err
	}
	k, err := c.AddNode(n.Key)
	if err != nil {
		return nil, err
	}
	v, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	fns := c.addMapFns(n.Map.Type().(*types.MapType))

	// Check if already was value, if so free
	pre := &strings.Builder{}
	name := c.GetTmp("hashchk")
	fmt.Fprintf(pre, "struct %s* %s = hashmap_get(%s->map, &(struct %s){ .key=%s });\n", fns.structName, name, m.Value, fns.structName, k.Value)
	fmt.Fprintf(pre, "if (%s != NULL) {\n", name)
	fmt.Fprintf(pre, "%s%s(%s);\n%s%s->len--;\n", c.Config.Tab, fns.freeFn, name, c.Config.Tab, m.Value)
	pre.WriteString("}\n")
	fmt.Fprintf(pre, "%s->len++;\n", m.Value)

	// Store
	fmt.Fprintf(pre, "hashmap_set(%s->map, &(struct %s){ .key=%s, .val=%s });", m.Value, fns.structName, k.Value, v.Value)

	// Grab
	grab := ""
	if isDynamic(n.Key.Type()) {
		grab = JoinCode(grab, c.GrabCode(k.Value, n.Key.Type()))
	}
	if isDynamic(n.Value.Type()) {
		grab = JoinCode(grab, c.GrabCode(v.Value, n.Value.Type()))
	}
	return &Code{
		Pre: JoinCode(m.Pre, k.Pre, v.Pre, grab, pre.String()),
	}, nil
}

func (c *CGen) addGet(n *ir.GetNode) (*Code, error) {
	m, err := c.AddNode(n.Map)
	if err != nil {
		return nil, err
	}
	k, err := c.AddNode(n.Key)
	if err != nil {
		return nil, err
	}

	// Store val
	name := c.GetTmp("hashval")
	fns := c.addMapFns(n.Map.Type().(*types.MapType))
	pre := fmt.Sprintf("%s %s = ((struct %s*)(hashmap_get(%s->map, &(struct %s){ .key=%s })))->val;", c.CType(n.Type()), name, fns.structName, m.Value, fns.structName, k.Value)
	return &Code{
		Pre:   JoinCode(m.Pre, k.Pre, pre),
		Value: name,
	}, nil
}

func (c *CGen) keysFns(t *types.MapType) {
	name := typName(t)
	_, exists := c.addedFns[name+"_keys"]
	if exists {
		return
	}

	fmt.Fprintf(c.globalfns, "bool %s_iter(const void *item, void *data) {\n", name)
	fmt.Fprintf(c.globalfns, "%sconst struct %s *v = item;\n", c.Config.Tab, name)
	fmt.Fprintf(c.globalfns, "%sarray* arr = *((array**)data);\n", c.Config.Tab)
	fmt.Fprintf(c.globalfns, "%sarray_append(arr, v->key);\n", c.Config.Tab)
	fmt.Fprintf(c.globalfns, "%sreturn true;\n", c.Config.Tab)
	c.globalfns.WriteString("}\n\n")

	// keys
	fmt.Fprintf(c.globalfns, "array* %s_keys(map* map) {\n", name)
	fmt.Fprintf(c.globalfns, "%sarray* arr = array_new(sizeof(struct %s), map->len);\n", c.Config.Tab, name)
	fmt.Fprintf(c.globalfns, "%shashmap_scan(map->map, %s_iter, &arr);\n", c.Config.Tab, name)
	fmt.Fprintf(c.globalfns, "%sreturn arr;\n", c.Config.Tab)
	c.globalfns.WriteString("}\n\n")

	c.addedFns[name+"_keys"] = struct{}{}
}

func (c *CGen) addKeys(n *ir.KeysNode) (*Code, error) {
	c.keysFns(n.Map.Type().(*types.MapType))
	m, err := c.AddNode(n.Map)
	if err != nil {
		return nil, err
	}
	name := c.GetTmp("keys")
	pre := fmt.Sprintf("array* %s = %s_keys(%s);", name, typName(n.Map.Type().(*types.MapType)), m.Value)
	c.stack.Add(c.FreeCode(name, n.Type()))
	return &Code{
		Pre:   JoinCode(m.Pre, pre),
		Value: name,
	}, nil
}

func (c *CGen) addExists(n *ir.ExistsNode) (*Code, error) {
	m, err := c.AddNode(n.Map)
	if err != nil {
		return nil, err
	}
	k, err := c.AddNode(n.Key)
	if err != nil {
		return nil, err
	}
	name := c.GetTmp("exists")
	pre := fmt.Sprintf("bool %s = hashmap_get(%s->map, &(struct %s){ .key=%s }) != NULL;", name, m.Value, typName(n.Map.Type().(*types.MapType)), k.Value)
	return &Code{
		Pre:   JoinCode(m.Pre, k.Pre, pre),
		Value: name,
	}, nil
}
