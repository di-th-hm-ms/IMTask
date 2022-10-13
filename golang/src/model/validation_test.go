package model

import "testing"

/**
* email
*/
func TestEmail(t *testing.T) {
	// 50 charas
	email := `0123456789012345678901234567890123abcd.e@gmail.com`
	if !ValidateEmail(email) {
		t.Error("miss")
	}
}

func TestEmail2(t *testing.T) {
	// all charas
	email := `1a._-B@-_2w.cOm`
	if !ValidateEmail(email) {
		t.Error()
	}
}

func TestEmail3(t *testing.T) {
	// 50 charas
	email := `0123456789012345678901234567890123abcd.e@gmail.com`
	if !ValidateEmail(email) {
		t.Error()
	}
}



func TestEmailErr(t *testing.T) {
	// 51 charas
	email := `012345678901234567890123456789012345678.e@gmail.com`
	if ValidateEmail(email) {
		t.Error("51 charas")
	}
}

func TestEmailErr2(t *testing.T) {
	// 51 charas
	email := `@gmail.com`
	if ValidateEmail(email) {
		t.Error("nothing before @")
	}
}

func TestEmailErr3(t *testing.T) {
	// 51 charas
	email := `abc123@.com`
	if ValidateEmail(email) {
		t.Error("no host")
	}
}

func TestEmailErr4(t *testing.T) {
	email := `abc123@gmail.`
	if ValidateEmail(email) {
		t.Error("no suffix")
	}
}

func TestEmailErr5(t *testing.T) {
	email := `abc123@gmail..`
	if ValidateEmail(email) {
		t.Error("suffix wrong chara")
	}
}

/**
* password
*/
func TestPw(t *testing.T) {
	password := `0Ab3dEf7`
	if !ValidatePassword(password) {
		t.Error("正常系エラー (8letters)")
	}
}

func TestPw100(t *testing.T) {
	password := `01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567aB`
	if !ValidatePassword(password) {
		// t.Error("正常系エラー (100letters), %d", len(password))
		t.Error(len(password))
	}
}

func TestPw101Err(t *testing.T) {
	password := `012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678aB`
	if ValidatePassword(password) {
		// t.Error("正常系エラー (100letters), %d", len(password))
		t.Error(len(password))
	}
}

func TestPwErr(t *testing.T) {
	// only number
	password := `12345678`
	if ValidatePassword(password) {
		t.Error("only number")
	}
}

func TestPwErr2(t *testing.T) {
	// only lower case
	password := `abcdefgh`
	if ValidatePassword(password) {
		t.Error("only lower case")
	}
}

func TestPwErr3(t *testing.T) {
	// only uppercase
	password := `ABCDEFGH`
	if ValidatePassword(password) {
		t.Error("only uppercase")
	}
}

func TestPwErr4(t *testing.T) {
	// no upper case
	password := `1234abcd`
	if ValidatePassword(password) {
		t.Error("no upper case")
	}
}

func TestPwErr5(t *testing.T) {
	// no number
	password := `abcdEFGH`
	if ValidatePassword(password) {
		t.Error("no number")
	}
}

func TestPwErr6(t *testing.T) {
	// no lower case
	password := `1234EFGH`
	if ValidatePassword(password) {
		t.Error("no lower case")
	}
}

func TestPwErr7(t *testing.T) {
	// no upper case
	password := `1234abC_d`
	if ValidatePassword(password) {
		t.Error("wrong character")
	}
}
