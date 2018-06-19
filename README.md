# vt-engine-tools
Engine related tools

## Building

```
make
```

This will produce the vt-engine-tools binary.  For usage, use the `help` command, e.g. `./vt-engine-tools help`



## merge-vlf

https://steel-ventures.atlassian.net/browse/VTN-8503

has the use case to use this command.

For usage:

```
./vt-engine-tools merge-vlf help
```


### Example usage

```
./vt-engine-tools merge-vlf samples/86940638_0_v-vlf.json samples/86940660_0_v-vlf.json -o res/output.json
```



### Sample file generation

`samples` directory has scripts:

 * `run.sh` to split a dual channel MP3 file to individual WAV files, transcoded to 16K mono WAV format.  It then submit to Veritone for transcribing.
 * `check-job.sh` can then be used subsequently to get the VLF and TTML file.

The VLF files can then be used in `merge-vlf` to produce a single engine output with speaker identifications