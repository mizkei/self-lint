package lint

import (
	"io"

	"gopkg.in/yaml.v2"
)

const (
	globalTarget = "global"
)

type Matter struct {
	Target []string            `yaml:"target"`
	Import []string            `yaml:"import"`
	Ref    map[string][]string `yaml:"ref"`
	Write  []string            `yaml:"write"`
}

type Config struct {
	ProhibitedMatter []Matter `yaml:"prohibited-matter"`
}

func (c Config) matterByTarget() map[string]*Matter {
	targetMap := make(map[string]*Matter)
	for _, mt := range c.ProhibitedMatter {
		for _, t := range mt.Target {
			if _, ok := targetMap[t]; !ok {
				targetMap[t] = &Matter{Ref: make(map[string][]string)}
			}
			targetMap[t].Import = append(targetMap[t].Import, mt.Import...)
			targetMap[t].Write = append(targetMap[t].Write, mt.Write...)
			for pkg, obj := range mt.Ref {
				targetMap[t].Ref[pkg] = append(targetMap[t].Ref[pkg], obj...)
			}
		}
	}

	gm, ok := targetMap[globalTarget]
	for _, m := range targetMap {
		if ok {
			m.Import = append(m.Import, gm.Import...)
			m.Write = append(m.Write, gm.Write...)
			for pkg, obj := range gm.Ref {
				m.Ref[pkg] = append(m.Ref[pkg], obj...)
			}
		}
		m.Import, m.Write = uniq(m.Import), uniq(m.Write)
		for pkg := range m.Ref {
			m.Ref[pkg] = uniq(m.Ref[pkg])
		}
	}

	return targetMap
}

func uniq(sl []string) []string {
	sm := make(map[string]struct{})
	for _, s := range sl {
		sm[s] = struct{}{}
	}

	res := make([]string, 0, len(sm))
	for s := range sm {
		res = append(res, s)
	}
	return res
}

func LoadConfig(r io.Reader) (Config, error) {
	var conf Config
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
