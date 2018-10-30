# Go Server Skeleton
A skeleton for making website servers in Go

My goal is to modify this so that it's a little be more semantically organized, and there is a clear delineation between dependencies/module type codeblocks and their examples

## Usage
Build and run server with:  
    `make && ./run`  
(you may need sudo access for ./run, depending on your ports)

## Versioning
Tags are used for versioning, to tag new version use:  
    `git tag -a <version>`  
[HINT: Use `git tag -l` to see past versions for reference]

If ever pushing with new tags, remember to:  
    `git push <remote> --tags`

To pull any new tags to local, use:  
    `git pull <remote>`

# Todo:

I basically ignored logging cause it's small and contained. Things aren't really built on it.
--- fix skeleton:  

- finish fleshing out tokens
		-- get it implemented (30 minutes)
- Get TLS in there  (30 minutes)
- rename packages (45 minutes)
- test it with curl and not https (30 min)
- did i miss something from keelins repo (30 min)
- look @ templates, whatever (45 minutes)
- document (30 min)
- pull request (20 minutes) // aj is done
- is error processing really sufficient? (we're done, keelin can take a look at this if he cares, integrate it with logging better)
- prevent collisions (key names, JWT, cookies, etc)
- whats the deal with the ClaimMaps and it's relationship with json
