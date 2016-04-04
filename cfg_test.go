package cfg

import (
	"testing"
)

func TestSimpleConfigString(t *testing.T) {
	cfg0s := "Title = Sample App Title\r\nDescription = Line 1\\\nLine2\nRating = 1 \"STAR\""
	cfg0, err := ParseString(cfg0s)
	if err != nil {
		t.Fatal(err)
	}
	if cfg0["Title"] != "Sample App Title" {
		t.Error(`cfg0["Title"]`, cfg0["Title"], " != Sample App Title")
	}
	if cfg0["Description"] != "Line 1\nLine2" {
		t.Error(`cfg0["Description"]`, cfg0["Description"], " != Line 1\nLine2")
	}
	if cfg0["Rating"] != "1 \"STAR\"" {
		t.Error(`cfg0["Rating"]`, cfg0["Rating"], " != 1 \"STAR\"")
	}
}

func TestSimpleConfigString2(t *testing.T) {
	cfg0s := `# AdminPagController
adminpag_newuser_ok_1        = Olá %v,\
O seu usuário Buenos Ayres está pronto para ser utilizado.
adminpag_newuser_ok_2        = \
Acesse este link %v para cadastrar uma senha.\
\
(este acesso expira em 24h)
`
	cfg0, err := ParseString(cfg0s)
	if err != nil {
		t.Fatal(err)
	}
	if cfg0["adminpag_newuser_ok_1"] != `Olá %v,
O seu usuário Buenos Ayres está pronto para ser utilizado.` {
		t.Error(`cfg0["adminpag_newuser_ok_1"]`, cfg0["adminpag_newuser_ok_1"], `Olá %v,
O seu usuário Buenos Ayres está pronto para ser utilizado.`)
	}
	if cfg0["adminpag_newuser_ok_2"] != `
Acesse este link %v para cadastrar uma senha.

(este acesso expira em 24h)` {
		t.Error(`cfg0["adminpag_newuser_ok_2"]`, "```", cfg0["adminpag_newuser_ok_2"], "```", `
Acesse este link %v para cadastrar uma senha.

(este acesso expira em 24h)`)
	}
}
