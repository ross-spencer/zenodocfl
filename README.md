# ZenodOFCL

Convert FHNW INK media server records to OCFL for publication on
Zenodo.

Consists of the following apps:

## Lister

Given a search term list records for appraisal.

### Allowlist

Lister's allowlist function provides finer grained ability to select the
records which will eventually be downloaded and added to Zenodo.

## Gather

Gather a list of ALL records and items based on the appraised selection.
Provides deduplication and a further check for appraisal.

### Allowlist

Gather accepts a further allowlist which lists specifically the records that
need to be downloaded for Zenodo and a further cross-check will be done against
this list for precision.

## Crater

Output a RO-Crate based on input data and optionally download the remainder
of the crate data.

## Example usage

Users of the ZenodOCFL workflow need to follow a basic workflow as follows:

1. Search by string `-search` or collection ID `-collection` if known, e.g.

```bash
./lister -collection 27 -o demo
```

2. The output will be two files:

```text
demo.manifest
demo.allowlist
```

3. Assuming no changes need to be made to the manifest, we will ignore the
allowlist in this example. Simply provide the manifest to `gather`:

```bash
./gather -download demo.manifest
```

All items will be output to a data folder.

4. A furhter RO-CRATE manifest must be created using `-list`. The RO-CRATE
manifest enables download of ALL further media items and anciliary records.

```bash
./gather -list -o demo
```

The new ro-crate manifest will be output:

```text
demo.collection
```

This can be provided to `crater` to complete the RO-CRATE generation process.

6. Create a metadata excerpt to be added to the RO-CRATE:

```json
{
  "identifier": "FHNW-01KHBM9PSRCH0KM7FGXSFCV5V8",
  "description": "Motet Cycles Research",
  "name": "Motet Cycles",
  "type": "Dataset",
  "data_published": "2018",
  "publisher": "https://ror.org/04mq2g308",
  "license": "https://creativecommons.org/publicdomain/zero/1.0/",
  "keywords": "renaissance, music, motet",
  "url": "https://ink.sammlung.cc/detail/motetcycles-research/"
}
```

> NB. Publisher should be provided by [ror.org][ror-1] if possible.

[ror-1]: https://ror.org/

5. Create RO-CRATE using the command line:

```bash
./crater -crate demo.collection-meta meta.json
```

> NB. a `-dry-run` flag is available to observe the output without downloading
the entire collection from INK.

6. Observe the output, e.g. for 10 records:

```text
output/
└── ro-crate-Motet-Cycles-1770997122
    ├── anciliary
    ├── media
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M001BeataProgenies.xml
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M002GloriosaeVirginis.xml
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M003SubTuamProtectionem.xml
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M004HortusConclusus.xml
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M006TotaPulchraEs.xml
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M015TuThronusEsSalomonis.xml
    │   ├── motet_cycles_data_test_M001.png
    │   ├── motet_cycles_data_test_motet_cycles_data_test_M002.png
    │   ├── motet_cycles_data_test_motet_cycles_data_test_M003.png
    │   └── motet_cycles_data_test_T001_Beata_progenies.pdf
    ├── posters
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M001BeataProgenies.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M002GloriosaeVirginis.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M003SubTuamProtectionem.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M004HortusConclusus.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M005DescendiInHortumMeum.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M006TotaPulchraEs.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M007OSacrumConvivium.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M008HocGaudiumEstSpiritus.png
    │   ├── motet_cycles_data_motet_cycles_data_MEI_files_motets_M015TuThronusEsSalomonis.png
    │   ├── motet_cycles_data_motet_cycles_data_motet_cycles_data.png
    │   ├── motet_cycles_data_test_M001.png
    │   ├── motet_cycles_data_test_motet_cycles_data_test_M002.png
    │   ├── motet_cycles_data_test_motet_cycles_data_test_M003.png
    │   └── motet_cycles_data_test_Website.png
    ├── records
    │   ├── motetcycle-0399.json
    │   ├── motetcycle-0400.json
    │   ├── motetcycle-0401.json
    │   ├── motetcycle-0955.json
    │   ├── motetcycle-0958.json
    │   ├── motetcycle-0960.json
    │   ├── motetcycle-0968.json
    │   ├── zotero2-2641719.5KZRPA2N.json
    │   ├── zotero2-2641719.R8Z7ZWIA.json
    │   └── zotero2-2641719.VF2JR42N.json
    └── ro-crate-metadata.json
```

## Preview

Preview is best done via the package [ro-crate-html][preview-1].

Two helper commands exist in the `justfile` to help with its use:

```just
    install-preview      # install ro-crate preview
    preview-rocrate dir  # preview ro-crate
```

[preview-1]: https://www.npmjs.com/package/ro-crate-html
