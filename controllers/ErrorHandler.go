package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

func init() {
	beego.ErrorHandler("403", Error403)
	beego.ErrorHandler("404", Error404)
	beego.ErrorHandler("501", Error501)
}

func Error403(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("403 Forbidden"))
}

func Error404(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("404 Not Found!"))
}

func Error501(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("501 Not Implemented"))
}
