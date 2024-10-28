package gorel

type Page interface {
	DecodePageFrom([]byte)
	Content() []byte
}
