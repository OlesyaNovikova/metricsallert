package main

import (
	"fmt"
	"net/http"
	"text/template"
)

const tplStr = `<table>
<thead>
	<tr>
		<th>Metric</th>
		<th>Value</th>
	</tr>
</thead>
<tbody>
	{{range $name, $value := . }}
		<tr>
			<td>{{ $name }}</td>
			<td>{{ $value }}</td>
		</tr>
	{{ end }}
</tbody>
</table>`

func getAllMems(res http.ResponseWriter, req *http.Request) {
	fmt.Print("Run getAllMems:\n")
	if req.Method != http.MethodGet {
		fmt.Print("Only GET requests are allowed!\n")
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	table := MemBase.S.GetAll()
	tpl, err := template.New("table").Parse(tplStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tpl.Execute(res, table)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	//res.WriteHeader(http.StatusOK)
}
