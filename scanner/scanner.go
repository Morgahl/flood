package scanner

import (
	"bufio"
	"fmt"
	"os"
)

type Scanner struct {
	file  *os.File
	buf   *bufio.Reader
	board [][]uint8
}

func New(filepath string) (*Scanner, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReaderSize(file, 64*1024)

	return &Scanner{
		file: file,
		buf:  buf,
	}, nil
}

func (s *Scanner) Scan() bool {
	s.board = make([][]uint8, 19)
	for i := 0; i < 19; i++ {
		s.board[i] = make([]uint8, 19)
		if _, err := fmt.Fscanf(s.buf, "%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d%1d\n",
			&s.board[i][0], &s.board[i][1], &s.board[i][2], &s.board[i][3], &s.board[i][4],
			&s.board[i][5], &s.board[i][6], &s.board[i][7], &s.board[i][8], &s.board[i][9],
			&s.board[i][10], &s.board[i][11], &s.board[i][12], &s.board[i][13], &s.board[i][14],
			&s.board[i][15], &s.board[i][16], &s.board[i][17], &s.board[i][18]); err != nil {
			return false // this should only fail at the double newline at the end of the file
		}
	}
	fmt.Fscanln(s.buf) // scan empty line after file

	return true
}

func (s *Scanner) Board() [][]uint8 {
	return s.board
}

func (s *Scanner) Close() error {
	return s.file.Close()
}
