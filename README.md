# hy: hierarchical files

## TODO
- Add options for default path names:
  - lowerCamelCase
  - CamelCase
  - snake_case
  - lowercaseonly
- Add support for reading FileTargets.
- Add support for auto-filling ID fields in map/slice elements on read.
  - Default field:  ID string
  - Default getter: ID() string
  - Default setter: SetID(string)
- On write, need to pick:
  - Fail if ID field not matching key or index?
  - Overwrite ID with current key or index?
  - Elide ID field from output altogether? (This should be the default, so
    it only matters in memory.)
  - Other?
- Add support for writing special maps with default fields/methods:
- Add support for writing actual files with a marshaller.
- Add support for reading actual files with a marshaller.
