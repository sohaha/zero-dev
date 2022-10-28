package tree

import (
	"fmt"
	"testing"

	"zlsapp/grbac/pkg/path"

	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

type RecordCase struct {
	args []string
	data interface{}
}

type QueryCase struct {
	args []string
	data interface{}
}

var (
	TestTree      *Tree
	TestQueryCase []QueryCase

	BenchTree           *Tree
	BenchForeachRecords []RecordCase
	BenchQueryCase      []QueryCase
)

func init() {
	defaultRecordCase := []RecordCase{
		{[]string{"*", "**", "*"}, "global category"},
		{[]string{"api-{prod,sit}.domain.com", "/article", "*"}, "article global category"},
		{[]string{"api-{prod,sit}.domain.com", "/article", "GET"}, "article get category"},
		{[]string{"api-{prod,sit}.domain.com", "/article", "POST"}, "article post category"},
		{[]string{"api-{prod,sit}.domain.com", "/article", "DELETE"}, "article delete category"},
		{[]string{"api-{prod,sit}.domain.com", "/login", "*"}, "login category"},
		{[]string{"api-{prod,sit}.domain.com", "/notice", "*"}, "notice category"},
		{[]string{"api-{prod,sit}.domain.com", "/query/*", "GET"}, "query category"},
		{[]string{"domain.com", "/login", "*"}, "login category"},
	}
	defaultQueryCase := []QueryCase{
		{[]string{"api-prod.domain.com", "/article", "GET"}, []interface{}{"global category", "article global category", "article get category"}},
		{[]string{"api-sit.domain.com", "/article", "DELETE"}, []interface{}{"global category", "article global category", "article delete category"}},
		{[]string{"api.domain.com", "/article", "POST"}, []interface{}{"global category"}},
		{[]string{"api-prod.domain.com", "/query/keywords", "GET"}, []interface{}{"global category", "query category"}},
	}

	TestTree = NewTree()
	TestQueryCase = defaultQueryCase
	for _, testCase := range defaultRecordCase {
		TestTree.Insert(testCase.args, testCase.data)
	}

	BenchForeachRecords = defaultRecordCase
	for i := 0; i < 1000; i++ {
		BenchForeachRecords = append(BenchForeachRecords, RecordCase{
			args: []string{
				"api-{prod,sit}.domain.com",
				"/" + faker.FirstName(),
				"*",
			},
		}, RecordCase{
			args: []string{
				faker.FirstName() + ".domain.com",
				"/" + faker.FirstName(),
				"GET",
			},
		}, RecordCase{
			args: []string{
				faker.FirstName(),
				fmt.Sprintf("%s/%s/%s/", faker.FirstName(), faker.FirstName(), faker.FirstName()),
				"GET",
			},
		})
	}

	BenchTree = NewTree()
	BenchQueryCase = defaultQueryCase
	for _, benchCase := range BenchForeachRecords {
		BenchTree.Insert(benchCase.args, benchCase.data)
	}
}

func TestTree_Query(t *testing.T) {
	tree := NewTree()
	conditions := []string{"layer1", "layer2", "layer3"}
	tree.Insert(conditions, "data1")
	tree.Insert(conditions, "data2")
	tree.Insert(conditions, "data3")
	tree.Insert(conditions, "data4")
	tree.Insert(conditions, "data5")
	data, err := tree.Query(conditions)
	assert.Equal(t, nil, err)
	assert.Equal(t, []interface{}{"data1", "data2", "data3", "data4", "data5"}, data)
}

func BenchmarkTree_Query(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, testCase := range BenchQueryCase {
				BenchTree.Query(testCase.args)
			}
		}
	})
}

func BenchmarkTree_Foreach_Query(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, treeCase := range BenchForeachRecords {
				for _, queryCase := range BenchQueryCase {
					for i, arg := range queryCase.args {
						matched := path.Match(treeCase.args[i], arg)
						if !matched {
							break
						}
					}
				}
			}
		}
	})
}
