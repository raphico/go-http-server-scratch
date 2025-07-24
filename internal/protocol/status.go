package protocol

const (
	StatusOk                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

var StatusText = map[int]string{
	StatusOk:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
}
