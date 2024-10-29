// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package test

import (
	"encoding/base64"
	"testing"

	"github.com/hootrhino/rhilex/component/security"
)

// go test -timeout 30s -run ^TestRSA_License github.com/hootrhino/rhilex/test -v -count=1
func TestRSA_License(t *testing.T) {
	if !security.CheckLicense() {
		security.InitSecurityLicense()
	}
	PasswordByte1, E1 := security.RSA_Encrypt([]byte("12345678"), "./.encrypt.pub")
	if E1 != nil {
		t.Fatal(E1)
	}
	Password := base64.StdEncoding.EncodeToString(PasswordByte1)
	t.Log("EncodeToString == ", Password)
	PasswordByte2, _ := base64.StdEncoding.DecodeString(Password)
	t.Log("PasswordByte1 == ", PasswordByte1)
	t.Log("PasswordByte2 == ", PasswordByte2)
	{
		B2, E2 := security.RSA_Decrypt(PasswordByte1, "./.encrypt.priv")
		if E2 != nil {
			t.Fatal(E1)
		}
		t.Log(B2, " ==== ", string(B2))
	}
	{
		B2, E2 := security.RSA_Decrypt(PasswordByte2, "./.encrypt.priv")
		if E2 != nil {
			t.Fatal(E1)
		}
		t.Log(B2, " ==== ", string(B2))
	}
}

// go test -timeout 30s -run ^TestRSA_LicenseDecrypt github.com/hootrhino/rhilex/test -v -count=1

func TestRSA_LicenseDecrypt(t *testing.T) {
	Password := `ySzFf1MjcZF1yEB6CsoeeZ7UVsItOpwyVtGRjLvW1cQnJFXYhmEYnwnp3iuXkmv3zRr3xtPRKjT2hlXvzkFJBMgOr5t5hZw7wPlrmaoRRISQXm5jmGhNdNZZXv5tG6EqId3iziX/j6Wm36nXKlWAOuX5TwKuEFa7s8mBxwRYkQKO2BP2hiuRJFhrXtVi3facv4eAHh65MxupV9GQJZ8W28wk9ICZpvaEJijAlMl+BPNH9b0an5TneCzZs2zIWZ4fcEuZttwyMsg1SDhEuEpdU2FM5/go0nzLolnJtQ9+0GmInWhOEHXxuKGDV6U9vkFW8WSTD7A5Fu4qqIRvIURVmw==`
	PasswordByte1, _ := base64.StdEncoding.DecodeString(Password)
	t.Log("PasswordByte1 == ", PasswordByte1)
	B2, E2 := security.RSA_Decrypt(PasswordByte1, "./.encrypt.priv")
	if E2 != nil {
		t.Fatal(E2)
	}
	t.Log("RSA_Decrypt == ", B2, " ==== ", string(B2))
}
