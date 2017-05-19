echo "getting external http dependancies..."
go get github.com/julienschmidt/httprouter

echo getting external json parsing dependancies...
go get github.com/mailru/easyjson/jlexer
go get github.com/pquerna/ffjson/ffjson
go get github.com/buger/jsonparser

echo "getting external tool for interactive go dev (not a requirememnt...)"
go get github.com/mkouhei/gosh

echo getting mongo db drivers
go get -v gopkg.in/mgo.v2

echo getting benchmarking tools...
go get github.com/rakyll/hey

echo getting validation packages
go get gopkg.in/go-playground/validator.v9
