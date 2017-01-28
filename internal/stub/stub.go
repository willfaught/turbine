/// {{$r := .Ident.Initial}}
/// {{$t := .Ident.Name}}

package __FIELD_Ident___

/// {{if .InterfacePackageName}}import "{{.InterfacePackageName}}"{{end}}

/// {{if false}}
type __VAR_t__ struct{} /// {{end}}

/// {{define "params"}}{{range $i, $p := .Type.Params}}{{if $i}}, {{end}}{{range $j, $n := $p.Idents}}{{if $j}}, {{end}}{{$n.Name}}{{end}} {{$p.Type}}{{end}}{{end}}
/// {{define "results"}}{{if .Type.Results}}{{if or (gt (len .Type.Results) 1) (index .Results 0).Idents}}({{end}}{{end}}{{range $i, $r := .Results}}{{if $i}}, {{end}}{{range $j, $n := $r.Idents}}{{if $j}}, {{end}}{{$n.Name}}{{end}} {{$r.Type}}{{end}}{{if .Results}}{{if or (gt (len .Results) 1) (index .Results 0).Idents}}){{end}}{{end}}{{end}}

/// {{range .Methods}}
/// {{range .Doc}}
/// {{.}}{{end}}
func (__VAR_r__ __VAR_t__) __FIELD_Ident_Name__( /** {{template "params" .}} **/ ) /** {{template "results" .}} **/ {
	panic("not implemented")
}

/// {{end}}
