//go:build generate

package generate

//go:generate go tool glua/tools/generate_ucd_tables -i testdata/DerivedCoreProperties.txt -p tables -o tables.go XID_Start XID_Continue
