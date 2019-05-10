/// {{$p := .ReceiverTypePointer}}
/// {{$r := .ReceiverIdentifier}}
/// {{$t := .ReceiverTypeName}}

package __FIELD_ReceiverTypePackageName__

/// {{if .InterfacePackageName}}import "{{.InterfacePackageName}}"{{end}}

/// {{if false}}
type __FIELD_InterfaceName__ interface{}
type __FIELD_InterfacePackageQualifier____FIELD_InterfaceName__ interface{}
type __FIELD_ReceiverTypeReference____VAR_t__ struct{}
type __FIELD_ReceiverTypeName__ struct{}
type __VAR_p____VAR_t__ struct{} /// {{end}}

/// {{if false}}
func (__VAR_p____VAR_t__) Called() int { return 0 } /// {{end}}

/// {{if .InterfacePackageQualifier}}
var _ __FIELD_InterfacePackageQualifier____FIELD_InterfaceName__ = __FIELD_ReceiverTypeReference____VAR_t__{} /// {{else}}
var _ __FIELD_InterfaceName__ = &__FIELD_ReceiverTypeName__{}                                                 /// {{end}}

/// {{define "params"}}{{range $i, $p := .ParamsGrouped}}{{if $i}}, {{end}}{{range $j, $n := $p.Names}}{{if $j}}, {{end}}{{$n}}{{end}} {{$p.Type}}{{end}}{{end}}
/// {{define "results"}}{{if .ResultsGrouped}}{{if or (gt (len .ResultsGrouped) 1) (index .ResultsGrouped 0).Names}}({{end}}{{end}}{{range $i, $r := .ResultsGrouped}}{{if $i}}, {{end}}{{range $j, $n := $r.Names}}{{if $j}}, {{end}}{{$n}}{{end}} {{$r.Type}}{{end}}{{if .ResultsGrouped}}{{if or (gt (len .ResultsGrouped) 1) (index .ResultsGrouped 0).Names}}){{end}}{{end}}{{end}}
/// {{define "args"}}{{range $i, $p := .ParamsFlat}}{{if $i}}, {{end}}{{$p.Name}}{{end}}{{end}}
/// {{define "body-none"}}{{.Interface.ReceiverIdentifier}} .Called( {{template "args" .}} ) {{end}}
/// {{define "body-one"}}return {{.Interface.ReceiverIdentifier}} .Called( {{template "args" .}} ) {{$t := (index .ResultsFlat 0).Type}} {{if ne $t "interface{}"}} . ( {{$t}} ) {{end}} {{end}}
/** {{define "body-many"}} {
	var result = {{.Interface.ReceiverIdentifier}} .Called( {{template "args" .}} )

	return {{range $i, $r := .ResultsFlat}} {{if $i}} , {{end}} result.Get( {{$i}} ) {{if ne $r.Type "interface{}"}} . ( {{$r.Type}} ) {{end}} {{end}}
} {{end}} **/
/// {{define "body"}} {{$n := len .ResultsFlat}} {{if eq $n 0}} {{template "body-none" .}} {{else if eq $n 1}} {{template "body-one" .}} {{else}} {{template "body-many" .}} {{end}} {{end}}

/// {{range .InterfaceMethods}}
/// {{range .Doc}}
/// {{.}}{{end}}
func (__VAR_r__ __VAR_p____VAR_t__) __FIELD_Name__( /** {{template "params" .}} **/ ) /** {{template "results" .}} **/ {
	/// {{template "body" .}}
}

/// {{end}}
