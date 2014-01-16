DailyMile.com History Puller
============================

Uses the Go compiler (http://golang.org/)

The command line help isn't working yet so here's a quick run-down...

The first time running it, the user needs to be set using this...

    dm -u <username> user <username>

After the first time, the default user can be changed using...

    dm user <username>

Pull default user's entries from dailymile.com...

    dm sync

The initial sync will probably take minutes depending on number of entries.  Future syncs will be incremental and should only take a few seconds.  It'll tell you if it was "Already up-to-date" or how many new entries it pulled down.  One thing it does NOT handle are deletes.  If you sync and then delete an entry on the website it will remain in your local copy of the data unless you delete and do a full sync. 

If you're on some flavor of unix, your data should be stored in ~/.dailymile_cli/<username>/entries.json.  You can use your favorite browser plugin or editor to view the JSON pretty-printed.  Or, if you have python 2.6+ installed...

    cat ~/.dailymile_cli/<username>/entries.json | python -mjson.tool | less

I also started adding basic search capability.  Still tinkering with this... not sure where or how far it's headed.  Not even sure it works correctly!  Here are the basics...

    dm find [-s start date] [-e end date] [-p regex pattern]

All of this year's entries in JSON:

    dm find -s 14/1/1

All of this year's entries in an abridged CSV format (needs work):

    dm find -s 14/1/1 -csv

All of this year's entries in a column layout (Linux):

    dm find -s 14/1/1 -csv | column -t -s,

All of last year's interval workouts in JSON:

    dm find -s 13/1/1 -e 13/12/31 -p "(?i)interval"

Search for patterns like "8 x 400", "10x800", "10 x 1600m":

    dm find -s 13/1/1 -e 13/12/31 -p "(?i)\d{1,2} *x *\d{3,4}"

Count ALL of your entries on Linux:

    dm find -csv | wc -l

The 'dm find' command can generate HTML for entries, formatted using a
user-defined template.  A simple example template, entries.html, is provided
with this source.  See http://golang.org/pkg/text/template for information
on the Go language template rules.

Format matching entries as HTML:

    dm find -s 14/1/1 -html entries.html
