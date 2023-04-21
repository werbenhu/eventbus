// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package eventbus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrString(t *testing.T) {
	c := err{
		Msg:  "test",
		Code: 0x1,
	}

	require.Equal(t, "test", c.String())
}

func TestErrErrorr(t *testing.T) {
	c := err{
		Msg:  "error",
		Code: 0x1,
	}

	require.Equal(t, "error", error(c).Error())
}
