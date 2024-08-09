# Shamir Secret Sharing

This is a practical example of [Shamir's secret sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing).

This tool allows you to divide up a secret string `S` into `n` shares where any `k` of the shares will allow you to reconstruct the original secret `S`.
If you have fewer than `k` of the shares, you will have no information (other than the length of the secret) about what that secret is.

## Usage

To distribute the secrets, use the command

`shamir distribute -S "<secret>" -n <number of shares to produce> -k <minimum number of shares to reconstruct the secret> -d "<output directory>"`

The shares will be saved as text files in the specified output directory.
The default directory is the current directory.

Way we run `shamir distribute -S="This is a super secret secret." -n=5 -k=3`.
Five text files will be saved in the current directory, each containing a share of the secret.
The first file might be called `shamir-IW6CPGV5COQ6FHMI-11d-1.txt`, and it contains the text

``` text
shamir-IW6CPGV5COQ6FHMI-11d-1-IyJneCEokIwf319Sl+sqeuLtBLwqWMxyd/m5Hj/U
```

The prefix `shamir-` indicates that this is a share in Shamir's secret sharing scheme.
`IW6CPGV5COQ6FHMI` is a randomly generated ID that allows you to correlate secrets and shares.
`11d` is the primitive polynomial used to construct the underlying Galois field.
`1` is the x coordinate of the share.
The remaining text is base64-encoded data.
Each byte corresponds to the value of a polynomial evaulated at the corresponding x-coordinate of the share.
(Note that each byte is encoded separately, each with a randomly-generated polynomial.)

To reconstruct the message, simply run `shamir reconstruct` in the same directory you ran `shamir distribute` command.
This will check any text files with the `.txt` file extension for shares.

Say we have the following text files:

``` text
shamir-IW6CPGV5COQ6FHMI-11d-1-IyJneCEokIwf319Sl+sqeuLtBLwqWMxyd/m5Hj/U
shamir-IW6CPGV5COQ6FHMI-11d-4-lKF0rNRgjQGXKJ0gbfKDMompvM6KGtAPw5TGAUGc
shamir-IW6CPGV5COQ6FHMI-11d-5-4+t6p9Uhbq3p17EHinzbaBgh2wDFNjwO0Q4Negpm
```

We then run the command `shamir reconstruct`, and we see that the secret is successfully reconstructed and printed to the console:

``` text
IW6CPGV5COQ6FHMI: This is a super secret secret.
```

If we try to only use shares 4 and 5, we cannot reconstruct the message, and we get gibberish:

``` text
IW6CPGV5COQ6FHMI: U�L��y&�r�-����G��=ѫ�G
����pS
```

Note that the file names themselves do not matter.
All of the information needed to reconstruct them is provided in the file itself.
