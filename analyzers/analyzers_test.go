package analyzers

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAnalyzersHaveUniqueNamesAndMetadata(t *testing.T) {
	seen := map[string]struct{}{}
	for _, analyzer := range allAnalyzersForUnitTests() {
		if analyzer.Name == "" {
			t.Fatal("analyzer name must not be empty")
		}
		if _, ok := seen[analyzer.Name]; ok {
			t.Fatalf("duplicate analyzer name %q", analyzer.Name)
		}
		seen[analyzer.Name] = struct{}{}
		if analyzer.Doc == "" {
			t.Fatalf("%s: doc must not be empty", analyzer.Name)
		}
		if analyzer.URL == "" {
			t.Fatalf("%s: URL must not be empty", analyzer.Name)
		}
		if analyzer.Run == nil {
			t.Fatalf("%s: Run must not be nil", analyzer.Name)
		}
	}
}

func TestEveryAnalyzerHasTestdata(t *testing.T) {
	for _, analyzer := range allAnalyzersForUnitTests() {
		path := filepath.Join("testdata", "src", analyzer.Name, analyzer.Name+".go")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("%s: expected testdata file %s: %v", analyzer.Name, path, err)
		}
	}
}

func allAnalyzersForUnitTests() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		NoPanic,
		ChanSize,
		EnumStart,
		NewRef,
		NilSlice,
		ZeroVar,
		StringBytes,
		GlobalPrefix,
		VarType,
		IfacePtr,
		AtomicStd,
		RawString,
		PublicEmbed,
		EmbedLayout,
		LocalVar,
		ZeroFields,
		MapInit,
		ConstPrintf,
		NakedParams,
		TimeField,
	}
}
