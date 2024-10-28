package gorel

type Page interface {
	DecodeFrom([]byte)
	Content() []byte
}
