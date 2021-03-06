==============
mixer.versions
==============

--------------------------------
Manage mix and upstream versions
--------------------------------

:Copyright: \(C) 2018 Intel Corporation, CC-BY-SA-3.0
:Manual section: 1


SYNOPSIS
========

``mixer versions``

``mixer versions [command]``


DESCRIPTION
===========

Manage mix and upstream versions. By itself the command will print the current
version of mix and upstream, and also report on the latest version of upstream
available. This command also allows the user to increment the mix version and
optionally the upstream version.


OPTIONS
=======

In addition to the globally recognized ``mixer`` flags (see ``mixer``\(1) for
more details), the following options are recognized.

-  ``-h, --help``

   Display subcommand help information and exit.


SUBCOMMANDS
===========

``update``

    Increment the mix version to generate a new release. By default the mix
    version is incremented by 10, following Clear Linux conventions to leave
    room for intermediate versions if necessary. The increment can be configured
    with the `--increment` flag. The ``update`` command also allows the user to
    update the upstream version.

    In addition to the global options ``mixer versions update`` takes the
    following options.

    - ``--clear-version {version}``

      Alias to `--upstream-version`

    - ``-c, --config {path}``

      Supply the `path` to the configuration file to use.

    - ``--increment {number}``

      Amount to increment the current mix version (default is 10).

    - ``--mix-version {version}``

      Set a specific mix `version`.

    - ``--upstream-version {version}``

      Set a specific next upstream version (either version number or the default
      "latest" string).

      This command will not update to an upstream version of a different format
      ("format bumps"). At the moment this needs to be handled manually.


EXIT STATUS
===========

On success, 0 is returned. A non-zero return code indicates a failure.

SEE ALSO
--------

* ``mixer``\(1)
