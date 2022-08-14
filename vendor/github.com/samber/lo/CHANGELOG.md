# Changelog

## 1.3.0 (2022-03-03)

Last and Nth return errors

## 1.2.0 (2022-03-03)

Adding `lop.Map` and `lop.ForEach`.

## 1.1.0 (2022-03-03)

Adding `i int` param to `lo.Map()`, `lo.Filter()`, `lo.ForEach()` and `lo.Reduce()` predicates.

## 1.0.0 (2022-03-02)

*Initial release*

Supported helpers for slices:

- Filter
- Map
- Reduce
- ForEach
- Uniq
- UniqBy
- GroupBy
- Chunk
- Flatten
- Shuffle
- Reverse
- Fill
- ToMap

Supported helpers for maps:

- Keys
- Values
- Entries
- FromEntries
- Assign (maps merge)

Supported intersection helpers:

- Contains
- Every
- Some
- Intersect
- Difference

Supported search helpers:

- IndexOf
- LastIndexOf
- Find
- Min
- Max
- Last
- Nth

Other functional programming helpers:

- Ternary (1 line if/else statement)
- If / ElseIf / Else
- Switch / Case / Default
- ToPtr
- ToSlicePtr

Constraints:

- Clonable
