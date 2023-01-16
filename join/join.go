package join

type Method struct {
	Between string
	Final   string
}

var CommaSeparatedJoin Method = Method{
	Between: ", ",
}

var ListJoin Method = Method{
	Between: ", ",
	Final:   " and ",
}
