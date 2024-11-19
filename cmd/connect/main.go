// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
)

func dontExit() {
	select {}
}

func main() {
	fmt.Println("Starting connect bot...")

	// connect.StartBotsAt("127.0.0.1:2106", 1)

	// dontExit()
}
