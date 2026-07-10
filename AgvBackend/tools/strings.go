package tools

import "fmt"

func CMDLink(text, link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", text, link)
}
