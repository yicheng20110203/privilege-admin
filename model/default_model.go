package model

type DefaultModel struct {

}

const (
	DefaultPage = 1
	DefaultSize = 10
)

func GetDefaultPage() int32 {
	return DefaultPage
}

func GetDefaultSize() int32 {
	return DefaultSize
}