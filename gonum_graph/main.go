package main

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	GA "github.com/dominikbraun/graph"
	batch "github.com/zqlpaopao/tool/batch_select/pkg"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type Nodes struct {
	graph.Node
	Data NodeM
}
type Edges struct {
	simple.Edge // 嵌入加权边
	Data        EdgeM
}

type NodeM struct {
	ID        int       `json:"id" gorm:"column:id"`
	Vid       string    `json:"vid" gorm:"column:vid"` // vid
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (m *NodeM) TableName() string {
	return "node"
}

type EdgeM struct {
	ID        int       `json:"id" gorm:"column:id"`
	SrcVid    string    `json:"src_vid" gorm:"column:src_vid"` // vid
	DstVid    string    `json:"dst_vid" gorm:"column:dst_vid"` // vid
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (m *EdgeM) TableName() string {
	return "edge"
}

func main() {

	InitDB()
	InitData()
	fmt.Println("data end.....")

	// 计算最短路径
	t := time.Now()
	path, err := GA.ShortestPath(baseGraph, "0", "500")
	if err != nil {
		panic(err)
	}
	fmt.Printf("路径: %v time : %v ,\n", path, time.Now().Sub(t).String())

}

var (
	url = "root:meimima123@tcp(127.0.0.1:3306)/test?charset=utf8&interpolateParams=true&&parseTime=True&loc=Local&&readTimeout=300s&writeTimeout=300s"

	MyBaseCli = &gorm.DB{}
	baseGraph = GA.New(GA.StringHash, GA.Weighted())

	//Graph = pkg.NewGraphHandler[GA.Graph[string, float64], Nodes, Edges](
	//	baseGraph,
	//	9999,
	//	9999,
	//)
)

func InitDB() {
	var err error
	if MyBaseCli, err = gorm.Open(mysql.Open(url), &gorm.Config{}); nil != err {
		panic(err)
	}
	MyBaseCli = MyBaseCli.Debug()
	if err = setParams(); nil != err {
		panic(err)
	}
	return
}

// setParams 设置连接
func setParams() (err error) {
	var db *sql.DB
	if db, err = MyBaseCli.DB(); nil != err {
		panic(err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Duration(86400) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(86100) * time.Second)
	return
}

func InitData() {
	var (
		err  error
		b    = batch.NewBatchSelect[NodeM]()
		e    = batch.NewBatchSelect[EdgeM]()
		node = "node"
		edge = "edge"
	)
	if err = b.InitTask(batch.InitTaskModel[NodeM]{
		TaskName: node,
		Opt: []batch.OptionInter[NodeM]{
			batch.WithHandleGoNum[NodeM](10),
			batch.WithHandleRevGoNum[NodeM](10),
			batch.WithLimit[NodeM](10000),
			batch.WithOrderColumn[NodeM]("id"),
			batch.WithTable[NodeM]("node"),
			batch.WithSqlWhere[NodeM]("", nil),
			batch.WithResChanSize[NodeM](20000),
			batch.WithMysqlSqlCli[NodeM](MyBaseCli),
			batch.WithCallFunc[NodeM](tidyNode)}}); nil != err {
		fmt.Println(err)
		return
	}
	if err = b.Run(node); nil != err {
		return
	}
	b.Wait()

	//spew.Dump(baseGraph)
	//os.Exit(1)

	if err = e.InitTask(batch.InitTaskModel[EdgeM]{
		TaskName: edge,
		Opt: []batch.OptionInter[EdgeM]{
			batch.WithHandleGoNum[EdgeM](10),
			batch.WithHandleRevGoNum[EdgeM](10),
			batch.WithLimit[EdgeM](10000),
			batch.WithOrderColumn[EdgeM]("id"),
			batch.WithTable[EdgeM]("edge"),
			batch.WithSqlWhere[EdgeM]("", nil),
			batch.WithResChanSize[EdgeM](20000),
			batch.WithMysqlSqlCli[EdgeM](MyBaseCli),
			batch.WithCallFunc[EdgeM](tidyEdge)}}); nil != err {
		fmt.Println(err)
		return
	}
	if err = e.Run(edge); nil != err {
		fmt.Println(err)
		return
	}
	e.Wait()
}

func tidyNode(data *[]NodeM) {
	for _, node := range *data {
		err := baseGraph.AddVertex(node.Vid)
		if err != nil {
			panic(err)
		}
		//Graph.GraphLock.Lock()
		//item := Nodes{
		//	Node: Graph.Graph.NewNode(),
		//	Data: node,
		//}
		//Graph.Graph.AddNode(item)
		//Graph.GraphLock.Unlock()
		//
		//Graph.NodeLock.Lock()
		//Graph.NodeCache[node.Vid] = item
		//Graph.NodeLock.Unlock()
	}

}

func tidyEdge(data *[]EdgeM) {
	for _, node := range *data {
		err := baseGraph.AddEdge(node.SrcVid, node.DstVid)
		if err != nil {
			panic(err)
		}
	}
}

func test() {
	WeightGraph()
	//graphData()
}

func GraphSimple() {
	g := simple.NewUndirectedGraph()

	nodes := make([]graph.Node, 4)
	for i := range nodes {
		nodes[i] = g.NewNode()
		g.AddNode(nodes[i])
	}

	g.SetEdge(g.NewEdge(nodes[0], nodes[1]))
	g.SetEdge(g.NewEdge(nodes[0], nodes[2]))
	g.SetEdge(g.NewEdge(nodes[1], nodes[3]))
	g.SetEdge(g.NewEdge(nodes[2], nodes[3]))

	shortest := path.DijkstraAllPaths(g)
	paths, _ := shortest.AllBetween(nodes[3].ID(), nodes[0].ID())

	fmt.Println("所有最短路径:")
	for i, p := range paths {
		fmt.Printf("路径%d: ", i+1)
		for _, n := range p {
			fmt.Printf("%d ", n)
		}
		fmt.Println()
	}

}

func WeightGraph() {
	//simple.NewDirectedGraph()
	//simple.NewUndirectedGraph()
	//
	//simple.NewWeightedDirectedGraph()
	//simple.NewWeightedUndirectedGraph()
	//
	//simple.NewDirectedMatrix()
	//simple.NewUndirectedMatrix()
	//
	//simple.NewDirectedMatrixFrom()
	//simple.NewUndirectedMatrixFrom()

	g := simple.NewWeightedDirectedGraph(8, 30)

	// 添加节点（示例使用4个节点）
	nodes := make([]graph.Node, 4)
	for i := range nodes {
		nodes[i] = g.NewNode()
		spew.Dump(g.NewNode())
		//os.Exit(1)
		g.AddNode(nodes[i])
	}

	// 设置带权边（构造两条不同路径）
	g.SetWeightedEdge(g.NewWeightedEdge(nodes[0], nodes[1], 1)) // 0→1 (权重1)
	g.SetWeightedEdge(g.NewWeightedEdge(nodes[0], nodes[2], 2)) // 0→2 (权重2)
	g.SetWeightedEdge(g.NewWeightedEdge(nodes[1], nodes[3], 3)) // 1→3 (权重3)
	g.SetWeightedEdge(g.NewWeightedEdge(nodes[2], nodes[3], 2)) // 2→3 (权重1)

	// 计算所有最短路径
	shortest := path.DijkstraAllPaths(g)
	paths, totalWeight := shortest.AllBetween(nodes[0].ID(), nodes[3].ID())

	// 打印所有路径详情
	fmt.Println("===== 所有最短路径 =====")
	for i, p := range paths {
		fmt.Printf("路径%d: ", i+1)
		var pathWeight float64
		for j := 0; j < len(p); j++ {
			fmt.Printf("%d", p[j].ID())
			if j < len(p)-1 {
				fmt.Print("→")
				// 计算当前边的权重
				edge := g.WeightedEdge(p[j].ID(), p[j+1].ID())
				pathWeight += edge.Weight()
			}
		}
		fmt.Printf(" | 总权重: %.1f (", pathWeight)
		// 显示权重组成
		for j := 0; j < len(p)-1; j++ {
			edge := g.WeightedEdge(p[j].ID(), p[j+1].ID())
			fmt.Printf("%.1f", edge.Weight())
			if j < len(p)-2 {
				fmt.Print("+")
			}
		}
		fmt.Println(")")
	}
	fmt.Printf("\n算法确认的最短路径权重: %.1f\n", totalWeight)

}

// 自定义节点类型
type MyNode struct {
	graph.Node        // 嵌入基础Node
	Name       string // 自定义字段
}

// 自定义边类型
type MyEdge struct {
	simple.WeightedEdge        // 嵌入加权边
	Label               string // 自定义字段
}

func graphData() {
	// 创建带权有向图，权重初始值为+Inf
	g := simple.NewWeightedDirectedGraph(0, math.Inf(1))

	// 1. 添加带标签的节点
	nodes := make([]*MyNode, 4)
	for i := range nodes {
		nodes[i] = &MyNode{
			Node: g.NewNode(), // 自动生成ID
			Name: fmt.Sprintf("NodeM%d", i),
		}
		g.AddNode(nodes[i])
	}

	// 2. 添加带标签的加权边
	edges := []struct {
		from, to int
		weight   float64
		label    string
	}{
		{0, 1, 1.0, "EdgeM-A"},
		{0, 2, 2.0, "EdgeM-B"},
		{1, 3, 3.0, "EdgeM-C"},
		{2, 3, 2.0, "EdgeM-D"},
	}

	for _, e := range edges {
		edge := &MyEdge{
			WeightedEdge: simple.WeightedEdge{
				F: nodes[e.from],
				T: nodes[e.to],
				W: e.weight,
			},
			Label: e.label,
		}
		g.SetWeightedEdge(edge)
	}

	//g.RemoveEdge()
	// 3. 计算所有最短路径
	shortest := path.DijkstraAllPaths(g)
	paths, totalWeight := shortest.AllBetween(nodes[0].ID(), nodes[3].ID())

	// 4. 打印带自定义数据的路径
	fmt.Printf("从 %s 到 %s 的最短路径权重: %.1f\n",
		nodes[0].Name, nodes[3].Name, totalWeight)

	for i, path := range paths {
		fmt.Printf("路径 %d:\n", i+1)
		for j, node := range path {
			// 类型断言获取自定义节点
			current := node.(*MyNode)
			fmt.Printf("  %s", current.Name)

			if j > 0 {
				prev := path[j-1].ID()
				// 获取边数据（需二次类型断言）
				edge := g.WeightedEdge(prev, node.ID()).(*MyEdge)
				fmt.Printf(" --[%s|%.1f]-->", edge.Label, edge.Weight())
			}
		}
		fmt.Println("\n" + strings.Repeat("-", 30))
	}
}
