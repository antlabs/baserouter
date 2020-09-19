package baserouter

import (
	"bytes"
	"fmt"
)

type path struct {
	originalPath []byte    //原始路径
	insertPath   []byte    //修改后的路径，单个变量变为: 所有变量变为*
	paramPath    []*handle //存放param
}

func genPath(p []byte) *path {
	p2 := &path{}
	p2.originalPath = p

	var paramName bytes.Buffer
	var insertPath bytes.Buffer

	foundParam := false
	wildcard := false
	maybeVar := false

	for i := 0; i < len(p); i++ {
		if !wildcard && !foundParam {
			if p[i] == '/' && !maybeVar {
				maybeVar = true
				insertPath.WriteByte('/')
				continue
			}
		}

		if maybeVar {
			maybeVar = false
			if !foundParam && !wildcard {

				if p[i] == ':' {
					foundParam = true
					insertPath.WriteString(":")
					continue
				}

				if p[i] == '*' {
					wildcard = true
					insertPath.WriteString("*")
					continue
				}
			}
		}

		if wildcard {
			if p[i] == '/' || foundParam {
				panic(fmt.Sprintf("catch-all routes are only allowed at the end of the path in path '%s'", p))
			}

			paramName.WriteByte(p[i])
			continue
		}

		if foundParam {

			if p[i] == '/' {
				foundParam = false
				maybeVar = true

				p2.checkParam(paramName)

				p2.addParamPath(insertPath, paramName)

				insertPath.WriteByte('/')

				paramName.Reset()
				continue
			}

			paramName.WriteByte(p[i])
			continue
		}

		insertPath.WriteByte(p[i])

	}

	if wildcard {

		p2.checkParam(paramName)

		p2.addParamPath(insertPath, paramName)
	}

	if foundParam {
		p2.checkParam(paramName)
		p2.addParamPath(insertPath, paramName)
	}

	if insertPath.Len() > 0 {
		p2.insertPath = insertPath.Bytes()
	}

	return p2
}

func (p *path) checkParam(paramName bytes.Buffer) {
	if paramName.Len() == 0 {
		panic(fmt.Sprintf("wildcards must be named with a non-empty name in path:%s",
			p.originalPath))
	}
}

func (p *path) addParamPath(insertPath, paramName bytes.Buffer) {
	if p.paramPath == nil {
		p.paramPath = make([]*handle, len(p.originalPath))
	}

	p.paramPath[insertPath.Len()-1] = &handle{paramName: paramName.String()}
}