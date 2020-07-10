# geoip

Simple command line utility to parse a local MaxMind GeoIP2 mmdb (leveraging the great [oschwald/geoip2-golang library](https://github.com/oschwald/geoip2-golang), 
and use it to enrich a CSV containing an IP address. It currently adds columns for `city`, `country`, `isAnonymousIP`,
and `isSatelliteProvider`. 

It's intentionally designed as a typical `*nix` cli tool, reading from `stdin` and writing to `stdout`. 

## Example

```sh
echo "8.8.8.8" | ./geoip                                                                                        
8.8.8.8,,United States,false,false
```

## Getting the GeoIP2 Database on Mac OS 

1. Sign up for the [free GeoLite2 DB access](https://www.maxmind.com/en/geolite2/signup)
2. Once you have an account, [configure a license key](https://www.maxmind.com/en/accounts/current/license-key) for the new database format
3. Once generated, either download the generated `.conf` file, or copy the license key and your account ID to a separate file
4. (Mac OS) `brew install geoipupdate`
5. (Optional) Configure `geoipupdate` by either copying the generated `.conf` file above to its default location (can be obtained with `-h`) or editing
the existing file with the copied values from step 3.
6. Run `geoipupdate`

Et voila! 
