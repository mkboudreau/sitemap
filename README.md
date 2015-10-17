# Sitemap Generator

## Summary
This little utility, given a starting URL, will crawl a website and find all the static assets and links on that site. 

## Design Goals
- Crawl an entire site and report on its structure
- Flexible output formats (i.e. json, tab, digraph)
- Customize performance characteristics 

## Design Decisions
- The utility will stay within the same domain
- THe utility, when it finds duplicate URLs, it will not traverse into its links, but still report on the links found.

## Features
- Ability to save results to a file
- Set number of worker threads/goroutines to crawl a site
- Set rate limiter, if desired
- Set inactivity timeout
- Read in saved results and redisplay in different formats

## How to get it

### (1) You have Docker installed
```bash
docker run mkboudreau/sitemap ....
```

### (2) You have Go installed
```bash
go get github.com/mkboudreau/sitemap 
make install
```

## Example Usage

### Crawl site with sensible defaults
`sitemap www.microsoft.com` 

### Crawl site with 50 workers
`sitemap -w 50 www.microsoft.com` 

### Crawl site with rate limiting turned off 
`sitemap -r 0s www.microsoft.com` 

### Crawl site and output JSON
`sitemap -f json www.microsoft.com` 

### Crawl site and output tabular format (default)
`sitemap -f tab www.microsoft.com` 

### Crawl site and output digraph (dot)
`sitemap -f digraph www.microsoft.com` 

### Crawl site and save results to file
`sitemap -o saved.json www.microsoft.com` 

### Use saved results and output as a digraph
`sitemap -i saved.json -f digraph` 

