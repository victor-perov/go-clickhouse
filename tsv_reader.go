package clickhouse

import "bufio"

type TsvParser struct {
	r         *bufio.Reader
	rawBuffer []byte
}

func newReader(r io.Reader) *TsvParser {
	return &TsvParser{
		r: bufio.NewReader(r),
	}
}

func (r *TsvParser) readLine() ([]byte, error) {
	line, err := r.r.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		r.rawBuffer = append(r.rawBuffer[:0], line...)
		for err == bufio.ErrBufferFull {
			line, err = r.r.ReadSlice('\n')
			r.rawBuffer = append(r.rawBuffer, line...)
		}
		line = r.rawBuffer
	}
	// drop trailing \n
	if n := len(line); n > 0 && line[n-1] == '\n' {
		line = line[:n-1]
	}
	// drop trailing \r
	if n := len(line); n > 0 && line[n-1] == '\r' {
		line = line[:n-1]
	}
	return line, err
}

func (r *TsvParser) Read() (record []string, err error) {
	line, errRead := r.readLine()
	if errRead != nil && errRead != io.EOF {
		return nil, errRead
	}
	return strings.Split(string(line), "\t"), errRead
}

