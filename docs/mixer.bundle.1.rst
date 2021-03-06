============
mixer.bundle
============

------------------------------------------------
Perform various configuration actions on bundles
------------------------------------------------

:Copyright: \(C) 2018 Intel Corporation, CC-BY-SA-3.0
:Manual section: 1


SYNOPSIS
========

``mixer bundle [command] [flags]``


DESCRIPTION
===========

Performs various configuration actions on upstream and local bundle definition
files. List and validate mix bundles. Validate local bundle definition files.


SUBCOMMANDS
===========

``add {bundle} [{bundle}...] [flags]``

    Add local or upstream bundles to your mix by modifying the mix bundles list
    at `<mixer/workspace>/mixbundles`. In addition to the global options ``mixer
    bundle add`` takes the following options.

    - ``--all-local``

      Add all local bundles defined in `<mixer/workspace>/local-bundles` to the
      mix. This command takes precedence over bundle list.

    - ``--all-upstream``

      Add all upstream bundles (cached in `<mixer/workspace>/upstream-bundles`
      to the mix. This command takes precedence over bundle list.

    - ``-c, --config {path}``

      Optionally tell ``mixer`` to use the configuration file at `path`. Uses
      the default `builder.conf` in the mixer workspace if this option is not
      provided.

    - ``--git``

      Automatically apply a new git commit after modifying the mix bundles list
      or bundle definitions. This is useful if your local bundles are kept under
      version control.

    - ``-h, --help``

      Display ``bundle add`` help information and exit.

``create``

    Create new bundles or copy existing bundles.
    This command will locate the bundle by first looking in local-bundles,
    and then in upstream-bundles. If the bundle is only found upstream,
    the bundle file will be copied to your local-bundles directory. If the bundle is
    not found anywhere, a blank template will be created with the correct name.

    Passing '--add' will also add the bundle(s) to your mix. Please note that
    bundles are added after all bundles are created, and thus will not be added
    if any errors are encountered earlier on.

    In addition to the global options ``mixer bundle create`` takes the following
    options.

    - ``--add``

      Add the bundle(s) to your mix after editing.

    - ``-c, --config {path}``

      Optionally tell ``mixer`` to use the configuration file at `path`. Uses
      the default `builder.conf` in the mixer workspace if this option is not
      provided.

    - ``--git``

      Automatically apply a new git commit after modifying the mix bundles list
      or bundle definitions. This is useful if your local bundles are kept under
      version control.

    - ``-h, --help``

      Display ``bundle create`` help information and exit.

``list [mix|local|upstream] [flags]``

    List the bundles in the mix, the available local bundles, or the available
    upstream bundles. In addition to the global options ``mixer bundle list``
    takes the following options.

    - ``mix``

      List the bundles in the mix.

    - ``local``

      List available locally-defined bundles.

    - ``upstream``

      List available upstream bundles.

    - ``-c, --config {path}``

      Optionally tell ``mixer`` to use the configuration file at `path`. Uses
      the default `builder.conf` in the mixer workspace if this option is not
      provided.

    - ``-h, --help``

      Display ``bundle list`` help information and exit.

    - ``--tree``

      Pretty-print the bundle list as a tree showing include information.

``remove``

    Removes bundles from your mix by modifying the mix bundle list (stored in
    the `<mixer/workspace>/mixbundles` file). The mix bundle list is parsed, the
    bundles are removed, and the resultant list is written back out in sorted
    order. If bundles do not exist in the mix, they are skipped. In addition to
    the global options ``mixer bundle remove`` takes the following options.

    - ``-c, --config {path}``

      Optionally tell ``mixer`` to use the configuration file at `path`. Uses
      the default `builder.conf` in the mixer workspace if this option is not
      provided.

    - ``--git``

      Automatically apply a new git commit after modifying the mix bundles list
      or bundle definitions. This is useful if your local bundles are kept under
      version control.

    - ``-h, --help``

      Display ``bundle remove`` help information and exit.

    - ``--local``

      Also remove the bundle file from local-bundles. This action is
      irrevocable.

    - ``--mix={bool}``

      Remove bundle from the mix bundle list. This defaults to true.

``validate``

    Checks bundle definition files for validity. Only local bundle files are
    checked; upstream bundles are trusted as valid. Valid bundles yield no
    output. Any invalid bundles will yield a non-zero return code.

    Basic validation includes checking syntax and structure, and that the bundle
    has a valid name. Commands like ``mixer bundle add`` run basic validation
    automatically.

    In addition to the global options ``mixer bundle remove`` takes the
    following options.

    - ``--all-local``

      Run validation against all local bundles.

    - ``-c, --config {path}``

      Optionally tell ``mixer`` to use the configuration file at `path`. Uses
      the default `builder.conf` in the mixer workspace if this option is not
      provided.

    - ``-h, --help``

      Display ``bundle validate`` help information and exit.

    - ``--strict``

      Perform strict validation to additionally check that the bundle header
      fields are parse-able and non-empty, and that the header 'Title' is itself
      valid and matches the bundle filename.


EXIT STATUS
===========

On success, 0 is returned. A non-zero return code indicates a failure.

SEE ALSO
--------

* ``mixer``\(1)
