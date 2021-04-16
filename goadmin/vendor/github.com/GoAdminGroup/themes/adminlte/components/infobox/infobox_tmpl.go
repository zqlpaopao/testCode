package infobox

var List = map[string]string{
	"infobox": `{{define "infobox"}}
<div class="info-box">
    {{if .IsHexColor}}
        <span class="info-box-icon" style="background-color: {{.Color}} !important;">
    {{else}}
        <span class="info-box-icon bg-{{.Color}}">
    {{end}}
    {{if .IsSvg}}
        {{.Icon}}
    {{else}}
        <i class="fa {{.Icon}}"></i>
    {{end}}
    </span>
    <div class="info-box-content">
        <span class="info-box-text">{{langHtml .Text}}</span>
        <span class="info-box-number">{{langHtml .Number}}</span>
        {{langHtml .Content}}
    </div>
</div>
{{end}}`,
}
