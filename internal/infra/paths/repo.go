package paths

// type Repo interface {
// 	Load() (*Paths, error)
// 	Save(*Paths) error
// }
//
// type pathsRepo struct {
// 	file string
// }
//
// // NewRepo creates a new paths repo
// func NewRepo() *pathsRepo {
// 	return &pathsRepo{
// 		file: filepath.Join(rootDir(), "paths.json"),
// 	}
// }
//
// // Init initializes paths with its repo
// func (r *pathsRepo) Init() (*Paths, error) {
// 	p, err := r.Load()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if p != nil {
// 		return p, nil
// 	}
//
// 	mcDir, err := detectOrPromptMinecraftDir()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	root := filepath.Dir(r.file)
// 	p = New(root, mcDir)
//
// 	if err := r.Save(p); err != nil {
// 		return nil, err
// 	}
//
// 	return p, nil
// }
//
// // Load loads paths from paths.json
// func (r *pathsRepo) Load() (*Paths, error) {
// 	data, err := os.ReadFile(r.file)
// 	if os.IsNotExist(err) {
// 		return nil, nil // first run
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var p Paths
// 	if err := json.Unmarshal(data, &p); err != nil {
// 		return nil, err
// 	}
// 	return &p, nil
// }
//
// // Save saves paths to paths.json
// func (r *pathsRepo) Save(p *Paths) error {
// 	data, err := json.MarshalIndent(p, "", " ")
// 	if err != nil {
// 		return nil
// 	}
//
// 	if err := os.MkdirAll(filepath.Dir(r.file), 0755); err != nil {
// 		return err
// 	}
//
// 	return os.WriteFile(r.file, data, 0644)
// }
